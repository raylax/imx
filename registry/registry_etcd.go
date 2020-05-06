package registry

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/raylax/imx/client"
	"github.com/raylax/imx/core"
	"go.etcd.io/etcd/clientv3"
	"log"
	"sync"
	"time"
)

const (
	dialTimeout = 5 * time.Second
	keepAlive   = 5
	keyPrefix   = "/imx"
	nodePrefix  = keyPrefix + "/node/"
	userPrefix  = keyPrefix + "/user/"
	groupPrefix = keyPrefix + "/group/"
)

type EtcdRegistry struct {
	endpoints []string
	node      core.Node
	cli       *clientv3.Client
	kv        clientv3.KV
	lease     clientv3.Lease
	leaseId   clientv3.LeaseID
	ctx       context.Context
	cancel    context.CancelFunc
	key       string
	users     map[string]core.User
	groups    map[string]core.Group
	userMux   sync.RWMutex
	groupMux  sync.RWMutex
}

func NewEtcdRegistry(endpoints []string, node core.Node) *EtcdRegistry {
	return &EtcdRegistry{
		endpoints: endpoints,
		node:      node,
	}
}

func (r *EtcdRegistry) Init() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   r.endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return err
	}
	r.kv = clientv3.NewKV(cli)
	r.lease = clientv3.NewLease(cli)
	r.cli = cli
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.key = nodePrefix + r.node.Key()
	r.users = make(map[string]core.User)
	return nil
}

func (r *EtcdRegistry) Reg() error {
	if err := r.reg(); err != nil {
		return err
	}
	r.watchNodes()
	return nil
}

func (r *EtcdRegistry) UnReg() {
	_, _ = r.cli.Revoke(r.ctx, r.leaseId)
	r.cancel()
}

func (r *EtcdRegistry) RegUser(u core.User) error {
	data, _ := json.Marshal(r.node)
	_, err := r.kv.Put(r.ctx, userPrefix+u.Id, string(data), clientv3.WithLease(r.leaseId))
	if err != nil {
		return err
	}
	r.userMux.Lock()
	r.users[u.Id] = u
	r.userMux.Unlock()
	return nil
}

func (r *EtcdRegistry) UnRegUser(u core.User) {
	r.userMux.Lock()
	r.users[u.Id] = u
	r.userMux.Unlock()
	_, _ = r.kv.Delete(r.ctx, userPrefix+u.Id)
}

func (r *EtcdRegistry) RegGroup(g core.Group, u core.User) error {
	_, err := r.kv.Put(r.ctx, groupPrefix+g.Id+"/"+u.Id, "", clientv3.WithLease(r.leaseId))
	if err != nil {
		return err
	}
	r.groupMux.Lock()
	r.groups[g.Id+"/"+u.Id] = core.Group{Users: []core.User{u}}
	r.groupMux.Unlock()
	return nil
}

func (r *EtcdRegistry) UnRegGroup(g core.Group, u core.User) {
	_, _ = r.kv.Delete(r.ctx, groupPrefix+g.Id+"/"+u.Id)
	r.groupMux.Lock()
	delete(r.groups, g.Id+"/"+u.Id)
	r.groupMux.Unlock()
}

func (r *EtcdRegistry) reRegGroups() {
	r.groupMux.RLock()
	for _, g := range r.groups {
		_ = r.RegGroup(g, g.Users[0])
	}
	r.groupMux.RUnlock()
}

func (r *EtcdRegistry) GetGroupUsers(gid string) ([]string, error) {
	resp, err := r.kv.Get(r.ctx, groupPrefix+gid+"/", clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}
	prefixLen := len(groupPrefix + gid + "/")
	users := make([]string, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		users[i] = string(kv.Key)[prefixLen:]
	}
	return users, nil
}

func (r *EtcdRegistry) LookupNodes(uid string) ([]core.Node, error) {
	resp, err := r.kv.Get(r.ctx, userPrefix+uid, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	nodes := make([]core.Node, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		nodes = append(nodes, core.NewNodeFromJSON(kv.Value))
	}
	return nodes, nil
}

// 注册服务
func (r *EtcdRegistry) reg() error {
	log.Printf("注册服务")
	err := r.reGrant()
	if err != nil {
		return err
	}
	if r.keepAlive() != nil {
		return err
	}
	if r.regNode() != nil {
		return err
	}
	log.Printf("服务注册完成")
	return nil
}

// 重新注册服务
func (r *EtcdRegistry) reReg() error {
	log.Printf("重新注册服务")
	r.ctx, r.cancel = context.WithCancel(context.Background())
	if err := r.reg(); err != nil {
		return err
	}
	r.regUsers()
	r.reRegGroups()
	log.Printf("服务重新注册完成")
	return nil
}

// 重新创建租约
func (r *EtcdRegistry) reGrant() error {
	log.Printf("创建租约")
	lease, err := r.lease.Grant(r.ctx, keepAlive)
	if err != nil {
		return err
	}
	r.leaseId = lease.ID
	log.Printf("租约创建完成 lease:%d", lease.ID)
	return nil
}

// 重新注册节点
func (r *EtcdRegistry) regNode() error {
	log.Printf("注册节点")
	val, err := json.Marshal(r.node)
	if err != nil {
		return err
	}
	resp, err := r.kv.Put(r.ctx, r.key, string(val), clientv3.WithLease(r.leaseId))
	if err != nil {
		return err
	}
	log.Printf("节点注册完成 revision:%d", resp.Header.Revision)
	return nil
}

// 重新注册用户
func (r *EtcdRegistry) regUsers() {
	r.userMux.RLock()
	for _, u := range r.users {
		_ = r.RegUser(u)
	}
	r.userMux.RUnlock()
}

func (r *EtcdRegistry) keepAlive() error {
	log.Printf("设置自动续约")
	respChan, err := r.lease.KeepAlive(r.ctx, r.leaseId)
	if err != nil {
		return err
	}
	// 监控租约到期
	go func() {
		for {
			select {
			case resp := <-respChan:
				if resp == nil {
					log.Printf("租约到期")
					// 租约到期重新注册
					_ = r.reReg()
					return
				}
			}
		}
	}()
	log.Printf("自动续约设置完成")
	return nil
}

func (r *EtcdRegistry) handleNodeChange(e *clientv3.Event) {
	node := core.NewNodeFromJSON(e.Kv.Value)
	// 排除自身节点
	if node.Endpoint() == r.node.Endpoint() {
		return
	}
	switch e.Type {
	// 监控到节点添加，添加RPC客户端
	case mvccpb.PUT:
		client.AddRpcClient(node)
	// 监控到节点删除，移除RPC客户端
	case mvccpb.DELETE:
		client.RemoveRpcClient(node)
	}
}

// 监控节点变化
func (r *EtcdRegistry) watchNodes() {
	watchChan := r.cli.Watch(r.ctx, nodePrefix, clientv3.WithPrefix(), clientv3.WithIgnoreValue())
	go func() {
		for {
			select {
			case resp := <-watchChan:
				for _, e := range resp.Events {
					r.handleNodeChange(e)
				}
			case <-r.ctx.Done():
				return
			}
		}
	}()
}

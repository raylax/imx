package registry

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/patrickmn/go-cache"
	"github.com/raylax/imx/core"
	"go.etcd.io/etcd/clientv3"
	"log"
	"sync"
	"time"
)

const (
	dialTimeout          = 5 * time.Second
	keepAlive            = 5
	keyPrefix            = "/imx"
	nodePrefix           = keyPrefix + "/node/"
	userPrefix           = keyPrefix + "/user/"
	cacheExpiration      = 1 * time.Minute
	cacheCleanupInterval = 30 * time.Second
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
	cache     *cache.Cache
	users     map[string]core.User
	m         sync.RWMutex
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
	r.cache = cache.New(cacheExpiration, cacheCleanupInterval)
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
	_, err := r.kv.Put(r.ctx, userPrefix+u.Id, r.node.Key(), clientv3.WithLease(r.leaseId))
	if err != nil {
		return err
	}
	r.m.Lock()
	r.users[u.Id] = u
	r.m.Unlock()
	return nil
}

func (r *EtcdRegistry) UnRegUser(u core.User) {
	r.m.Lock()
	r.users[u.Id] = u
	r.m.Unlock()
	_, _ = r.kv.Delete(r.ctx, userPrefix+u.Id)
}

func (r *EtcdRegistry) Lookup(u core.User) ([]core.Node, error) {
	v, ok := r.cache.Get(u.Id)
	if ok {
		return v.([]core.Node), nil
	}
	resp, err := r.kv.Get(r.ctx, userPrefix+u.Id, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	nodes := make([]core.Node, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		nodes = append(nodes, core.NewNodeFromKey(string(kv.Value)))
	}
	r.cache.SetDefault(u.Id, nodes)
	return nodes, nil
}

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

func (r *EtcdRegistry) reReg() error {
	log.Printf("重新注册服务")
	r.ctx, r.cancel = context.WithCancel(context.Background())
	if err := r.reg(); err != nil {
		return err
	}
	r.regUsers()
	log.Printf("服务重新注册完成")
	return nil
}

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

func (r *EtcdRegistry) regUsers()  {
	r.m.RLock()
	for _, u := range r.users {
		_ = r.RegUser(u)
	}
	r.m.RUnlock()
}

func (r *EtcdRegistry) keepAlive() error {
	log.Printf("设置自动续约")
	respChan, err := r.lease.KeepAlive(r.ctx, r.leaseId)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case resp := <-respChan:
				if resp == nil {
					log.Printf("租约到期")
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
	log.Printf("watchNodes => " + e.Type.String() + " , " + string(e.Kv.Key) + " , " + string(e.Kv.Value))
	switch e.Type {
	case mvccpb.PUT:

	case mvccpb.DELETE:

	}
}

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

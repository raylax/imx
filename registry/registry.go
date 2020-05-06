package registry

import "github.com/raylax/imx/core"

type Registry interface {

	// 注册节点
	Reg() error

	// 取消注册节点
	UnReg()

	// 注册用户
	RegUser(u core.User) error

	// 取消注册用户
	UnRegUser(u core.User)

	// 注册群组
	RegGroup(g core.Group, u core.User) error

	// 取消注册群组
	UnRegGroup(g core.Group, u core.User)

	// 获取群组下所有在线用户
	GetGroupUsers(gid string) ([]string, error)

	// 查找用户所在节点
	LookupNodes(uid string) ([]core.Node, error)
}

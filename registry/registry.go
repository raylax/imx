package registry

import "github.com/raylax/imx/core"

type Registry interface {
	Reg() error

	UnReg()

	RegUser(u core.User) error

	UnRegUser(u core.User)

	RegGroup(g core.Group, u core.User) error

	UnRegGroup(g core.Group, u core.User)

	GetGroupUsers(gid string) ([]string, error)

	LookupNodes(uid string) ([]core.Node, error)
}

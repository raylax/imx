package registry

import "github.com/raylax/imx/core"

type Registry interface {
	Reg() error

	UnReg()

	RegUser(u core.User) error

	UnRegUser(u core.User)

	LookupNode(id string) ([]core.Node, error)
}

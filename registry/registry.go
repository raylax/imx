package registry

import "github.com/raylax/imx/core"

type Registry interface {
	Reg() error

	UnReg()

	RegUser(u core.User) error

	UnRegUser(u core.User)

	Lookup(u core.User) ([]core.Node, error)
}

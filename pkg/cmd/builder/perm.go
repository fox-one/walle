package builder

import (
	"github.com/fox-one/walle/core"
)

type customPerm struct {
	Builder
	perm core.Perm
}

func (c customPerm) Perm() core.Perm {
	return c.perm
}

func WithPerm(b Builder, p core.Perm) Builder {
	return customPerm{
		Builder: b,
		perm:    p,
	}
}

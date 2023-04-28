package gen

import (
	"sync/atomic"

	"github.com/hawkingrei/gsqlancer/pkg/errors"
)

var GlobalStatue = &globalStatue{}

type globalStatue struct {
	errorIgnore errors.ErrorIgnore
	databaseID  atomic.Uint64
}

func (g *globalStatue) SetErrorIgnore(e errors.ErrorIgnore) {
	g.errorIgnore = e
}

func (g *globalStatue) GetDatabaseID() uint64 {
	return g.databaseID.Add(1)
}

func (g *globalStatue) IsIgnoreError(err error) bool {
	return g.errorIgnore.Contains(err) || g.errorIgnore.Regexp(err)
}

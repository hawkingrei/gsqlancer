package executor

import "sync/atomic"

var GlobalStatue = &globalStatue{}

type globalStatue struct {
	databaseID atomic.Uint64
}

func (g *globalStatue) GetDatabaseID() uint64 {
	return g.databaseID.Add(1)
}

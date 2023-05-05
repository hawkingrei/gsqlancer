package gen

import (
	"fmt"
	"math/rand"
	"sync/atomic"

	"golang.org/x/exp/maps"
)

type Action int

const (
	CREATE_TABLE Action = iota
	CREATE_INDEX
	INSERT
	ALTER_TABLE
	ANALYZE_TABLE
)

type TiDBState struct {
	tableID    map[uint32]uint32
	databaseID uint64
	tableIDGen atomic.Uint32
}

func NewTiDBState() *TiDBState {
	return &TiDBState{
		tableID: make(map[uint32]uint32),
	}
}

func (t *TiDBState) GenTableID() uint32 {
	return t.tableIDGen.Add(1)
}

func (t *TiDBState) GenColumn(tid uint32) uint32 {
	if _, ok := t.tableID[tid]; !ok {
		t.tableID[tid] = 0
	}
	t.tableID[tid] += 1
	return t.tableID[tid]
}

func (t *TiDBState) GenDatabaseID() uint64 {
	t.databaseID = GlobalStatue.GetDatabaseID()
	return t.databaseID
}

func (t *TiDBState) RandGetTableID() string {
	ids := maps.Keys(t.tableID)
	if len(ids) == 0 {
		return ""
	}
	return fmt.Sprintf("t%d", ids[rand.Intn(len(ids))])
}

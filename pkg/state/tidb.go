package state

import "sync/atomic"

type Action int

const (
	CREATE_TABLE Action = iota
	CREATE_INDEX
	INSERT
	ALTER_TABLE
	ANALYZE_TABLE
)

type TiDBState struct {
	databaseID atomic.Uint32
	tableIDGen atomic.Uint32
	tableID    map[uint32]uint32
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

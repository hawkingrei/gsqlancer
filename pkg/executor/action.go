package executor

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
	tableID    map[uint32]uint32
	databaseID uint64
	tableIDGen atomic.Uint32
}

func NewTiDBState() *TiDBState {
	databaseID := GlobalStatue.GetDatabaseID()
	return &TiDBState{
		databaseID: databaseID,
		tableID:    make(map[uint32]uint32),
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

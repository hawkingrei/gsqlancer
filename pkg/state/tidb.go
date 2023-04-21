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
	tableID    map[uint32]uint32
}

func NewTiDBState() *TiDBState {
	return &TiDBState{
		tableID: make(map[uint32]uint32),
	}
}

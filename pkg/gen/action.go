package gen

import (
	"fmt"
	"math/rand"
	"sync/atomic"

	"github.com/hawkingrei/gsqlancer/pkg/model"
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
	tableID       map[string]uint32 // tableID  -> columnID
	tableMeta     map[string]*model.Table
	databaseID    uint64
	tableIDGen    atomic.Uint32
	resultTable   []*model.Table
	tmpTableIDGen atomic.Uint32
	tableAlias    map[string]string
}

func NewTiDBState() *TiDBState {
	return &TiDBState{
		tableID: make(map[string]uint32),
	}
}

func (t *TiDBState) GenTableID() uint32 {
	return t.tableIDGen.Add(1)
}

func (t *TiDBState) GenColumn(tid uint32) uint32 {
	tbl := fmt.Sprintf("t%d", tid)
	if _, ok := t.tableID[tbl]; !ok {
		t.tableID[tbl] = 0
	}
	t.tableID[tbl] += 1
	return t.tableID[tbl]
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
	return ids[rand.Intn(len(ids))]
}

func (t *TiDBState) GetRandTableList() []string {
	ids := maps.Keys(t.tableID)
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	return ids
}

func (t *TiDBState) findTableByName(name string) *model.Table {
	meta, ok := t.tableMeta[name]
	if !ok {
		return nil
	}
	return meta
}

func (t *TiDBState) SetResultTable(tables []*model.Table) {
	t.resultTable = tables
}

func (t *TiDBState) AppendResultTable(table *model.Table) {
	t.resultTable = append(t.resultTable, table)
}

func (t *TiDBState) GetResultTable() []*model.Table {
	return t.resultTable
}

func (t *TiDBState) CreateTmpTable() string {
	id := t.tmpTableIDGen.Add(1)
	return fmt.Sprintf("tmp%d", id)
}

func (t *TiDBState) SetTableAlias(table string, alias string) {
	t.tableAlias[table] = alias
}

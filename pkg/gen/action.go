package gen

import (
	"fmt"
	"math/rand"
	"sync/atomic"

	"github.com/hawkingrei/gsqlancer/pkg/connection"
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
	tableID         map[string]uint32 // tableID  -> columnID
	tableMeta       map[string]*model.Table
	databaseID      uint64
	tableIDGen      atomic.Uint32
	tmpColIndex     atomic.Uint32
	resultTable     map[string]*model.Table
	InUsedTable     map[string]*model.Table
	tmpTableIDGen   atomic.Uint32
	TableAlias      map[string]string
	PivotRows       map[string]*connection.QueryItem
	unwrapPivotRows map[string]interface{}
}

func NewTiDBState() *TiDBState {
	return &TiDBState{
		tableID:     make(map[string]uint32),
		tableMeta:   make(map[string]*model.Table),
		TableAlias:  make(map[string]string),
		resultTable: make(map[string]*model.Table),
		InUsedTable: make(map[string]*model.Table),
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
	//logging.StatusLog().Debug("set table", zap.Any("table", tables))
	for _, table := range tables {
		t.resultTable[table.Name()] = table
	}
}

func (t *TiDBState) AppendResultTable(table *model.Table) {
	//logging.StatusLog().Debug("add table", zap.String("table", table.Name()), zap.String("alias", table.AliasName.String()))
	t.resultTable[table.Name()] = table
}

func (t *TiDBState) GetResultTable() []*model.Table {
	return maps.Values(t.resultTable)
}

func (t *TiDBState) CreateTmpTable() string {
	id := t.tmpTableIDGen.Add(1)
	return fmt.Sprintf("tmp%d", id)
}

func (t *TiDBState) SetTableAlias(table string, alias string) {
	t.TableAlias[table] = alias
}

func (t *TiDBState) GetInUsedTable() []*model.Table {
	return maps.Values(t.InUsedTable)
}

func (t *TiDBState) SetInUsedTable(tables []*model.Table) {
	for _, table := range tables {
		t.InUsedTable[table.Name()] = table
	}
}

func (t *TiDBState) TableMeta(name string) (*model.Table, bool) {
	tbl, ok := t.tableMeta[name]
	return tbl, ok
}

func (t *TiDBState) AddTableMeta(name string, tbl *model.Table) {
	t.tableMeta[name] = tbl
}

func (t *TiDBState) CreateTmpColumn() string {
	index := t.tmpColIndex.Add(1)
	return fmt.Sprintf("col_%d", index)
}

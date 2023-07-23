package executor

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	"github.com/hawkingrei/gsqlancer/pkg/gen"
	"github.com/hawkingrei/gsqlancer/pkg/gen/setgen"
	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/pingcap/tidb/parser/ast"
)

type generator struct {
	tableGen              *gen.TiDBTableGenerator
	analyzeGen            *gen.TiDBAnalyzeTable
	tiflashReplicaStmtGen *gen.TiflashReplicaStmtGen
	selectGen             *gen.TiDBSelectStmtGen
	insertGen             *gen.TiDBInsertGenerator
	setGen                *setgen.SetVariableGen
}

func newGenerator(c *config.Config, g *gen.TiDBState) generator {
	return generator{
		tableGen:   gen.NewTiDBTableGenerator(c, g),
		analyzeGen: gen.NewAnalyzeTable(c, g),
		selectGen:  gen.NewTiDBSelectStmtGen(c, g),
		insertGen:  gen.NewTiDBInsertGenerator(c, g),
	}
}

func (g *generator) CreateTable() (*model.SQL, *model.Table, error) {
	return g.tableGen.GenerateDDLCreateTable()
}

func (g *generator) AnalyzeTable() *model.SQL {
	return g.analyzeGen.GenerateAnalyzeTables()
}

func (g *generator) TiflashReplicaStmt(table string, replicaNum int) *model.SQL {
	return g.tiflashReplicaStmtGen.Build(table, replicaNum).Generate()
}

func (g *generator) SelectTable(usedTables []*model.Table) (selectStmtNode *ast.SelectStmt, sql string, err error) {
	return g.selectGen.GenSelectStmt(usedTables)
}

func (g *generator) InsertTable(table string) (*model.SQL, error) {
	return g.insertGen.GenerateDMLInsertByTable(table)
}

func (g *generator) GenPQSSelectStmt(pivotRows map[string]*connection.QueryItem, usedTables []*model.Table) (
	selectStmtNode *ast.SelectStmt, sql string, columnInfos []model.Column, updatedPivotRows map[string]*connection.QueryItem, err error) {
	return g.selectGen.GenPQSSelectStmt(pivotRows, usedTables)
}

func (g *generator) SetVariable() *model.SQL {
	return g.setGen.Gen()
}

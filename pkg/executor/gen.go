package executor

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	"github.com/hawkingrei/gsqlancer/pkg/gen"
	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/pingcap/tidb/parser/ast"
)

type generator struct {
	tableGen              *gen.TiDBTableGenerator
	analyzeGen            *gen.TiDBAnalyzeTable
	tiflashReplicaStmtGen *gen.TiflashReplicaStmtGen
	selectGen             *gen.TiDBSelectStmtGen
}

func newGenerator(c *config.Config, g *gen.TiDBState) generator {
	return generator{
		tableGen:   gen.NewTiDBTableGenerator(c, g),
		analyzeGen: gen.NewAnalyzeTable(c, g),
		selectGen:  gen.NewTiDBSelectStmtGen(c, g),
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

func (g *generator) SelectTable() (selectStmtNode *ast.SelectStmt, sql string, columnInfos []types.Column, updatedPivotRows map[string]*connection.QueryItem, err error) {
	return g.selectGen.Gen()
}

func (g *generator) GenPQSSelectStmt(pivotRows map[string]*connection.QueryItem,
	usedTables []*model.Table) (selectStmtNode *ast.SelectStmt, sql string, columnInfos []types.Column, updatedPivotRows map[string]*connection.QueryItem, err error) {
	return g.selectGen.GenPQSSelectStmt(pivotRows, usedTables)
}

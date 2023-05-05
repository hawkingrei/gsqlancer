package executor

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/gen"
	"github.com/hawkingrei/gsqlancer/pkg/model"
)

type generator struct {
	tableGen              *gen.TiDBTableGenerator
	analyzeGen            *gen.TiDBAnalyzeTable
	tiflashReplicaStmtGen *gen.TiflashReplicaStmtGen
}

func newGenerator(c *config.Config, g *gen.TiDBState) generator {
	return generator{
		tableGen:   gen.NewTiDBTableGenerator(c, g),
		analyzeGen: gen.NewAnalyzeTable(c, g),
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

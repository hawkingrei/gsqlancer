package executor

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/gen"
)

type generator struct {
	tableGen *gen.TiDBTableGenerator
}

func newGenerator(c *config.Config, g *gen.TiDBState) generator {
	return generator{
		tableGen: gen.NewTiDBTableGenerator(c, g),
	}
}

func (g *generator) CreateTable() {
	g.tableGen.GenerateDDLCreateTable()
}

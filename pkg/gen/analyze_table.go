package gen

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/model"
)

type TiDBAnalyzeTable struct {
	c           *config.Config
	globalState *TiDBState
}

func NewAnalyzeTable(c *config.Config, globalState *TiDBState) *TiDBAnalyzeTable {
	return &TiDBAnalyzeTable{
		c:           c,
		globalState: globalState,
	}
}

func (a *TiDBAnalyzeTable) GenerateAnalyzeTables() *model.SQL {
	if t := a.globalState.RandGetTableID(); t != "" {
		return &model.SQL{
			SQLType:  model.SQLTypeAnalyzeTable,
			SQLTable: t,
			SQLStmt:  "analyze table " + t,
		}

	}
	return nil
}

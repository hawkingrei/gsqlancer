package gen

import (
	"fmt"

	"github.com/hawkingrei/gsqlancer/pkg/model"
)

const tiflashReplicaStmt = "ALTER TABLE %s SET TIFLASH REPLICA count %d;"

type TiflashReplicaStmtGen struct {
	table      string
	replicaNum int
}

func (t *TiflashReplicaStmtGen) Build(table string, replicaNum int) *TiflashReplicaStmtGen {
	t.table = table
	t.replicaNum = replicaNum
	return t
}

func (t *TiflashReplicaStmtGen) Generate() *model.SQL {
	return &model.SQL{
		SQLType: model.SQLTypeDDLAlterTable,
		SQLStmt: fmt.Sprintf(tiflashReplicaStmt, t.table, t.replicaNum),
	}
}

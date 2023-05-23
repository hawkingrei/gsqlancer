package gen

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/parser/ast"
)

type TiDBInsertGenerator struct {
	c           *config.Config
	globalState *TiDBState
}

func NewTiDBInsertGenerator(c *config.Config, globalState *TiDBState) *TiDBInsertGenerator {
	return &TiDBInsertGenerator{
		c:           c,
		globalState: globalState,
	}
}

func (i *TiDBInsertGenerator) GenerateDMLInsertByTable(table string) (*types.SQL, error) {
	tree := insertStmt()
	stmt, err := i.walkInsertStmtForTable(tree, table)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &types.SQL{
		SQLType:  types.SQLTypeDMLInsert,
		SQLStmt:  stmt,
		SQLTable: table,
	}, nil
}

func insertStmt() *ast.InsertStmt {
	insertStmtNode := ast.InsertStmt{
		Table: &ast.TableRefsClause{
			TableRefs: &ast.Join{
				Left: &ast.TableName{},
			},
		},
		Lists:   [][]ast.ExprNode{},
		Columns: []*ast.ColumnName{},
	}
	return &insertStmtNode
}

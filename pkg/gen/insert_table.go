package gen

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
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

func (i *TiDBInsertGenerator) walkInsertStmtForTable(node *ast.InsertStmt, tableName string) (string, error) {
	table, ok := i.globalState.TableMeta(tableName)
	if !ok {
		return "", errors.Errorf("table %s not exist", tableName)
	}
	node.Table.TableRefs.Left.(*ast.TableName).Name = model.NewCIStr(table.Name())
	columns := e.walkColumns(&node.Columns, table)
	e.walkLists(&node.Lists, columns)
	return BufferOut(node)
}

func (i *TiDBInsertGenerator) walkColumns(columns *[]*ast.ColumnName, table *types.Table) []types.Column {
	cols := make([]types.Column, 0, len(table.Columns))
	for _, column := range table.Columns {
		if column.Name.HasPrefix("id_") || column.Name.EqString("id") {
			continue
		}
		*columns = append(*columns, &ast.ColumnName{
			Table: table.Name.ToModel(),
			Name:  column.Name.ToModel(),
		})
		cols = append(cols, column.Clone())
	}
	return cols
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

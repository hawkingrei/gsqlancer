package gen

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	gmodel "github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/hawkingrei/gsqlancer/pkg/util"
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
	columns := i.walkColumns(&node.Columns, table)
	i.walkLists(&node.Lists, columns)
	return BufferOut(node)
}

func (i *TiDBInsertGenerator) walkColumns(columns *[]*ast.ColumnName, table *gmodel.Table) []gmodel.Column {
	cols := make([]gmodel.Column, 0, len(table.Columns()))
	for _, column := range table.Columns() {
		*columns = append(*columns, &ast.ColumnName{
			Table: model.NewCIStr(table.Name()),
			Name:  model.NewCIStr(column.Name()),
		})
		cols = append(cols, column.Clone())
	}
	return cols
}

func (i *TiDBInsertGenerator) walkLists(lists *[][]ast.ExprNode, columns []gmodel.Column) {
	count := int(util.RdRange(10, 20))
	for i := 0; i < count; i++ {
		*lists = append(*lists, randList(columns))
	}
	// *lists = append(*lists, randor0(columns)...)
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

func randList(columns []gmodel.Column) []ast.ExprNode {
	var list []ast.ExprNode
	for _, column := range columns {
		// GenerateEnumDataItem
		list = append(list, ast.NewValueExpr(GenerateEnumDataItem(column), "", ""))
	}
	return list
}

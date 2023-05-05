package gen

import (
	"fmt"
	"math/rand"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/parser/ast"
	parsermodel "github.com/pingcap/tidb/parser/model"
	parserTypes "github.com/pingcap/tidb/parser/types"
)

var allColumnTypes = []string{"int", "float", "varchar"}

type TiDBTableGenerator struct {
	c           *config.Config
	globalState *TiDBState
}

func NewTiDBTableGenerator(c *config.Config, g *TiDBState) *TiDBTableGenerator {
	return &TiDBTableGenerator{
		c:           c,
		globalState: g,
	}
}

// GenerateDDLCreateTable rand create table statement
func (e *TiDBTableGenerator) GenerateDDLCreateTable() (*model.SQL, *model.Table, error) {
	tableBuilder := model.NewTableBuilder()
	tree := createTableStmt(e.c)
	stmt, table, tableMeta, err := e.walkDDLCreateTable(1, tree, tableBuilder)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	return &model.SQL{
		SQLType:  model.SQLTypeDDLCreateTable,
		SQLTable: table,
		SQLStmt:  stmt,
	}, tableMeta, nil
}

func (e *TiDBTableGenerator) walkDDLCreateTable(index int,
	node *ast.CreateTableStmt, tableBuilder *model.TableBuilder) (sql, table string, tableMeta *model.Table, err error) {

	tid := e.globalState.GenTableID()
	table = fmt.Sprintf("%s%d", "t", tid)
	tableBuilder.SetName(table)
	idColName := fmt.Sprintf("idx%d", index)

	idFieldType := parserTypes.NewFieldType(Type2Tp("bigint"))
	idFieldType.SetFlen(DataType2Len("bigint"))
	idCol := &ast.ColumnDef{
		Name:    &ast.ColumnName{Name: parsermodel.NewCIStr(idColName)},
		Tp:      idFieldType,
		Options: []*ast.ColumnOption{randColumnOptionAuto()},
	}
	tableBuilder.AddColumn(idCol)
	node.Cols = append(node.Cols, idCol)
	makeConstraintPrimaryKey(node, idColName)

	node.Table.Name = parsermodel.NewCIStr(table)
	min := 2
	max := 10
	columnsNum := rand.Intn(max-min+1) + min
	for i := 0; i < columnsNum; i++ {
		col := model.RandGenColumn(fmt.Sprintf("c%d", e.globalState.GenColumn(tid)))
		node.Cols = append(node.Cols, col.GetAst())
	}

	// Auto_random should only have one primary key
	// The columns in expressions of partition table must be included
	// in primary key/unique constraints
	// So no partition table if there is auto_random column
	if idCol.Options[0].Tp == ast.ColumnOptionAutoRandom {
		node.Partition = nil
	}
	if node.Partition != nil {
		//if colType := e.walkPartition(index, node.Partition, colTypes); colType != "" {
		//	makeConstraintPrimaryKey(node, fmt.Sprintf("col_%s_%d", colType, index))
		//} else {
		//	node.Partition = nil
		//}
		node.Partition = nil
	}
	tableMeta = tableBuilder.Build()
	sql, err = BufferOut(node)
	if err != nil {
		return "", "", nil, errors.Trace(err)
	}
	return sql, table, tableMeta, errors.Trace(err)
}

func createTableStmt(c *config.Config) *ast.CreateTableStmt {
	createTableNode := ast.CreateTableStmt{
		Table:       &ast.TableName{},
		Cols:        []*ast.ColumnDef{},
		Constraints: []*ast.Constraint{},
		Options:     []*ast.TableOption{},
	}
	// TODO: config for enable partition
	// partitionStmt is disabled
	if rand.Intn(2) == 0 && c.EnablePartition() {
		createTableNode.Partition = partitionStmt()
	}

	return &createTableNode
}

func partitionStmt() *ast.PartitionOptions {
	return &ast.PartitionOptions{
		PartitionMethod: ast.PartitionMethod{
			ColumnNames: []*ast.ColumnName{},
		},
		Definitions: []*ast.PartitionDefinition{},
	}
}

func makeConstraintPrimaryKey(node *ast.CreateTableStmt, column string) {
	for _, constraint := range node.Constraints {
		if constraint.Tp == ast.ConstraintPrimaryKey {
			constraint.Keys = append(constraint.Keys, &ast.IndexPartSpecification{
				Column: &ast.ColumnName{
					Name: parsermodel.NewCIStr(column),
				},
			})
			return
		}
	}
	node.Constraints = append(node.Constraints, &ast.Constraint{
		Tp: ast.ConstraintPrimaryKey,
		Keys: []*ast.IndexPartSpecification{
			{
				Column: &ast.ColumnName{
					Name: parsermodel.NewCIStr(column),
				},
			},
		},
	})
}

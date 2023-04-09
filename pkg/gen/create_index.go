package gen

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/parser/ast"
	parsermodel "github.com/pingcap/tidb/parser/model"
	parserTypes "github.com/pingcap/tidb/parser/types"
)

type TiDBTableGenerator struct {
	c *config.Config
}

// GenerateDDLCreateTable rand create table statement
func (e *TiDBTableGenerator) GenerateDDLCreateTable() (*model.SQL, error) {
	tree := createTableStmt(e.c)

	stmt, table, err := e.walkDDLCreateTable(index, tree, colTypes)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &model.SQL{
		SQLType:  model.SQLTypeDDLCreateTable,
		SQLTable: table,
		SQLStmt:  stmt,
	}, nil
}

func (e *TiDBTableGenerator) walkDDLCreateTable(index int, node *ast.CreateTableStmt, colTypes []string) (sql string, table string, err error) {
	table = fmt.Sprintf("%s_%s", "table", strings.Join(colTypes, "_"))
	idColName := fmt.Sprintf("id_%d", index)

	idFieldType := parserTypes.NewFieldType(Type2Tp("bigint"))
	idFieldType.SetFlen(DataType2Len("bigint"))
	idCol := &ast.ColumnDef{
		Name:    &ast.ColumnName{Name: parsermodel.NewCIStr(idColName)},
		Tp:      idFieldType,
		Options: []*ast.ColumnOption{randColumnOptionAuto()},
	}
	node.Cols = append(node.Cols, idCol)
	makeConstraintPrimaryKey(node, idColName)

	node.Table.Name = parsermodel.NewCIStr(table)
	for _, colType := range colTypes {
		fieldType := parserTypes.NewFieldType(Type2Tp(colType))
		fieldType.SetFlen(DataType2Len(colType))
		node.Cols = append(node.Cols, &ast.ColumnDef{
			Name: &ast.ColumnName{Name: parsermodel.NewCIStr(fmt.Sprintf("col_%s_%d", colType, index))},
			Tp:   fieldType,
		})
	}

	// Auto_random should only have one primary key
	// The columns in expressions of partition table must be included
	// in primary key/unique constraints
	// So no partition table if there is auto_random column
	if idCol.Options[0].Tp == ast.ColumnOptionAutoRandom {
		node.Partition = nil
	}
	if node.Partition != nil {
		if colType := e.walkPartition(index, node.Partition, colTypes); colType != "" {
			makeConstraintPrimaryKey(node, fmt.Sprintf("col_%s_%d", colType, index))
		} else {
			node.Partition = nil
		}
	}
	sql, err = BufferOut(node)
	if err != nil {
		return "", "", err
	}
	return sql, table, errors.Trace(err)
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
	if rand.Intn(2) == 0 && c.EnablePartition {
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
					Name: model.NewCIStr(column),
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
					Name: model.NewCIStr(column),
				},
			},
		},
	})
}

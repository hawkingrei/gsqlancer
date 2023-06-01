package model

import (
	"fmt"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/types"
)

type Column struct {
	ColumnDef  *ast.ColumnDef
	Table      model.CIStr
	AliasTable model.CIStr
	AliasName  model.CIStr
}

func NewColumn(def *ast.ColumnDef) *Column {
	return &Column{
		ColumnDef: def,
	}
}

func (c *Column) Type() *types.FieldType {
	return c.ColumnDef.Tp
}

func (c *Column) Name() string {
	return c.ColumnDef.Name.String()
}

func (c *Column) String() string {
	return fmt.Sprintf("%s.%s", c.Table, c.Name)
}

// Clone makes a replica of column
func (c *Column) Clone() Column {
	return Column{
		ColumnDef:  &*(c.ColumnDef),
		Table:      c.Table,
		AliasTable: c.AliasTable,
		AliasName:  c.AliasName,
	}
}

func (c *Column) HasOption(opt ast.ColumnOptionType) bool {
	for _, option := range c.ColumnDef.Options {
		if option.Tp == opt {
			return true
		}
	}
	return false
}

// ToModel converts to ast model
func (c *Column) ToColumnNameExpr(tableName string) *ast.ColumnNameExpr {
	table := c.Table
	if c.AliasTable.String() != "" {
		table = c.AliasTable
	}
	name := c.Name()
	if c.AliasTable.String() != "" {
		name = c.AliasName.String()
	}
	if table.String() == "" {
		table = model.NewCIStr(tableName)
	}
	return &ast.ColumnNameExpr{
		Name: &ast.ColumnName{
			Table: table,
			Name:  model.NewCIStr(name),
		},
	}
}

type Table struct {
	Refer     *ast.ReferenceDef // Used for foreign key.
	Option    *ast.IndexOption  // Index Options
	name      string
	AliasName model.CIStr
	columns   []*Column
	indexes   []*ast.IndexPartSpecification // Used for PRIMARY KEY, UNIQUE, ......
}

func (t *Table) Name() string {
	if t.name == "" {
		return t.AliasName.String()
	}
	return t.name
}

func (t *Table) Index() []*ast.IndexPartSpecification {
	return t.indexes
}

func (t *Table) Rename(name string) *Table {
	table := NewTableBuilder().SetName(name)
	for _, column := range t.columns {
		table.AddColumn(column)
	}
	for _, index := range t.indexes {
		table.AddIndex(index)
	}
	return table.Build()
}

func (t *Table) Columns() []*Column {
	return t.columns
}

type TableBuilder struct {
	table *Table
}

func NewTableBuilder() *TableBuilder {
	return &TableBuilder{
		table: &Table{
			columns: make([]*Column, 0),
			indexes: make([]*ast.IndexPartSpecification, 0),
		},
	}
}

func (t *TableBuilder) SetName(name string) *TableBuilder {
	t.table.name = name
	t.table.AliasName = model.NewCIStr(name)
	return t
}

func (t *TableBuilder) AddColumn(column *Column) *TableBuilder {
	t.table.columns = append(t.table.columns, column)
	return t
}

func (t *TableBuilder) AddIndex(index *ast.IndexPartSpecification) *TableBuilder {
	t.table.indexes = append(t.table.indexes, index)
	return t
}

func (t *TableBuilder) Build() *Table {
	return t.table
}

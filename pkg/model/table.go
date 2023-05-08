package model

import "github.com/pingcap/tidb/parser/ast"

type Table struct {
	Refer   *ast.ReferenceDef // Used for foreign key.
	Option  *ast.IndexOption  // Index Options
	name    string
	columns []*ast.ColumnDef
	indexes []*ast.IndexPartSpecification // Used for PRIMARY KEY, UNIQUE, ......
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Index() []*ast.IndexPartSpecification {
	return t.indexes
}

type TableBuilder struct {
	table *Table
}

func NewTableBuilder() *TableBuilder {
	return &TableBuilder{
		table: &Table{
			columns: make([]*ast.ColumnDef, 0),
			indexes: make([]*ast.IndexPartSpecification, 0),
		},
	}
}

func (t *TableBuilder) SetName(name string) *TableBuilder {
	t.table.name = name
	return t
}

func (t *TableBuilder) AddColumn(column *ast.ColumnDef) *TableBuilder {
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

package gen

import (
	"math/rand"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	"github.com/hawkingrei/gsqlancer/pkg/gen/hint"
	gmodel "github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/parser/mysql"
	tidb_types "github.com/pingcap/tidb/types"
	parser_driver "github.com/pingcap/tidb/types/parser_driver"
	"go.uber.org/zap"
)

type TiDBSelectStmtGen struct {
	c           *config.Config
	globalState *TiDBState
	gen         *Generator
}

func NewTiDBSelectStmtGen(c *config.Config, g *TiDBState) *TiDBSelectStmtGen {
	return &TiDBSelectStmtGen{
		c:           c,
		globalState: g,
		gen:         NewGenerator(c, g),
	}
}

func (t *TiDBSelectStmtGen) GenPQSSelectStmt(pivotRows map[string]*connection.QueryItem,
	usedTables []*gmodel.Table) (*ast.SelectStmt, string, []gmodel.Column, map[string]*connection.QueryItem, error) {
	t.globalState.PivotRows = pivotRows
	t.globalState.InUsedTable = usedTables
	return t.Gen()
}

func (t *TiDBSelectStmtGen) Gen() (selectStmtNode *ast.SelectStmt, sql string, columnInfos []gmodel.Column, updatedPivotRows map[string]*connection.QueryItem, err error) {
	selectStmtNode = &ast.SelectStmt{
		SelectStmtOpts: &ast.SelectStmtOpts{
			SQLCache: true,
		},
		Fields: &ast.FieldList{
			Fields: []*ast.SelectField{},
		},
	}
	selectStmtNode.From = t.TableRefsClause()
	t.walkTableRefs(selectStmtNode.From.TableRefs)
	selectStmtNode.Where = t.ConditionClause()
	selectStmtNode.TableHints = t.tableHintsExpr(t.globalState.InUsedTable)
	columnInfos, updatedPivotRows = t.walkResultFields(selectStmtNode)
	sql, err = BufferOut(selectStmtNode)
	return selectStmtNode, sql, columnInfos, updatedPivotRows, err
}

// TableRefsClause generates linear joins:
//
//	          Join
//	    Join        t4
//	Join     t3
//
// t1    t2
func (t *TiDBSelectStmtGen) TableRefsClause() *ast.TableRefsClause {
	clause := &ast.TableRefsClause{TableRefs: &ast.Join{
		Left:  &ast.TableName{},
		Right: &ast.TableName{},
	}}
	usedTables := t.globalState.GetRandTableList()
	var node = clause.TableRefs
	// TODO: it works, but need to refactor
	if len(usedTables) == 1 {
		clause.TableRefs = &ast.Join{
			Left: &ast.TableSource{
				Source: &ast.TableName{
					Name: model.NewCIStr(usedTables[0]),
				},
			},
			Right: nil,
		}
	}
	for i := len(usedTables) - 1; i >= 1; i-- {
		var tp ast.JoinType
		if !t.c.EnableLeftRightJoin {
			tp = ast.CrossJoin
		} else {
			switch rand.Intn(3) {
			case 0:
				tp = ast.RightJoin
			case 1:
				tp = ast.LeftJoin
			default:
				tp = ast.CrossJoin
			}
		}
		node.Right = &ast.TableSource{
			Source: &ast.TableName{
				Name: model.NewCIStr(usedTables[i]),
			},
		}
		node.Tp = tp
		if i == 1 {
			node.Left = &ast.TableSource{
				Source: &ast.TableName{
					Name: model.NewCIStr(usedTables[i-1]),
				},
			}
		} else {
			node.Left = &ast.Join{}
			node = node.Left.(*ast.Join)
		}
	}
	return clause
}

func (t *TiDBSelectStmtGen) tableHintsExpr(usedTables []*gmodel.Table) []*ast.TableOptimizerHint {
	hints := make([]*ast.TableOptimizerHint, 0)
	if !t.c.Hint {
		return hints
	}
	// avoid duplicated hints
	enabledHints := make(map[string]bool)
	length := rand.Intn(4)
	for i := 0; i < length; i++ {
		to := hint.GenerateHintExpr(usedTables)
		if to == nil {
			continue
		}
		if _, ok := enabledHints[to.HintName.String()]; !ok {
			hints = append(hints, to)
			enabledHints[to.HintName.String()] = true
		}
	}
	return hints
}

func (t *TiDBSelectStmtGen) walkResultFields(node *ast.SelectStmt) ([]gmodel.Column, map[string]*connection.QueryItem) {
	if t.c.EnableNoRECApproach {
		exprNode := &parser_driver.ValueExpr{}
		tp := tidb_types.NewFieldType(mysql.TypeLonglong)
		tp.SetFlag(128)
		exprNode.TexprNode.SetType(tp)
		exprNode.Datum.SetInt64(1)
		countField := ast.SelectField{
			Expr: &ast.AggregateFuncExpr{
				F: "count",
				Args: []ast.ExprNode{
					exprNode,
				},
			},
		}
		node.Fields.Fields = append(node.Fields.Fields, &countField)
		return nil, nil
	}
	// only used for tlp mode now, may effective on more cases in future
	if !t.c.EnablePQSApproach {
		// rand if there is aggregation func in result field
		if rand.Intn(6) > 3 {
			selfComposableAggs := []string{"sum", "min", "max"}
			for i := 0; i < rand.Intn(3)+1; i++ {
				child, _, err := t.generateExpr(types.TypeNumberLikeArg, 2)
				if err != nil {
					logging.StatusLog().Error("generate child expr of aggeration in result field error", zap.Error(err))
				} else {
					aggField := ast.SelectField{
						Expr: &ast.AggregateFuncExpr{
							F: selfComposableAggs[rand.Intn(len(selfComposableAggs))],
							Args: []ast.ExprNode{
								child,
							},
						},
					}
					node.Fields.Fields = append(node.Fields.Fields, &aggField)
				}
			}

			if len(node.Fields.Fields) != 0 {
				return nil, nil
			}
		}
	}
	columns := make([]gmodel.Column, 0)
	row := make(map[string]*connection.QueryItem)
	for _, table := range t.globalState.GetResultTable() {
		for _, column := range table.Columns() {
			asname := t.globalState.CreateTmpColumn()
			selectField := ast.SelectField{
				Expr:   column.ToColumnNameExpr(),
				AsName: model.NewCIStr(asname),
			}
			node.Fields.Fields = append(node.Fields.Fields, &selectField)
			col := column.Clone()
			col.AliasName = model.NewCIStr(asname)
			columns = append(columns, col)
			row[asname] = t.globalState.PivotRows[column.String()]
		}
	}
	return columns, row
}

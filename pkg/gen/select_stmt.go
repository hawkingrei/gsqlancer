package gen

import (
	"math/rand"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	"github.com/hawkingrei/gsqlancer/pkg/gen/hint"
	gmodel "github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
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
	usedTables []*gmodel.Table) (*ast.SelectStmt, string, []types.Column, map[string]*connection.QueryItem, error) {
	t.globalState.PivotRows = pivotRows
	t.globalState.InUsedTable = usedTables
	return t.Gen()
}

func (t *TiDBSelectStmtGen) Gen() (selectStmtNode *ast.SelectStmt, sql string, columnInfos []types.Column, updatedPivotRows map[string]*connection.QueryItem, err error) {
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
	return
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

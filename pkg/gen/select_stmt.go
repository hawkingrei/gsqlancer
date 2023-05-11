package gen

import (
	"math/rand"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
)

type TiDBSelectStmtGen struct {
	c           *config.Config
	globalState *TiDBState
}

func NewTiDBSelectStmtGen(c *config.Config, g *TiDBState) *TiDBSelectStmtGen {
	return &TiDBSelectStmtGen{
		c:           c,
		globalState: g,
	}
}

func (*TiDBSelectStmtGen) Gen() (selectStmtNode *ast.SelectStmt) {
	selectStmtNode = &ast.SelectStmt{
		SelectStmtOpts: &ast.SelectStmtOpts{
			SQLCache: true,
		},
		Fields: &ast.FieldList{
			Fields: []*ast.SelectField{},
		},
	}
	return selectStmtNode
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

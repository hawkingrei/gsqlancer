package gen

import (
	gmodel "github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"go.uber.org/zap"
)

func (g *TiDBSelectStmtGen) walkTableRefs(node *ast.Join) {
	if node.Right == nil {
		if node, ok := node.Left.(*ast.TableSource); ok {
			if tn, ok := node.Source.(*ast.TableName); ok {
				if table := g.globalState.findTableByName(tn.Name.L); table != nil {
					g.globalState.AppendResultTable(table)
					return
				}
			}
		}
		panic("unreachable")
	}

	if right, ok := node.Right.(*ast.TableSource); ok {
		var (
			leftTables []*gmodel.Table
			rightTable *gmodel.Table
		)
		switch node := node.Left.(type) {
		case *ast.Join:
			g.walkTableRefs(node)
			leftTables = g.globalState.GetResultTable()
		case *ast.TableSource:
			if tn, ok := node.Source.(*ast.TableName); ok {
				if table := g.globalState.findTableByName(tn.Name.L); table != nil {
					tmpTable := g.globalState.CreateTmpTable()
					node.AsName = model.NewCIStr(tmpTable)
					leftTables = []*gmodel.Table{table.Rename(tmpTable)}
					g.globalState.SetTableAlias(table.Name(), tmpTable)
					break
				}
			}
		default:
			panic("unreachable")
		}
		if table := g.globalState.findTableByName(right.Source.(*ast.TableName).Name.L); table != nil {
			tmpTable := g.globalState.CreateTmpTable()
			right.AsName = model.NewCIStr(tmpTable)
			rightTable = table.Rename(tmpTable)
			g.globalState.SetTableAlias(table.Name(), tmpTable)
		} else {
			panic("unreachable")
		}
		allTables := append(leftTables, rightTable)
		// usedTables := genCtx.UsedTables
		g.globalState.SetResultTable(allTables)
		// genCtx.UsedTables = allTables
		// defer func() {
		// 	genCtx.UsedTables = usedTables
		// }()
		node.On = &ast.OnCondition{}
		// for _, table := range genCtx.ResultTables {
		// 	fmt.Println(table.Name, table.AliasName)
		// }
		node.On.Expr = g.ConditionClause(1)
		return
	}

	panic("unreachable")
}

// ConditionClause is to generate a ConditionClause
func (g *TiDBSelectStmtGen) ConditionClause(depth int) ast.ExprNode {
	// TODO: support subquery
	// TODO: more ops
	exprType := types.TypeNumberLikeArg
	var err error
	retry := 2
	node, val, err := g.generateExpr(ctx, exprType, depth)
	for err != nil && retry > 0 {
		log.L().Error("generate where expr error", zap.Error(err))
		node, val, err = g.generateExpr(ctx, exprType, depth)
		retry--
	}
	if err != nil {
		panic("retry times exceed 3")
	}

	return g.rectifyCondition(node, val)
}

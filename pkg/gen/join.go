package gen

import (
	"math/rand"

	gmodel "github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/hawkingrei/gsqlancer/pkg/util"
	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/parser/opcode"
	tidb_types "github.com/pingcap/tidb/types"
	parser_driver "github.com/pingcap/tidb/types/parser_driver"
	"go.uber.org/zap"
)

func (t *TiDBSelectStmtGen) walkTableRefs(node *ast.Join) {
	if node.Right == nil {
		if node, ok := node.Left.(*ast.TableSource); ok {
			if tn, ok := node.Source.(*ast.TableName); ok {
				if table := t.globalState.findTableByName(tn.Name.L); table != nil {
					t.globalState.AppendResultTable(table)
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
			t.walkTableRefs(node)
			leftTables = t.globalState.GetResultTable()
		case *ast.TableSource:
			if tn, ok := node.Source.(*ast.TableName); ok {
				if table := t.globalState.findTableByName(tn.Name.L); table != nil {
					tmpTable := t.globalState.CreateTmpTable()
					node.AsName = model.NewCIStr(tmpTable)
					leftTables = []*gmodel.Table{table.Rename(tmpTable)}
					t.globalState.SetTableAlias(table.Name(), tmpTable)
					break
				}
			}
		default:
			panic("unreachable")
		}
		if table := t.globalState.findTableByName(right.Source.(*ast.TableName).Name.L); table != nil {
			tmpTable := t.globalState.CreateTmpTable()
			right.AsName = model.NewCIStr(tmpTable)
			rightTable = table.Rename(tmpTable)
			t.globalState.SetTableAlias(table.Name(), tmpTable)
		} else {
			panic("unreachable")
		}
		allTables := append(leftTables, rightTable)
		t.globalState.SetResultTable(allTables)
		node.On = &ast.OnCondition{}
		node.On.Expr = t.ConditionClause()
		return
	}
	panic("unreachable")
}

// ConditionClause is to generate a ConditionClause
func (t *TiDBSelectStmtGen) ConditionClause() ast.ExprNode {
	// TODO: support subquery
	// TODO: more ops
	exprType := types.TypeNumberLikeArg
	var err error
	retry := 2
	node, val, err := t.generateExpr(exprType, t.c.Depth)
	for err != nil && retry > 0 {
		log.L().Error("generate where expr error", zap.Error(err))
		node, val, err = t.generateExpr(exprType, t.c.Depth)
		retry--
	}
	if err != nil {
		panic("retry times exceed 3")
	}

	return t.rectifyCondition(node, val)
}

func (*TiDBSelectStmtGen) rectifyCondition(node ast.ExprNode, val parser_driver.ValueExpr) ast.ExprNode {
	// pthese := ast.ParenthesesExpr{}
	// pthese.Expr = node
	switch val.Kind() {
	case tidb_types.KindNull:
		node = &ast.IsNullExpr{
			// Expr: &pthese,
			Expr: node,
			Not:  false,
		}
	default:
		// make it true
		zero := parser_driver.ValueExpr{}
		zero.SetValue(false)
		if util.CompareValueExpr(val, zero) == 0 {
			node = &ast.UnaryOperationExpr{
				Op: opcode.Not,
				// V:  &pthese,
				V: node,
			}
		} else {
			// make it return 1 as true through add IS NOT TRUE
			node = &ast.IsNullExpr{
				// Expr: &pthese,
				Expr: node,
				Not:  true,
			}
		}
	}
	// remove useless parenthese
	if n, ok := node.(*ast.ParenthesesExpr); ok {
		node = n.Expr
	}
	return node
}

func (t *TiDBSelectStmtGen) generateExpr(valueTp uint64, depth int) (ast.ExprNode, parser_driver.ValueExpr, error) {
	if valueTp == 0 {
		log.L().Warn("required return type is empty")
		e := parser_driver.ValueExpr{}
		e.SetValue(nil)
		return ast.NewValueExpr(nil, "", ""), e, nil
	}
	// select one type from valueTp randomly
	// TODO: make different types have different chance
	// e.g. INT/FLOAT appears more often than NumberLikeString...
	possibleTypes := make([]uint64, 0)
	for valueTp != 0 {
		possibleTypes = append(possibleTypes, valueTp&^(valueTp-1))
		valueTp &= valueTp - 1
	}
	tp := possibleTypes[rand.Intn(len(possibleTypes))]

	// TODO: fix me!
	// when tp is NumberLikeString, no function can return it
	// but in ALMOST cases, it can be replaced by INT/FLOAT
	if tp == types.TypeNumberLikeStringArg {
		if util.RdBool() {
			tp = types.TypeFloatArg
		} else {
			tp = types.TypeIntArg
		}
	}

	// generate leaf node
	if depth == 0 {
		if rand.Intn(3) > 1 {
			exprNode, value, err := t.gen.columnExpr(tp)
			if err == nil {
				return exprNode, value, nil
			}
		}
		exprNode, value := t.gen.constValueExpr(tp)
		return exprNode, value, nil
	}

	// select a function with return type tp
	fn, err := util.OpFuncGroupByRet.RandOpFn(tp)
	if err != nil {
		log.L().Warn("generate fn or op failed", zap.Error(err))
		// if no op/fn satisfied requirement
		// generate const instead
		exprNode, value := t.gen.constValueExpr(tp)
		return exprNode, value, nil
	}

	pthese := ast.ParenthesesExpr{}
	var value parser_driver.ValueExpr
	pthese.Expr, value, err = fn.Node(func(childTp uint64) (ast.ExprNode, parser_driver.ValueExpr, error) {
		return t.generateExpr(childTp, depth-1)
	}, tp)
	if err != nil {
		return nil, parser_driver.ValueExpr{}, errors.Trace(err)
	}
	return &pthese, value, nil
}

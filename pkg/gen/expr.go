// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gen

import (
	"fmt"
	"math/rand"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	gmodel "github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/hawkingrei/gsqlancer/pkg/util"
	"github.com/juju/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/parser/mysql"
	tidbtypes "github.com/pingcap/tidb/types"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
	"go.uber.org/zap"
)

type Generator struct {
	c           *config.Config
	globalState *TiDBState
}

func NewGenerator(c *config.Config, g *TiDBState) *Generator {
	return &Generator{
		c:           c,
		globalState: g,
	}
}

func (*Generator) constValueExpr(tp uint64) (ast.ValueExpr, parserdriver.ValueExpr) {
	val := parserdriver.ValueExpr{}
	// return null randomly
	if rand.Intn(6) == 0 {
		val.SetNull()
		return ast.NewValueExpr(nil, "", ""), val
	}

	// choose one type from tp
	possibleTypes := make([]uint64, 0)
	for tp != 0 {
		possibleTypes = append(possibleTypes, tp&^(tp-1))
		tp &= tp - 1
	}
	tp = possibleTypes[rand.Intn(len(possibleTypes))]

	switch tp {
	case types.TypeFloatArg:
		switch rand.Intn(6) {
		case 0, 1:
			val.SetFloat64(0)
			return ast.NewValueExpr(0.0, "", ""), val
		case 2:
			val.SetFloat64(1.0)
			return ast.NewValueExpr(1.0, "", ""), val
		case 3:
			val.SetFloat64(-1.0)
			return ast.NewValueExpr(-1.0, "", ""), val
		default:
			f := rand.Float64()
			val.SetFloat64(f)
			return ast.NewValueExpr(f, "", ""), val
		}
	case types.TypeIntArg:
		switch rand.Intn(6) {
		case 0, 1:
			val.SetInt64(0)
			return ast.NewValueExpr(0, "", ""), val
		case 2:
			val.SetInt64(1)
			return ast.NewValueExpr(1, "", ""), val
		case 3:
			val.SetInt64(-1)
			return ast.NewValueExpr(-1, "", ""), val
		default:
			i := util.RdInt64()
			val.SetInt64(i)
			return ast.NewValueExpr(i, "", ""), val
		}
	case types.TypeNonFormattedStringArg:
		switch rand.Intn(3) {
		case 0, 1:
			s := util.RdString(rand.Intn(10))
			val.SetString(s, "")
			return ast.NewValueExpr(s, "", ""), val
		default:
			val.SetString("", "")
			return ast.NewValueExpr("", "", ""), val
		}
	case types.TypeNumberLikeStringArg:
		var s = fmt.Sprintf("%f", rand.Float64())
		if util.RdBool() {
			s = fmt.Sprintf("%d", util.RdInt64())
		}
		val.SetString(s, "")
		return ast.NewValueExpr(s, "", ""), val
	case types.TypeDatetimeArg:
		if util.RdBool() {
			// is TIMESTAMP
			t := tidbtypes.NewTime(tidbtypes.FromGoTime(util.RdTimestamp()), mysql.TypeTimestamp, rand.Intn(7))
			val.SetMysqlTime(t)
			return ast.NewValueExpr(t, "", ""), val
		}
		// is DATETIME
		t := tidbtypes.NewTime(tidbtypes.FromGoTime(util.RdDate()), mysql.TypeDatetime, 0)
		val.SetMysqlTime(t)
		return ast.NewValueExpr(t, "", ""), val
	case types.TypeDatetimeLikeStringArg:
		s := util.RdTimestamp().Format("2006-01-02 15:04:05")
		val.SetString(s, "")
		return ast.NewValueExpr(s, "", ""), val
	default:
		log.L().Error("this type is not implemented", zap.Uint64("tp", tp))
		panic("unreachable")
	}
}

func (g *Generator) columnExpr(argTp uint64) (*ast.ColumnNameExpr, parserdriver.ValueExpr, error) {
	val := parserdriver.ValueExpr{}
	tables := g.globalState.GetResultTable()
	if g.c.IsInUpdateDeleteStmt || g.c.IsInExprIndex {
		tables = g.globalState.GetInUsedTable()
	}
	randTable := tables[rand.Intn(len(tables))]
	tempCols := make([]*gmodel.Column, 0)
	for _, column := range randTable.Columns() {
		if util.TransStringType(column.Type().String())&argTp != 0 {
			tempCols = append(tempCols, column)
		}
	}
	if len(tempCols) == 0 {
		return nil, val, fmt.Errorf("no valid column as arg %d table %s", argTp, randTable.Name)
	}
	randColumn := tempCols[rand.Intn(len(tempCols))]
	colName, typeStr := randColumn.Name(), randColumn.Type().String()
	col := new(ast.ColumnNameExpr)
	col.Name = &ast.ColumnName{
		Table: model.NewCIStr(randTable.Name()),
		Name:  model.NewCIStr(colName),
	}
	col.SetType(tidbtypes.NewFieldType(util.TransToMysqlType(util.TransStringType(typeStr))))

	// no pivot rows
	// so value is generate by random
	if !g.c.EnablePQSApproach || g.c.EnableNoRECApproach || g.c.IsInUpdateDeleteStmt {
		tp := util.TransStringType(typeStr)
		_, val = g.constValueExpr(tp)
		return col, val, nil
	}

	// find the column value in pivot rows
	for key, value := range g.globalState.unwrapPivotRows {
		originTableName := col.Name.Table.L
		for k, v := range g.globalState.TableAlias {
			if v == originTableName {
				originTableName = k
				break
			}
		}
		originColumnName := col.Name.Name.L
		if key == fmt.Sprintf("%s.%s", originTableName, originColumnName) {
			val.SetValue(value)
			if tmpTable, ok := g.globalState.TableAlias[col.Name.Table.L]; ok {
				col.Name.Table = model.NewCIStr(tmpTable)
			}
			break
		}
	}

	return col, val, nil
}

func (g *Generator) generateExpr(valueTp uint64, depth int) (ast.ExprNode, parserdriver.ValueExpr, error) {
	if valueTp == 0 {
		log.L().Warn("required return type is empty")
		e := parserdriver.ValueExpr{}
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
			exprNode, value, err := g.columnExpr(tp)
			if err == nil {
				return exprNode, value, nil
			}
		}
		exprNode, value := g.constValueExpr(tp)
		return exprNode, value, nil
	}

	// select a function with return type tp
	fn, err := util.OpFuncGroupByRet.RandOpFn(tp)
	if err != nil {
		log.L().Warn("generate fn or op failed", zap.Error(err))
		// if no op/fn satisfied requirement
		// generate const instead
		exprNode, value := g.constValueExpr(tp)
		return exprNode, value, nil
	}

	pthese := ast.ParenthesesExpr{}
	var value parserdriver.ValueExpr
	pthese.Expr, value, err = fn.Node(func(childTp uint64) (ast.ExprNode, parserdriver.ValueExpr, error) {
		return g.generateExpr(childTp, depth-1)
	}, tp)
	if err != nil {
		return nil, parserdriver.ValueExpr{}, errors.Trace(err)
	}
	return &pthese, value, nil
}

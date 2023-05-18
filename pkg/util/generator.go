package util

import (
	"fmt"
	"strings"

	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/pingcap/tidb/parser/charset"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/sessionctx/stmtctx"
	parser_driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/pingcap/tidb/util/collate"
)

// OpFuncGroupByRet is a map
var OpFuncGroupByRet = make(types.OpFuncIndex)

// RegisterToOpFnIndex is to register OpFuncEval
func RegisterToOpFnIndex(o types.OpFuncEval) {
	returnType := o.GetPossibleReturnType()

	for returnType != 0 {
		// i must be n-th power of 2: 1001100 => i(0000100)
		i := returnType &^ (returnType - 1)
		if _, ok := OpFuncGroupByRet[i]; !ok {
			OpFuncGroupByRet[i] = make(map[string]types.OpFuncEval)
		}
		if _, ok := OpFuncGroupByRet[i][o.GetName()]; !ok {
			OpFuncGroupByRet[i][o.GetName()] = o
		}
		// make the last non-zero bit to be zero: 1001100 => 1001000
		returnType &= (returnType - 1)
	}
}

// CompareValueExpr TODO(mahjonp) because of CompareDatum can tell us whether there exists an truncate, we can use it in `comparisionValidator`
func CompareValueExpr(a, b parser_driver.ValueExpr) int {
	charset.GetCollationByName(b.Collation())
	statementContext := &stmtctx.StatementContext{AllowInvalidDate: true}
	statementContext.IgnoreTruncate.Store(true)
	res, _ := a.Compare(statementContext, &b.Datum, collate.GetCollator(b.GetType().GetCollate()))
	return res
}

// TransStringType TODO: decimal NOT Equals to float
func TransStringType(s string) uint64 {
	s = strings.ToLower(s)
	switch {
	case strings.Contains(s, "string"), strings.Contains(s, "char"), strings.Contains(s, "text"), strings.Contains(s, "json"):
		return types.TypeStringArg
	case strings.Contains(s, "int"), strings.Contains(s, "long"), strings.Contains(s, "short"), strings.Contains(s, "tiny"):
		return types.TypeIntArg
	case strings.Contains(s, "float"), strings.Contains(s, "decimal"), strings.Contains(s, "double"):
		return types.TypeFloatArg
	case strings.Contains(s, "time"), strings.Contains(s, "date"):
		return types.TypeDatetimeArg
	default:
		panic(fmt.Sprintf("no implement for type: %s", s))
	}
}

// TransToMysqlType is to be transferred into mysql type
func TransToMysqlType(i uint64) byte {
	switch i {
	case types.TypeIntArg:
		return mysql.TypeLong
	case types.TypeFloatArg:
		return mysql.TypeDouble
	case types.TypeDatetimeArg:
		return mysql.TypeDatetime
	case types.TypeStringArg:
		return mysql.TypeVarchar
	default:
		panic(fmt.Sprintf("no implement this type: %d", i))
	}
}

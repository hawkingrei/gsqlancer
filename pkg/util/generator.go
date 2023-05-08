package util

import (
	"github.com/pingcap/tidb/parser/charset"
	"github.com/pingcap/tidb/sessionctx/stmtctx"
	parser_driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/pingcap/tidb/util/collate"
)

// CompareValueExpr TODO(mahjonp) because of CompareDatum can tell us whether there exists an truncate, we can use it in `comparisionValidator`
func CompareValueExpr(a, b parser_driver.ValueExpr) int {
	charset.GetCollationByName(b.Collation())
	statementContext := &stmtctx.StatementContext{AllowInvalidDate: true}
	statementContext.IgnoreTruncate.Store(true)
	res, _ := a.Compare(statementContext, &b.Datum, collate.GetCollator(b.GetType().GetCollate()))
	return res
}

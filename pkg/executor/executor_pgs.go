package executor

import (
	"github.com/hawkingrei/gsqlancer/pkg/knownbugs"
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
	"github.com/pingcap/tidb/parser/ast"
	"go.uber.org/zap"
)

func (e *Executor) DoPGS() (success bool) {
	logging.StatusLog().Info("do pgs")
	rawPivotRows, usedTables, err := e.ChoosePivotedRow()
	if err != nil {
		logging.StatusLog().Fatal("choose pivot row failed", zap.Error(err))
	}
	if len(rawPivotRows) == 0 {
		logging.StatusLog().Info("no pivot row, restart")
		return false
	}
	//logging.StatusLog().Info("rawPivotRows", zap.Any("rows", rawPivotRows), zap.Any("tables", usedTables))
	selectAST, selectSQL, columns, pivotRows, err := e.gen.GenPQSSelectStmt(rawPivotRows, usedTables)
	if err != nil {
		logging.StatusLog().Fatal("generate PQS statement error", zap.Error(err))
	}
	resultRows, err := e.conn.Select(e.ctx, selectSQL)
	if err != nil {
		logging.StatusLog().Info("select", zap.Error(err))
		return false
	}

	correct := e.verifyPQS(pivotRows, columns, resultRows)
	if correct {
		dust := knownbugs.NewDustbin([]ast.Node{selectAST}, pivotRows)
		if dust.IsKnownBug() {
			// TODO: add known bug log
			return true
		}
		return true
	}
	return true
}

package executor

import (
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/transformer"
	"github.com/hawkingrei/gsqlancer/pkg/util"
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
	"github.com/pingcap/tidb/parser/ast"
	tidbmodel "github.com/pingcap/tidb/parser/model"
	"go.uber.org/zap"
)

func (e *Executor) DoNoRECAndTLP(approach testingApproach) bool {
	usedTables := e.state.GetRandTables()
	if len(usedTables) > 2 {
		cnt := mathutil.Min(rand.Intn(len(usedTables)-2)+2, e.cfg.SelectDepth())
		usedTables = util.ChoiceSubset(usedTables, cnt)
	}
	var columns []*model.Column
	for _, table := range usedTables {
		column := util.Choice(table.Columns())
		column.Table = tidbmodel.NewCIStr(table.Name())
		columns = append(columns, column)
	}

	//log.Info("DoNoRECAndTLP", zap.Any("usedTables", usedTables))
	selectStmtNode, _, err := e.gen.SelectTable(columns, usedTables)
	if err != nil {
		logging.StatusLog().Error("generate normal SQL statement failed", zap.Error(err))
	}
	var transformers []transformer.Transformer
	switch approach {
	case approachNoREC:
		transformers = append(transformers, transformer.NoREC)
	case approachTLP:
		transformers = append(
			transformers,
			&transformer.TLPTrans{
				Expr: &ast.ParenthesesExpr{Expr: e.gen.selectGen.ConditionClause()},
				Tp:   transformer.RandTLPType(),
			},
		)
	default:
		logging.StatusLog().Fatal("unknown approach")
	}

	nodesArr := transformer.RandTransformer(transformers...).Transform([]ast.ResultSetNode{selectStmtNode})
	if len(nodesArr) < 2 {
		sql, _ := util.BufferOut(selectStmtNode)
		logging.StatusLog().Warn("no enough sqls were generated", zap.String("error sql", sql), zap.Int("node length", len(nodesArr)))
		return true
	}
	sqlInOneGroup := make([]string, 0)
	resultSet := make([][]connection.QueryItems, 0)
	logging.StatusLog().Info("DoNoRECAndTLP", zap.Int("nodesArr", len(nodesArr)))
	var sql string
	var resultRows []connection.QueryItems
	for _, node := range nodesArr {
		sql, err = util.BufferOut(node)
		if err != nil {
			logging.StatusLog().Error("err on restoring", zap.Error(err))
		} else {
			resultRows, err = e.conn.Select(e.ctx, sql)
			if err != nil {
				logging.StatusLog().Error("execSelect failed", zap.Error(err), zap.String("sql", sql))
				return false
			}
			resultSet = append(resultSet, resultRows)
		}
		sqlInOneGroup = append(sqlInOneGroup, sql)
	}
	correct := checkResultSet(resultSet, true)
	if !correct {
		logging.StatusLog().Error("last round SQLs", zap.Strings("", sqlInOneGroup))
		logging.StatusLog().Fatal("NoREC/TLP data verified failed")
	}
	if correct {
		variableGen := e.gen.setGen.Rand()
		logging.StatusLog().Info("start checking variable", zap.String("name", variableGen.Name()))
		e.conn.MustExec(e.ctx, variableGen.GenVariableOn())
		resultRows1, err := e.conn.Select(e.ctx, sql)
		if err != nil {
			logging.StatusLog().Error("execSelect failed", zap.Error(err), zap.String("sql", sql))
			return false
		}
		e.conn.MustExec(e.ctx, variableGen.GenVariableOff())
		resultRows2, err := e.conn.Select(e.ctx, sql)
		if err != nil {
			logging.StatusLog().Error("execSelect failed", zap.Error(err), zap.String("sql", sql))
			return false
		}
		correct := checkResultSet([][]connection.QueryItems{resultRows1, resultRows2}, true)
		if !correct {
			logging.StatusLog().Error("last round SQLs", zap.Strings("", sqlInOneGroup))
			logging.StatusLog().Fatal("NoREC/TLP data verified failed")
		}
	}
	logging.StatusLog().Info("check finished", zap.Int("batch", e.batch), zap.Int("round", e.roundInBatch), zap.Bool("result", correct))
	//log.L().Info("check finished", zap.String("approach", "NoREC"), zap.Int("batch", p.batch), zap.Int("round", p.roundInBatch), zap.Bool("result", correct))
	return true
}

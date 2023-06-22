package executor

import (
	"fmt"
	"math/rand"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	"github.com/hawkingrei/gsqlancer/pkg/connection/realdb"
	"github.com/hawkingrei/gsqlancer/pkg/gen"
	"github.com/hawkingrei/gsqlancer/pkg/knownbugs"
	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/report"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/hawkingrei/gsqlancer/pkg/util"
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/dumpling/context"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/util/mathutil"
	"go.uber.org/zap"
)

type testingApproach = int

const (
	approachPQS testingApproach = iota
	approachNoREC
	approachTLP
)

// Executor define test executor
type Executor struct {
	cfg         *config.Config
	state       *gen.TiDBState
	exitCh      chan struct{}
	ctx         *context.Context
	conn        *realdb.DBConn
	sessionMeta *model.SessionMeta
	gen         generator
	status      *ExecutorStat
	id          int
	action      ActionType
	useDatabase string
	report      *report.Reporter

	// copy from
	batch        int
	roundInBatch int
}

func NewExecutor(id int, cfg *config.Config, exitCh chan struct{}, conn *realdb.DBConn, report *report.Reporter) *Executor {
	action := gen.NewTiDBState()
	return &Executor{
		action:      ActionCreateTableStmt,
		ctx:         context.Background(),
		state:       action,
		cfg:         cfg,
		conn:        conn,
		exitCh:      exitCh,
		id:          id,
		sessionMeta: model.NewSessionMeta(),
		status:      NewExecutorStat(),
		gen:         newGenerator(cfg, action),
		report:      report,
	}
}

func (e *Executor) Run() {
	err := e.initDatabase()
	if err != nil {
		logging.StatusLog().Error("fail to init database", zap.Error(err))
	}
	logging.StatusLog().Info("init database success", zap.Int("id", e.id))
	for {
		select {
		case <-e.exitCh:
			return
		default:
		}
		for !e.Do() {
			continue
		}
		logging.StatusLog().Info("Done")
		e.Next()
	}
}

// return whether it continues or not.
func (e *Executor) progress() bool {
	var approaches []testingApproach
	// Because we creates view just in time with process(we creates views on first ViewCount rounds)
	// and current implementation depends on the PQS to ensures that there exists at lease one row in that view
	// so we must choose approachPQS in this scenario
	if e.roundInBatch < e.cfg.ViewCount {
		approaches = []testingApproach{approachPQS}
	} else {
		if e.cfg.EnableNoRECApproach {
			approaches = append(approaches, approachNoREC)
		}
		if e.cfg.EnablePQSApproach {
			approaches = append(approaches, approachPQS)
		}
		if e.cfg.EnableTLPApproach {
			approaches = append(approaches, approachTLP)
		}
	}
	approach := approaches[rand.Intn(len(approaches))]
	var succeed bool
	switch approach {
	case approachPQS:
		succeed = e.DoPGS()
	case approachNoREC, approachTLP:

	}
	e.state.Reset()
	return succeed
}

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
	log.Info("start query")
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
	log.Info("pass")
	return true
}

// ChoosePivotedRow choose a row
// it may move to another struct
func (e *Executor) ChoosePivotedRow() (map[string]*connection.QueryItem, []*model.Table, error) {
	result := make(map[string]*connection.QueryItem)
	usedTables := e.state.GetRandTables()
	if len(usedTables) > 2 {
		usedTables = util.ChoiceSubset(usedTables, mathutil.Min(rand.Intn(len(usedTables)-2)+2, 5))
	}
	var reallyUsed []*model.Table

	for _, i := range usedTables {
		sql := fmt.Sprintf("SELECT * FROM %s ORDER BY RAND() LIMIT 1;", i.Name())
		exeRes, err := e.conn.Select(e.ctx, sql)
		if err != nil {
			logging.SQLLOG().Error(sql, zap.Error(err))
			return nil, nil, err
		}
		if len(exeRes) > 0 {
			for idx := range exeRes[0] {
				tableColumn := types.Column{Table: types.CIStr(i.Name()), Name: types.CIStr(exeRes[0][idx].ValType.Name())}
				result[tableColumn.String()] = exeRes[0][idx]
			}
			reallyUsed = append(reallyUsed, i)
		} else {
			logging.StatusLog().Info("no row", zap.String("table", i.Name()))
		}
	}
	//log.Info("ChoosePivotedRow", zap.Any("result", result), zap.Any("usedTables", usedTables))
	return result, reallyUsed, nil
}

func (e *Executor) verifyPQS(pivotRows map[string]*connection.QueryItem, columns []model.Column, resultSets []connection.QueryItems) bool {
	for _, row := range resultSets {
		if e.checkRow(pivotRows, columns, row) {
			return true
		}
	}
	return false
}

func (e *Executor) checkRow(pivotRows map[string]*connection.QueryItem, columns []model.Column, resultSet connection.QueryItems) bool {
	for i, c := range columns {
		//logging.StatusLog().Info(fmt.Sprintf("i: %d, column: %+v, left: %+v, right: %+v", i, c, pivotRows[c.AliasName.String()], resultSet[i]))
		_, ok := pivotRows[c.AliasName.String()]
		if !ok {
			logging.StatusLog().Fatal("pivotRows[c.AliasName.String()] not exist")
		}
		if !compareQueryItem(pivotRows[c.AliasName.String()], resultSet[i]) {
			return false
		}
	}
	return true
}

func compareQueryItem(left *connection.QueryItem, right *connection.QueryItem) bool {
	// if left.ValType.Name() != right.ValType.Name() {
	// 	return false
	// }
	if left == nil {
		logging.StatusLog().Fatal("left is nil")
	}
	if right == nil {
		logging.StatusLog().Fatal("right is nil")
	}
	if left.Null != right.Null {
		return false
	}

	return (left.Null && right.Null) || (left.ValString == right.ValString)
}

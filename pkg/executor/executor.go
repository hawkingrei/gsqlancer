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
	"github.com/pingcap/tidb/dumpling/context"
	"github.com/pingcap/tidb/parser/ast"
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
		e.Do()
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
	switch approach {
	case approachPQS:
		rawPivotRows, usedTables, err := e.ChoosePivotedRow()
		if err != nil {
			logging.StatusLog().Fatal("choose pivot row failed", zap.Error(err))
		}
		selectAST, selectSQL, columns, pivotRows, err := e.gen.GenPQSSelectStmt(rawPivotRows, usedTables)
		if err != nil {
			logging.StatusLog().Fatal("generate PQS statement error", zap.Error(err))
		}
		resultRows, err := e.conn.Select(e.ctx, selectSQL)
		if err != nil {
			return true
		}
		correct := e.verifyPQS(pivotRows, columns, resultRows)
		if correct {
			dust := knownbugs.NewDustbin([]ast.Node{selectAST}, pivotRows)
			if dust.IsKnownBug() {
				// TODO: add known bug log
				return true
			}
			return false
		}
	case approachNoREC, approachTLP:

	}
	e.state.Reset()
	return true
}

// ChoosePivotedRow choose a row
// it may move to another struct
func (e *Executor) ChoosePivotedRow() (map[string]*connection.QueryItem, []*model.Table, error) {
	result := make(map[string]*connection.QueryItem)
	usedTables := e.state.GetRandTables()
	if len(usedTables) > 2 {
		usedTables = util.ChoiceSubset(usedTables, rand.Intn(len(usedTables)-2)+2)
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
			for _, c := range exeRes[0] {
				// panic(fmt.Sprintf("no rows in table %s", i.Column))
				tableColumn := types.Column{Table: types.CIStr(i.Name()), Name: types.CIStr(c.ValType.Name())}
				result[tableColumn.String()] = c
			}
			reallyUsed = append(reallyUsed, i)
		}
	}
	return result, reallyUsed, nil
}

func (e *Executor) verifyPQS(originRow map[string]*connection.QueryItem, columns []model.Column, resultSets []connection.QueryItems) bool {
	for _, row := range resultSets {
		if e.checkRow(originRow, columns, row) {
			return true
		}
	}
	return false
}

func (e *Executor) checkRow(originRow map[string]*connection.QueryItem, columns []model.Column, resultSet connection.QueryItems) bool {
	for i, c := range columns {
		// fmt.Printf("i: %d, column: %+v, left: %+v, right: %+v", i, c, originRow[c], resultSet[i])
		if !compareQueryItem(originRow[c.AliasName.String()], resultSet[i]) {
			return false
		}
	}
	return true
}

func compareQueryItem(left *connection.QueryItem, right *connection.QueryItem) bool {
	// if left.ValType.Name() != right.ValType.Name() {
	// 	return false
	// }
	if left.Null != right.Null {
		return false
	}
	return (left.Null && right.Null) || (left.ValString == right.ValString)
}

package executor

import (
	"fmt"
	"math/rand"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	"github.com/hawkingrei/gsqlancer/pkg/gen"
	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/types"
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
	"github.com/pingcap/tidb/dumpling/context"
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
	conn        *connection.DBConn
	sessionMeta *model.SessionMeta
	gen         generator
	status      *ExecutorStat
	id          int
	action      ActionType
	useDatabase string
	// copy from
	batch        int
	roundInBatch int
}

func NewExecutor(id int, cfg *config.Config, exitCh chan struct{}, conn *connection.DBConn) *Executor {
	action := gen.NewTiDBState()
	return &Executor{
		ctx:         context.Background(),
		state:       action,
		cfg:         cfg,
		conn:        conn,
		exitCh:      exitCh,
		id:          id,
		sessionMeta: model.NewSessionMeta(),
		status:      NewExecutorStat(),
		gen:         newGenerator(cfg, action),
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
		e.Next()
		e.Do()
	}
}

func (e *Executor) progress() {
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
	case approachNoREC, approachTLP:
	}
}

// ChoosePivotedRow choose a row
// it may move to another struct
func (e *Executor) ChoosePivotedRow() (map[string]*connection.QueryItem, []*model.Table, error) {
	result := make(map[string]*connection.QueryItem)
	usedTables := e.state.GetResultTable()
	var reallyUsed []*model.Table

	for _, i := range usedTables {
		sql := fmt.Sprintf("SELECT * FROM %s ORDER BY RAND() LIMIT 1;", i.Name())
		exeRes, err := e.conn.Select(e.ctx, sql)
		if err != nil {
			panic(err)
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

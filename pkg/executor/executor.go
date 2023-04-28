package executor

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	gen2 "github.com/hawkingrei/gsqlancer/pkg/gen"
	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/dumpling/context"
	"go.uber.org/zap"
)

// Executor define test executor
type Executor struct {
	cfg    *config.Config
	action *gen2.TiDBState
	exitCh chan struct{}
	ctx    *context.Context
	conn   *connection.DBConn

	tables map[string]model.Table
	gen    generator
	id     int
}

func NewExecutor(id int, cfg *config.Config, exitCh chan struct{}, conn *connection.DBConn) *Executor {
	action := gen2.NewTiDBState()
	return &Executor{
		ctx:    context.Background(),
		action: action,
		cfg:    cfg,
		conn:   conn,
		exitCh: exitCh,
		id:     id,
		gen:    newGenerator(cfg, action),
	}
}

func (e *Executor) Run() {
	err := e.initDatabase()
	if err != nil {
		log.Error("fail to init database", zap.Error(err))
	}
	log.Info("init database success", zap.Int("id", e.id))
}

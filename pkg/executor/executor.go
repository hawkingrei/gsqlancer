package executor

import (
	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/dumpling/context"
	"go.uber.org/zap"
)

// Executor define test executor
type Executor struct {
	id     int
	cfg    *config.Config
	action *TiDBState
	exitCh chan struct{}
	ctx    *context.Context
	conn   *connection.DBConn
}

func NewExecutor(id int, cfg *config.Config, exitCh chan struct{}, conn *connection.DBConn) *Executor {
	return &Executor{
		ctx:    context.Background(),
		action: NewTiDBState(),
		cfg:    cfg,
		conn:   conn,
		exitCh: exitCh,
		id:     id,
	}
}

func (e *Executor) Run() {
	err := e.initDatabase()
	if err != nil {
		log.Error("fail to init database", zap.Error(err))
	}
	log.Info("init database success", zap.Int("id", e.id))
}

package executor

import "github.com/hawkingrei/gsqlancer/pkg/config"

// Executor define test executor
type Executor struct {
	id     int
	cfg    *config.Config
	action *TiDBState
	exitCh chan struct{}
}

func NewExecutor(id int, cfg *config.Config, exitCh chan struct{}) *Executor {
	return &Executor{
		id:     id,
		cfg:    cfg,
		exitCh: exitCh,
		action: NewTiDBState(),
	}
}

func (e *Executor) Run() {
}

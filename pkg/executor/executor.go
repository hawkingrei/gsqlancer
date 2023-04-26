package executor

import "github.com/hawkingrei/gsqlancer/pkg/config"

// Executor define test executor
type Executor struct {
	id     int
	db     int32
	cfg    *config.Config
	exitCh chan struct{}
}

func NewExecutor(id int, cfg *config.Config, exitCh chan struct{}) *Executor {
	return &Executor{
		id:     id,
		cfg:    cfg,
		exitCh: exitCh,
	}
}

func (e *Executor) Run() {

}

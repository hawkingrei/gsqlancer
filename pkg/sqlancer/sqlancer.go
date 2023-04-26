package sqlancer

import (
	"sync"

	"github.com/hawkingrei/gsqlancer/pkg/config"
)

type SQLancer struct {
	wg  sync.WaitGroup
	cfg *config.Config

	exitCh chan struct{}
}

func NewSQLancer(cfg *config.Config) *SQLancer {
	return &SQLancer{
		cfg: cfg,
	}
}

func (s *SQLancer) Run() {

}

func (s *SQLancer) Stop() {
	close(s.exitCh)
	s.wg.Wait()
}

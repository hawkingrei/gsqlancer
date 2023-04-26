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
		cfg:    cfg,
		exitCh: make(chan struct{}),
	}
}

func (s *SQLancer) Run() {
	for i := 0; i < int(s.cfg.Concurrency()); i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
		}()
	}
}

func (s *SQLancer) Stop() {
	close(s.exitCh)
	s.wg.Wait()
}

package sqlancer

import (
	"time"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	tidbutil "github.com/pingcap/tidb/util"
)

type SQLancer struct {
	wg  tidbutil.WaitGroupWrapper
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
	s.wg.Run(s.tick)
}

func (s *SQLancer) Stop() {
	close(s.exitCh)
	s.wg.Wait()
}

func (s *SQLancer) tick() {
	ticker := time.NewTicker(s.cfg.MaxTestTime())
	defer ticker.Stop()
	for {
		select {
		case <-s.exitCh:
			return
		case <-ticker.C:
			close(s.exitCh)
			return
		}
	}
}

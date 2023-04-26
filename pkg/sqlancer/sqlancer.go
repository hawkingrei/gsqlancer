package sqlancer

import (
	"time"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/connection"
	"github.com/hawkingrei/gsqlancer/pkg/executor"
	"github.com/pingcap/log"
	tidbutil "github.com/pingcap/tidb/util"
	"go.uber.org/zap"
)

type SQLancer struct {
	wg     tidbutil.WaitGroupWrapper
	cfg    *config.Config
	dbConn *connection.DBConnect

	exitCh chan struct{}
}

func NewSQLancer(cfg *config.Config) *SQLancer {
	return &SQLancer{
		cfg:    cfg,
		exitCh: make(chan struct{}),
		dbConn: connection.NewDBConnect(cfg.DBConfig()),
	}
}

func (s *SQLancer) Run() {
	for i := 0; i < int(s.cfg.Concurrency()); i++ {
		conn, err := s.dbConn.GetConnection()
		if err != nil {
			log.Fatal("failed to get connection", zap.Error(err))
		}
		exec := executor.NewExecutor(i, s.cfg, s.exitCh, conn)
		s.wg.Run(exec.Run)
	}
	s.wg.Run(s.tick)
}

func (s *SQLancer) Stop() {
	close(s.exitCh)
	s.wg.Wait()
}

func (s *SQLancer) tick() {
	ticker := time.NewTicker(s.cfg.MaxTestTime())
	pingTicker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-s.exitCh:
			return
		case <-ticker.C:
			close(s.exitCh)
			return
		case <-pingTicker.C:
			s.dbConn.Ping()
		}
	}
}

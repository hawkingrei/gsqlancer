package executor

import (
	"fmt"

	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
	"go.uber.org/zap"
)

func (e *Executor) initDatabase() error {
	id := e.state.GenDatabaseID()
	e.useDatabase = fmt.Sprintf("database%d", id)
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s", e.useDatabase)
	err := e.conn.ExecContext(e.ctx, sql)
	if err != nil {
		return err
	}

	sql = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", e.useDatabase)
	err = e.conn.ExecContext(e.ctx, sql)
	if err != nil {
		return err
	}
	logging.SQLLOG().Info("success to create database", zap.String("sql", sql))
	err = e.conn.ExecContext(e.ctx, "use "+e.useDatabase)
	return err
}

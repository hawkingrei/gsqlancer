package executor

import (
	"fmt"

	"github.com/pingcap/log"
	"go.uber.org/zap"
)

func (e *Executor) initDatabase() error {
	id := e.state.GenDatabaseID()
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS database%d", id)
	_, err := e.conn.ExecContext(e.ctx, sql)
	if err != nil {
		return err
	}

	sql = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS database%d", id)
	_, err = e.conn.ExecContext(e.ctx, sql)
	if err != nil {
		return err
	}
	log.Info("success to create database", zap.String("sql", sql))
	return err
}

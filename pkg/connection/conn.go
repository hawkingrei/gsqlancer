package connection

import (
	"context"
	"database/sql"

	"github.com/pingcap/log"
	"go.uber.org/zap"
)

type DBConn struct {
	conn *sql.Conn
}

func (c *DBConn) Close() error {
	return c.conn.Close()
}

func (c *DBConn) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := c.conn.ExecContext(ctx, query, args...)
	if err != nil {
		log.Error("fail to execute sql", zap.String("sql", query), zap.Error(err))
		return nil, err
	}
	log.Info("success to execute sql", zap.String("sql", query))
	return result, err
}

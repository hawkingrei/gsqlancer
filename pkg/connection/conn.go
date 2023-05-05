package connection

import (
	"context"
	"database/sql"

	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
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
		logging.SQLLOG().Error("fail to execute sql", zap.String("sql", query), zap.Error(err))
		return nil, err
	}
	logging.SQLLOG().Info("success to execute sql", zap.String("sql", query))
	return result, err
}

func (c *DBConn) MustExec(ctx context.Context, query *model.SQL) error {
	_, err := c.conn.ExecContext(ctx, query.SQLStmt)
	if err != nil {
		logging.SQLLOG().Error("fail to execute sql",
			zap.String("type", query.SQLType.String()), zap.String("sql", query.SQLStmt), zap.Error(err))
		return err
	}
	logging.SQLLOG().Info("success to execute sql",
		zap.String("type", query.SQLType.String()), zap.String("sql", query.SQLStmt))
	return nil
}

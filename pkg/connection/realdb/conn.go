package realdb

import (
	"context"
	"database/sql"
	"time"

	"github.com/hawkingrei/gsqlancer/pkg/connection"
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

func (c *DBConn) ExecContext(ctx context.Context, query string, args ...interface{}) error {
	if query == "" {
		logging.SQLLOG().Error("empty sql", zap.Stack("stack"))
		return nil
	}
	_, err := c.conn.ExecContext(ctx, query, args...)
	if err != nil {
		logging.SQLLOG().Error("fail to execute sql", zap.String("sql", query), zap.Error(err))
		return err
	}
	logging.SQLLOG().Info("success to execute sql", zap.String("sql", query))
	return err
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

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (c *DBConn) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if query == "" {
		logging.SQLLOG().Error("empty sql", zap.Stack("stack"))
		return nil, nil
	}
	rows, err := c.conn.QueryContext(ctx, query, args...)
	if err != nil {
		logging.SQLLOG().Error("fail to query sql", zap.String("sql", query), zap.Error(err), zap.Stack("stack"))
		return nil, err
	}
	if rows.Err() != nil {
		logging.SQLLOG().Error("fail to query sql", zap.String("sql", query), zap.Error(rows.Err()), zap.Stack("stack"))
		return nil, err
	}
	logging.SQLLOG().Info("success to query sql", zap.String("sql", query))
	return rows, err
}

// Select run select statement and return query result
func (c *DBConn) Select(ctx context.Context, stmt string, args ...interface{}) ([]connection.QueryItems, error) {
	if stmt == "" {
		logging.StatusLog().Error("empty sql", zap.Stack("stack"))
		return nil, nil
	}
	start := time.Now()
	rows, err := c.conn.QueryContext(ctx, stmt, args...)
	if err != nil {
		logging.StatusLog().Fatal("fail to query sql", zap.String("sql", stmt), zap.Error(err), zap.Duration("time", time.Since(start)), zap.Stack("stack"))
		return []connection.QueryItems{}, err
	}
	if rows.Err() != nil {
		logging.StatusLog().Fatal("fail to query sql", zap.String("sql", stmt), zap.Error(rows.Err()), zap.Duration("time", time.Since(start)), zap.Stack("stack"))
		return []connection.QueryItems{}, rows.Err()
	}

	columnTypes, _ := rows.ColumnTypes()
	var result []connection.QueryItems

	for rows.Next() {
		var (
			rowResultSets []interface{}
			resultRow     connection.QueryItems
		)
		for range columnTypes {
			rowResultSets = append(rowResultSets, new(interface{}))
		}
		if err := rows.Scan(rowResultSets...); err != nil {
			logging.StatusLog().Error("fail to scan sql", zap.String("sql", stmt), zap.Error(err))
		}
		for index, resultItem := range rowResultSets {
			r := *resultItem.(*interface{})
			item := connection.QueryItem{
				ValType: columnTypes[index],
			}
			if r != nil {
				bytes := r.([]byte)
				item.ValString = string(bytes)
			} else {
				item.Null = true
			}
			resultRow = append(resultRow, &item)
		}
		result = append(result, resultRow)
	}
	logging.StatusLog().Info("success to query sql", zap.String("sql", stmt), zap.Duration("duration", time.Since(start)))
	return result, nil
}

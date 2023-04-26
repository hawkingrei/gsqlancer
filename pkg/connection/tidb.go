package connection

import (
	"database/sql"
	"sync"

	"github.com/pingcap/log"
	"go.uber.org/zap"
)

// DBConnect wraps db
type DBConnect struct {
	mu  sync.Mutex
	dsn string
	db  *sql.DB
}

func NewDBConnect(config Config) {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(config.MaxLifetime)
}

func (c *DBConnect) Ping() {
	c.db.Ping()
}

func (c *DBConnect) MustExec(sql string, args ...interface{}) error {
	_, err := c.db.Exec(sql, args...)
	if err != nil {
		log.Error("exec sql failed", zap.String("sql", sql), zap.Error(err))
	}
	return err
}

func (c *DBConnect) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(sql, args...)
}

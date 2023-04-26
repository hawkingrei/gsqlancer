package connection

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/tidb/dumpling/context"
)

// DBConnect wraps db
type DBConnect struct {
	ctx context.Context
	mu  sync.Mutex
	dsn string
	db  *sql.DB
}

func NewDBConnect(config *Config) *DBConnect {
	ctx := context.Background()
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(config.MaxLifetime)
	return &DBConnect{
		ctx: *ctx,
		dsn: config.DSN,
		db:  db,
	}
}

func (c *DBConnect) Ping() {
	c.db.Ping()
}

func (c *DBConnect) Close() {
	c.db.Close()
}

func (c *DBConnect) GetConnection() (*DBConn, error) {
	conn, err := c.db.Conn(c.ctx)
	if err != nil {
		return nil, err
	}
	return &DBConn{conn: conn}, err
}

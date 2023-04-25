package connection

import (
	"database/sql"
	"sync"
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

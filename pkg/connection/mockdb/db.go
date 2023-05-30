package mockdb

import (
	"context"

	"github.com/pingcap/tidb/testkit"
)

type DBTestConn struct {
	tk *testkit.TestKit
}

func NewDBTestConn(tk *testkit.TestKit) *DBTestConn {
	return &DBTestConn{tk: tk}
}

func (c *DBTestConn) ExecContext(ctx context.Context, query string, args ...interface{}) error {
	_, err := c.tk.ExecWithContext(ctx, query, args...)
	return err
}

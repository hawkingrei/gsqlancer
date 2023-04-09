package ddltest

import (
	"testing"

	"github.com/pingcap/tidb/testkit"
)

func TestDDL(t *testing.T) {
	store := testkit.CreateMockStore(t)

	tk := testkit.NewTestKit(t, store)
	tk.MustExec("use test")
}

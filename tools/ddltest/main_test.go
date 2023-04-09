package ddltest

import (
	"testing"

	"github.com/hawkingrei/gsqlancer/pkg/config"
	"github.com/hawkingrei/gsqlancer/pkg/gen"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/testkit"
	"github.com/stretchr/testify/require"
)

func TestDDL(t *testing.T) {
	store := testkit.CreateMockStore(t)

	tk := testkit.NewTestKit(t, store)
	tk.MustExec("use test")
	generator := gen.NewTiDBTableGenerator(config.DefaultConfig())
	sql, err := generator.GenerateDDLCreateTable()
	require.NoError(t, err)
	tk.MustExec(sql.SQLStmt)
}

func TestGenerateDDL(t *testing.T) {
	generator := gen.NewTiDBTableGenerator(config.DefaultConfig())
	for i := 0; i < 200; i++ {
		sql, err := generator.GenerateDDLCreateTable()
		require.NoError(t, err)
		log.Info(sql.SQLStmt)
	}
}

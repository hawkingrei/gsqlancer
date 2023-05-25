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
	generator := gen.NewTiDBTableGenerator(config.DefaultConfig(), gen.NewTiDBState())
	sql, _, err := generator.GenerateDDLCreateTable()
	require.NoError(t, err)
	tk.MustExec(sql.SQLStmt)
}

func TestGenerateDDL(t *testing.T) {
	generator := gen.NewTiDBTableGenerator(config.DefaultConfig(), gen.NewTiDBState())
	for i := 0; i < 200; i++ {
		sql, _, err := generator.GenerateDDLCreateTable()
		require.NoError(t, err)
		log.Info(sql.SQLStmt)
	}
}

func TestGenerateInsert(t *testing.T) {
	cfg := config.DefaultConfig()
	state := gen.NewTiDBState()
	insertGen := gen.NewTiDBInsertGenerator(cfg, state)
	TableGen := gen.NewTiDBTableGenerator(cfg, state)
	sql, _, err := TableGen.GenerateDDLCreateTable()
	require.NoError(t, err)
	log.Info(sql.SQLStmt)
	for i := 0; i < 200; i++ {
		sql, err := insertGen.GenerateDMLInsertByTable("t1")
		require.NoError(t, err)
		log.Info(sql.SQLStmt)
	}
}

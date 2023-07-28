package executor

import (
	"fmt"

	"github.com/hawkingrei/gsqlancer/pkg/util"
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
	"go.uber.org/zap"
)

type ActionType int

const (
	ActionAnalyzeStmt ActionType = iota
	ActionCreateTableStmt
	ActionInsertTableStmt
	ActionSelectStmt
	ActionDeleteTable

	defaultTiflashReplicasCnt = 1
)

var AllAction = []ActionType{
	ActionAnalyzeStmt,
	ActionCreateTableStmt,
	ActionInsertTableStmt,
	ActionSelectStmt,
	ActionDeleteTable,
}

func (e *Executor) Next() {
	switch e.action {
	case ActionAnalyzeStmt:
		e.action = ActionSelectStmt
	case ActionCreateTableStmt:
		e.action = ActionInsertTableStmt
	case ActionInsertTableStmt:
		e.action = ActionAnalyzeStmt
	case ActionSelectStmt:
		e.action = util.Choice([]ActionType{ActionCreateTableStmt, ActionDeleteTable})
	case ActionDeleteTable:
		e.action = ActionInsertTableStmt
	}
}

func (e *Executor) Do() bool {
	switch e.action {
	case ActionAnalyzeStmt:
		e.gen.AnalyzeTable()
		logging.StatusLog().Info("ActionAnalyzeStmt")
	case ActionCreateTableStmt:
		logging.StatusLog().Info("ActionCreateTableStmt")
		if e.state.TableMetaCnt() >= e.cfg.MaxCreateTable {
			return true
		}
		sql, table, err := e.gen.CreateTable()
		if err != nil {
			logging.StatusLog().Info("fail to create table", zap.Error(err))
		}
		err = e.conn.MustExec(e.ctx, sql)
		if err != nil {
			panic(err)
		}
		e.state.AddTableMeta(table.Name(), table)
		if e.cfg.EnableTiflashReplicas() {
			e.gen.TiflashReplicaStmt(table.Name(), defaultTiflashReplicasCnt)
		}
		if err != nil {
			panic(err)
		}
		logging.StatusLog().Info("ActionCreateTableStmt", zap.String("name", table.Name()))
		e.state.SetRecentCreateTableName(table.Name())
	case ActionInsertTableStmt:
		table := ""
		if t := e.state.GetRecentCreateTableName(); t != "" {
			table = t
		} else {
			tbl := e.state.GetRandTable()
			table = tbl.Name()
		}
		sql, err := e.gen.InsertTable(table)
		if err != nil {
			logging.StatusLog().Error("fail to gen insert table", zap.Error(err))
		}
		e.conn.MustExec(e.ctx, sql)
		logging.StatusLog().Info("ActionInsertTableStmt")
	case ActionSelectStmt:
		if e.state.TableMetaCnt() > 2 {
			logging.StatusLog().Info("ActionSelectStmt")
			return e.progress()
		}
	case ActionDeleteTable:
		logging.StatusLog().Info("ActionDeleteTable")
		for _, tbl := range e.state.GetAllTableName() {
			sql := "select count(*) from " + tbl
			var line uint64
			e.conn.QueryRowContext(e.ctx, sql).Scan(&line)
			if line < 25 {
				continue
			}
			sql = fmt.Sprintf("delete from %s limit %d", tbl, line/2)
			err := e.conn.ExecContext(e.ctx, sql)
			if err != nil {
				logging.StatusLog().Error("fail to delete table", zap.Error(err))
				return false
			} else {
				logging.StatusLog().Debug("delete data", zap.String("table", tbl))
			}
		}
	}
	return true
}

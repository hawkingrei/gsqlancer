package executor

import (
	"github.com/hawkingrei/gsqlancer/pkg/util/logging"
	"go.uber.org/zap"
)

type ActionType int

const (
	ActionAnalyzeStmt ActionType = iota
	ActionCreateTableStmt
	ActionInsertTableStmt
	ActionSelectStmt

	defaultTiflashReplicasCnt = 1
)

func (e *Executor) Next() {
	switch e.action {
	case ActionAnalyzeStmt:
		e.action = ActionCreateTableStmt
	case ActionCreateTableStmt:
		e.action = ActionAnalyzeStmt
	case ActionInsertTableStmt:
		e.action = ActionSelectStmt
	case ActionSelectStmt:
		e.action = ActionAnalyzeStmt
	}
}

func (e *Executor) Do() {
	switch e.action {
	case ActionAnalyzeStmt:
		e.gen.AnalyzeTable()
	case ActionCreateTableStmt:
		sql, table, err := e.gen.CreateTable()
		if err != nil {
			logging.StatusLog().Error("fail to create table", zap.Error(err))
		}
		err = e.conn.MustExec(e.ctx, sql)
		if err != nil {
			panic(err)
		}
		e.state.AddTableMeta(table.Name(), table)
		if e.cfg.EnableTiflashReplicas() {
			e.gen.TiflashReplicaStmt(table.Name(), defaultTiflashReplicasCnt)
		}

	}
}

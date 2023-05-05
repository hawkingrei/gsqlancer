package executor

type ActionType int

const (
	ActionAnalyzeStmt ActionType = iota
	ActionCreateTableStmt
)

func (e *Executor) Next() {
	switch e.action {
	case ActionAnalyzeStmt:

	case ActionCreateTableStmt:

	}
}

func (e *Executor) Do(gen *generator) {
	switch e.action {
	case ActionAnalyzeStmt:
		gen.AnalyzeTable()
	case ActionCreateTableStmt:
		gen.CreateTable()
	}
}

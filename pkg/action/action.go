package action

import "github.com/hawkingrei/gsqlancer/pkg/config"

type ActionType int

const (
	ActionAnalyze ActionType = iota
	ActionCreateTable
)

type Action struct {
	cfg    *config.Config
	action ActionType
}

func (a *Action) Next() {
	switch a.action {
	case ActionAnalyze:
	case ActionCreateTable:

	}
}

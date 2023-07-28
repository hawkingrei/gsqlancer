package setgen

import (
	"github.com/hawkingrei/gsqlancer/pkg/gen/setgen/variablegen"
	"github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/util"
)

type VariableGen interface {
	Name() string
	GenVariableOn() *model.SQL
	GenVariableOff() *model.SQL
	IsSwitch() bool
	Gen() *model.SQL
}

var GlobalVariableGen = []VariableGen{
	variablegen.NewBoolGen("tidb_enable_non_prepared_plan_cache", true),
}

type SetVariableGen struct{}

func NewSetVariableGen() *SetVariableGen {
	return &SetVariableGen{}
}

func (s *SetVariableGen) Gen() *model.SQL {
	return util.Choice(GlobalVariableGen).Gen()
}

func (s *SetVariableGen) Rand() VariableGen {
	return util.Choice(GlobalVariableGen)
}

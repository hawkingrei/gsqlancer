package setgen

import (
	"github.com/hawkingrei/gsqlancer/pkg/gen/setgen/variablegen"
	"github.com/hawkingrei/gsqlancer/pkg/util"
)

type VariableGen interface {
	GenVariableOn() string
	GenVariableOff() string
	IsSwitch() bool
	Gen() string
}

var GlobalVariableGen = []VariableGen{
	variablegen.NewBoolGen("tidb_enable_non_prepared_plan_cache", true),
}

type SetVariableGen struct{}

func NewSetVariableGen() *SetVariableGen {
	return &SetVariableGen{}
}

func (s *SetVariableGen) Gen() string {
	return util.Choice(GlobalVariableGen).Gen()
}

package variablegen

import (
	"fmt"

	"github.com/hawkingrei/gsqlancer/pkg/util"
)

type BoolGen struct {
	name         string
	defaultStats bool
	status       bool
}

func NewBoolGen(name string, defaultStats bool) *BoolGen {
	return &BoolGen{
		name:         name,
		defaultStats: defaultStats,
		status:       defaultStats,
	}
}

func (b *BoolGen) GenVariableOn() string {
	b.status = true
	return fmt.Sprintf("set %s=1", b.name)
}

func (b *BoolGen) GenVariableOff() string {
	b.status = false
	return fmt.Sprintf("set %s=0", b.name)
}

func (b *BoolGen) IsSwitch() bool {
	return true
}

func (b *BoolGen) Gen() string {
	value := util.RdBool()
	defer func() {
		b.status = value
	}()
	if value {
		return b.GenVariableOn()
	}
	return b.GenVariableOff()
}

func (b *BoolGen) Name() string {
	return b.name
}

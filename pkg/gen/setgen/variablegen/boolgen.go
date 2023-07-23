package variablegen

import "fmt"

type BoolGen struct {
	name         string
	defaultStats bool
}

func NewBoolGen(name string, defaultStats bool) *BoolGen {
	return &BoolGen{
		name:         name,
		defaultStats: defaultStats,
	}
}

func (b *BoolGen) GenVariableOn() string {
	return fmt.Sprintf("set %s=1", b.name)
}

func (b *BoolGen) GenVariableOff() string {
	return fmt.Sprintf("set %s=0", b.name)
}

func (b *BoolGen) IsSwitch() bool {
	return true
}

func (b *BoolGen) Gen() string {
	if b.defaultStats {
		return b.GenVariableOn()
	}
	return b.GenVariableOff()
}

func (b *BoolGen) Name() string {
	return b.name
}

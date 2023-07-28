package variablegen

import (
	"fmt"

	"github.com/hawkingrei/gsqlancer/pkg/model"
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

func (b *BoolGen) GenVariableOn() *model.SQL {
	b.status = true
	return &model.SQL{
		SQLType: model.SQLTypeSetVariable,
		SQLStmt: fmt.Sprintf("set %s=1", b.name),
	}
}

func (b *BoolGen) GenVariableOff() *model.SQL {
	b.status = false
	return &model.SQL{
		SQLType: model.SQLTypeSetVariable,
		SQLStmt: fmt.Sprintf("set %s=1", b.name),
	}
}

func (b *BoolGen) IsSwitch() bool {
	return true
}

func (b *BoolGen) Gen() *model.SQL {
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

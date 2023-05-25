package gen

import (
	"math/rand"

	gmodel "github.com/hawkingrei/gsqlancer/pkg/model"
	"github.com/hawkingrei/gsqlancer/pkg/util"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/types"
)

var (
	stringEnums = []string{"", "0", "-0", "1", "-1", "true", "false", " ", "NULL", "2020-02-02 02:02:00", "0000-00-00 00:00:00"}
	intEnums    = []int{0, 1, -1, 65535}
	floatEnums  = []float64{0, 0.1, 1.0000, -0.1, -1.0000, 0.5, 1.5}
	timeEnums   = []types.CoreTime{
		types.FromDate(0, 0, 0, 0, 0, 0, 0),
		types.FromDate(1, 1, 1, 0, 0, 0, 0),
		types.FromDate(2020, 2, 2, 2, 2, 2, 0),
		types.FromDate(9999, 9, 9, 9, 9, 9, 0),
	}
)

// GenerateEnumDataItem gets enum data interface with given type
func GenerateEnumDataItem(column gmodel.Column) interface{} {
	var res interface{}
	if !column.HasOption(ast.ColumnOptionNotNull) && util.RdRange(0, 3) == 0 {
		return nil
	}
	switch column.Type().String() {
	case "varchar":
		res = stringEnums[rand.Intn(len(stringEnums))]
	case "text":
		res = stringEnums[rand.Intn(len(stringEnums))]
	case "int":
		res = intEnums[rand.Intn(len(intEnums))]
	case "datetime":
		res = types.NewTime(timeEnums[rand.Intn(len(timeEnums))], mysql.TypeDatetime, 0)
	case "timestamp":
		res = types.NewTime(timeEnums[rand.Intn(len(timeEnums))], mysql.TypeDatetime, 0)
	case "float":
		res = floatEnums[rand.Intn(len(floatEnums))]
	}
	return res
}

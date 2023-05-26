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
	if !column.HasOption(ast.ColumnOptionNotNull) && util.RdRange(0, 5) == 0 {
		return nil
	}
	switch column.Type().GetType() {
	case mysql.TypeVarchar:
		res = stringEnums[rand.Intn(len(stringEnums))]
	case mysql.TypeString:
		res = stringEnums[rand.Intn(len(stringEnums))]
	case mysql.TypeTiny:
		return util.RdRange(-128, 127)
	case mysql.TypeShort:
		return util.RdRange(-32768, 32767)
	case mysql.TypeInt24:
		return util.RdRange(-8388608, 8388607)
	case mysql.TypeLong, mysql.TypeLonglong:
		return util.RdRange(-2147483648, 2147483647)
	case mysql.TypeDatetime:
		res = types.NewTime(timeEnums[rand.Intn(len(timeEnums))], mysql.TypeDatetime, 0)
	case mysql.TypeTimestamp:
		res = types.NewTime(timeEnums[rand.Intn(len(timeEnums))], mysql.TypeDatetime, 0)
	case mysql.TypeDouble, mysql.TypeFloat:
		res = floatEnums[rand.Intn(len(floatEnums))]
	}
	return res
}

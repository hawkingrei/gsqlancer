package gendata

import (
	"strings"

	"github.com/google/uuid"
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

var uuidGen = uuid.New()

func GetUUID() string {
	return strings.ToUpper(uuidGen.String())
}

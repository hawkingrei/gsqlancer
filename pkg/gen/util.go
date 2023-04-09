package gen

import (
	"bytes"
	"math/rand"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	"github.com/pingcap/tidb/parser/mysql"
)

// DataType2Len is to convert data type into length
func DataType2Len(t string) int {
	switch t {
	case "int":
		return 16
	case "bigint":
		return 64
	case "varchar":
		return 511
	case "timestamp":
		return 255
	case "datetime":
		return 255
	case "text":
		return 511
	case "float":
		return 53
	}
	return 16
}

var (
	autoOptTypes = []ast.ColumnOptionType{ast.ColumnOptionAutoIncrement, ast.ColumnOptionAutoRandom}
)

func randColumnOptionAuto() *ast.ColumnOption {
	opt := &ast.ColumnOption{Tp: autoOptTypes[rand.Intn(len(autoOptTypes))]}
	if opt.Tp == ast.ColumnOptionAutoRandom {
		// TODO: support auto_random with random shard bits
		opt.AutoRandOpt.ShardBits = 5
		opt.AutoRandOpt.RangeBits = 64
	}
	return opt
}

// BufferOut parser ast node to SQL string
func BufferOut(node ast.Node) (string, error) {
	out := new(bytes.Buffer)
	err := node.Restore(format.NewRestoreCtx(format.RestoreStringDoubleQuotes, out))
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Type2Tp conver type string to tp byte
// TODO: complete conversion map
func Type2Tp(t string) byte {
	switch t {
	case "int":
		return mysql.TypeLong
	case "bigint":
		return mysql.TypeLonglong
	case "varchar":
		return mysql.TypeVarchar
	case "timestamp":
		return mysql.TypeTimestamp
	case "datetime":
		return mysql.TypeDatetime
	case "text":
		return mysql.TypeBlob
	case "float":
		return mysql.TypeFloat
	}
	return mysql.TypeNull
}

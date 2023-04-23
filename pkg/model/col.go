package model

import (
	"math/rand"

	"github.com/hawkingrei/gsqlancer/pkg/util"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/parser/ast"
	parsermodel "github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/parser/mysql"
	parserTypes "github.com/pingcap/tidb/parser/types"
	"go.uber.org/zap"
)

var supportKind = []int{
	util.KindTINYINT,
	util.KindSMALLINT,
	util.KindMEDIUMINT,
	util.KindInt32,
	util.KindBigInt,
	util.KindFloat,
	util.KindDouble,
}

func Type2Tp(t int) byte {
	switch t {
	case util.KindDouble:
		return mysql.TypeDouble
	case util.KindMEDIUMINT:
		return mysql.TypeInt24
	case util.KindInt32:
		return mysql.TypeLong
	case util.KindTINYINT:
		return mysql.TypeTiny
	case util.KindSMALLINT:
		return mysql.TypeShort
	case util.KindBigInt:
		return mysql.TypeLonglong
	case util.KindVarChar:
		return mysql.TypeVarchar
	case util.KindTIMESTAMP:
		return mysql.TypeTimestamp
	case util.KindDATETIME:
		return mysql.TypeDatetime
	case util.KindBLOB:
		return mysql.TypeBlob
	case util.KindFloat:
		return mysql.TypeFloat
	}
	log.Fatal("not support type", zap.Int("type", t))
	return mysql.TypeNull
}

type Column struct {
	k    int
	name string
}

func RandGenColumn(name string) *Column {
	k := supportKind[rand.Intn(len(supportKind))]
	return &Column{
		k,
		name,
	}
}

func (c *Column) GetAst() *ast.ColumnDef {
	fieldType := parserTypes.NewFieldType(Type2Tp(c.k))
	//fieldType.SetFlen(DataType2Len(colType))
	return &ast.ColumnDef{
		Name: &ast.ColumnName{Name: parsermodel.NewCIStr(c.name)},
		Tp:   fieldType,
	}
}

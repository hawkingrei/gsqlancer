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

type _Column struct {
	defaultValue any
	name         string
	setValue     []string //for enum , set data type
	kind         int
	filedTypeM   int //such as:  VARCHAR(10) ,    filedTypeM = 10
	filedTypeD   int //such as:  DECIMAL(10,5) ,  filedTypeD = 5
}

func RandGenColumn(name string) *_Column {
	k := supportKind[rand.Intn(len(supportKind))]
	return &_Column{
		kind: k,
		name: name,
	}
}

func (c *_Column) GetAst() *ast.ColumnDef {
	fieldType := parserTypes.NewFieldType(Type2Tp(c.kind))
	//fieldType.SetFlen(DataType2Len(colType))
	return &ast.ColumnDef{
		Name: &ast.ColumnName{Name: parsermodel.NewCIStr(c.name)},
		Tp:   fieldType,
	}
}

func (c *_Column) canHaveDefaultValue() bool {
	switch c.kind {
	case util.KindBLOB, util.KindTINYBLOB, util.KindMEDIUMBLOB, util.KindLONGBLOB, util.KindTEXT,
		util.KindTINYTEXT, util.KindMEDIUMTEXT, util.KindLONGTEXT, util.KindJSON:
		return false
	default:
		return true
	}
}

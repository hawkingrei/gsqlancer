package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"sync/atomic"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/hawkingrei/gsqlancer/pkg/util"
)

const (
	ddlTestValueNull string = "NULL"
)

type ddlColumn struct {
	defaultValue any
	mValue       map[string]any
	dependency   *ddlColumn
	rows         *arraylist.List
	fieldType    string

	name      string
	nameOfGen string

	dependenciedCols []*ddlColumn

	setValue []string //for enum , set data type

	filedTypeM      int //such as:  VARCHAR(10) ,    filedTypeM = 10
	filedTypeD      int //such as:  DECIMAL(10,5) ,  filedTypeD = 5
	filedPrecision  int
	k               int
	indexReferences int

	renamed      int32
	deleted      atomic.Bool
	isPrimaryKey bool
}

func (col *ddlColumn) isDeleted() bool {
	return col.deleted.Load()
}

func (col *ddlColumn) setDeleted() {
	col.deleted.Load()
}

func (col *ddlColumn) isRenamed() bool {
	return atomic.LoadInt32(&col.renamed) != 0
}

func (col *ddlColumn) setRenamed() {
	atomic.StoreInt32(&col.renamed, 1)
}

func (col *ddlColumn) delRenamed() {
	atomic.StoreInt32(&col.renamed, 0)
}

func (col *ddlColumn) setDeletedRecover() {
	col.deleted.Store(false)
}

func (col *ddlColumn) getMatchedColumnDescriptor(descriptors []*ddlColumnDescriptor) *ddlColumnDescriptor {
	for _, d := range descriptors {
		if d.column == col {
			return d
		}
	}
	return nil
}

func (col *ddlColumn) getDefinition() string {
	if col.isPrimaryKey {
		if col.canHaveDefaultValue() {
			return fmt.Sprintf("%s DEFAULT %v", col.fieldType, getDefaultValueString(col.k, col.defaultValue))
		}
		return col.fieldType
	}

	if col.isGenerated() {
		return fmt.Sprintf("%s AS (JSON_EXTRACT(`%s`,'$.%s'))", col.fieldType, col.dependency.name, col.nameOfGen)
	}

	if col.canHaveDefaultValue() {
		return fmt.Sprintf("%s NULL DEFAULT %v", col.fieldType, getDefaultValueString(col.k, col.defaultValue))
	} else {
		return fmt.Sprintf("%s NULL", col.fieldType)
	}

}

func (col *ddlColumn) getSelectName() string {
	if col.k == util.KindBit {
		return fmt.Sprintf("bin(`%s`)", col.name)
	} else {
		return fmt.Sprintf("`%s`", col.name)
	}
}

// BLOB, TEXT, GEOMETRY or JSON column 'b' can't have a default value")
func (col *ddlColumn) canHaveDefaultValue() bool {
	switch col.k {
	case util.KindBLOB, util.KindTINYBLOB, util.KindMEDIUMBLOB, util.KindLONGBLOB, util.KindTEXT,
		util.KindTINYTEXT, util.KindMEDIUMTEXT, util.KindLONGTEXT, util.KindJSON:
		return false
	default:
		return true
	}
}

func (col *ddlColumn) isGenerated() bool {
	return col.dependency != nil
}

func (col *ddlColumn) notGenerated() bool {
	return col.dependency == nil
}

func (col *ddlColumn) hasGenerateCol() bool {
	return len(col.dependenciedCols) > 0
}

// randValue return a rand value of the column
func (col *ddlColumn) randValue() interface{} {
	switch col.k {
	case util.KindTINYINT:
		return rand.Int31n(1<<8) - 1<<7
	case util.KindSMALLINT:
		return rand.Int31n(1<<16) - 1<<15
	case util.KindMEDIUMINT:
		return rand.Int31n(1<<24) - 1<<23
	case util.KindInt32:
		return rand.Int63n(1<<32) - 1<<31
	case util.KindBigInt:
		if rand.Intn(2) == 1 {
			return rand.Int63()
		}
		return -1 - rand.Int63()
	case util.KindBit:
		if col.filedTypeM >= 64 {
			return fmt.Sprintf("%b", rand.Uint64())
		} else {
			m := col.filedTypeM
			if col.filedTypeM > 7 { // it is a bug
				m = m - 1
			}
			n := (int64)((1 << (uint)(m)) - 1)
			return fmt.Sprintf("%b", rand.Int63n(n))
		}
	case util.KindFloat:
		return rand.Float32() + 1
	case util.KindDouble:
		return rand.Float64() + 1
	case util.KindDECIMAL:
		return util.RandDecimal(col.filedTypeM, col.filedTypeD)
	case util.KindChar, util.KindVarChar, util.KindBLOB, util.KindTINYBLOB, util.KindMEDIUMBLOB, util.KindLONGBLOB, util.KindTEXT, util.KindTINYTEXT, util.KindMEDIUMTEXT, util.KindLONGTEXT:
		if col.filedTypeM == 0 {
			return ""
		} else {
			if col.isGenerated() {
				if col.filedTypeM <= 2 {
					return ""
				}
				return util.RandSeq(rand.Intn(col.filedTypeM - 2))
			}
			return util.RandSeq(rand.Intn(col.filedTypeM))
		}
	case util.KindBool:
		return rand.Intn(2)
	case util.KindDATE:
		randTime := time.Unix(util.MinDATETIME.Unix()+rand.Int63n(util.GapDATETIMEUnix), 0)
		return randTime.Format(util.TimeFormatForDATE)
	case util.KindTIME:
		randTime := time.Unix(util.MinTIMESTAMP.Unix()+rand.Int63n(util.GapTIMESTAMPUnix), 0)
		return randTime.Format(util.TimeFormatForTIME)
	case util.KindDATETIME:
		randTime := randTime(util.MinDATETIME, util.GapDATETIMEUnix)
		return randTime.Format(util.TimeFormat)
	case util.KindTIMESTAMP:
		randTime := randTime(util.MinTIMESTAMP, util.GapTIMESTAMPUnix)
		return randTime.Format(util.TimeFormat)
	case util.KindYEAR:
		return rand.Intn(254) + 1901 //1901 ~ 2155
	case util.KindJSON:
		return col.randJsonValue()
	case util.KindEnum:
		setLen := len(col.setValue)
		if setLen == 0 {
			// Type is change.
			return col.randValue()
		}
		i := rand.Intn(len(col.setValue))
		return col.setValue[i]
	case util.KindSet:
		var l int
		for l == 0 {
			setLen := len(col.setValue)
			if setLen == 0 {
				// Type is change.
				return col.randValue()
			}
			l = rand.Intn(len(col.setValue))
		}
		idxs := make([]int, l)
		m := make(map[int]struct{})
		for i := 0; i < l; i++ {
			idx := rand.Intn(len(col.setValue))
			_, ok := m[idx]
			for ok {
				idx = rand.Intn(len(col.setValue))
				_, ok = m[idx]
			}
			m[idx] = struct{}{}
			idxs[i] = idx
		}
		sort.Ints(idxs)
		s := ""
		for i := range idxs {
			if i > 0 {
				s += ","
			}
			s += col.setValue[idxs[i]]
		}
		return s
	default:
		return nil
	}
}

// getDefaultValueString returns a string representation of `defaultValue` according
// to the type `kind` of column.
func getDefaultValueString(k int, defaultValue interface{}) string {
	if k == util.KindBit {
		return fmt.Sprintf("b'%v'", defaultValue)
	} else {
		return fmt.Sprintf("'%v'", defaultValue)
	}
}

func randTime(minTime time.Time, gap int64) time.Time {
	// https://github.com/chronotope/chrono-tz/issues/23
	// see all invalid time: https://timezonedb.com/time-zones/Asia/Shanghai
	var randTime time.Time
	for {
		randTime = time.Unix(minTime.Unix()+rand.Int63n(gap), 0).In(util.Local)
		if util.NotAmbiguousTime(randTime) {
			break
		}
	}
	return randTime
}

func (col *ddlColumn) randJsonValue() string {
	for _, dCol := range col.dependenciedCols {
		col.mValue[dCol.nameOfGen] = dCol.randValue()
	}
	jsonRow, _ := json.Marshal(col.mValue)
	return string(jsonRow)
}

type ddlColumnDescriptor struct {
	column *ddlColumn
	value  interface{}
}

func (ddlt *ddlColumnDescriptor) getValueString() string {
	// make bit data visible
	if ddlt.column.k == util.KindBit {
		return fmt.Sprintf("b'%v'", ddlt.value)
	} else {
		return fmt.Sprintf("'%v'", ddlt.value)
	}
}

func (ddlt *ddlColumnDescriptor) buildConditionSQL() string {
	var sql string
	if ddlt.value == ddlTestValueNull || ddlt.value == nil {
		sql += fmt.Sprintf("`%s` IS NULL", ddlt.column.name)
	} else {
		switch ddlt.column.k {
		case util.KindFloat:
			sql += fmt.Sprintf("abs(`%s` - %v) < 0.0000001", ddlt.column.name, ddlt.getValueString())
		case util.KindDouble:
			sql += fmt.Sprintf("abs(`%s` - %v) < 0.0000000000000001", ddlt.column.name, ddlt.getValueString())
		default:
			sql += fmt.Sprintf("`%s` = %v", ddlt.column.name, ddlt.getValueString())
		}
	}
	return sql
}

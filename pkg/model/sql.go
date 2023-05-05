package model

// SQLType enums for SQL types
type SQLType int

// SQLTypeDMLSelect
const (
	SQLTypeUnknown SQLType = iota
	SQLTypeReloadSchema
	SQLTypeDMLSelect
	SQLTypeDMLSelectForUpdate
	SQLTypeDMLUpdate
	SQLTypeDMLInsert
	SQLTypeDMLDelete
	SQLTypeDDLCreateTable
	SQLTypeDDLAlterTable
	SQLTypeDDLCreateIndex
	SQLTypeTxnBegin
	SQLTypeTxnCommit
	SQLTypeTxnRollback
	SQLTypeExec
	SQLTypeExit
	SQLTypeSleep
	SQLTypeCreateDatabase
	SQLTypeDropDatabase
	SQLTypeAnalyzeTable
)

// SQL struct
type SQL struct {
	SQLStmt  string `json:"sql_stmt"`
	SQLTable string
	SQLType  SQLType
	// ExecTime is for avoid lock watched interference before time out
	// useful for sleep statement
	ExecTime int
	Fail     bool `json:"fail"`
}

func (s *SQL) SetFail() {
	s.Fail = true
}

func (t SQLType) String() string {
	switch t {
	case SQLTypeReloadSchema:
		return "SQLTypeReloadSchema"
	case SQLTypeDMLSelect:
		return "SQLTypeDMLSelect"
	case SQLTypeDMLSelectForUpdate:
		return "SQLTypeDMLSelectForUpdate"
	case SQLTypeDMLUpdate:
		return "SQLTypeDMLUpdate"
	case SQLTypeDMLInsert:
		return "SQLTypeDMLInsert"
	case SQLTypeDMLDelete:
		return "SQLTypeDMLDelete"
	case SQLTypeDDLCreateTable:
		return "SQLTypeDDLCreateTable"
	case SQLTypeDDLAlterTable:
		return "SQLTypeDDLAlterTable"
	case SQLTypeDDLCreateIndex:
		return "SQLTypeDDLCreateIndex"
	case SQLTypeTxnBegin:
		return "SQLTypeTxnBegin"
	case SQLTypeTxnCommit:
		return "SQLTypeTxnCommit"
	case SQLTypeTxnRollback:
		return "SQLTypeTxnRollback"
	case SQLTypeExec:
		return "SQLTypeExec"
	case SQLTypeExit:
		return "SQLTypeExit"
	case SQLTypeSleep:
		return "SQLTypeSleep"
	case SQLTypeCreateDatabase:
		return "SQLTypeCreateDatabase"
	case SQLTypeDropDatabase:
		return "SQLTypeDropDatabase"
	default:
		return "SQLTypeUnknown"
	}
}

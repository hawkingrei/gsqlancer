package report

import "time"

type ReportResultBuilder struct {
}

type ReportResult struct {
	// prepare
	Process []string `json:"process"`
	// test Method
	Method string `json:"method"`
	// Error Message
	ErrorMessage string `json:"error_message"`
	// Error Stack if it is in panic.
	Stack string `json:"stack"`
	// Error sql is the sql that cause error.s
	ErrorSql string `json:"error_sql"`
	// Database version
	DatabaseVersion string `json:"database_version"`
	// time
	Timestamp time.Time `json:"timestamp"`
	// environment variables
	EnvironmentVariables map[string]string `json:"environment_variables"`
}

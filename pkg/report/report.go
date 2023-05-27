package report

type Report struct {
	// prepare
	Process []string
	// test Method
	Method string
	// Error Message
	ErrorMessage string
	// Error Stack if it is in panic.
	Stack string
	// Error sql is the sql that cause error.s
	ErrorSql string
}

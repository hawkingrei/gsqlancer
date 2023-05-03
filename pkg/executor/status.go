package executor

type ExecutorStat struct {
	sqlCnt int64
}

func NewExecutorStat() *ExecutorStat {
	return &ExecutorStat{}
}

func (s *ExecutorStat) AddSQLCnt() {
	s.sqlCnt++
}

func (s *ExecutorStat) SQLCnt() int64 {
	return s.sqlCnt
}

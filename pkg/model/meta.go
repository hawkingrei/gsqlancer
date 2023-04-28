package model

import "golang.org/x/exp/maps"

type SessionMeta struct {
	tables map[string]*Table
}

func NewSessionMeta() *SessionMeta {
	return &SessionMeta{
		tables: make(map[string]*Table),
	}
}

func (s *SessionMeta) Reset() {
	maps.Clear(s.tables)
}

func (s *SessionMeta) CreateTable(name string, table *Table) {
	s.tables[name] = table
}

func (s *SessionMeta) GetTable(name string) (*Table, bool) {
	table, ok := s.tables[name]
	return table, ok
}

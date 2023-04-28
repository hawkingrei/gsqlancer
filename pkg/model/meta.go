package model

type SessionMeta struct {
	tables map[string]Table
}

func NewSessionMeta() *SessionMeta {
	return &SessionMeta{
		tables: make(map[string]Table),
	}
}

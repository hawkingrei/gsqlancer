package model

type Variable interface {
	ToSQL() string
}

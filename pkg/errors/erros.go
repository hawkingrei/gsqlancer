package errors

import (
	"regexp"
	"strings"
)

type ErrorIgnore interface {
	Contains(err error) bool
	Regexp(err error) bool
}
type TiDBErrorIgnore struct {
	containFilter []string
	regexpFilter  []regexp.Regexp
}

func NewTiDBErrorIgnore() *TiDBErrorIgnore {
	return &TiDBErrorIgnore{}
}

func (e *TiDBErrorIgnore) Contains(err error) bool {
	es := err.Error()
	for _, filter := range e.containFilter {
		if strings.Contains(es, filter) {
			return true
		}
	}
	return false
}

func (e *TiDBErrorIgnore) Regexp(err error) bool {
	es := err.Error()
	for _, filter := range e.regexpFilter {
		if filter.MatchString(es) {
			return true
		}
	}
	return false
}

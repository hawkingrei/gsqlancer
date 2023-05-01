package logging

import "sync/atomic"

var globalLogger atomic.Pointer[Logger]

type Logger interface {
}

package logger

import (
	"log"

	"github.com/avareum/avareum-hubble-signer/pkg/logger/types"
)

type LocalLogger struct {
	types.Logger
}

func NewLocalLogger() *LocalLogger {
	return &LocalLogger{}
}

func (l *LocalLogger) Info(a ...any) {
	log.Println(append([]any{"INFO:"}, a...)...)
}

func (l *LocalLogger) Warn(a ...any) {
	log.Println(append([]any{"WARN:"}, a...)...)
}

func (l *LocalLogger) Err(a ...any) {
	log.Println(append([]any{"ERR:"}, a...)...)
}

package logger

import (
	"fmt"

	"github.com/avareum/avareum-hubble-signer/pkg/logger/types"
)

type LocalLogger struct {
	types.Logger
}

func NewLocalLogger() *LocalLogger {
	return &LocalLogger{}
}

func (l *LocalLogger) Info(a ...any) {
	fmt.Println(a...)
}

func (l *LocalLogger) Warn(a ...any) {
	fmt.Println(a...)
}

func (l *LocalLogger) Err(a ...any) {
	fmt.Println(a...)
}

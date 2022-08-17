package types

type Logger interface {
	Info(a ...any)
	Warn(a ...any)
	Err(a ...any)
}

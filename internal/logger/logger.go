package logger

type Logger interface {
	Logf(format string, args ...any)
}

type NoopLogger struct{}

func (NoopLogger) Logf(format string, args ...any) {}

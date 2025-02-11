package testutil

import (
	"testing"

	"github.com/nxt-fwd/kinde-go/internal/logger"
)

func NewTestLogger(t *testing.T) logger.Logger {
	return &testLogger{t}
}

type testLogger struct {
	t *testing.T
}

func (l testLogger) Logf(format string, args ...any) {
	l.t.Helper()
	l.t.Logf(format, args...)
}

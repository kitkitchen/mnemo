package mnemo

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

const (
	Debug LogLevel = iota
	Info
	Warn
	Fatal
	Panic
)

var logger = NewCharmLogger("Mnemo: ")

type (
	LogLevel int
	Logger   interface {
		Info(string)
		Debug(string)
		Warn(string)
		Error(string)
		Fatal(string)
	}
	CharmLogger struct {
		*log.Logger
	}
)

func NewCharmLogger(prefix string) *CharmLogger {
	l := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          prefix,
	})
	return &CharmLogger{l}
}

func (l *CharmLogger) Info(msg string) {
	l.Logger.Info(msg)
}

func (l *CharmLogger) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l *CharmLogger) Warn(msg string) {
	l.Logger.Warn(msg)
}

func (l *CharmLogger) Error(msg string) {
	l.Logger.Error(msg)
}

func (l *CharmLogger) Fatal(msg string) {
	l.Logger.Fatal(msg)
}

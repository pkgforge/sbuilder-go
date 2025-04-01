package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var Log *log.Logger = NewLogger()

func NewLogger() *log.Logger {
	return log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
	})
}

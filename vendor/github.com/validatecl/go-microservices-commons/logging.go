package commons

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

//ConfigureLogger configura un JSON logger por defecto con levels habilitados debug, info, warn, error
func ConfigureLogger(loggingLevel string) log.Logger {
	var logger log.Logger
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	switch loggingLevel {
	case "info":
		logger = level.NewFilter(logger, level.AllowInfo())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	logger = log.WithPrefix(logger, "time", log.TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		"2006-01-02T15:04:05:0000",
	))

	return log.WithPrefix(logger, "defaultCaller", log.DefaultCaller)
}

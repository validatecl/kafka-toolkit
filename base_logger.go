package kafka_toolkit

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Log logger base global
var Log BaseLogger

//BaseLogger interface
type BaseLogger interface {
	Debug(keyvals ...interface{})
	Info(keyvals ...interface{})
	Warn(keyvals ...interface{})
	Error(keyvals ...interface{})
	Panic(keyvals ...interface{})
	Fatal(msg string, args ...interface{})
}

type baseLogger struct {
	logger kitlog.Logger
}

func init() {
	// Inicia logging verboso de Sarama si variable KAFKA_VERBOSE esta en true
	if os.Getenv("KAFKA_VERBOSE") == "true" {
		sarama.Logger = log.New(os.Stdout, "kafka-verbose", log.LstdFlags|log.LUTC|log.Lmicroseconds)
	}
}

//NewBaseLogger construct
func NewBaseLogger(logger kitlog.Logger) BaseLogger {
	if logger == nil {
		logger = kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	}

	Log = &baseLogger{
		logger: logger,
	}
	return Log
}

func (b *baseLogger) Debug(keyvals ...interface{}) {
	level.Debug(b.logger).Log(keyvals...)
}

func (b *baseLogger) Info(keyvals ...interface{}) {
	level.Info(b.logger).Log(keyvals...)
}

func (b *baseLogger) Warn(keyvals ...interface{}) {
	level.Warn(b.logger).Log(keyvals...)
}

func (b *baseLogger) Error(keyvals ...interface{}) {
	level.Error(b.logger).Log(keyvals...)
}
func (b *baseLogger) Panic(keyvals ...interface{}) {
	level.Error(b.logger).Log("panic", "true", keyvals)
}

func (b *baseLogger) Fatal(msg string, args ...interface{}) {
	log.Fatalf(msg, args...)
}

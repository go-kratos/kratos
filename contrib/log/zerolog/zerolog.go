package zerolog

import (
	"github.com/rs/zerolog"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log *zerolog.Logger
}

func NewLogger(logger *zerolog.Logger) log.Logger {
	return &Logger{
		log: logger,
	}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) (err error) {
	var event *zerolog.Event
	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}

	switch level {
	case log.LevelDebug:
		event = l.log.Debug()
	case log.LevelInfo:
		event = l.log.Info()
	case log.LevelWarn:
		event = l.log.Warn()
	case log.LevelError:
		event = l.log.Error()
	case log.LevelFatal:
		event = l.log.Fatal()
	default:
		event = l.log.Debug()
	}

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		event = event.Any(key, keyvals[i+1])
	}
	event.Send()
	return
}

package stack

import "github.com/go-kratos/kratos/v2/log"

var _ log.Logger = (*stackLogger)(nil)

type stackLogger struct {
	loggers          []log.Logger
	ignoreExceptions bool
}

type Option func(*stackLogger)

func IgnoreExceptions() Option {
	return func(logger *stackLogger) {
		logger.ignoreExceptions = true
	}
}

func New(loggers []log.Logger, opts ...Option) log.Logger {
	logger := stackLogger{
		loggers: loggers,
	}

	for _, opt := range opts {
		opt(&logger)
	}

	return logger
}

func (s stackLogger) Log(level log.Level, keyvals ...interface{}) error {
	for _, logger := range s.loggers {
		if err := logger.Log(level, keyvals...); err != nil {
			if !s.ignoreExceptions {
				return err
			}
		}
	}

	return nil
}

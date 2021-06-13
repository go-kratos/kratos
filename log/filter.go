package log

import (
	"fmt"
)

// Filter is a logger filter.
type Filter struct {
	logger  Logger
	options options
}

// Option is tracing option.
type Option func(*options)

type options struct {
	filterLevel map[Level]struct{}
	filterKey   map[string]struct{}
	filterValue map[string]struct{}
	filterFunc  func(level Level, keyvals ...interface{}) bool
}

// WithFilterLevel with filter level
func FilterLevel(level map[Level]struct{}) Option {
	return func(opts *options) {
		opts.filterLevel = level
	}
}

// FilterKey with filter key
func FilterKey(key map[string]struct{}) Option {
	return func(opts *options) {
		opts.filterKey = key
	}
}

// WithFilterValue with filter value
func FilterValue(value map[string]struct{}) Option {
	return func(opts *options) {
		opts.filterValue = value
	}
}

// WithFilterFunc with filter func
func FilterFunc(f func(level Level, keyvals ...interface{}) bool) Option {
	return func(opts *options) {
		opts.filterFunc = f
	}
}

// NewFilter new a logger filter.
func NewFilter(logger Logger, opts ...Option) *Filter {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	return &Filter{
		logger:  logger,
		options: options,
	}
}

// Log .
func (h *Filter) Log(level Level, keyvals ...interface{}) {
	if f := filter(level, h.options, keyvals); f {
		return
	}
	h.logger.Log(level, keyvals...)
}

// Debug logs a message at debug level.
func (h *Filter) Debug(a ...interface{}) {
	if f := filter(LevelDebug, h.options, a); f {
		return
	}
	h.logger.Log(LevelDebug, "msg", fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Filter) Debugf(format string, a ...interface{}) {
	if f := filter(LevelDebug, h.options, a); f {
		return
	}
	h.logger.Log(LevelDebug, "msg", fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Filter) Debugw(keyvals ...interface{}) {
	if f := filter(LevelDebug, h.options, keyvals); f {
		return
	}
	h.logger.Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func (h *Filter) Info(a ...interface{}) {
	if f := filter(LevelInfo, h.options, a); f {
		return
	}
	h.logger.Log(LevelInfo, "msg", fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *Filter) Infof(format string, a ...interface{}) {
	if f := filter(LevelInfo, h.options, a); f {
		return
	}
	h.logger.Log(LevelInfo, "msg", fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Filter) Infow(keyvals ...interface{}) {
	if f := filter(LevelInfo, h.options, keyvals); f {
		return
	}
	h.logger.Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func (h *Filter) Warn(a ...interface{}) {
	if f := filter(LevelWarn, h.options, a); f {
		return
	}
	h.logger.Log(LevelWarn, "msg", fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Filter) Warnf(format string, a ...interface{}) {
	if f := filter(LevelWarn, h.options, a); f {
		return
	}
	h.logger.Log(LevelWarn, "msg", fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Filter) Warnw(keyvals ...interface{}) {
	if f := filter(LevelWarn, h.options, keyvals); f {
		return
	}
	h.logger.Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func (h *Filter) Error(a ...interface{}) {
	if f := filter(LevelError, h.options, a); f {
		return
	}
	h.logger.Log(LevelError, "msg", fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *Filter) Errorf(format string, a ...interface{}) {
	if f := filter(LevelError, h.options, a); f {
		return
	}
	h.logger.Log(LevelError, "msg", fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Filter) Errorw(keyvals ...interface{}) {
	if f := filter(LevelError, h.options, keyvals); f {
		return
	}
	h.logger.Log(LevelError, keyvals...)
}

func filter(level Level, options options, keyvals []interface{}) bool {
	if _, ok := options.filterLevel[level]; ok {
		return true
	}
	if len(keyvals)%2 == 0 {
		for i := 0; i < len(keyvals); i += 2 {
			if options.filterKey != nil {
				if _, ok := options.filterKey[keyvals[i].(string)]; ok {
					keyvals[i+1] = "***"
				}
			}
			if options.filterValue != nil {
				if _, ok := options.filterValue[keyvals[i].(string)]; ok {
					keyvals[i] = "***"
				}
				if _, ok := options.filterValue[keyvals[i+1].(string)]; ok {
					keyvals[i+1] = "***"
				}
			}
		}
	} else {
		for i := 0; i < len(keyvals); i++ {
			if options.filterValue != nil {
				if _, ok := options.filterValue[keyvals[i].(string)]; ok {
					keyvals[i] = "***"
				}
			}
		}
	}
	if options.filterFunc != nil {
		return options.filterFunc(level, keyvals...)
	}
	return false
}

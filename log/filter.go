package log

import (
	"fmt"
)

// Filter is a logger filter.
type Filter struct {
	logger      Logger
	filterLevel Level
	filterKey   map[string]struct{}
	filterValue map[string]struct{}
	filterFunc  func(level Level, keyvals ...interface{}) bool
}

// FilterOption is filter option.
type FilterOption func(*Filter)

// FilterLevel with filter level
func FilterLevel(level Level) FilterOption {
	return func(opts *Filter) {
		opts.filterLevel = level
	}
}

// FilterKey with filter key
func FilterKey(key ...string) FilterOption {
	return func(opts *Filter) {
		filterKey := make(map[string]struct{})
		for _, v := range key {
			filterKey[v] = struct{}{}
		}
		opts.filterKey = filterKey
	}
}

// FilterValue with filter value
func FilterValue(value ...string) FilterOption {
	return func(opts *Filter) {
		filterValue := make(map[string]struct{})
		for _, v := range value {
			filterValue[v] = struct{}{}
		}
		opts.filterValue = filterValue
	}
}

// FilterFunc with filter func
func FilterFunc(f func(level Level, keyvals ...interface{}) bool) FilterOption {
	return func(opts *Filter) {
		opts.filterFunc = f
	}
}

// NewFilter new a logger filter.
func NewFilter(logger Logger, opts ...FilterOption) *Filter {
	options := Filter{logger: logger}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

// Log .
func (h *Filter) Log(level Level, keyvals ...interface{}) {
	if f := h.filter(level, keyvals); f {
		return
	}
	h.logger.Log(level, keyvals...)
}

// Debug logs a message at debug level.
func (h *Filter) Debug(a ...interface{}) {
	if f := h.filter(LevelDebug, a); f {
		return
	}
	h.logger.Log(LevelDebug, "msg", fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Filter) Debugf(format string, a ...interface{}) {
	if f := h.filter(LevelDebug, a); f {
		return
	}
	h.logger.Log(LevelDebug, "msg", fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Filter) Debugw(keyvals ...interface{}) {
	if f := h.filter(LevelDebug, keyvals); f {
		return
	}
	h.logger.Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func (h *Filter) Info(a ...interface{}) {
	if f := h.filter(LevelInfo, a); f {
		return
	}
	h.logger.Log(LevelInfo, "msg", fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *Filter) Infof(format string, a ...interface{}) {
	if f := h.filter(LevelInfo, a); f {
		return
	}
	h.logger.Log(LevelInfo, "msg", fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Filter) Infow(keyvals ...interface{}) {
	if f := h.filter(LevelInfo, keyvals); f {
		return
	}
	h.logger.Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func (h *Filter) Warn(a ...interface{}) {
	if f := h.filter(LevelWarn, a); f {
		return
	}
	h.logger.Log(LevelWarn, "msg", fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Filter) Warnf(format string, a ...interface{}) {
	if f := h.filter(LevelWarn, a); f {
		return
	}
	h.logger.Log(LevelWarn, "msg", fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Filter) Warnw(keyvals ...interface{}) {
	if f := h.filter(LevelWarn, keyvals); f {
		return
	}
	h.logger.Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func (h *Filter) Error(a ...interface{}) {
	if f := h.filter(LevelError, a); f {
		return
	}
	h.logger.Log(LevelError, "msg", fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *Filter) Errorf(format string, a ...interface{}) {
	if f := h.filter(LevelError, a); f {
		return
	}
	h.logger.Log(LevelError, "msg", fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Filter) Errorw(keyvals ...interface{}) {
	if f := h.filter(LevelError, keyvals); f {
		return
	}
	h.logger.Log(LevelError, keyvals...)
}

func (h *Filter) filter(level Level, keyvals []interface{}) bool {
	if h.filterLevel > level {
		return true
	}
	if len(keyvals)%2 == 0 {
		for i := 0; i < len(keyvals); i += 2 {
			if h.filterKey != nil {
				if _, ok := h.filterKey[keyvals[i].(string)]; ok {
					keyvals[i+1] = "***"
				}
			}
			if h.filterValue != nil {
				if _, ok := h.filterValue[keyvals[i].(string)]; ok {
					keyvals[i] = "***"
				}
				if _, ok := h.filterValue[keyvals[i+1].(string)]; ok {
					keyvals[i+1] = "***"
				}
			}
		}
	} else {
		for i := 0; i < len(keyvals); i++ {
			if h.filterValue != nil {
				if _, ok := h.filterValue[keyvals[i].(string)]; ok {
					keyvals[i] = "***"
				}
			}
		}
	}
	if h.filterFunc != nil {
		return h.filterFunc(level, keyvals...)
	}
	return false
}

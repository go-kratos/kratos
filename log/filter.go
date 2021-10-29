package log

import "context"

// FilterOption is filter option.
type FilterOption func(*Filter)

const fuzzyStr = "***"

// FilterLevel with filter level.
func FilterLevel(level Level) FilterOption {
	return func(opts *Filter) {
		opts.level = level
	}
}

// FilterKey with filter key.
func FilterKey(key ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range key {
			o.key[v] = struct{}{}
		}
	}
}

// FilterValue with filter value.
func FilterValue(value ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range value {
			o.value[v] = struct{}{}
		}
	}
}

// FilterFunc with filter func.
func FilterFunc(f func(level Level, keyvals ...interface{}) bool) FilterOption {
	return func(o *Filter) {
		o.filter = f
	}
}

// Filter is a logger filter.
type Filter struct {
	logger Logger
	level  Level
	key    map[interface{}]struct{}
	value  map[interface{}]struct{}
	filter func(level Level, keyvals ...interface{}) bool
}

// NewFilter new a logger filter.
func NewFilter(logger Logger, opts ...FilterOption) *Filter {
	options := &Filter{
		logger: logger,
		key:    make(map[interface{}]struct{}),
		value:  make(map[interface{}]struct{}),
	}
	options = options.WithContext(context.Background())

	for _, o := range opts {
		o(options)
	}

	addSkipDepth(options, 1)
	return options
}

// WithContext returns a shallow copy of h with its context changed
// to ctx. The provided ctx must be non-nil.
func (f *Filter) WithContext(ctx context.Context) *Filter {
	switch lgr := f.logger.(type) {
	case *logger:
		return &Filter{
			logger: WithContext(ctx, lgr),
			key:    f.key,
			value:  f.value,
		}
	case *Filter:
		return &Filter{
			logger: lgr.WithContext(ctx),
			key:    f.key,
			value:  f.value,
		}
	}
	panic("f.logger must be *logger and *Filter")
}

// Log Print log by level and keyvals.
func (f *Filter) Log(level Level, keyvals ...interface{}) error {
	if level < f.level {
		return nil
	}
	if f.filter != nil && f.filter(level, keyvals...) {
		return nil
	}
	if len(f.key) > 0 || len(f.value) > 0 {
		for i := 0; i < len(keyvals); i += 2 {
			v := i + 1
			if v >= len(keyvals) {
				continue
			}
			if _, ok := f.key[keyvals[i]]; ok {
				keyvals[v] = fuzzyStr
			}
			if _, ok := f.value[keyvals[v]]; ok {
				keyvals[v] = fuzzyStr
			}
		}
	}
	return f.logger.Log(level, keyvals...)
}

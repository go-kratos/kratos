package log

// FilterOption is filter option.
type FilterOption func(*Filter)

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
	options := Filter{
		logger: logger,
		key:    make(map[interface{}]struct{}),
		value:  make(map[interface{}]struct{}),
	}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

// Log Print log by level and keyvals.
func (f *Filter) Log(level Level, keyvals ...interface{}) error {
	if f.level > level {
		return nil
	}
	if f.filter != nil && f.filter(level, keyvals...) {
		return nil
	}
	for i := 0; i < len(keyvals); i += 2 {
		iv := i + 1
		if iv >= len(keyvals) {
			continue
		}
		k := keyvals[i]
		v := keyvals[iv]
		if _, ok := f.key[k]; ok {
			keyvals[i+1] = "***"
		}
		if _, ok := f.value[v]; ok {
			keyvals[i+1] = "***"
		}
	}
	return f.logger.Log(level, keyvals...)
}

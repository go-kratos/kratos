package log

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
func FilterFunc(f func(level Level, keyvals ...any) bool) FilterOption {
	return func(o *Filter) {
		o.filter = f
	}
}

// Filter is a logger filter.
type Filter struct {
	logger Logger
	level  Level
	key    map[any]struct{}
	value  map[any]struct{}
	filter func(level Level, keyvals ...any) bool
}

// NewFilter new a logger filter.
func NewFilter(logger Logger, opts ...FilterOption) *Filter {
	options := Filter{
		logger: logger,
		key:    make(map[any]struct{}),
		value:  make(map[any]struct{}),
	}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

// Log Print log by level and keyvals.
func (f *Filter) Log(level Level, keyvals ...any) error {
	if level < f.level {
		return nil
	}
	// prefixkv contains the slice of arguments defined as prefixes during the log initialization
	var prefixkv []any
	l, ok := f.logger.(*logger)
	if ok && len(l.prefix) > 0 {
		prefixkv = make([]any, 0, len(l.prefix))
		prefixkv = append(prefixkv, l.prefix...)
	}

	if f.filter != nil && (f.filter(level, prefixkv...) || f.filter(level, keyvals...)) {
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

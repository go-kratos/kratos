package log

// Filter is a logger filter.
type Filter struct {
	logger     Logger
	level      Level
	key        map[string]struct{}
	value      map[string]struct{}
	filterHook func(level Level, keyvals ...interface{}) bool
}

// FilterOption is filter option.
type FilterOption func(*Filter)

// FilterLevel with filter level
func FilterLevel(level Level) FilterOption {
	return func(opts *Filter) {
		opts.level = level
	}
}

// FilterKeys with filter key
func FilterKeys(key ...string) FilterOption {
	return func(opts *Filter) {
		filterKey := make(map[string]struct{})
		for _, v := range key {
			filterKey[v] = struct{}{}
		}
		opts.key = filterKey
	}
}

// FilterValues with filter value
func FilterValues(value ...string) FilterOption {
	return func(opts *Filter) {
		filterValue := make(map[string]struct{})
		for _, v := range value {
			filterValue[v] = struct{}{}
		}
		opts.value = filterValue
	}
}

// FilterHook with filter func
func FilterHook(f func(level Level, keyvals ...interface{}) bool) FilterOption {
	return func(opts *Filter) {
		opts.filterHook = f
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
func (h *Filter) Log(level Level, keyvals ...interface{}) error {
	if f := h.filter(level, keyvals); f {
		return nil
	}
	return h.logger.Log(level, keyvals...)
}

func (h *Filter) filter(level Level, keyvals []interface{}) bool {
	if h.level > level {
		return true
	}
	if len(keyvals)%2 == 0 {
		for i := 0; i < len(keyvals); i += 2 {
			if h.key != nil {
				if v, ok := keyvals[i].(string); ok {
					if _, ok := h.key[v]; ok {
						keyvals[i+1] = "***"
					}
				}
			}
			if h.value != nil {
				if v, ok := keyvals[i+1].(string); ok {
					if _, ok := h.value[v]; ok {
						keyvals[i+1] = "***"
					}
				}
			}
		}
	} else {
		for i := 0; i < len(keyvals); i++ {
			if v, ok := keyvals[i].(string); ok {
				if _, ok := h.value[v]; ok {
					keyvals[i] = "***"
				}
			}
		}
	}
	if h.filterHook != nil {
		return h.filterHook(level, keyvals...)
	}
	return false
}

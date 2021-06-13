package log

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
func (h *Filter) Log(level Level, keyvals ...interface{}) error {
	if f := h.filter(level, keyvals); f {
		return nil
	}
	return h.logger.Log(level, keyvals...)
}

func (h *Filter) filter(level Level, keyvals []interface{}) bool {
	if h.filterLevel > level {
		return true
	}
	if len(keyvals)%2 == 0 {
		for i := 0; i < len(keyvals); i += 2 {
			if h.filterKey != nil {
				if v, ok := keyvals[i].(string); ok {
					if _, ok := h.filterKey[v]; ok {
						keyvals[i+1] = "***"
					}
				}
			}
			if h.filterValue != nil {
				if v, ok := keyvals[i+1].(string); ok {
					if _, ok := h.filterValue[v]; ok {
						keyvals[i+1] = "***"
					}
				}
			}
		}
	} else {
		for i := 0; i < len(keyvals); i++ {
			if v, ok := keyvals[i].(string); ok {
				if _, ok := h.filterValue[v]; ok {
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

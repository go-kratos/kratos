package log

type filter struct {
	level Level
	log   Logger
}

// NewFilter new a filter with level.
func NewFilter(l Logger, level Level) Logger {
	return &filter{log: l, level: level}
}

func (f *filter) Print(kv ...interface{}) {
	for i := 1; i < len(kv); i += 2 {
		if v, ok := kv[i].(Level); ok {
			if v < f.level {
				return
			}
			break
		}
	}
	f.log.Print(kv...)
}

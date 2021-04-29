package log

type filter struct {
	level Level
	log   Logger
}

// NewFilter new a filter with level.
func NewFilter(l Logger, level Level) Logger {
	return &filter{log: l, level: level}
}

func (f *filter) Print(kvs ...interface{}) {
	for i := 1; i < len(kvs); i += 2 {
		if v, ok := kvs[i].(Level); ok {
			if v < f.level {
				return
			}
			break
		}
	}
	f.log.Print(kvs...)
}

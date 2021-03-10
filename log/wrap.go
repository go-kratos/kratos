package log

type wrapper []Logger

func (w wrapper) Print(level Level, pairs ...interface{}) {
	for _, p := range w {
		p.Print(level, pairs...)
	}
}

// Wrap wraps multi logger.
func Wrap(l ...Logger) Logger {
	return wrapper(l)
}

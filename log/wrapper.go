package log

type wrapper []Logger

func (w wrapper) Print(pairs ...interface{}) {
	for _, p := range w {
		p.Print(pairs...)
	}
}

// Wrap wraps multi logger.
func Wrap(l ...Logger) Logger {
	return wrapper(l)
}

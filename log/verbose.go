package log

import "context"

var verbose int

// SetVerbose set the logger verbose.
func SetVerbose(v int) {
	verbose = v
}

// Verbose is a boolean type that implements logger.
type Verbose bool

// Info logs a message at info level.
func (v Verbose) Info(ctx context.Context, a ...interface{}) {
	if v {
		defaultLogger.Print(ctx, LevelInfo, a)
	}
}

// Infof logs a message at info level.
func (v Verbose) Infof(ctx context.Context, format string, a ...interface{}) {
	if v {
		defaultLogger.Printf(ctx, LevelInfo, format, a)
	}
}

// Infow logs a message at info level.
func (v Verbose) Infow(ctx context.Context, kvpair ...interface{}) {
	if v {
		defaultLogger.Printw(ctx, LevelInfo, kvpair)
	}
}

// V reports whether verbosity at the call site is at least the requested level.
func V(v int) Verbose {
	return defaultLogger.Verbose(v)
}

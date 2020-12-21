package log

import "context"

var verbose Verbose

// SetVerbose set the logger verbose.
func SetVerbose(v int) {
	verbose = Verbose(v)
}

// Verbose is a boolean type that implements logger.
type Verbose int

// Info logs a message at info level.
func (v Verbose) Info(ctx context.Context, a ...interface{}) {
	if v >= verbose {
		//		defaultLogger.Printf(ctx, LevelInfo, "", a)
	}
}

// Infof logs a message at info level.
func (v Verbose) Infof(ctx context.Context, format string, a ...interface{}) {
	if v >= verbose {
		//		defaultLogger.Printf(ctx, LevelInfo, format, a)
	}
}

// Infow logs a message at info level.
func (v Verbose) Infow(ctx context.Context, kvpair ...interface{}) {
	if v >= verbose {
		//		defaultLogger.Printf(ctx, LevelInfo, "", kvpair)
	}
}

// V reports whether verbosity at the call site is at least the requested level.
func V(v int) Verbose {
	return Verbose(v)
}

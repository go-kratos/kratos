package log

// Verbose is a verbose type that implements Logger Print.
type Verbose struct {
	log     Logger
	level   Level
	enabled bool
}

// NewVerbose new a verbose with level.
func NewVerbose(log Logger, level Level) Verbose {
	return Verbose{log: log, level: level}
}

// Enabled will return true if this log level is enabled, guarded by the value of v.
func (v Verbose) Enabled(level Level) bool {
	return v.level <= level
}

// V reports whether verbosity at the call site is at least the requested level.
func (v Verbose) V(level Level) Verbose {
	return Verbose{log: v.log, enabled: v.Enabled(level)}
}

// Print is equivalent to the Print function, guarded by the value of v.
func (v Verbose) Print(a ...interface{}) {
	if v.enabled {
		v.log.Print(a...)
	}
}

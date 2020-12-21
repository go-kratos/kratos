package log

// Option is a logger option.
type Option func(*Options)

// Options is a logger options.
type Options struct {
	Level   Level
	Verbose int
}

// WithLevel with level option.
func WithLevel(l Level) Option {
	return func(o *Options) {
		o.Level = l
	}
}

// WithVerbose with verbose option.
func WithVerbose(v int) Option {
	return func(o *Options) {
		o.Verbose = v
	}
}

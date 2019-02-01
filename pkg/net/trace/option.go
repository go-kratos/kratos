package trace

var defaultOption = option{}

type option struct {
	Debug bool
}

// Option dapper Option
type Option func(*option)

// EnableDebug enable debug mode
func EnableDebug() Option {
	return func(opt *option) {
		opt.Debug = true
	}
}

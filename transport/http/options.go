package http

// ServerOption is HTTP server option.
type ServerOption func(o *ServerOptions)

// ServerOptions is HTTP server options.
type ServerOptions struct {
	Address string
}

// WithAddress is address option.
func WithAddress(a string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = a
	}
}

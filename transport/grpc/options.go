package grpc

// ServerOption is gRPC server option.
type ServerOption func(o *ServerOptions)

// ServerOptions is gRPC server options.
type ServerOptions struct {
	Address string
}

// WithAddress is address option.
func WithAddress(a string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = a
	}
}

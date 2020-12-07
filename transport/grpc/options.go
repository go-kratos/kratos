package grpc

// ServerOption is gRPC server option.
type ServerOption func(o *serverOptions)

type serverOptions struct {
	Address string
}

// WithAddress is bind address option.
func WithAddress(a string) ServerOption {
	return func(o *serverOptions) {
		o.Address = a
	}
}

package selector

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

var globalSelector Builder

// GlobalSelector returns global selector builder.
func GlobalSelector() Builder {
	return globalSelector
}

// SetGlobalSelector set global selector builder.
func SetGlobalSelector(builder Builder) {
	globalSelector = builder
	grpc.RegisterGlobalBalancerSelector(builder)
}

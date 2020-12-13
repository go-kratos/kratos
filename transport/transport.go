package transport

import (
	"context"
)

// Server is a server interface.
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

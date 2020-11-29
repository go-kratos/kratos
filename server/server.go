package transport

import (
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
)

// Server is a server interface.
type Server interface {
	Start() error
	Stop() error

	Use(...middleware.Middleware)
}

// ServerInfo is the server request infomation.
type ServerInfo interface {
	Path() string
	Method() string
	ContentType() string
	Metadata() metadata.MD
}

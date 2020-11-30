package transport

import (
	"github.com/go-kratos/kratos/v2/metadata"
)

// Server is a server interface.
type Server interface {
	Start() error
	Stop() error
}

// ServerInfo is the server request infomation.
type ServerInfo interface {
	Path() string
	Method() string
	ContentType() string
	Metadata() metadata.MD
}

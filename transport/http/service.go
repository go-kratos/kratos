package http

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true

type methodHandler func(srv interface{}, ctx context.Context, req *http.Request, dec func(interface{}) error, m middleware.Middleware) (out interface{}, err error)

// MethodDesc represents a Proto service's method specification.
type MethodDesc struct {
	Path    string
	Method  string
	Handler methodHandler
}

// ServiceDesc represents a Proto service's specification.
type ServiceDesc struct {
	ServiceName string
	Methods     []MethodDesc
	Metadata    interface{}
}

// ServiceRegistrar wraps a single method that supports service registration.
type ServiceRegistrar interface {
	RegisterService(desc *ServiceDesc, impl interface{})
}

// RegisterService .
func (s *Server) RegisterService(desc *ServiceDesc, impl interface{}) {
	for _, m := range desc.Methods {
		h := m.Handler
		s.router.HandleFunc(m.Path, func(res http.ResponseWriter, req *http.Request) {
			out, err := h(impl, req.Context(), req, func(v interface{}) error {
				return s.requestDecoder(req, v)
			}, s.middleware)
			if err != nil {
				s.errorEncoder(res, req, err)
				return
			}
			if err := s.responseEncoder(res, req, out); err != nil {
				s.errorEncoder(res, req, err)
			}
		}).Methods(m.Method)
	}
}

package http

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/gorilla/mux"
)

// SupportPackageIsVersion1 These constants should not be referenced from any other code.
const SupportPackageIsVersion1 = true

// ServiceRegistrar wraps a single method that supports service registration.
type ServiceRegistrar interface {
	RegisterService(desc *ServiceDesc, impl interface{})
}

// ServiceDesc represents a HTTP service's specification.
type ServiceDesc struct {
	ServiceName string
	HandlerType interface{}
	Methods     []MethodDesc
	Metadata    interface{}
}

type methodHandler func(srv interface{}, ctx context.Context, m Marshaler) (interface{}, error)

// MethodDesc represents a HTTP service's method specification.
type MethodDesc struct {
	Path    string
	Method  string
	Handler methodHandler
}

// Server is a HTTP server wrapper.
type Server struct {
	router      *mux.Router
	opts        serverOptions
	middlewares map[interface{}]middleware.Middleware
}

// NewServer creates a HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	options := serverOptions{
		errorHandler:    DefaultErrorHandler,
		responseHandler: DefaultResponseHandler,
	}
	for _, o := range opts {
		o(&options)
	}
	return &Server{
		opts:        options,
		router:      mux.NewRouter(),
		middlewares: make(map[interface{}]middleware.Middleware),
	}
}

// Use .
func (s *Server) Use(srv interface{}, m middleware.Middleware) {
	s.middlewares[srv] = m
}

// Handle registers a new route with a matcher for the URL path.
func (s *Server) Handle(path string, handler http.Handler) {
	s.router.Handle(path, handler)
}

// HandleFunc registers a new route with a matcher for the URL path.
func (s *Server) HandleFunc(path string, h func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(path, h)
}

// ServeHTTP should write reply headers and data to the ResponseWriter and then return.
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(res, req)
}

// RegisterService registers a service and its implementation to the HTTP server.
func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) {
	for _, method := range sd.Methods {
		s.registerHandle(ss, method)
	}
}

func (s *Server) registerHandle(srv interface{}, md MethodDesc) {
	s.router.HandleFunc(md.Path, func(res http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		m, err := codecForReq(req)
		if err != nil {
			s.opts.errorHandler(ctx, err, m, res)
			return
		}

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return md.Handler(srv, ctx, m)
		}
		if m, ok := s.middlewares[srv]; ok {
			handler = m(handler)
		}
		if s.opts.middleware != nil {
			handler = s.opts.middleware(handler)
		}

		reply, err := handler(ctx, req)
		if err != nil {
			s.opts.errorHandler(ctx, err, m, res)
			return
		}

		s.opts.responseHandler(ctx, reply, m, res)

	}).Methods(md.Method)
}

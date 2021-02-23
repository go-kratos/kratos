package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/internal/host"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/gorilla/mux"
)

const loggerName = "transport/http"

var _ transport.Server = (*Server)(nil)

// DecodeRequestFunc deocder request func.
type DecodeRequestFunc func(req *http.Request, v interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(res http.ResponseWriter, req *http.Request, v interface{}) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(res http.ResponseWriter, req *http.Request, err error)

// ServerOption is HTTP server option.
type ServerOption func(*Server)

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(loggerName, logger)
	}
}

// Middleware with server middleware option.
func Middleware(m middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.middleware = m
	}
}

// ErrorEncoder with error handler option.
func ErrorEncoder(fn EncodeErrorFunc) ServerOption {
	return func(s *Server) {
		s.errorEncoder = fn
	}
}

// Server is a HTTP server wrapper.
type Server struct {
	*http.Server
	lis             net.Listener
	network         string
	address         string
	timeout         time.Duration
	middleware      middleware.Middleware
	requestDecoder  DecodeRequestFunc
	responseEncoder EncodeResponseFunc
	errorEncoder    EncodeErrorFunc
	router          *mux.Router
	log             *log.Helper
}

// NewServer creates a HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network:         "tcp",
		address:         ":0",
		timeout:         time.Second,
		requestDecoder:  defaultRequestDecoder,
		responseEncoder: defaultResponseEncoder,
		errorEncoder:    defaultErrorEncoder,
		middleware:      recovery.Recovery(),
		log:             log.NewHelper(loggerName, log.DefaultLogger),
	}
	for _, o := range opts {
		o(srv)
	}
	srv.router = mux.NewRouter()
	srv.Server = &http.Server{Handler: srv}
	return srv
}

// RouteGroup returns a new route group for the URL path prefix.
func (s *Server) RouteGroup(prefix string) *RouteGroup {
	return &RouteGroup{prefix: prefix, router: s.router}
}

// Handle registers a new route with a matcher for the URL path.
func (s *Server) Handle(path string, h http.Handler) {
	s.router.Handle(path, h)
}

// HandleFunc registers a new route with a matcher for the URL path.
func (s *Server) HandleFunc(path string, h http.HandlerFunc) {
	s.router.HandleFunc(path, h)
}

// PrefixHanlde  registers a new route with a matcher for the URL path prefix.
func (s *Server) PrefixHanlde(prefix string, h http.Handler) {
	s.router.PathPrefix(prefix).Handler(h)
}

// ServeHTTP should write reply headers and data to the ResponseWriter and then return.
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), s.timeout)
	defer cancel()
	ctx = transport.NewContext(ctx, transport.Transport{Kind: "HTTP"})
	ctx = NewServerContext(ctx, ServerInfo{Request: req, Response: res})

	h := func(ctx context.Context, req interface{}) (interface{}, error) {
		s.router.ServeHTTP(res, req.(*http.Request))
		return res, nil
	}
	if s.middleware != nil {
		h = s.middleware(h)
	}
	if _, err := h(ctx, req.WithContext(ctx)); err != nil {
		s.errorEncoder(res, req, err)
	}
}

// Endpoint return a real address to registry endpoint.
// examples:
//   http://127.0.0.1:8000?isSecure=false
func (s *Server) Endpoint() (string, error) {
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s", addr), nil
}

// Start start the HTTP server.
func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.lis = lis
	s.log.Infof("[HTTP] server listening on: %s", lis.Addr().String())
	if err := s.Serve(lis); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stop the HTTP server.
func (s *Server) Stop() error {
	s.log.Info("[HTTP] server stopping")
	return s.Shutdown(context.Background())
}

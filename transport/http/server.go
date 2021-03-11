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
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/gorilla/mux"
)

const loggerName = "transport/http"

var _ transport.Server = (*Server)(nil)

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

// Server is a HTTP server wrapper.
type Server struct {
	*http.Server
	lis     net.Listener
	network string
	address string
	timeout time.Duration
	router  *mux.Router
	log     *log.Helper
}

// NewServer creates a HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: time.Second,
		log:     log.NewHelper(loggerName, log.DefaultLogger),
	}
	for _, o := range opts {
		o(srv)
	}
	srv.router = mux.NewRouter()
	srv.Server = &http.Server{Handler: srv}
	return srv
}

// Handle registers a new route with a matcher for the URL path.
func (s *Server) Handle(path string, h http.Handler) {
	s.router.Handle(path, h)
}

// HandlePrefix registers a new route with a matcher for the URL path prefix.
func (s *Server) HandlePrefix(prefix string, h http.Handler) {
	s.router.PathPrefix(prefix).Handler(h)
}

// HandleFunc registers a new route with a matcher for the URL path.
func (s *Server) HandleFunc(path string, h http.HandlerFunc) {
	s.router.HandleFunc(path, h)
}

// ServeHTTP should write reply headers and data to the ResponseWriter and then return.
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), s.timeout)
	defer cancel()
	ctx = transport.NewContext(ctx, transport.Transport{Kind: transport.KindHTTP})
	ctx = NewServerContext(ctx, ServerInfo{Request: req, Response: res})
	s.router.ServeHTTP(res, req.WithContext(ctx))
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

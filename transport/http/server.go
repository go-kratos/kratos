package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	ic "github.com/go-kratos/kratos/v2/internal/context"
	"github.com/go-kratos/kratos/v2/internal/host"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/metadata/builtin"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/gorilla/mux"
)

var _ transport.Server = (*Server)(nil)
var _ transport.Endpointer = (*Server)(nil)

// ServerOption is an HTTP server option.
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
		s.log = log.NewHelper(logger)
	}
}

func MetadataExtractor(mdExtrator ExtractMetadataFunc) ServerOption {
	return func(o *Server) {
		o.mdExtractor = mdExtrator
	}
}

type ExtractMetadataFunc func(metadata.Metadata, http.Header)

func MetadataBuilder(builder metadata.Builder) ServerOption {
	return func(o *Server) {
		o.mdBuilder = builder
	}
}

// Server is an HTTP server wrapper.
type Server struct {
	*http.Server
	ctx         context.Context
	lis         net.Listener
	once        sync.Once
	err         error
	network     string
	address     string
	endpoint    *url.URL
	timeout     time.Duration
	router      *mux.Router
	log         *log.Helper
	mdExtractor ExtractMetadataFunc
	mdBuilder   metadata.Builder
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network:   "tcp",
		address:   ":0",
		timeout:   1 * time.Second,
		log:       log.NewHelper(log.DefaultLogger),
		mdBuilder: &builtin.Builder{},
		mdExtractor: func(md metadata.Metadata, header http.Header) {
			for key, values := range header {
				key = strings.ToLower(key)
				if len(values) == 1 && strings.HasPrefix(key, "x-md-") {
					key = strings.TrimLeft(key, "x-md-")
					md.Set(key, values[0])
				}
			}
		},
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
	ctx, cancel := ic.Merge(req.Context(), s.ctx)
	defer cancel()
	ctx = transport.NewContext(ctx, transport.Transport{Kind: transport.KindHTTP, Endpoint: s.endpoint.String()})
	ctx = NewServerContext(ctx, ServerInfo{Request: req, Response: res})
	md := s.mdBuilder.Build()
	s.mdExtractor(md, req.Header)
	ctx = metadata.NewContext(ctx, md)

	if s.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
	}
	s.router.ServeHTTP(res, req.WithContext(ctx))
}

// Endpoint return a real address to registry endpoint.
// examples:
//   http://127.0.0.1:8000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	s.once.Do(func() {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return
		}
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			lis.Close()
			s.err = err
			return
		}
		s.lis = lis
		s.endpoint = &url.URL{Scheme: "http", Host: addr}
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	if _, err := s.Endpoint(); err != nil {
		return err
	}
	s.ctx = ctx
	s.log.Infof("[HTTP] server listening on: %s", s.lis.Addr().String())
	if err := s.Serve(s.lis); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("[HTTP] server stopping")
	return s.Shutdown(context.Background())
}

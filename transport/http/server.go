package http

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-kratos/kratos/v2/internal/endpoint"
	"github.com/go-kratos/kratos/v2/internal/host"
	"github.com/go-kratos/kratos/v2/internal/matcher"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
	_ http.Handler         = (*Server)(nil)
)

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

// Endpoint with server address.
func Endpoint(endpoint *url.URL) ServerOption {
	return func(s *Server) {
		s.endpoint = endpoint
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// Logger with server logger.
// Deprecated: use global logger instead.
func Logger(log.Logger) ServerOption {
	return func(*Server) {}
}

// Middleware with service middleware option.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.middleware.Use(m...)
	}
}

// Filter with HTTP middleware option.
func Filter(filters ...FilterFunc) ServerOption {
	return func(o *Server) {
		o.filters = filters
	}
}

// RequestVarsDecoder with request decoder.
func RequestVarsDecoder(dec DecodeRequestFunc) ServerOption {
	return func(o *Server) {
		o.decVars = dec
	}
}

// RequestQueryDecoder with request decoder.
func RequestQueryDecoder(dec DecodeRequestFunc) ServerOption {
	return func(o *Server) {
		o.decQuery = dec
	}
}

// RequestDecoder with request decoder.
func RequestDecoder(dec DecodeRequestFunc) ServerOption {
	return func(o *Server) {
		o.decBody = dec
	}
}

// ResponseEncoder with response encoder.
func ResponseEncoder(en EncodeResponseFunc) ServerOption {
	return func(o *Server) {
		o.enc = en
	}
}

// ErrorEncoder with error encoder.
func ErrorEncoder(en EncodeErrorFunc) ServerOption {
	return func(o *Server) {
		o.ene = en
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(o *Server) {
		o.tlsConf = c
	}
}

// StrictSlash is with mux's StrictSlash
// If true, when the path pattern is "/path/", accessing "/path" will
// redirect to the former and vice versa.
func StrictSlash(strictSlash bool) ServerOption {
	return func(o *Server) {
		o.strictSlash = strictSlash
	}
}

// Listener with server lis
func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// PathPrefix with mux's PathPrefix, router will be replaced by a subrouter that start with prefix.
func PathPrefix(prefix string) ServerOption {
	return func(s *Server) {
		s.router = s.router.PathPrefix(prefix).Subrouter()
	}
}

func NotFoundHandler(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.router.NotFoundHandler = handler
	}
}

func MethodNotAllowedHandler(handler http.Handler) ServerOption {
	return func(s *Server) {
		s.router.MethodNotAllowedHandler = handler
	}
}

// Server is an HTTP server wrapper.
type Server struct {
	*http.Server
	lis         net.Listener
	tlsConf     *tls.Config
	endpoint    *url.URL
	err         error
	network     string
	address     string
	timeout     time.Duration
	filters     []FilterFunc
	middleware  matcher.Matcher
	decVars     DecodeRequestFunc
	decQuery    DecodeRequestFunc
	decBody     DecodeRequestFunc
	enc         EncodeResponseFunc
	ene         EncodeErrorFunc
	strictSlash bool
	router      *mux.Router
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network:     "tcp",
		address:     ":0",
		timeout:     1 * time.Second,
		middleware:  matcher.New(),
		decVars:     DefaultRequestVars,
		decQuery:    DefaultRequestQuery,
		decBody:     DefaultRequestDecoder,
		enc:         DefaultResponseEncoder,
		ene:         DefaultErrorEncoder,
		strictSlash: true,
		router:      mux.NewRouter(),
	}
	srv.router.NotFoundHandler = http.DefaultServeMux
	srv.router.MethodNotAllowedHandler = http.DefaultServeMux
	for _, o := range opts {
		o(srv)
	}
	srv.router.StrictSlash(srv.strictSlash)
	srv.router.Use(srv.filter())
	srv.Server = &http.Server{
		Handler:   FilterChain(srv.filters...)(srv.router),
		TLSConfig: srv.tlsConf,
	}
	return srv
}

// Use uses a service middleware with selector.
// selector:
//   - '/*'
//   - '/helloworld.v1.Greeter/*'
//   - '/helloworld.v1.Greeter/SayHello'
func (s *Server) Use(selector string, m ...middleware.Middleware) {
	s.middleware.Add(selector, m...)
}

// WalkRoute walks the router and all its sub-routers, calling walkFn for each route in the tree.
func (s *Server) WalkRoute(fn WalkRouteFunc) error {
	return s.router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		methods, err := route.GetMethods()
		if err != nil {
			return nil // ignore no methods
		}
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		for _, method := range methods {
			if err := fn(RouteInfo{Method: method, Path: path}); err != nil {
				return err
			}
		}
		return nil
	})
}

// WalkHandle walks the router and all its sub-routers, calling walkFn for each route in the tree.
func (s *Server) WalkHandle(handle func(method, path string, handler http.HandlerFunc)) error {
	return s.WalkRoute(func(r RouteInfo) error {
		handle(r.Method, r.Path, s.ServeHTTP)
		return nil
	})
}

// Route registers an HTTP router.
func (s *Server) Route(prefix string, filters ...FilterFunc) *Router {
	return newRouter(prefix, s, filters...)
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

// HandleHeader registers a new route with a matcher for the header.
func (s *Server) HandleHeader(key, val string, h http.HandlerFunc) {
	s.router.Headers(key, val).Handler(h)
}

// ServeHTTP should write reply headers and data to the ResponseWriter and then return.
func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.Handler.ServeHTTP(res, req)
}

func (s *Server) filter() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var (
				ctx    context.Context
				cancel context.CancelFunc
			)
			if s.timeout > 0 {
				ctx, cancel = context.WithTimeout(req.Context(), s.timeout)
			} else {
				ctx, cancel = context.WithCancel(req.Context())
			}
			defer cancel()

			pathTemplate := req.URL.Path
			if route := mux.CurrentRoute(req); route != nil {
				// /path/123 -> /path/{id}
				pathTemplate, _ = route.GetPathTemplate()
			}

			tr := &Transport{
				operation:    pathTemplate,
				pathTemplate: pathTemplate,
				reqHeader:    headerCarrier(req.Header),
				replyHeader:  headerCarrier(w.Header()),
				request:      req,
				response:     w,
			}
			if s.endpoint != nil {
				tr.endpoint = s.endpoint.String()
			}
			tr.request = req.WithContext(transport.NewServerContext(ctx, tr))
			next.ServeHTTP(w, tr.request)
		})
	}
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	https://127.0.0.1:8000
//	Legacy: http://127.0.0.1:8000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	log.Infof("[HTTP] server listening on: %s", s.lis.Addr().String())
	var err error
	if s.tlsConf != nil {
		err = s.ServeTLS(s.lis, "", "")
	} else {
		err = s.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Info("[HTTP] server stopping")
	return s.Shutdown(ctx)
}

func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("http", s.tlsConf != nil), addr)
	}
	return s.err
}

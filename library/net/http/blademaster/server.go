package blademaster

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/conf/dsn"
	"go-common/library/log"
	"go-common/library/net/ip"
	"go-common/library/net/metadata"
	"go-common/library/stat"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

var (
	_     IRouter = &Engine{}
	stats         = stat.HTTPServer

	_httpDSN string
)

func init() {
	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
	v := os.Getenv("HTTP")
	if v == "" {
		v = "tcp://0.0.0.0:8000/?timeout=1s"
	}
	fs.StringVar(&_httpDSN, "http", v, "listen http dsn, or use HTTP env variable.")
}

func parseDSN(rawdsn string) *ServerConfig {
	conf := new(ServerConfig)
	d, err := dsn.Parse(rawdsn)
	if err != nil {
		panic(errors.Wrapf(err, "blademaster: invalid dsn: %s", rawdsn))
	}
	if _, err = d.Bind(conf); err != nil {
		panic(errors.Wrapf(err, "blademaster: invalid dsn: %s", rawdsn))
	}
	return conf
}

// Handler responds to an HTTP request.
type Handler interface {
	ServeHTTP(c *Context)
}

// HandlerFunc http request handler function.
type HandlerFunc func(*Context)

// ServeHTTP calls f(ctx).
func (f HandlerFunc) ServeHTTP(c *Context) {
	f(c)
}

// ServerConfig is the bm server config model
type ServerConfig struct {
	Network string `dsn:"network"`
	// FIXME: rename to Address
	Addr         string         `dsn:"address"`
	Timeout      xtime.Duration `dsn:"query.timeout"`
	ReadTimeout  xtime.Duration `dsn:"query.readTimeout"`
	WriteTimeout xtime.Duration `dsn:"query.writeTimeout"`
}

// MethodConfig is
type MethodConfig struct {
	Timeout xtime.Duration
}

// Start listen and serve bm engine by given DSN.
func (engine *Engine) Start() error {
	conf := engine.conf
	l, err := net.Listen(conf.Network, conf.Addr)
	if err != nil {
		errors.Wrapf(err, "blademaster: listen tcp: %s", conf.Addr)
		return err
	}

	log.Info("blademaster: start http listen addr: %s", conf.Addr)
	server := &http.Server{
		ReadTimeout:  time.Duration(conf.ReadTimeout),
		WriteTimeout: time.Duration(conf.WriteTimeout),
	}
	go func() {
		if err := engine.RunServer(server, l); err != nil {
			if errors.Cause(err) == http.ErrServerClosed {
				log.Info("blademaster: server closed")
				return
			}
			panic(errors.Wrapf(err, "blademaster: engine.ListenServer(%+v, %+v)", server, l))
		}
	}()

	return nil
}

// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
// Create an instance of Engine, by using New() or Default()
type Engine struct {
	RouterGroup

	lock sync.RWMutex
	conf *ServerConfig

	address string

	mux       *http.ServeMux                    // http mux router
	server    atomic.Value                      // store *http.Server
	metastore map[string]map[string]interface{} // metastore is the path as key and the metadata of this path as value, it export via /metadata

	pcLock        sync.RWMutex
	methodConfigs map[string]*MethodConfig

	injections []injection
}

type injection struct {
	pattern  *regexp.Regexp
	handlers []HandlerFunc
}

// New returns a new blank Engine instance without any middleware attached.
//
// Deprecated: please use NewServer.
func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		address: ip.InternalIP(),
		conf: &ServerConfig{
			Timeout: xtime.Duration(time.Second),
		},
		mux:           http.NewServeMux(),
		metastore:     make(map[string]map[string]interface{}),
		methodConfigs: make(map[string]*MethodConfig),
		injections:    make([]injection, 0),
	}
	engine.RouterGroup.engine = engine
	// NOTE add prometheus monitor location
	engine.addRoute("GET", "/metrics", monitor())
	engine.addRoute("GET", "/metadata", engine.metadata())
	startPerf()
	return engine
}

// NewServer returns a new blank Engine instance without any middleware attached.
func NewServer(conf *ServerConfig) *Engine {
	if conf == nil {
		if !flag.Parsed() {
			fmt.Fprint(os.Stderr, "[blademaster] please call flag.Parse() before Init warden server, some configure may not effect.\n")
		}
		conf = parseDSN(_httpDSN)
	} else {
		fmt.Fprintf(os.Stderr, "[blademaster] config will be deprecated, argument will be ignored. please use -http flag or HTTP env to configure http server.\n")
	}

	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		address:       ip.InternalIP(),
		mux:           http.NewServeMux(),
		metastore:     make(map[string]map[string]interface{}),
		methodConfigs: make(map[string]*MethodConfig),
	}
	if err := engine.SetConfig(conf); err != nil {
		panic(err)
	}
	engine.RouterGroup.engine = engine
	// NOTE add prometheus monitor location
	engine.addRoute("GET", "/metrics", monitor())
	engine.addRoute("GET", "/metadata", engine.metadata())
	startPerf()
	return engine
}

// SetMethodConfig is used to set config on specified path
func (engine *Engine) SetMethodConfig(path string, mc *MethodConfig) {
	engine.pcLock.Lock()
	engine.methodConfigs[path] = mc
	engine.pcLock.Unlock()
}

// DefaultServer returns an Engine instance with the Recovery, Logger and CSRF middleware already attached.
func DefaultServer(conf *ServerConfig) *Engine {
	engine := NewServer(conf)
	engine.Use(Recovery(), Trace(), Logger(), CSRF(), Mobile())
	return engine
}

// Default returns an Engine instance with the Recovery, Logger and CSRF middleware already attached.
//
// Deprecated: please use DefaultServer.
func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Trace(), Logger(), CSRF(), Mobile())
	return engine
}

func (engine *Engine) addRoute(method, path string, handlers ...HandlerFunc) {
	if path[0] != '/' {
		panic("blademaster: path must begin with '/'")
	}
	if method == "" {
		panic("blademaster: HTTP method can not be empty")
	}
	if len(handlers) == 0 {
		panic("blademaster: there must be at least one handler")
	}
	if _, ok := engine.metastore[path]; !ok {
		engine.metastore[path] = make(map[string]interface{})
	}
	engine.metastore[path]["method"] = method
	engine.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		c := &Context{
			Context:  nil,
			engine:   engine,
			index:    -1,
			handlers: nil,
			Keys:     nil,
			method:   "",
			Error:    nil,
		}

		c.Request = req
		c.Writer = w
		c.handlers = handlers
		c.method = method

		engine.handleContext(c)
	})
}

// SetConfig is used to set the engine configuration.
// Only the valid config will be loaded.
func (engine *Engine) SetConfig(conf *ServerConfig) (err error) {
	if conf.Timeout <= 0 {
		return errors.New("blademaster: config timeout must greater than 0")
	}
	if conf.Network == "" {
		conf.Network = "tcp"
	}
	engine.lock.Lock()
	engine.conf = conf
	engine.lock.Unlock()
	return
}

func (engine *Engine) methodConfig(path string) *MethodConfig {
	engine.pcLock.RLock()
	mc := engine.methodConfigs[path]
	engine.pcLock.RUnlock()
	return mc
}

func (engine *Engine) handleContext(c *Context) {
	var cancel func()
	req := c.Request
	ctype := req.Header.Get("Content-Type")
	switch {
	case strings.Contains(ctype, "multipart/form-data"):
		req.ParseMultipartForm(defaultMaxMemory)
	default:
		req.ParseForm()
	}
	// get derived timeout from http request header,
	// compare with the engine configured,
	// and use the minimum one
	engine.lock.RLock()
	tm := time.Duration(engine.conf.Timeout)
	engine.lock.RUnlock()
	// the method config is preferred
	if pc := engine.methodConfig(c.Request.URL.Path); pc != nil {
		tm = time.Duration(pc.Timeout)
	}
	if ctm := timeout(req); ctm > 0 && tm > ctm {
		tm = ctm
	}
	md := metadata.MD{
		metadata.Color:      color(req),
		metadata.RemoteIP:   remoteIP(req),
		metadata.RemotePort: remotePort(req),
		metadata.Caller:     caller(req),
		metadata.Mirror:     mirror(req),
	}
	ctx := metadata.NewContext(context.Background(), md)
	if tm > 0 {
		c.Context, cancel = context.WithTimeout(ctx, tm)
	} else {
		c.Context, cancel = context.WithCancel(ctx)
	}
	defer cancel()
	c.Next()
}

// Router return a http.Handler for using http.ListenAndServe() directly.
func (engine *Engine) Router() http.Handler {
	return engine.mux
}

// Server is used to load stored http server.
func (engine *Engine) Server() *http.Server {
	s, ok := engine.server.Load().(*http.Server)
	if !ok {
		return nil
	}
	return s
}

// Shutdown the http server without interrupting active connections.
func (engine *Engine) Shutdown(ctx context.Context) error {
	server := engine.Server()
	if server == nil {
		return errors.New("blademaster: no server")
	}
	return errors.WithStack(server.Shutdown(ctx))
}

// UseFunc attachs a global middleware to the router. ie. the middleware attached though UseFunc() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) UseFunc(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.UseFunc(middleware...)
	return engine
}

// Use attachs a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...Handler) IRoutes {
	engine.RouterGroup.Use(middleware...)
	return engine
}

// Ping is used to set the general HTTP ping handler.
func (engine *Engine) Ping(handler HandlerFunc) {
	engine.GET("/monitor/ping", handler)
}

// Register is used to export metadata to discovery.
func (engine *Engine) Register(handler HandlerFunc) {
	engine.GET("/register", handler)
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) Run(addr ...string) (err error) {
	address := resolveAddress(addr)
	server := &http.Server{
		Addr:    address,
		Handler: engine.mux,
	}
	engine.server.Store(server)
	if err = server.ListenAndServe(); err != nil {
		err = errors.Wrapf(err, "addrs: %v", addr)
	}
	return
}

// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunTLS(addr, certFile, keyFile string) (err error) {
	server := &http.Server{
		Addr:    addr,
		Handler: engine.mux,
	}
	engine.server.Store(server)
	if err = server.ListenAndServeTLS(certFile, keyFile); err != nil {
		err = errors.Wrapf(err, "tls: %s/%s:%s", addr, certFile, keyFile)
	}
	return
}

// RunUnix attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunUnix(file string) (err error) {
	os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		err = errors.Wrapf(err, "unix: %s", file)
		return
	}
	defer listener.Close()
	server := &http.Server{
		Handler: engine.mux,
	}
	engine.server.Store(server)
	if err = server.Serve(listener); err != nil {
		err = errors.Wrapf(err, "unix: %s", file)
	}
	return
}

// RunServer will serve and start listening HTTP requests by given server and listener.
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunServer(server *http.Server, l net.Listener) (err error) {
	server.Handler = engine.mux
	engine.server.Store(server)
	if err = server.Serve(l); err != nil {
		err = errors.Wrapf(err, "listen server: %+v/%+v", server, l)
		return
	}
	return
}

func (engine *Engine) metadata() HandlerFunc {
	return func(c *Context) {
		c.JSON(engine.metastore, nil)
	}
}

// Inject is
func (engine *Engine) Inject(pattern string, handlers ...HandlerFunc) {
	engine.injections = append(engine.injections, injection{
		pattern:  regexp.MustCompile(pattern),
		handlers: handlers,
	})
}

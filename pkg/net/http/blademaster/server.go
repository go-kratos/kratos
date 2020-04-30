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

	"github.com/go-kratos/kratos/pkg/conf/dsn"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/net/criticality"
	"github.com/go-kratos/kratos/pkg/net/ip"
	"github.com/go-kratos/kratos/pkg/net/metadata"
	xtime "github.com/go-kratos/kratos/pkg/time"

	"github.com/pkg/errors"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

var (
	_ IRouter = &Engine{}

	_httpDSN       string
	default405Body = []byte("405 method not allowed")
	default404Body = []byte("404 page not found")
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
	Network      string         `dsn:"network"`
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
		return errors.Wrapf(err, "blademaster: listen tcp: %s", conf.Addr)
	}

	log.Info("blademaster: start http listen addr: %s", l.Addr().String())
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

	trees     methodTrees
	server    atomic.Value                      // store *http.Server
	metastore map[string]map[string]interface{} // metastore is the path as key and the metadata of this path as value, it export via /metadata

	pcLock        sync.RWMutex
	methodConfigs map[string]*MethodConfig

	injections []injection

	// If enabled, the url.RawPath will be used to find parameters.
	UseRawPath bool

	// If true, the path value will be unescaped.
	// If UseRawPath is false (by default), the UnescapePathValues effectively is true,
	// as url.Path gonna be used, which is already unescaped.
	UnescapePathValues bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	allNoRoute  []HandlerFunc
	allNoMethod []HandlerFunc
	noRoute     []HandlerFunc
	noMethod    []HandlerFunc

	pool sync.Pool
}

type injection struct {
	pattern  *regexp.Regexp
	handlers []HandlerFunc
}

// NewServer returns a new blank Engine instance without any middleware attached.
func NewServer(conf *ServerConfig) *Engine {
	if conf == nil {
		if !flag.Parsed() {
			fmt.Fprint(os.Stderr, "[blademaster] please call flag.Parse() before Init blademaster server, some configure may not effect.\n")
		}
		conf = parseDSN(_httpDSN)
	}
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		address:                ip.InternalIP(),
		trees:                  make(methodTrees, 0, 9),
		metastore:              make(map[string]map[string]interface{}),
		methodConfigs:          make(map[string]*MethodConfig),
		HandleMethodNotAllowed: true,
		injections:             make([]injection, 0),
	}
	if err := engine.SetConfig(conf); err != nil {
		panic(err)
	}
	engine.pool.New = func() interface{} {
		return engine.newContext()
	}
	engine.RouterGroup.engine = engine
	// NOTE add prometheus monitor location
	engine.addRoute("GET", "/metrics", monitor())
	engine.addRoute("GET", "/metadata", engine.metadata())
	engine.NoRoute(func(c *Context) {
		c.Bytes(404, "text/plain", default404Body)
		c.Abort()
	})
	engine.NoMethod(func(c *Context) {
		c.Bytes(405, "text/plain", []byte(http.StatusText(405)))
		c.Abort()
	})
	startPerf(engine)
	return engine
}

// SetMethodConfig is used to set config on specified path
func (engine *Engine) SetMethodConfig(path string, mc *MethodConfig) {
	engine.pcLock.Lock()
	engine.methodConfigs[path] = mc
	engine.pcLock.Unlock()
}

// DefaultServer returns an Engine instance with the Recovery and Logger middleware already attached.
func DefaultServer(conf *ServerConfig) *Engine {
	engine := NewServer(conf)
	engine.Use(Recovery(), Trace(), Logger())
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
	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}

	prelude := func(c *Context) {
		c.method = method
		c.RoutePath = path
	}
	handlers = append([]HandlerFunc{prelude}, handlers...)
	root.addRoute(path, handlers)
}

func (engine *Engine) prepareHandler(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.Request.URL.EscapedPath()) > 0 {
		rPath = c.Request.URL.EscapedPath()
		unescape = engine.UnescapePathValues
	}
	rPath = cleanPath(rPath)

	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		handlers, params, _ := root.getValue(rPath, c.Params, unescape)
		if handlers != nil {
			c.handlers = handlers
			c.Params = params
			return
		}
		break
	}

	if engine.HandleMethodNotAllowed {
		for _, tree := range engine.trees {
			if tree.method == httpMethod {
				continue
			}
			if handlers, _, _ := tree.root.getValue(rPath, nil, unescape); handlers != nil {
				c.handlers = engine.allNoMethod
				return
			}
		}
	}
	c.handlers = engine.allNoRoute
	return
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
		metadata.RemoteIP:    remoteIP(req),
		metadata.RemotePort:  remotePort(req),
		metadata.Criticality: string(criticality.Critical),
	}
	parseMetadataTo(req, md)
	ctx := metadata.NewContext(context.Background(), md)
	if tm > 0 {
		c.Context, cancel = context.WithTimeout(ctx, tm)
	} else {
		c.Context, cancel = context.WithCancel(ctx)
	}
	defer cancel()
	engine.prepareHandler(c)
	c.Next()
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

// Router return a http.Handler for using http.ListenAndServe() directly.
func (engine *Engine) Router() http.Handler {
	return engine
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
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}

// Use attachs a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...Handler) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}

// Ping is used to set the general HTTP ping handler.
func (engine *Engine) Ping(handler HandlerFunc) {
	engine.GET("/ping", handler)
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
		Handler: engine,
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
		Handler: engine,
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
		Handler: engine,
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
	server.Handler = engine
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

// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.Request = req
	c.Writer = w
	c.reset()

	engine.handleContext(c)
	engine.pool.Put(c)
}

//newContext for sync.pool
func (engine *Engine) newContext() *Context {
	return &Context{engine: engine}
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

// NoMethod sets the handlers called when... TODO.
func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}

func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}

package rpc

import (
	"bufio"
	ctx "context"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"go-common/library/conf/dsn"
	"go-common/library/conf/env"
	xlog "go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/net/rpc/context"
	"go-common/library/net/rpc/interceptor"
	"go-common/library/net/trace"

	pkgerr "github.com/pkg/errors"
)

var (
	_gorpcDSN string
)

func init() {
	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
	v := os.Getenv("GORPC")
	if v == "" {
		if env.GORPCPort != "" {
			v = "tcp://0.0.0.0:" + env.GORPCPort
		} else {
			v = "tcp://0.0.0.0:8099"
		}
	}
	fs.StringVar(&_gorpcDSN, "gorpc", v, "listen go rpc dsn, or use GORPC env variable.")
}

func parseDSN(rawdsn string) *ServerConfig {
	conf := new(ServerConfig)
	d, err := dsn.Parse(rawdsn)
	if err != nil {
		panic(pkgerr.WithMessage(err, "net/rpc: invalid dsn"))
	}
	if _, err = d.Bind(conf); err != nil {
		panic(pkgerr.WithMessage(err, "net/rpc: invalid dsn"))
	}
	return conf
}

// ServerConfig rpc server settings.
type ServerConfig struct {
	Proto string `dsn:"network"`
	Addr  string `dsn:"address"`
}

// NewServer new a rpc server.
func NewServer(c *ServerConfig) *Server {
	if c == nil {
		if !flag.Parsed() {
			fmt.Fprint(os.Stderr, "[net/rpc] please call flag.Parse() before Init go rpc server, some configure may not effect.\n")
		}
		c = parseDSN(_gorpcDSN)
	} else {
		fmt.Fprintf(os.Stderr, "[net/rpc] config will be deprecated, argument will be ignored. please use -gorpc flag or GORPC env to configure go rpc server.\n")
	}

	s := newServer()
	s.Interceptor = interceptor.NewInterceptor("")
	go rpcListen(c, s)
	return s
}

// rpcListen start rpc listen.
func rpcListen(c *ServerConfig, s *Server) {
	l, err := net.Listen(c.Proto, c.Addr)
	if err != nil {
		xlog.Error("net.Listen(rpcAddr:(%v)) error(%v)", c.Addr, err)
		panic(err)
	}
	// if process exit, then close the rpc bind
	defer func() {
		xlog.Info("rpc addr:(%s) close", c.Addr)
		if err := l.Close(); err != nil {
			xlog.Error("listener.Close() error(%v)", err)
		}
	}()
	xlog.Info("start rpc listen addr: %s", c.Addr)
	s.Accept(l)
}

var (
	// Precompute the reflect type for error. Can't use error directly
	// because Typeof takes an empty interface value. This is annoying.
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
	ctxType     = reflect.TypeOf((*context.Context)(nil)).Elem()
)

// methodType 方法类型
type methodType struct {
	method    reflect.Method //方法
	ArgType   reflect.Type   //参数类型
	ReplyType reflect.Type   //回复类型
}

// service 服务
type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods 注册的方法
}

// Request is a header written before every RPC call. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Request struct {
	Color         string        // color
	RemoteIP      string        // remoteIP
	Timeout       time.Duration // timeout
	ServiceMethod string        // format: "Service.Method"
	Seq           uint64        // sequence number chosen by client
	Trace         TraceInfo     // trace info

	ctx context.Context
}

// Auth handshake struct.
type Auth struct {
	User string
}

// Response is a header written before every RPC return. It is used internally
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Response struct {
	ServiceMethod string // echoes that of the Request
	Seq           uint64 // echoes that of the request
	Error         string // error, if any.
}

// Interceptor interface.
type Interceptor interface {
	Rate(context.Context) error
	Stat(context.Context, interface{}, error)
	Auth(context.Context, net.Addr, string) error // ip, token
}

// Server represents an RPC Server.
type Server struct {
	lis         net.Listener
	serviceMap  map[string]*service
	Interceptor Interceptor //拦截器
}

// newServer returns a new Server.
func newServer() *Server {
	return &Server{serviceMap: make(map[string]*service)}
}

// DefaultServer is the default instance of *Server.
var DefaultServer = newServer()

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//	- exported method of exported type
//	- two arguments, both of exported type
//	- the second argument is a pointer
//	- one return value, of type error
// It returns an error if the receiver is not an exported type or has
// no suitable methods. It also logs the error using package log.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.
func (server *Server) Register(rcvr interface{}) (err error) {
	if err = server.register(rcvr, "", false); err != nil {
		return
	}
	return server.register(new(pinger), _service, true)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (server *Server) RegisterName(name string, rcvr interface{}) (err error) {
	if err = server.register(rcvr, name, true); err != nil {
		return
	}
	return server.register(new(pinger), _service, true)
}

func (server *Server) register(rcvr interface{}, name string, useName bool) error {
	if server.serviceMap == nil {
		server.serviceMap = make(map[string]*service)
	}
	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if useName {
		sname = name
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + s.typ.String()
		log.Print(s)
		return errors.New(s)
	}
	if !isExported(sname) && !useName {
		s := "rpc.Register: type " + sname + " is not exported"
		log.Print(s)
		return errors.New(s)
	}
	if _, present := server.serviceMap[sname]; present {
		return errors.New("rpc: service already defined: " + sname)
	}
	s.name = sname
	// Install the methods
	s.method = suitableMethods(s.typ, true)

	if len(s.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PtrTo(s.typ), false)
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		log.Print(str)
		return errors.New(str)
	}
	server.serviceMap[s.name] = s
	return nil
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using log if reportErr is true.
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs ins: receiver, context, *arg, *reply.
		if mtype.NumIn() != 4 {
			if reportErr {
				log.Println("method", mname, "has wrong number of ins:", mtype.NumIn())
			}
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !argType.Implements(ctxType) {
			if reportErr {
				log.Println(mname, "argument type must implements:", ctxType)
			}
			continue
		}
		// Second arg need not be a pointer.
		argType = mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			if reportErr {
				log.Println(mname, "argument type not exported:", argType)
			}
			continue
		}
		// Thrid arg must be a pointer.
		replyType := mtype.In(3)
		if replyType.Kind() != reflect.Ptr {
			if reportErr {
				log.Println("method", mname, "reply type not a pointer:", replyType)
			}
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if reportErr {
				log.Println("method", mname, "reply type not exported:", replyType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			if reportErr {
				log.Println("method", mname, "has wrong number of outs:", mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			if reportErr {
				log.Println("method", mname, "returns", returnType.String(), "not error")
			}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}

// A value sent as a placeholder for the server's response value when the server
// receives an invalid request. It is never decoded by the client since the Response
// contains an error when it is used.
var invalidRequest = struct{}{}

func (server *Server) sendResponse(c context.Context, codec *serverCodec, reply interface{}, errmsg string) {
	var (
		err  error
		ts   Response
		resp = &codec.resp
	)
	if errmsg != "" {
		reply = invalidRequest
	}
	ts.ServiceMethod = c.ServiceMethod()
	ts.Seq = c.Seq()
	ts.Error = errmsg
	codec.sending.Lock()
	// NOTE must keep resp goroutine safe
	*resp = ts
	// Encode the response header
	if err = codec.writeResponse(reply); err != nil {
		log.Println("rpc: writing response:", err)
	}
	codec.sending.Unlock()
}

func (s *service) call(c context.Context, server *Server, mtype *methodType, argv, replyv reflect.Value, codec *serverCodec) {
	var (
		err          error
		errmsg       string
		errInter     interface{}
		cv           reflect.Value
		returnValues []reflect.Value
	)
	t, _ := trace.FromContext(c)
	defer func() {
		if err1 := recover(); err1 != nil {
			err = err1.(error)
			errmsg = err.Error()
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			xlog.Error("rpc call panic: %v \n%s", err1, buf)
			server.sendResponse(c, codec, replyv.Interface(), errmsg)
			if server.Interceptor != nil {
				server.Interceptor.Stat(c, argv.Interface(), err)
			}
			if t != nil {
				t.Finish(&err)
			}
		}
	}()
	// rate limit
	if server.Interceptor != nil {
		if err = server.Interceptor.Rate(c); err != nil {
			errmsg = err.Error()
		}
	}
	if err == nil {
		// Invoke the method, providing a new value for the reply.
		cv = reflect.New(ctxType)
		*cv.Interface().(*context.Context) = c
		returnValues = mtype.method.Func.Call([]reflect.Value{s.rcvr, cv.Elem(), argv, replyv})
		// The return value for the method is an error.
		if errInter = returnValues[0].Interface(); errInter != nil {
			err = errInter.(error)
			errmsg = pkgerr.Cause(err).Error()
		}
	}
	server.sendResponse(c, codec, replyv.Interface(), errmsg)
	// stat
	if server.Interceptor != nil {
		server.Interceptor.Stat(c, argv.Interface(), err)
	}
	if t != nil {
		t.Finish(&err)
	}
}

type serverCodec struct {
	sending sync.Mutex
	resp    Response
	req     Request
	auth    Auth

	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
	addr   net.Addr
	closed bool
}

func (c *serverCodec) readRequestHeader() error {
	return pkgerr.WithStack(c.dec.Decode(&c.req))
}

func (c *serverCodec) readRequestBody(body interface{}) error {
	return pkgerr.WithStack(c.dec.Decode(body))
}

func (c *serverCodec) writeResponse(body interface{}) (err error) {
	if err = c.enc.Encode(&c.resp); err != nil {
		err = pkgerr.WithStack(err)
		if c.encBuf.Flush() == nil {
			// Gob couldn't encode the header. Should not happen, so if it does,
			// shut down the connection to signal that the connection is broken.
			log.Println("rpc: gob error encoding response:", err)
			c.close()
		}
		return
	}
	if err = c.enc.Encode(body); err != nil {
		err = pkgerr.WithStack(err)
		if c.encBuf.Flush() == nil {
			// Was a gob problem encoding the body but the header has been written.
			// Shut down the connection to signal that the connection is broken.
			log.Println("rpc: gob error encoding body:", err)
			c.close()
		}
		return
	}
	return pkgerr.WithStack(c.encBuf.Flush())
}

func (c *serverCodec) close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

// ServeConn runs the server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
// ServeConn uses the gob wire format (see package gob) on the
// connection. To use an alternate codec, use ServeCodec.
func (server *Server) ServeConn(conn net.Conn) {
	buf := bufio.NewWriter(conn)
	srv := &serverCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(buf),
		encBuf: buf,
		addr:   conn.RemoteAddr(),
	}
	server.serveCodec(srv)
}

// serveCodec is like ServeConn but uses the specified codec to
// decode requests and encode responses.
func (server *Server) serveCodec(codec *serverCodec) {
	req := &codec.req
	for {
		// serve request
		service, mtype, argv, replyv, err := server.readRequest(codec)
		if err != nil {
			if err != io.EOF {
				log.Println("rpc:", err)
			}
			if req.ctx == nil {
				break
			}
			errmsg := err.Error()
			if req.ServiceMethod == _authServiceMethod {
				errmsg = ""
			}
			server.sendResponse(req.ctx, codec, invalidRequest, errmsg)
			continue
		}
		if req.ServiceMethod == _authServiceMethod {
			codec.auth = *(argv.Interface().(*Auth))
			req.ctx = context.NewContext(ctx.Background(), req.ServiceMethod, codec.auth.User, req.Seq)
			server.sendResponse(req.ctx, codec, invalidRequest, "")
			continue
		}
		go service.call(req.ctx, server, mtype, argv, replyv, codec)
	}
	codec.close()
}

func (server *Server) readRequest(codec *serverCodec) (service *service, mtype *methodType, argv, replyv reflect.Value, err error) {
	var req = &codec.req
	*req = Request{}
	if service, mtype, err = server.readRequestHeader(codec); err != nil {
		// keepreading
		if req.ctx == nil {
			return
		}
		// discard body
		codec.readRequestBody(nil)
		return
	}

	// Decode the argument value.
	argIsValue := false // if true, need to indirect before calling.
	if mtype.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(mtype.ArgType.Elem())
	} else {
		argv = reflect.New(mtype.ArgType)
		argIsValue = true
	}
	// argv guaranteed to be a pointer now.
	if err = codec.readRequestBody(argv.Interface()); err != nil {
		return
	}
	if argIsValue {
		argv = argv.Elem()
	}

	replyv = reflect.New(mtype.ReplyType.Elem())
	return
}

func (server *Server) readRequestHeader(codec *serverCodec) (service *service, mtype *methodType, err error) {
	var t trace.Trace
	req := &codec.req
	if err = codec.readRequestHeader(); err != nil {
		return
	}
	if t, _ = trace.Extract(nil, &req.Trace); t == nil {
		t = trace.New(req.ServiceMethod)
	}
	t.SetTitle(req.ServiceMethod)
	t.SetTag(trace.String(trace.TagAddress, codec.addr.String()))
	md := metadata.MD{
		metadata.Trace:    t,
		metadata.Color:    req.Color,
		metadata.RemoteIP: req.RemoteIP,
		metadata.Caller:   req.Trace.Caller,
	}
	// FIXME(maojian) Timeout?
	c1 := metadata.NewContext(ctx.Background(), md)

	caller := codec.auth.User
	if caller == "" {
		caller = req.Trace.Caller
	}

	// NOTE ctx not nil then keepreading
	req.ctx = context.NewContext(c1, req.ServiceMethod, caller, req.Seq)
	// We read the header successfully. If we see an error now,
	// we can still recover and move on to the next request.
	dot := strings.LastIndex(req.ServiceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc: service/method request ill-formed: " + req.ServiceMethod)
		return
	}
	serviceName := req.ServiceMethod[:dot]
	methodName := req.ServiceMethod[dot+1:]

	// Look up the request.
	service = server.serviceMap[serviceName]
	if service == nil {
		err = errors.New("rpc: can't find service " + req.ServiceMethod)
		return
	}
	mtype = service.method[methodName]
	if mtype == nil {
		err = errors.New("rpc: can't find method " + req.ServiceMethod)
	}
	return
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection. Accept blocks until the listener
// returns a non-nil error. The caller typically invokes Accept in a
// go statement.
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Print("rpc.Serve: accept:", err.Error())
			return
		}
		go server.ServeConn(conn)
	}
}

// Close stop the rpc server.
func (server *Server) Close() error {
	if server.lis != nil {
		return server.lis.Close()
	}
	return nil
}

// Register publishes the receiver's methods in the DefaultServer.
func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func RegisterName(name string, rcvr interface{}) error {
	return DefaultServer.RegisterName(name, rcvr)
}

// ServeConn runs the DefaultServer on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
// ServeConn uses the gob wire format (see package gob) on the
// connection. To use an alternate codec, use ServeCodec.
func ServeConn(conn net.Conn) {
	DefaultServer.ServeConn(conn)
}

// Accept accepts connections on the listener and serves requests
// to DefaultServer for each incoming connection.
// Accept blocks; the caller typically invokes it in a go statement.
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

const (
	_authServiceMethod = "inner.Auth"
	_pingServiceMethod = "inner.Ping"
	_service           = "inner"
)

var (
	_pingArg = &struct{}{}
)

// pinger rpc ping service
type pinger struct {
}

// Ping rpc ping.
func (p *pinger) Ping(c context.Context, arg *struct{}, reply *struct{}) error {
	return nil
}

// Auth
func (p *pinger) Auth(c context.Context, arg *Auth, reply *struct{}) error {
	return nil
}

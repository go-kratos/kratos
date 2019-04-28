package rpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-common/library/net/rpc/context"
)

var (
	newSvr, newInterceptorServer                        *Server
	serverAddr, newServerAddr, newServerInterceptorAddr string
	once, newOnce, newInterceptorOnce                   sync.Once
	testInterceptor                                     *TestInterceptor
	statCount                                           uint64
)

const (
	_testToken = "test_token"
)

type TestInterceptor struct {
	Token, RateMethod string
}

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

// Some of Arith's methods have value args, some have pointer args. That's deliberate.

func (t *TestInterceptor) Rate(c context.Context) error {
	log.Printf("Interceptor rate method: %s, current: %s", t.RateMethod, c.ServiceMethod())
	if t.RateMethod == c.ServiceMethod() {
		return fmt.Errorf("Interceptor rate method: %s, time: %s", c.ServiceMethod(), c.Now())
	}
	return nil
}

func (t *TestInterceptor) Auth(c context.Context, addr net.Addr, token string) error {
	if t.Token != token {
		return fmt.Errorf("Interceptor auth token: %s, ip: %s seq: %d failed", token, addr, c.Seq())
	}
	log.Printf("Interceptor auth token: %s, ip: %s, seq: %d ok", token, addr, c.Seq())
	return nil
}

func (t *TestInterceptor) Stat(c context.Context, args interface{}, err error) {
	atomic.AddUint64(&statCount, 1)
}

func (t *Arith) Auth(c context.Context, args Auth, reply *Reply) error {
	return nil
}

func (t *Arith) Add(c context.Context, args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Mul(c context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Div(c context.Context, args Args, reply *Reply) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	reply.C = args.A / args.B
	return nil
}

func (t *Arith) String(c context.Context, args *Args, reply *string) error {
	*reply = fmt.Sprintf("%d+%d=%d", args.A, args.B, args.A+args.B)
	return nil
}

func (t *Arith) Scan(c context.Context, args string, reply *Reply) (err error) {
	_, err = fmt.Sscan(args, &reply.C)
	return
}

func (t *Arith) Error(c context.Context, args *Args, reply *Reply) error {
	panic("ERROR")
}

type hidden int

func (t *hidden) Exported(c context.Context, args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

type Embed struct {
	hidden
}

// NOTE listen and start the server

func listenTCP() (net.Listener, string) {
	l, e := net.Listen("tcp", "127.0.0.1:0") // any available address
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}
	return l, l.Addr().String()
}

func startServer() {
	Register(new(Arith))
	Register(new(Embed))
	RegisterName("net.rpc.Arith", new(Arith))

	var l net.Listener
	l, serverAddr = listenTCP()
	log.Println("Test RPC server listening on", serverAddr)
	go Accept(l)
}

func startNewServer() {
	newSvr = newServer()
	newSvr.Register(new(Arith))
	newSvr.Register(new(Embed))
	newSvr.RegisterName("net.rpc.Arith", new(Arith))
	newSvr.RegisterName("newServer.Arith", new(Arith))

	var l net.Listener
	l, newServerAddr = listenTCP()
	log.Println("NewServer test RPC server listening on", newServerAddr)
	go newSvr.Accept(l)
}

func startNewInterceptorServer() {
	testInterceptor = &TestInterceptor{Token: _testToken}
	newInterceptorServer = newServer()
	newInterceptorServer.Register(new(Arith))
	newInterceptorServer.Register(new(Embed))
	newInterceptorServer.RegisterName("net.rpc.Arith", new(Arith))
	newInterceptorServer.RegisterName("newServer.Arith", new(Arith))
	newInterceptorServer.Interceptor = testInterceptor

	var l net.Listener
	l, newServerInterceptorAddr = listenTCP()
	log.Println("NewInterceptorServer test RPC server listening on", newServerAddr)
	go newInterceptorServer.Accept(l)
}

// NOTE test rpc call with check expected

func TestServerRPC(t *testing.T) {
	once.Do(startServer)
	newOnce.Do(startNewServer)
	newInterceptorOnce.Do(startNewInterceptorServer)
	testRPC(t, serverAddr)
	testRPC(t, newServerAddr)
	testNewServerRPC(t, newServerAddr)
	testNewServerAuthRPC(t, newServerInterceptorAddr)
	testNewInterceptorServerRPC(t, newServerInterceptorAddr)
	testNewInterceptorServerRateRPC(t, newServerInterceptorAddr)
}

func testRPC(t *testing.T, addr string) {
	client, err := dial("tcp", addr, time.Second)
	if err != nil {
		t.Fatal("dialing", err)
	}
	defer client.Close()

	// Synchronous calls
	args := &Args{7, 8}
	reply := new(Reply)
	err = client.Call("Arith.Add", args, reply)
	if err != nil {
		t.Errorf("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		t.Errorf("Add: expected %d got %d", reply.C, args.A+args.B)
	}

	// Methods exported from unexported embedded structs
	args = &Args{7, 0}
	reply = new(Reply)
	err = client.Call("Embed.Exported", args, reply)
	if err != nil {
		t.Errorf("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		t.Errorf("Add: expected %d got %d", reply.C, args.A+args.B)
	}

	// Nonexistent method
	args = &Args{7, 0}
	reply = new(Reply)
	err = client.Call("Arith.BadOperation", args, reply)
	// expect an error
	if err == nil {
		t.Error("BadOperation: expected error")
	} else if !strings.HasPrefix(err.Error(), "rpc: can't find method ") {
		t.Errorf("BadOperation: expected can't find method error; got %q", err)
	}

	// Unknown service
	args = &Args{7, 8}
	reply = new(Reply)
	err = client.Call("Arith.Unknown", args, reply)
	if err == nil {
		t.Error("expected error calling unknown service")
	} else if !strings.Contains(err.Error(), "method") {
		t.Error("expected error about method; got", err)
	}

	// Out of order.
	args = &Args{7, 8}
	mulReply := new(Reply)
	mulCall := client.Go("Arith.Mul", args, mulReply, nil)
	addReply := new(Reply)
	addCall := client.Go("Arith.Add", args, addReply, nil)

	addCall = <-addCall.Done
	if addCall.Error != nil {
		t.Errorf("Add: expected no error but got string %q", addCall.Error.Error())
	}
	if addReply.C != args.A+args.B {
		t.Errorf("Add: expected %d got %d", addReply.C, args.A+args.B)
	}

	mulCall = <-mulCall.Done
	if mulCall.Error != nil {
		t.Errorf("Mul: expected no error but got string %q", mulCall.Error.Error())
	}
	if mulReply.C != args.A*args.B {
		t.Errorf("Mul: expected %d got %d", mulReply.C, args.A*args.B)
	}

	// Error test
	args = &Args{7, 0}
	reply = new(Reply)
	err = client.Call("Arith.Div", args, reply)
	// expect an error: zero divide
	if err == nil {
		t.Error("Div: expected error")
	} else if err.Error() != "divide by zero" {
		t.Error("Div: expected divide by zero error; got", err)
	}

	// Bad type.
	reply = new(Reply)
	err = client.Call("Arith.Add", reply, reply) // args, reply would be the correct thing to use
	if err == nil {
		t.Error("expected error calling Arith.Add with wrong arg type")
	} else if !strings.Contains(err.Error(), "type") {
		t.Error("expected error about type; got", err)
	}

	// Non-struct argument
	const Val = 12345
	str := fmt.Sprint(Val)
	reply = new(Reply)
	err = client.Call("Arith.Scan", &str, reply)
	if err != nil {
		t.Errorf("Scan: expected no error but got string %q", err.Error())
	} else if reply.C != Val {
		t.Errorf("Scan: expected %d got %d", Val, reply.C)
	}

	// Non-struct reply
	args = &Args{27, 35}
	str = ""
	err = client.Call("Arith.String", args, &str)
	if err != nil {
		t.Errorf("String: expected no error but got string %q", err.Error())
	}
	expect := fmt.Sprintf("%d+%d=%d", args.A, args.B, args.A+args.B)
	if str != expect {
		t.Errorf("String: expected %s got %s", expect, str)
	}

	args = &Args{7, 8}
	reply = new(Reply)
	err = client.Call("Arith.Mul", args, reply)
	if err != nil {
		t.Errorf("Mul: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A*args.B {
		t.Errorf("Mul: expected %d got %d", reply.C, args.A*args.B)
	}

	// ServiceName contain "." character
	args = &Args{7, 8}
	reply = new(Reply)
	err = client.Call("net.rpc.Arith.Add", args, reply)
	if err != nil {
		t.Errorf("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		t.Errorf("Add: expected %d got %d", reply.C, args.A+args.B)
	}
}

func testNewServerRPC(t *testing.T, addr string) {
	client, err := dial("tcp", addr, time.Second)
	if err != nil {
		t.Fatal("dialing", err)
	}
	defer client.Close()

	// Synchronous calls
	args := &Args{7, 8}
	reply := new(Reply)
	err = client.Call("newServer.Arith.Add", args, reply)
	if err != nil {
		t.Errorf("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		t.Errorf("Add: expected %d got %d", reply.C, args.A+args.B)
	}
}

func testNewInterceptorServerRPC(t *testing.T, addr string) {
	client, err := dial("tcp", addr, time.Second)
	if err != nil {
		t.Fatal("authing", err)
	}
	defer client.Close()

	// Synchronous calls
	args := &Args{7, 8}
	reply := new(Reply)
	err = client.Call("newServer.Arith.Add", args, reply)
	if err != nil {
		t.Errorf("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		t.Errorf("Add: expected %d got %d", reply.C, args.A+args.B)
	}
}

func testNewInterceptorServerRateRPC(t *testing.T, addr string) {
	client, err := dial("tcp", addr, time.Second)
	if err != nil {
		t.Fatal("authing", err)
	}
	defer client.Close()

	// Synchronous calls
	args := &Args{7, 8}
	reply := new(Reply)
	err = client.Call("Arith.Add", args, reply)
	if err != nil {
		t.Errorf("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		t.Errorf("Add: expected %d got %d", reply.C, args.A+args.B)
	}

	// check rate the method
	testInterceptor.RateMethod = "Arith.Add"
	args = &Args{7, 8}
	reply = new(Reply)
	err = client.Call("Arith.Add", args, reply)
	if err == nil {
		t.Errorf("Add: expected error this rate method")
	}
}

func testNewServerAuthRPC(t *testing.T, addr string) {
	_, err := dial("tcp", addr, time.Second)
	if err != nil {
		t.Errorf("Auth: expected error %q", err.Error())
	}
}

// NOTE test no point structs for registration

type ReplyNotPointer int
type ArgNotPublic int
type ReplyNotPublic int
type NeedsPtrType int
type local struct{}

func (t *ReplyNotPointer) ReplyNotPointer(c context.Context, args *Args, reply Reply) error {
	return nil
}

func (t *ArgNotPublic) ArgNotPublic(c context.Context, args *local, reply *Reply) error {
	return nil
}

func (t *ReplyNotPublic) ReplyNotPublic(c context.Context, args *Args, reply *local) error {
	return nil
}

func (t *NeedsPtrType) NeedsPtrType(c context.Context, args *Args, reply *Reply) error {
	return nil
}

// Check that registration handles lots of bad methods and a type with no suitable methods.
func TestRegistrationError(t *testing.T) {
	err := Register(new(ReplyNotPointer))
	if err == nil {
		t.Error("expected error registering ReplyNotPointer")
	}
	err = Register(new(ArgNotPublic))
	if err == nil {
		t.Error("expected error registering ArgNotPublic")
	}
	err = Register(new(ReplyNotPublic))
	if err == nil {
		t.Error("expected error registering ReplyNotPublic")
	}
	err = Register(NeedsPtrType(0))
	if err == nil {
		t.Error("expected error registering NeedsPtrType")
	} else if !strings.Contains(err.Error(), "pointer") {
		t.Error("expected hint when registering NeedsPtrType")
	}
}

// NOTE test multiple call methods

func dialDirect() (*client, error) {
	return dial("tcp", serverAddr, time.Second)
}

func countMallocs(dial func() (*client, error), t *testing.T) float64 {
	once.Do(startServer)
	client, err := dial()
	if err != nil {
		t.Fatal("error dialing", err)
	}
	defer client.Close()

	args := &Args{7, 8}
	reply := new(Reply)
	return testing.AllocsPerRun(100, func() {
		err := client.Call("Arith.Add", args, reply)
		if err != nil {
			t.Errorf("Add: expected no error but got string %q", err.Error())
		}
		if reply.C != args.A+args.B {
			t.Errorf("Add: expected %d got %d", reply.C, args.A+args.B)
		}
	})
}

func TestCountMallocs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping malloc count in short mode")
	}
	if runtime.GOMAXPROCS(0) > 1 {
		t.Skip("skipping; GOMAXPROCS>1")
	}
	fmt.Printf("mallocs per rpc round trip: %v\n", countMallocs(dialDirect, t))
}

func TestTCPClose(t *testing.T) {
	once.Do(startServer)

	client, err := dialDirect()
	if err != nil {
		t.Fatalf("dialing: %v", err)
	}
	defer client.Close()

	args := Args{17, 8}
	var reply Reply
	err = client.Call("Arith.Mul", args, &reply)
	if err != nil {
		t.Fatal("arith error:", err)
	}
	t.Logf("Arith: %d*%d=%d\n", args.A, args.B, reply)
	if reply.C != args.A*args.B {
		t.Errorf("Add: expected %d got %d", reply.C, args.A*args.B)
	}
}

func TestErrorAfterClientClose(t *testing.T) {
	once.Do(startServer)

	client, err := dialDirect()
	if err != nil {
		t.Fatalf("dialing: %v", err)
	}
	err = client.Close()
	if err != nil {
		t.Fatal("close error:", err)
	}
	err = client.Call("Arith.Add", &Args{7, 9}, new(Reply))
	if err != ErrShutdown {
		t.Errorf("Forever: expected ErrShutdown got %v", err)
	}
}

// Tests the fix to issue 11221. Without the fix, this loops forever or crashes.
func TestAcceptExitAfterListenerClose(t *testing.T) {
	newSvr = newServer()
	newSvr.Register(new(Arith))
	newSvr.RegisterName("net.rpc.Arith", new(Arith))
	newSvr.RegisterName("newServer.Arith", new(Arith))

	var l net.Listener
	l, newServerAddr = listenTCP()
	l.Close()
	newSvr.Accept(l)
}

func TestParseDSN(t *testing.T) {
	c := parseDSN("tcp://127.0.0.1:8099")
	if c.Proto != "tcp" {
		t.Error("parse dsn proto not equal tcp")
	}
	if c.Addr != "127.0.0.1:8099" {
		t.Error("parse dsn addr not equal")
	}
}

func benchmarkEndToEnd(dial func() (*client, error), b *testing.B) {
	once.Do(startServer)
	client, err := dial()
	if err != nil {
		b.Fatal("error dialing:", err)
	}
	defer client.Close()

	// Synchronous calls
	args := &Args{7, 8}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		reply := new(Reply)
		for pb.Next() {
			err := client.Call("Arith.Add", args, reply)
			if err != nil {
				b.Fatalf("rpc error: Add: expected no error but got string %q", err.Error())
			}
			if reply.C != args.A+args.B {
				b.Fatalf("rpc error: Add: expected %d got %d", reply.C, args.A+args.B)
			}
		}
	})
}

func benchmarkEndToEndAsync(dial func() (*client, error), b *testing.B) {
	if b.N == 0 {
		return
	}
	const MaxConcurrentCalls = 100
	once.Do(startServer)
	client, err := dial()
	if err != nil {
		b.Fatal("error dialing:", err)
	}
	defer client.Close()

	// Asynchronous calls
	args := &Args{7, 8}
	procs := 4 * runtime.GOMAXPROCS(-1)
	send := int32(b.N)
	recv := int32(b.N)
	var wg sync.WaitGroup
	wg.Add(procs)
	gate := make(chan bool, MaxConcurrentCalls)
	res := make(chan *Call, MaxConcurrentCalls)
	b.ResetTimer()

	for p := 0; p < procs; p++ {
		go func() {
			for atomic.AddInt32(&send, -1) >= 0 {
				gate <- true
				reply := new(Reply)
				client.Go("Arith.Add", args, reply, res)
			}
		}()
		go func() {
			for call := range res {
				A := call.Args.(*Args).A
				B := call.Args.(*Args).B
				C := call.Reply.(*Reply).C
				if A+B != C {
					return
				}
				<-gate
				if atomic.AddInt32(&recv, -1) == 0 {
					close(res)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkEndToEnd(b *testing.B) {
	benchmarkEndToEnd(dialDirect, b)
}

func BenchmarkEndToEndAsync(b *testing.B) {
	benchmarkEndToEndAsync(dialDirect, b)
}

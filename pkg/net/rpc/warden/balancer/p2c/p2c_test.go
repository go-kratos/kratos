package p2c

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bilibili/kratos/pkg/conf/env"

	nmd "github.com/bilibili/kratos/pkg/net/metadata"
	wmeta "github.com/bilibili/kratos/pkg/net/rpc/warden/internal/metadata"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

var serverNum int
var cliNum int
var concurrency int
var extraLoad int64
var extraDelay int64
var extraWeight uint64

func init() {
	flag.IntVar(&serverNum, "snum", 6, "-snum 6")
	flag.IntVar(&cliNum, "cnum", 12, "-cnum 12")
	flag.IntVar(&concurrency, "concurrency", 10, "-cc 10")
	flag.Int64Var(&extraLoad, "exload", 3, "-exload 3")
	flag.Int64Var(&extraDelay, "exdelay", 250, "-exdelay 250")
	flag.Uint64Var(&extraWeight, "extraWeight", 50, "-exdelay 50")
}

type testSubConn struct {
	addr resolver.Address
	wait chan struct{}
	//statics
	reqs      int64
	usage     int64
	cpu       int64
	prevReq   int64
	prevUsage int64
	//control params
	loadJitter  int64
	delayJitter int64
}

func newTestSubConn(addr string, weight uint64, color string) (sc *testSubConn) {
	sc = &testSubConn{
		addr: resolver.Address{
			Addr: addr,
			Metadata: wmeta.MD{
				Weight: weight,
				Color:  color,
			},
		},
		wait: make(chan struct{}, 1000),
	}
	go func() {
		for {
			for i := 0; i < 210; i++ {
				<-sc.wait
			}
			time.Sleep(time.Millisecond * 20)
		}
	}()

	return
}

func (s *testSubConn) connect(ctx context.Context) {
	time.Sleep(time.Millisecond * 15)
	//add qps counter when request come in
	atomic.AddInt64(&s.reqs, 1)
	select {
	case <-ctx.Done():
		return
	case s.wait <- struct{}{}:
		atomic.AddInt64(&s.usage, 1)
	}
	load := atomic.LoadInt64(&s.loadJitter)
	if load > 0 {
		for i := 0; i <= rand.Intn(int(load)); i++ {
			select {
			case <-ctx.Done():
				return
			case s.wait <- struct{}{}:
				atomic.AddInt64(&s.usage, 1)
			}
		}
	}
	delay := atomic.LoadInt64(&s.delayJitter)
	if delay > 0 {
		delay = rand.Int63n(delay)
		time.Sleep(time.Millisecond * time.Duration(delay))
	}
}

func (s *testSubConn) UpdateAddresses([]resolver.Address) {

}

// Connect starts the connecting for this SubConn.
func (s *testSubConn) Connect() {

}

func TestBalancerPick(t *testing.T) {
	scs := map[resolver.Address]balancer.SubConn{}
	sc1 := &testSubConn{
		addr: resolver.Address{
			Addr: "test1",
			Metadata: wmeta.MD{
				Weight: 8,
			},
		},
	}
	sc2 := &testSubConn{
		addr: resolver.Address{
			Addr: "test2",
			Metadata: wmeta.MD{
				Weight: 4,
				Color:  "red",
			},
		},
	}
	sc3 := &testSubConn{
		addr: resolver.Address{
			Addr: "test3",
			Metadata: wmeta.MD{
				Weight: 2,
				Color:  "red",
			},
		},
	}
	sc4 := &testSubConn{
		addr: resolver.Address{
			Addr: "test4",
			Metadata: wmeta.MD{
				Weight: 2,
				Color:  "purple",
			},
		},
	}
	scs[sc1.addr] = sc1
	scs[sc2.addr] = sc2
	scs[sc3.addr] = sc3
	scs[sc4.addr] = sc4
	b := &p2cPickerBuilder{}
	picker := b.Build(scs)
	res := []string{"test1", "test1", "test1", "test1"}
	for i := 0; i < 3; i++ {
		conn, _, err := picker.Pick(context.Background(), balancer.PickOptions{})
		if err != nil {
			t.Fatalf("picker.Pick failed!idx:=%d", i)
		}
		sc := conn.(*testSubConn)
		if sc.addr.Addr != res[i] {
			t.Fatalf("the subconn picked(%s),but expected(%s)", sc.addr.Addr, res[i])
		}
	}

	ctx := nmd.NewContext(context.Background(), nmd.New(map[string]interface{}{"color": "black"}))
	for i := 0; i < 4; i++ {
		conn, _, err := picker.Pick(ctx, balancer.PickOptions{})
		if err != nil {
			t.Fatalf("picker.Pick failed!idx:=%d", i)
		}
		sc := conn.(*testSubConn)
		if sc.addr.Addr != res[i] {
			t.Fatalf("the (%d) subconn picked(%s),but expected(%s)", i, sc.addr.Addr, res[i])
		}
	}

	env.Color = "purple"
	ctx2 := context.Background()
	for i := 0; i < 4; i++ {
		conn, _, err := picker.Pick(ctx2, balancer.PickOptions{})
		if err != nil {
			t.Fatalf("picker.Pick failed!idx:=%d", i)
		}
		sc := conn.(*testSubConn)
		if sc.addr.Addr != "test4" {
			t.Fatalf("the (%d) subconn picked(%s),but expected(%s)", i, sc.addr.Addr, res[i])
		}
	}

}

func Benchmark_Wrr(b *testing.B) {
	scs := map[resolver.Address]balancer.SubConn{}
	for i := 0; i < 50; i++ {
		addr := resolver.Address{
			Addr:     fmt.Sprintf("addr_%d", i),
			Metadata: wmeta.MD{Weight: 10},
		}
		scs[addr] = &testSubConn{addr: addr}
	}
	wpb := &p2cPickerBuilder{}
	picker := wpb.Build(scs)
	opt := balancer.PickOptions{}
	ctx := context.Background()
	for idx := 0; idx < b.N; idx++ {
		_, done, err := picker.Pick(ctx, opt)
		if err != nil {
			done(balancer.DoneInfo{})
		}
	}
}

func TestChaosPick(t *testing.T) {
	flag.Parse()
	fmt.Printf("start chaos test!svrNum:%d cliNum:%d concurrency:%d exLoad:%d exDelay:%d\n", serverNum, cliNum, concurrency, extraLoad, extraDelay)
	c := newController(serverNum, cliNum)
	c.launch(concurrency)
	go c.updateStatics()
	go c.control(extraLoad, extraDelay)
	time.Sleep(time.Second * 50)
}

func newController(svrNum int, cliNum int) *controller {
	//new servers
	servers := []*testSubConn{}
	var weight uint64 = 10
	if extraWeight > 0 {
		weight = extraWeight
	}
	for i := 0; i < svrNum; i++ {
		weight += extraWeight
		sc := newTestSubConn(fmt.Sprintf("addr_%d", i), weight, "")
		servers = append(servers, sc)
	}
	//new clients
	var clients []balancer.Picker
	scs := map[resolver.Address]balancer.SubConn{}
	for _, v := range servers {
		scs[v.addr] = v
	}
	for i := 0; i < cliNum; i++ {
		wpb := &p2cPickerBuilder{}
		picker := wpb.Build(scs)
		clients = append(clients, picker)
	}

	c := &controller{
		servers: servers,
		clients: clients,
	}
	return c
}

type controller struct {
	servers []*testSubConn
	clients []balancer.Picker
}

func (c *controller) launch(concurrency int) {
	opt := balancer.PickOptions{}
	bkg := context.Background()
	for i := range c.clients {
		for j := 0; j < concurrency; j++ {
			picker := c.clients[i]
			go func() {
				for {
					ctx, cancel := context.WithTimeout(bkg, time.Millisecond*250)
					sc, done, _ := picker.Pick(ctx, opt)
					server := sc.(*testSubConn)
					server.connect(ctx)
					var err error
					if ctx.Err() != nil {
						err = status.Errorf(codes.DeadlineExceeded, "dead")
					}
					cancel()
					cpu := atomic.LoadInt64(&server.cpu)
					md := make(map[string]string)
					md[wmeta.CPUUsage] = strconv.FormatInt(cpu, 10)
					done(balancer.DoneInfo{Trailer: metadata.New(md), Err: err})
					time.Sleep(time.Millisecond * 10)
				}
			}()
		}
	}
}

func (c *controller) updateStatics() {
	for {
		time.Sleep(time.Millisecond * 500)
		for _, sc := range c.servers {
			usage := atomic.LoadInt64(&sc.usage)
			avgCpu := (usage - sc.prevUsage) * 2
			atomic.StoreInt64(&sc.cpu, avgCpu)
			sc.prevUsage = usage
		}
	}

}

func (c *controller) control(extraLoad, extraDelay int64) {
	var chaos int
	for {
		fmt.Printf("\n")
		//make some chaos
		n := rand.Intn(3)
		chaos = n + 1
		for i := 0; i < chaos; i++ {
			if extraLoad > 0 {
				degree := rand.Int63n(extraLoad)
				degree++
				atomic.StoreInt64(&c.servers[i].loadJitter, degree)
				fmt.Printf("set addr_%d load:%d ", i, degree)
			}
			if extraDelay > 0 {
				degree := rand.Int63n(extraDelay)
				atomic.StoreInt64(&c.servers[i].delayJitter, degree)
				fmt.Printf("set addr_%d delay:%dms ", i, degree)
			}
		}
		fmt.Printf("\n")
		sleep := int64(5)
		time.Sleep(time.Second * time.Duration(sleep))
		for _, sc := range c.servers {
			req := atomic.LoadInt64(&sc.reqs)
			qps := (req - sc.prevReq) / sleep
			wait := len(sc.wait)
			sc.prevReq = req
			fmt.Printf("%s qps:%d waits:%d\n", sc.addr.Addr, qps, wait)
		}
		for _, picker := range c.clients {
			p := picker.(*p2cPicker)
			p.printStats()
		}
		fmt.Printf("\n")
		//reset chaos
		for i := 0; i < chaos; i++ {
			atomic.StoreInt64(&c.servers[i].loadJitter, 0)
			atomic.StoreInt64(&c.servers[i].delayJitter, 0)
		}
		chaos = 0
	}
}

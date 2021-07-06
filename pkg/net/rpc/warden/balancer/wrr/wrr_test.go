package wrr

import (
	"context"
	"fmt"
	"google.golang.org/grpc/balancer/base"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-kratos/kratos/pkg/conf/env"
	nmd "github.com/go-kratos/kratos/pkg/net/metadata"
	wmeta "github.com/go-kratos/kratos/pkg/net/rpc/warden/internal/metadata"
	"github.com/go-kratos/kratos/pkg/stat/metric"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

type testSubConn struct {
	addr resolver.Address
}

func (s *testSubConn) UpdateAddresses([]resolver.Address) {

}

// Connect starts the connecting for this SubConn.
func (s *testSubConn) Connect() {
	fmt.Println(s.addr.Addr)
}

func TestBalancerPick(t *testing.T) {
	addr1 := resolver.Address{Addr: "test1", Metadata: wmeta.MD{Weight: 8}}
	addr2 := resolver.Address{Addr: "test2", Metadata: wmeta.MD{Weight: 4, Color: "red"}}
	addr3 := resolver.Address{Addr: "test3", Metadata: wmeta.MD{Weight: 2, Color: "red"}}
	bi := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			&testSubConn{addr: addr1}: {Address: addr1},
			&testSubConn{addr: addr2}: {Address: addr2},
			&testSubConn{addr: addr3}: {Address: addr3},
		},
	}
	b := &wrrPickerBuilder{}
	picker := b.Build(bi)
	res := []string{"test1", "test1", "test1", "test1"}
	for i := 0; i < 3; i++ {
		pr, err := picker.Pick(balancer.PickInfo{Ctx: context.Background()})
		if err != nil {
			t.Fatalf("picker.Pick failed!idx:=%d", i)
		}
		sc := pr.SubConn.(*testSubConn)
		if sc.addr.Addr != res[i] {
			t.Fatalf("the subconn picked(%s),but expected(%s)", sc.addr.Addr, res[i])
		}
	}
	res2 := []string{"test2", "test3", "test2", "test2", "test3", "test2"}
	ctx := nmd.NewContext(context.Background(), nmd.New(map[string]interface{}{"color": "red"}))
	for i := 0; i < 6; i++ {
		pr, err := picker.Pick(balancer.PickInfo{Ctx: ctx})
		if err != nil {
			t.Fatalf("picker.Pick failed!idx:=%d", i)
		}
		sc := pr.SubConn.(*testSubConn)
		if sc.addr.Addr != res2[i] {
			t.Fatalf("the (%d) subconn picked(%s),but expected(%s)", i, sc.addr.Addr, res2[i])
		}
	}
	ctx = nmd.NewContext(context.Background(), nmd.New(map[string]interface{}{"color": "black"}))
	for i := 0; i < 4; i++ {
		pr, err := picker.Pick(balancer.PickInfo{Ctx: ctx})
		if err != nil {
			t.Fatalf("picker.Pick failed!idx:=%d", i)
		}
		sc := pr.SubConn.(*testSubConn)
		if sc.addr.Addr != res[i] {
			t.Fatalf("the (%d) subconn picked(%s),but expected(%s)", i, sc.addr.Addr, res[i])
		}
	}

	// test for env color
	ctx = context.Background()
	env.Color = "red"
	for i := 0; i < 6; i++ {
		pr, err := picker.Pick(balancer.PickInfo{Ctx: ctx})
		if err != nil {
			t.Fatalf("picker.Pick failed!idx:=%d", i)
		}
		sc := pr.SubConn.(*testSubConn)
		if sc.addr.Addr != res2[i] {
			t.Fatalf("the (%d) subconn picked(%s),but expected(%s)", i, sc.addr.Addr, res2[i])
		}
	}
}

func TestBalancerDone(t *testing.T) {
	addr := resolver.Address{Addr: "test1", Metadata: wmeta.MD{Weight: 8}}
	b := &wrrPickerBuilder{}
	bi := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			&testSubConn{addr: addr}: {Address: addr},
		},
	}
	picker := b.Build(bi)

	pr, _ := picker.Pick(balancer.PickInfo{Ctx: context.Background()})
	time.Sleep(100 * time.Millisecond)
	pr.Done(balancer.DoneInfo{Err: status.Errorf(codes.Unknown, "test")})
	err, req := picker.(*wrrPicker).subConns[0].errSummary()
	assert.Equal(t, int64(0), err)
	assert.Equal(t, int64(1), req)

	latency, count := picker.(*wrrPicker).subConns[0].latencySummary()
	expectLatency := float64(100*time.Millisecond) / 1e5
	if latency < expectLatency || latency > (expectLatency+500) {
		t.Fatalf("latency is less than 100ms or greater than 150ms, %f", latency)
	}
	assert.Equal(t, int64(1), count)

	pr, _ = picker.Pick(balancer.PickInfo{Ctx: context.Background()})
	pr.Done(balancer.DoneInfo{Err: status.Errorf(codes.Aborted, "test")})
	err, req = picker.(*wrrPicker).subConns[0].errSummary()
	assert.Equal(t, int64(1), err)
	assert.Equal(t, int64(2), req)
}

func TestErrSummary(t *testing.T) {
	sc := &subConn{
		err: metric.NewRollingCounter(metric.RollingCounterOpts{
			Size:           10,
			BucketDuration: time.Millisecond * 100,
		}),
		latency: metric.NewRollingGauge(metric.RollingGaugeOpts{
			Size:           10,
			BucketDuration: time.Millisecond * 100,
		}),
	}
	for i := 0; i < 10; i++ {
		sc.err.Add(0)
		sc.err.Add(1)
	}
	err, req := sc.errSummary()
	assert.Equal(t, int64(10), err)
	assert.Equal(t, int64(20), req)
}

func TestLatencySummary(t *testing.T) {
	sc := &subConn{
		err: metric.NewRollingCounter(metric.RollingCounterOpts{
			Size:           10,
			BucketDuration: time.Millisecond * 100,
		}),
		latency: metric.NewRollingGauge(metric.RollingGaugeOpts{
			Size:           10,
			BucketDuration: time.Millisecond * 100,
		}),
	}
	for i := 1; i <= 100; i++ {
		sc.latency.Add(int64(i))
	}
	latency, count := sc.latencySummary()
	assert.Equal(t, 50.50, latency)
	assert.Equal(t, int64(100), count)
}

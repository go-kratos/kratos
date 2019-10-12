package wrr

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"

	"github.com/bilibili/kratos/pkg/conf/env"
	nmd "github.com/bilibili/kratos/pkg/net/metadata"
	wmeta "github.com/bilibili/kratos/pkg/net/rpc/warden/internal/metadata"
	"github.com/bilibili/kratos/pkg/stat/metric"

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
	scs[sc1.addr] = sc1
	scs[sc2.addr] = sc2
	scs[sc3.addr] = sc3
	b := &wrrPickerBuilder{}
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
	res2 := []string{"test2", "test3", "test2", "test2", "test3", "test2"}
	ctx := nmd.NewContext(context.Background(), nmd.New(map[string]interface{}{"color": "red"}))
	for i := 0; i < 6; i++ {
		conn, _, err := picker.Pick(ctx, balancer.PickOptions{})
		if err != nil {
			t.Fatalf("picker.Pick failed!idx:=%d", i)
		}
		sc := conn.(*testSubConn)
		if sc.addr.Addr != res2[i] {
			t.Fatalf("the (%d) subconn picked(%s),but expected(%s)", i, sc.addr.Addr, res2[i])
		}
	}
	ctx = nmd.NewContext(context.Background(), nmd.New(map[string]interface{}{"color": "black"}))
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

	// test for env color
	ctx = context.Background()
	env.Color = "red"
	for i := 0; i < 6; i++ {
		conn, _, err := picker.Pick(ctx, balancer.PickOptions{})
		if err != nil {
			t.Fatalf("picker.Pick failed!idx:=%d", i)
		}
		sc := conn.(*testSubConn)
		if sc.addr.Addr != res2[i] {
			t.Fatalf("the (%d) subconn picked(%s),but expected(%s)", i, sc.addr.Addr, res2[i])
		}
	}
}

func TestBalancerDone(t *testing.T) {
	scs := map[resolver.Address]balancer.SubConn{}
	sc1 := &testSubConn{
		addr: resolver.Address{
			Addr: "test1",
			Metadata: wmeta.MD{
				Weight: 8,
			},
		},
	}
	scs[sc1.addr] = sc1
	b := &wrrPickerBuilder{}
	picker := b.Build(scs)

	_, done, _ := picker.Pick(context.Background(), balancer.PickOptions{})
	time.Sleep(100 * time.Millisecond)
	done(balancer.DoneInfo{Err: status.Errorf(codes.Unknown, "test")})
	err, req := picker.(*wrrPicker).subConns[0].errSummary()
	assert.Equal(t, int64(0), err)
	assert.Equal(t, int64(1), req)

	latency, count := picker.(*wrrPicker).subConns[0].latencySummary()
	expectLatency := float64(100*time.Millisecond) / 1e5
	if latency < expectLatency || latency > (expectLatency+500) {
		t.Fatalf("latency is less than 100ms or greater than 150ms, %f", latency)
	}
	assert.Equal(t, int64(1), count)

	_, done, _ = picker.Pick(context.Background(), balancer.PickOptions{})
	done(balancer.DoneInfo{Err: status.Errorf(codes.Aborted, "test")})
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

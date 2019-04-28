package wrr

import (
	"context"
	"fmt"
	"testing"

	nmd "go-common/library/net/metadata"
	wmeta "go-common/library/net/rpc/warden/metadata"

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
}

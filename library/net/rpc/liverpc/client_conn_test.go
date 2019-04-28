package liverpc

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"go-common/library/conf/env"
	"go-common/library/net/metadata"

	. "github.com/smartystreets/goconvey/convey"
)

type testConn struct {
}

func (t *testConn) Read(p []byte) (n int, err error) {
	return
}

func (t *testConn) Write(p []byte) (n int, err error) {
	return
}

func (t *testConn) Close() (err error) {
	return
}

func TestClientConn(t *testing.T) {
	var (
		req  protoReq
		resp protoResp
	)
	req.Header.magic = _magic
	req.Header.checkSum = 0
	req.Header.seq = rand.Uint32()
	req.Header.timestamp = uint32(time.Now().Unix())
	req.Header.reserved = 1
	req.Header.version = 1
	// command: {message_type}controller.method
	req.Header.cmd = make([]byte, 32)
	req.Header.cmd[0] = _cmdReqType
	// serviceMethod: Room.room_init
	copy(req.Header.cmd[1:], []byte("Room.room_init"))

	codec := &ClientConn{rwc: &testConn{}}
	if err := codec.writeRequest(context.TODO(), &req); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := codec.readResponse(context.TODO(), &resp); err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("request header:%+v body:%s", req.Header, req.Body)
	t.Logf("response header:%+v body:%s", resp.Header, resp.Body)
}

func TestGroups(t *testing.T) {
	Convey("groups", t, func() {
		client := Client{}

		// env before empty
		env.Color = "red"
		client.addrs.Store(GroupAddrs{
			"red":  []string{"r1", "r2"},
			"":     {"empty1", "empty2"},
			"dark": {"d1", "d2"},
			"blue": []string{"b1", "b2"},
		})
		addrs, _ := client.fetchAddrs(context.TODO(), nil)
		So(addrs, ShouldResemble, []string{"r1", "r2"})

		// can get env
		env.Color = ""
		addrs, _ = client.fetchAddrs(context.TODO(), nil)
		So(addrs, ShouldResemble, []string{"empty1", "empty2"})

		// from metadata.color
		env.Color = "red"
		addrs, _ = client.fetchAddrs(
			metadata.NewContext(context.TODO(), metadata.MD{metadata.Color: "blue"}), nil)
		So(addrs, ShouldResemble, []string{"b1", "b2"})

		// header source_group always first
		addrs, _ = client.fetchAddrs(
			metadata.NewContext(context.TODO(), metadata.MD{metadata.Color: "blue"}),
			&Args{Header: &Header{SourceGroup: "dark"}})
		So(addrs, ShouldResemble, []string{"d1", "d2"})
	})

}

package v1

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/broadcast/model"
	"go-common/library/log"
	"go-common/library/naming/discovery"
	"go-common/library/net/netutil/breaker"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/rpc/warden/resolver"
	xtime "go-common/library/time"
)

func testInit() ZergClient {
	log.Init(nil)
	conf := &warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 10),
		Timeout: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(3 * time.Second),
			Sleep:   xtime.Duration(3 * time.Second),
			Bucket:  10,
			Ratio:   0.3,
			Request: 20,
		},
	}
	wc := warden.NewClient(conf)
	resolver.Register(discovery.New(nil))
	conn, err := wc.Dial(context.TODO(), "discovery://default/push.interface.broadcast")
	if err != nil {
		panic(err)
	}
	return NewZergClient(conn)
}

func TestPushMsg(t *testing.T) {
	client := testInit()
	time.Sleep(10 * time.Second)
	client.PushMsg(context.Background(), &PushMsgReq{
		Keys:    []string{"test"},
		ProtoOp: model.OpSendMsg,
		Proto: &model.Proto{
			Ver:       0,
			SeqId:     0,
			Operation: model.OpSendMsgReply,
			Body:      []byte("{\"test1111111\"}"),
		},
	})
}

/*
func TestBroadcastMsg(t *testing.T) {
	client := testInit()
	client.Broadcast(context.Background(), 102, &model.Proto{
		Ver:       0,
		SeqId:     0,
		Operation: define.OP_SEND_SMS_REPLY,
		Body:      []byte("{\"test broadcast 104\"}"),
	})
}

func TestBroadcastRoom(t *testing.T) {
	client := testInit()
	client.BroadcastRoom(context.Background(), "test_room", &model.Proto{
		Ver:       0,
		SeqId:     0,
		Operation: define.OP_SEND_SMS_REPLY,
		Body:      []byte("{\"test broadcast\"}"),
	})
}

func TestRooms(t *testing.T) {
	client := testInit()
	rooms, err := client.Rooms(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("TestRooms.rooms:%v", rooms)
}
*/

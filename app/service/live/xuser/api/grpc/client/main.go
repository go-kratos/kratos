package main

import (
	"context"
	"flag"
	"github.com/davecgh/go-spew/spew"
	"log"
	"time"

	"go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
)

var name, addr string

func init() {
	flag.StringVar(&name, "name", "lily", "name")
	flag.StringVar(&addr, "addr", "127.0.0.1:9004", "server addr")
}

func main() {
	flag.Parse()
	cfg := &warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 3),
		Timeout: xtime.Duration(time.Second * 3),
	}
	cc, err := warden.NewClient(cfg).Dial(context.Background(), addr)
	if err != nil {
		log.Fatalf("new client failed!err:=%v", err)
		return
	}

	client := v1.NewRoomAdminClient(cc)

	//resp, _ := client.IsAny(context.Background(), &v1.RoomAdminShowEntryReq{
	//	Uid: 1001,
	//})
	//spew.Dump("IsAny", resp)
	//
	//resp2, _ := client.GetByUid(context.Background(), &v1.RoomAdminGetByUidReq{
	//	Uid:  1001,
	//	Page: 1,
	//})
	//spew.Dump("GetByUid", resp2)
	//
	//resp3, err := client.SearchForAdmin(context.Background(), &v1.RoomAdminSearchForAdminReq{
	//	Uid:     28007796,
	//	KeyWord: "1001",
	//})
	//spew.Dump("SearchForAdmin", resp3, err)
	//
	//resp4, err := client.GetByAnchor(context.Background(), &v1.RoomAdminGetByAnchorReq{
	//	Page: 0,
	//	Uid:  28007796,
	//})
	//spew.Dump("GetByAnchor", resp4, err)
	//
	//resp5, err := client.Dismiss(context.Background(), &v1.RoomAdminDismissAdminReq{
	//	Uid:      1001,
	//	AnchorId: 28007796,
	//})
	//spew.Dump("Dismiss", resp5, err)
	//
	//resp6, err := client.Appoint(context.Background(), &v1.RoomAdminAddReq{
	//	Uid:      1001,
	//	AnchorId: 28007796,
	//})
	//spew.Dump("Appoint", resp6, err)

	resp7, err := client.IsAdminShort(context.Background(), &v1.RoomAdminIsAdminShortReq{
		Uid:    2101323,
		Roomid: 5171,
	})
	spew.Dump("IsAdminShort", resp7, err)
}

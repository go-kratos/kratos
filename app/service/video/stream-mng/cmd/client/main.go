package main

import (
	"context"
	"fmt"

	"go-common/app/service/video/stream-mng/api/v1"
	"go-common/library/ecode"
	"go-common/library/net/rpc/warden"
)

func main() {

	cc, err := warden.NewClient(nil).Dial(context.Background(), "127.0.0.1:9000")

	if err != nil {
		return
	}

	client := v1.NewStreamClient(cc)

	resp, err := client.GetRoomIDByStreamName(context.Background(), &v1.GetRoomIDByStreamNameReq{
		StreamName: "",
	})

	fmt.Printf("%v=%v", resp, err)

	e := ecode.Cause(err)

	fmt.Printf("%v=%v=%v=%v", e, e.Code(), e.Message(), e.Details())
}

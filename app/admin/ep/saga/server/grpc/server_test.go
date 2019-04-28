package grpc

import (
	"context"
	"testing"

	"go-common/app/admin/ep/saga/api/grpc/v1"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/rpc/warden/resolver"
	"go-common/library/net/rpc/warden/resolver/direct"
)

func TestGRPC(t *testing.T) {
	resolver.Register(direct.New())
	conn, err := warden.NewClient(nil).Dial(context.Background(), "direct://default/127.0.0.1:9000")
	if err != nil {
		t.Fatal(err)
	}
	client := v1.NewSagaAdminClient(conn)
	if _, err = client.PushMsg(context.Background(), &v1.PushMsgReq{Username: []string{"wuwei"}, Content: "test"}); err != nil {
		t.Error(err)
	}
}

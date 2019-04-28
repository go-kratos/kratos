package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID discovery appid.
const AppID = "main.appsvr.shareservice"

// NewClient new grpc client.
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (ShareClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	return NewShareClient(conn), nil
}

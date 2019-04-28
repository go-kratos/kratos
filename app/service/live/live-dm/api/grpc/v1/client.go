package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

//AppID 弹幕服务discoverID
const AppID = "live.livedm"

//NewClient 弹幕服务client创建
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (DMClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}

	return NewDMClient(conn), nil
}

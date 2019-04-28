package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

//AppID 直播风控服务discoverID
const AppID = "live.riskcontrol"

//NewClient 直播风控服务client创建
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (IsForbiddenClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://detault/"+AppID)
	if err != nil {
		return nil, err
	}

	return NewIsForbiddenClient(conn), nil
}

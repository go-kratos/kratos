package v1

import (
	"context"
	"fmt"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID .
const AppID = "app.resource"

// NewClient new grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (AppResourceClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AppID))
	if err != nil {
		return nil, err
	}
	return NewAppResourceClient(cc), nil
}

//go:generate $GOPATH/src/go-common/app/tool/warden/protoc.sh

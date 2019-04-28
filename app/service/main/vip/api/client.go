package v1

import (
	"fmt"

	"go-common/library/net/rpc/warden"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// AppID .
const AppID = "account.service.vip"

// NewClient new grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (VipClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AppID))
	if err != nil {
		return nil, err
	}
	return NewVipClient(cc), nil
}

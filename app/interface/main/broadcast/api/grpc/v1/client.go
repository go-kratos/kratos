package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// NewClient new a client.
func NewClient(target string, cfg *warden.ClientConfig, opts ...grpc.DialOption) (ZergClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), target)
	if err != nil {
		return nil, err
	}
	return NewZergClient(cc), nil
}

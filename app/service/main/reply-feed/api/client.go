package v1

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

// AppID unique app id for service discovery
const AppID = "community.reply.feed"

// NewClient new identify grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (ReplyFeedClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	return NewReplyFeedClient(conn), nil
}

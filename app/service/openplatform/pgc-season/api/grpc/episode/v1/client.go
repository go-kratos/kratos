package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// DiscoveryID season
const DiscoveryID = "season.service"

// NewClient new identify grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (EpisodeClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+DiscoveryID)
	if err != nil {
		return nil, err
	}
	return NewEpisodeClient(conn), nil
}

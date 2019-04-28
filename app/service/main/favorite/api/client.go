package api

import (
	"context"
	"fmt"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID .
const AppID = "community.service.favorite"

// New fav service client
func New(c *warden.ClientConfig, opts ...grpc.DialOption) (FavoriteClient, error) {
	client := warden.NewClient(c, opts...)
	conn, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AppID))
	return NewFavoriteClient(conn), err
}

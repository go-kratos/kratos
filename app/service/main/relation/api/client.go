package api

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

// AppID AppID
const AppID = "account.service.relation"

// NewClient new member grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (RelationClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	return NewRelationClient(conn), nil
}

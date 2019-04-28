package v1

import (
	"context"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID unique app id for service discovery
const AppID = "live.xuser"

// Client client
type Client struct {
	UserExpClient
	VipClient
	GuardClient
}

// NewXuserRoomAdminClient new member grpc client
func NewXuserRoomAdminClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (RoomAdminClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	//conn, err := client.Dial(context.Background(), "127.0.0.1:9000")

	if err != nil {
		return nil, err
	}
	return NewRoomAdminClient(conn), nil
}

// NewClient new resource gRPC client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (*Client, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	cli := &Client{}
	cli.UserExpClient = NewUserExpClient(conn)
	cli.GuardClient = NewGuardClient(conn)
	cli.VipClient = NewVipClient(conn)
	return cli, nil
}

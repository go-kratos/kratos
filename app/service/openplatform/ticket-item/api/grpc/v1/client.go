package v1

import (
	"context"

	"go-common/library/log"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/rpc/warden/balancer/wrr"

	"google.golang.org/grpc"
)

// AppID unique app id for service diSCovery
const AppID = "ticket.service.item"

//Client 客户端枚举
type Client struct {
	IC ItemClient
	GC GuestClient
	BC BulletinClient
	VC VenueClient
	PC PlaceClient
	AC AreaClient
	SC SeatClient
}

// New 新建客户端实例
func New(c *warden.ClientConfig) (*Client, error) {
	client := warden.NewClient(c, grpc.WithBalancerName(wrr.Name))
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		log.Error("client can not connect server: %v", err)
		return nil, err
	}
	return &Client{
		IC: NewItemClient(conn),
		GC: NewGuestClient(conn),
		BC: NewBulletinClient(conn),
		VC: NewVenueClient(conn),
		PC: NewPlaceClient(conn),
		AC: NewAreaClient(conn),
		SC: NewSeatClient(conn),
	}, nil
}

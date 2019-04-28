package v1

import (
	"context"
	"fmt"

	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID 本服务的discoveryID
const AppID = "ticket.service.sales"

// TradeSalesClient .
type TradeSalesClient interface {
	TradeClient
	PromotionClient
}

var _ TradeSalesClient = client{}

type client struct {
	TradeClient
	PromotionClient
}

// NewClient include TradeClient adn SalesClient
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (TradeSalesClient, error) {
	cc, err := warden.NewClient(cfg, opts...).Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AppID))
	if err != nil {
		return nil, err
	}
	return client{TradeClient: NewTradeClient(cc), PromotionClient: NewPromotionClient(cc)}, nil
}

// Package v1 .
// NOTE: need registery discovery resovler into grpc before use this client
/*
import (
	"go-common/library/naming/discovery"
	"go-common/library/net/rpc/warden/resolver"
)

func main() {
	resolver.Register(discovery.New(nil))
}
*/
package v1

import (
	"context"

	"google.golang.org/grpc"

	"go-common/library/net/rpc/warden"
)

// AppID unique app id for service discovery
const AppID = "passport.service.identify"

// NewClient new identify grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (IdentifyClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
	if err != nil {
		return nil, err
	}
	return NewIdentifyClient(conn), nil
}

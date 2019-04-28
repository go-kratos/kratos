package livezk

import (
	"context"
	"fmt"
	"net"

	"go-common/library/naming"
	lz "go-common/library/naming/livezk"
	"go-common/library/net/ip"
)

// Register self grpc service to live zookeeper
func Register(config *lz.Zookeeper, addr string, discoveryID string) (context.CancelFunc, error) {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	z, err := lz.New(config)
	if err != nil {
		return nil, err
	}
	internalIP := ip.InternalIP()
	ins := &naming.Instance{
		AppID: discoveryID,
		Addrs: []string{fmt.Sprintf("grpc://%s:%s", internalIP, port)},
	}
	return z.Register(context.Background(), ins)
}

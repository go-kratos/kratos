package lancergateway

import (
	"time"
	"errors"
	"fmt"

	"go-common/library/net/rpc/warden/resolver"
	"go-common/library/net/rpc/warden/balancer/wrr"
	"go-common/library/naming/discovery"
	xtime "go-common/library/time"

	"google.golang.org/grpc"
)

type Config struct {
	AppId   string        `toml:"appId"`
	Timeout xtime.Duration `toml:"timeout"`
	Subset  int           `toml:"subset"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of LancerGateway can't be nil")
	}

	if c.AppId == "" {
		c.AppId = "datacenter.lancer.gateway2-server"
	}

	if c.Timeout == 0 {
		c.Timeout = xtime.Duration(time.Second * 5)
	}

	if c.Subset == 0 {
		c.Subset = 5
	}
	return nil
}

func init() {
	resolver.Register(discovery.Builder())
}

// NewClient new member grpc client
func NewClient(c *Config) (Gateway2ServerClient, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBalancerName(wrr.Name),
	}

	if c.Timeout != 0 {
		opts = append(opts, grpc.WithTimeout(time.Duration(c.Timeout)))
	}

	conn, err := grpc.Dial(fmt.Sprintf("discovery://default/%s?subset=%d", c.AppId, c.Subset), opts...)
	if err != nil {
		return nil, err
	}
	return NewGateway2ServerClient(conn), nil
}

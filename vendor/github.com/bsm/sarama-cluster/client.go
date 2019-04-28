package cluster

import (
	"errors"
	"sync/atomic"

	"github.com/Shopify/sarama"
)

var errClientInUse = errors.New("cluster: client is already used by another consumer")

// Client is a group client
type Client struct {
	sarama.Client
	config Config

	inUse uint32
}

// NewClient creates a new client instance
func NewClient(addrs []string, config *Config) (*Client, error) {
	if config == nil {
		config = NewConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	client, err := sarama.NewClient(addrs, &config.Config)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client, config: *config}, nil
}

// ClusterConfig returns the cluster configuration.
func (c *Client) ClusterConfig() *Config {
	cfg := c.config
	return &cfg
}

func (c *Client) claim() bool {
	return atomic.CompareAndSwapUint32(&c.inUse, 0, 1)
}

func (c *Client) release() {
	atomic.CompareAndSwapUint32(&c.inUse, 1, 0)
}

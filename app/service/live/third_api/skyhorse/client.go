package skyhorse

import "go-common/library/net/http/blademaster"

type Client struct {
	conf *blademaster.ClientConfig
}

func (c *Client) getConf() *blademaster.ClientConfig {
	return c.conf
}

func New(c *blademaster.ClientConfig) *Client {
	return &Client{conf: c}
}

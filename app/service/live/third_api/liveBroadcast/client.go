package liveBroadcast

import "go-common/library/net/http/blademaster"

//Client  客户端
type Client struct {
	conf *blademaster.ClientConfig
}

func (c *Client) getConf() *blademaster.ClientConfig {
	return c.conf
}

//New 创建
func New(c *blademaster.ClientConfig) *Client {
	return &Client{conf: c}
}

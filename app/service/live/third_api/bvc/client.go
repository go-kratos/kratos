package bvc

import (
	"go-common/library/net/http/blademaster"
)

type Client struct {
	conf    *blademaster.ClientConfig
	bvcHost string
	mock    string
}

func (c *Client) getConf() *blademaster.ClientConfig {
	return c.conf
}
func (c *Client) getBvcHost(def string) string {
	if c.bvcHost == "" {
		return def
	}
	return c.bvcHost
}

func New(c *blademaster.ClientConfig, bvcHost string, mock string) *Client {
	return &Client{
		conf:    c,
		bvcHost: bvcHost,
		mock:    mock,
	}
}

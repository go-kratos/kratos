package http

import (
	"go-common/app/interface/video/portal/conf"
	bm "go-common/library/net/http/blademaster"
)

func txStreamForwardingConf(c *bm.Context) {
	thisconf := *conf.Conf
	c.JSONMap(map[string]interface{}{"message": "ok", "data": thisconf.StreamForward}, nil)
}

// bvcStreamForwardingConf bvc转推白名单
func bvcStreamForwardingConf(c *bm.Context) {
	thisconf := *conf.Conf
	c.JSONMap(map[string]interface{}{"message": "ok", "data": thisconf.BvcStreamForward}, nil)
}

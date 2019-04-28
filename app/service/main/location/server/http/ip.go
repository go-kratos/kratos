package http

import (
	"strings"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// info ip info.
func info(c *bm.Context) {
	var (
		ip    string
		query = c.Request.Form
	)
	if ip = query.Get("ip"); ip == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.Info(c, ip))
}

// infos ip info.
func infos(c *bm.Context) {
	var (
		ips   string
		query = c.Request.Form
	)
	if ips = query.Get("ips"); ips == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.Infos(c, strings.Split(ips, ",")))
}

// infoComplete get whole ip info.
func infoComplete(c *bm.Context) {
	var (
		ip    string
		query = c.Request.Form
	)
	if ip = query.Get("ip"); ip == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.InfoComplete(c, ip))
}

// infosComplete get whole ip infos.
func infosComplete(c *bm.Context) {
	var (
		ips   string
		query = c.Request.Form
	)
	if ips = query.Get("ips"); ips == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.InfosComplete(c, strings.Split(ips, ",")))
}

// anonym ip info.
func anonym(c *bm.Context) {
	var (
		ip    string
		query = c.Request.Form
	)
	if ip = query.Get("ip"); ip == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(svr.Anonym(ip))
}

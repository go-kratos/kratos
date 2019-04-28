package http

import (
	"strconv"

	"go-common/library/log"
	"go-common/library/net/http/blademaster"
)

// setContextMid 把form中的mid写入context中，用以调用interface的http接口
func setContextMid(c *blademaster.Context) {
	var _, ok = c.Get("mid")
	if ok {
		return
	}
	var midstr = c.Request.Form.Get("mid")
	if midstr == "" {
		return
	}

	var mid, err = strconv.ParseInt(midstr, 10, 64)
	if err != nil {
		log.Error("parse mid fail, midstr=%s, err=%v", midstr, err)
		return
	}

	c.Set("mid", mid)
}

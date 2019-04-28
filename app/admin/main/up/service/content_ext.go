package service

import (
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"net/http"
)

//BmHTTPErrorWithMsg return customed message to client
func BmHTTPErrorWithMsg(c *bm.Context, err error, msg string) {
	if c.IsAborted() {
		return
	}
	c.Error = err
	bcode := ecode.Cause(err)
	if msg == "" {
		msg = err.Error()
	}
	c.Render(http.StatusOK, render.JSON{
		Code:    bcode.Code(),
		Message: msg,
		Data:    nil,
	})
}

//BmGetStringOrDefault util to get string from context
func BmGetStringOrDefault(c *bm.Context, key string, defaul string) (value string, exist bool) {
	i, exist := c.Get(key)

	if !exist {
		value = defaul
		return
	}

	value, exist = i.(string)
	if !exist {
		value = defaul
	}
	return
}

//BmGetInt64OrDefault util to get int64 from context
func BmGetInt64OrDefault(c *bm.Context, key string, defaul int64) (value int64, exist bool) {
	i, exist := c.Get(key)

	if !exist {
		value = defaul
		return
	}

	value, exist = i.(int64)
	if !exist {
		value = defaul
	}
	return
}

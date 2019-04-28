package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/stat/prom"
)

func dynamicNew(c *bm.Context) {
	params := c.Request.Form
	data, err := externalSvc.DynamicNew(c, params.Encode())
	dynamicResult(c, data, err)
}

func dynamicCount(c *bm.Context) {
	params := c.Request.Form
	data, err := externalSvc.DynamicCount(c, params.Encode())
	dynamicResult(c, data, err)
}

func dynamicHistory(c *bm.Context) {
	params := c.Request.Form
	data, err := externalSvc.DynamicHistory(c, params.Encode())
	dynamicResult(c, data, err)
}

func dynamicResult(c *bm.Context, data json.RawMessage, err error) {
	params := c.Request.Form
	path := c.Request.URL.Path
	code := ecode.OK
	if err != nil {
		code = ecode.Int(-22)
		c.JSON(nil, code)
	} else {
		if !bytes.Contains(data, []byte(`"code":0`)) {
			var res struct {
				Code int `json:"code"`
			}
			json.Unmarshal(data, &res)
			code = ecode.Int(res.Code)
		}
		c.Bytes(http.StatusOK, "text/json; charset=utf-8", data)
	}
	prom.HTTPServer.Incr("no_user", path[1:], code.Error())
	prom.HTTPServer.Timing("no_user", int64(time.Since(time.Now())/time.Millisecond), path[1:])
	log.Info("method:%s,mid:%v,ip:%s,params:%s,ret:%d[%s] stack:%+v", path, params.Get("uid"), metadata.String(c, metadata.RemoteIP), params.Encode(), code.Code(), code.Message(), err)
}

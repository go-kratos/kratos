package http

import (
	"strconv"
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

func readPing(c *bm.Context) {
	var (
		buvid  = buvid(c)
		aid    int64
		mid    int64
		ip     = metadata.String(c, metadata.RemoteIP)
		cur    = time.Now().Unix()
		source = c.Request.Form.Get("source")
		err    error
	)
	if aid, err = strconv.ParseInt(c.Request.Form.Get("aid"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if i, ok := c.Get("mid"); ok {
		mid = i.(int64)
	} else {
		mid = 0
	}
	c.JSON(nil, artSrv.ReadPing(c, buvid, aid, mid, ip, cur, source))
}

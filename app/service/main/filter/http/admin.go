package http

import (
	"strings"

	"go-common/app/service/main/filter/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func filterTest(c *bm.Context) {
	var (
		areas  []string
		params = c.Request.Form
	)
	msgStr := params.Get("msg")
	if msgStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	areaStr := params.Get("area")
	if areaStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	areas = strings.Split(strings.TrimSpace(areaStr), ",")

	ws, bs := svc.FilterTest(c, areas, msgStr)
	var data = map[string]interface{}{
		"whits": ws,
		"hits":  bs,
	}
	c.JSON(data, nil)
}

func testKey(c *bm.Context) {
	var (
		err    error
		rs     []*model.KeyTestResult
		params = c.Request.Form
	)
	keyStr := params.Get("key")
	areaStr := params.Get("area")
	msgStr := params.Get("msg")
	if areaStr == "" || keyStr == "" {
		log.Error("strconv.ParseInt(%s,%s) err(%+v)", keyStr, areaStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	areas := strings.Split(areaStr, ",")
	if rs, err = svc.KeyTest(c, keyStr, msgStr, areas); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var data = map[string]interface{}{
		"rules": rs,
	}
	c.JSON(data, nil)
}

package http

import (
	"strconv"
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_headerBuvid    = "Buvid"
	_headerDeviceID = "Device-ID"
)

// interest get interest list
func interest(c *bm.Context) {
	var (
		header = c.Request.Header
		params = c.Request.Form
	)
	mobiApp := params.Get("mobi_app")
	buvid := header.Get(_headerBuvid)
	c.JSON(guideSvc.Interest(mobiApp, buvid, time.Now()), nil)
}

// interest2 get interest list
func interest2(c *bm.Context) {
	var (
		header = c.Request.Header
		params = c.Request.Form
	)
	mobiApp := params.Get("mobi_app")
	device := params.Get("device")
	buildStr := params.Get("build")
	buvid := header.Get(_headerBuvid)
	if buvid == "" || len(buvid) <= 5 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	dvcid := header.Get(_headerDeviceID)
	c.JSON(guideSvc.Interest2(mobiApp, device, dvcid, buvid, build, time.Now()), nil)
}

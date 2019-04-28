package http

import (
	"strconv"

	bm "go-common/library/net/http/blademaster"
)

func authors(c *bm.Context) {
	var (
		mid    int64
		author int64
		params = c.Request.Form
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	author, _ = strconv.ParseInt(params.Get("author"), 10, 64)
	authors := artSrv.Authors(c, mid, author)
	c.JSONMap(map[string]interface{}{"authors": authors}, nil)
}

// func authors(c *bm.Context) {
// 	var (
// 		mid      int64
// 		userID   int64
// 		request  = c.Request
// 		params   = request.Form
// 		clientIP = metadata.String(c, metadata.RemoteIP)
// 		err      error
// 		authors  *artmdl.RecommendAuthors
// 	)
// 	author := params.Get("author")
// 	if author == "" {
// 		err = ecode.RequestErr
// 		c.JSON(nil, err)
// 		return
// 	}
// 	if mid, err = strconv.ParseInt(author, 10, 64); err != nil {
// 		c.JSON(nil, err)
// 		return
// 	}
// 	mobiApp := params.Get("mobi_app")
// 	device := params.Get("device")
// 	plat := artmdl.Plat(mobiApp, device)
// 	platform := artmdl.Client(plat)
// 	buildStr := params.Get("build")
// 	build, _ := strconv.Atoi(buildStr)
// 	buvid := buvid(c)
// 	if midInter, ok := c.Get("mid"); ok {
// 		userID = midInter.(int64)
// 	}

// 	authors, err = artSrv.RecommendAuthors(c, platform, mobiApp, device, build, clientIP, userID, buvid, mid)
// 	c.JSON(authors, err)
// }

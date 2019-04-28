package http

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// @params FeedListReq
// @router get /bbq/app-bbq/feed/list/
// @response FeedListResponse
func feedList(c *bm.Context) {
	arg := &v1.FeedListRequest{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	mid, _ := c.Get("mid")
	arg.MID = mid.(int64)
	dev, _ := c.Get("device")
	arg.Device = dev.(*bm.Device)
	b, _ := c.Get("BBQBase")
	arg.BUVID = b.(*v1.Base).BUVID
	resp, err := srv.FeedList(c, arg)
	c.JSON(resp, err)

	// 埋点
	if err != nil {
		return
	}
	var svidList []int64
	for _, v := range resp.List {
		svidList = append(svidList, v.SVID)
	}
	ext := &struct {
		Svids []int64
	}{
		Svids: svidList,
	}
	uiLog(c, model.ActionFeedList, ext)
}

// @params mid
// @router get /bbq/app-bbq/feed/update_num
// @response FeedUpdateNumResponse
func feedUpdateNum(c *bm.Context) {
	mid, _ := c.Get("mid")
	c.JSON(srv.FeedUpdateNum(c, mid.(int64)))
}

// @params SpaceSvListRequest
// @router get /bbq/app-bbq/space/sv/list/
// @response SpaceSvListResponse
func spaceSvList(c *bm.Context) {
	arg := &v1.SpaceSvListRequest{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	mid, exists := c.Get("mid")
	if exists {
		arg.MID = mid.(int64)
	}
	dev, _ := c.Get("device")
	arg.Device = dev.(*bm.Device)
	arg.Size = model.SpaceListLen

	c.JSON(srv.SpaceSvList(c, arg))
}

// @params SpaceSvListRequest
// @router get /bbq/app-bbq/detail/sv/list/
// @response SpaceSvListResponse
func detailSvList(c *bm.Context) {
	arg := &v1.SpaceSvListRequest{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	mid, exists := c.Get("mid")
	if exists {
		arg.MID = mid.(int64)
	}
	dev, _ := c.Get("device")
	arg.Device = dev.(*bm.Device)
	// 暂时设置size=3
	arg.Size = 3
	arg.CursorNext = ""
	arg.CursorPrev = ""
	c.JSON(srv.SpaceSvList(c, arg))
}

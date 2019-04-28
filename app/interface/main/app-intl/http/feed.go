package http

import (
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card"
	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/feed"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_headerBuvid     = "Buvid"
	_headerDisplayID = "Display-ID"
	_headerDeviceID  = "Device-ID"
)

func feedIndex(c *bm.Context) {
	var mid int64
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	header := c.Request.Header
	buvid := header.Get(_headerBuvid)
	disid := header.Get(_headerDisplayID)
	dvcid := header.Get(_headerDeviceID)
	param := &feed.IndexParam{}
	// get params
	if err := c.Bind(param); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	column, ok := cdm.Columnm[param.Column]
	if !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// 兼容老的style逻辑，3为新单列
	style := int(cdm.Columnm[param.Column])
	if style == 1 {
		style = 3
	}
	// check params
	plat := model.Plat(param.MobiApp, param.Device)
	now := time.Now()
	// index
	data, userFeature, isRcmd, newUser, code, autoPlay, feedclean, autoPlayInfoc, err := feedSvc.Index(c, buvid, mid, plat, param, now, style)
	autoplayCard := struct {
		Column          cdm.ColumnStatus `json:"column"`
		AutoplayCard    int8             `json:"autoplay_card"`
		FeedCleanAbtest int8             `json:"feed_clean_abtest"`
	}{Column: column, AutoplayCard: autoPlay, FeedCleanAbtest: feedclean}
	c.JSON(struct {
		Item   []card.Handler `json:"items"`
		Config interface{}    `json:"config"`
	}{Item: data, Config: autoplayCard}, err)
	if err != nil {
		return
	}
	// infoc
	items := make([]*ai.Item, 0, len(data))
	for _, item := range data {
		items = append(items, item.Get().Rcmd)
	}
	feedSvc.IndexInfoc(c, mid, plat, param.Build, buvid, disid, "/x/intl/feed/index", userFeature, style, code, items, isRcmd, param.Pull, newUser, now, "", dvcid, param.Network, param.Flush, autoPlayInfoc, param.DeviceType)
}

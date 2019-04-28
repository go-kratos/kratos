package http

import (
	"fmt"
	"net/http"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func queryBlacklist(c *bm.Context) {
	v := new(struct {
		FromTime int64  `form:"from_time"`
		ToTime   int64  `form:"to_time"`
		Type     int    `form:"type"`
		Reason   int    `form:"reason"`
		MID      int64  `form:"mid"`
		Nickname string `form:"nickname"`
		AID      int64  `form:"aid"`
		From     int    `form:"from" validate:"min=0" default:"0"`
		Limit    int    `form:"limit" validate:"min=1" default:"20"`
		Sort     string `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	total, blacklist, err := svr.QueryBlacklist(v.FromTime, v.ToTime, v.Type, v.Reason, v.MID, v.Nickname, v.AID, v.From, v.Limit, v.Sort)
	if err != nil {
		log.Error("growup svr.QueryBlacklist error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    blacklist,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func exportBlack(c *bm.Context) {
	v := new(struct {
		FromTime int64  `form:"from_time"`
		ToTime   int64  `form:"to_time"`
		Type     int    `form:"type"`
		Reason   int    `form:"reason"`
		MID      int64  `form:"mid"`
		Nickname string `form:"nickname"`
		AID      int64  `form:"aid"`
		From     int    `form:"from" validate:"min=0" default:"0"`
		Limit    int    `form:"limit" validate:"min=1" default:"20"`
		Sort     string `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	content, err := svr.ExportBlacklist(v.FromTime, v.ToTime, v.Type, v.Reason, v.MID, v.Nickname, v.AID, v.From, v.Limit, v.Sort)
	if err != nil {
		log.Error("growup svr.ExportBlacklist error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "blacklist"),
	})
}

func recoverBlacklist(c *bm.Context) {
	v := new(struct {
		AID  int64 `form:"aid"`
		Type int   `form:"type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.RecoverBlacklist(v.AID, v.Type)
	if err != nil {
		log.Error("growup svr.RecoverBlacklist error(%v)", err)
	}
	c.JSON(nil, err)
}

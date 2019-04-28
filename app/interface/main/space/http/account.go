package http

import (
	"strconv"
	"strings"

	"go-common/app/interface/main/space/conf"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func navNum(c *bm.Context) {
	var (
		vmid, mid int64
		err       error
	)
	midStr := c.Request.Form.Get("mid")
	if vmid, err = strconv.ParseInt(midStr, 10, 64); err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(spcSvc.NavNum(c, mid, vmid), nil)
}

func upStat(c *bm.Context) {
	var (
		mid int64
		err error
	)
	midStr := c.Request.Form.Get("mid")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(spcSvc.UpStat(c, mid))
}

func myInfo(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(spcSvc.MyInfo(c, mid))
}

func notice(c *bm.Context) {
	v := new(struct {
		Mid int64 `form:"mid" validate:"gt=0"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(spcSvc.Notice(c, v.Mid))
}

func setNotice(c *bm.Context) {
	v := new(struct {
		Notice string `form:"notice" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	notice := strings.Trim(v.Notice, " ")
	if len([]rune(notice)) > conf.Conf.Rule.MaxNoticeLen {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, spcSvc.SetNotice(c, mid, notice))
}

func accTags(c *bm.Context) {
	v := new(struct {
		Mid int64 `form:"mid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(spcSvc.AccTags(c, v.Mid))
}

func setAccTags(c *bm.Context) {
	v := new(struct {
		Tags string `form:"tags" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	tags := strings.Split(v.Tags, ",")
	var addTags []string
	for _, v := range tags {
		if tag := strings.TrimSpace(v); tag != "" {
			addTags = append(addTags, tag)
		}
	}
	if len(addTags) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spcSvc.SetAccTags(c, strings.Join(addTags, ","), c.Request.Header.Get("Cookie")))
}

func accInfo(c *bm.Context) {
	var (
		mid, vmid int64
		err       error
	)
	vmidStr := c.Request.Form.Get("mid")
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(spcSvc.AccInfo(c, mid, vmid))
}

func lastPlayGame(c *bm.Context) {
	var (
		mid, vmid int64
		err       error
	)
	vmidStr := c.Request.Form.Get("mid")
	if vmid, err = strconv.ParseInt(vmidStr, 10, 64); err != nil || vmid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(spcSvc.LastPlayGame(c, mid, vmid))
}

func themeList(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(spcSvc.ThemeList(c, mid))
}

func themeActive(c *bm.Context) {
	var (
		themeID int64
		err     error
	)
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	themeIDStr := c.Request.Form.Get("theme_id")
	if themeID, err = strconv.ParseInt(themeIDStr, 10, 64); err != nil || themeID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spcSvc.ThemeActive(c, mid, themeID))
}

func relation(c *bm.Context) {
	v := new(struct {
		Vmid int64 `form:"mid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(spcSvc.Relation(c, mid, v.Vmid), nil)
}

func clearCache(c *bm.Context) {
	v := new(struct {
		Msg string `form:"msg" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, spcSvc.ClearCache(c, v.Msg))
}

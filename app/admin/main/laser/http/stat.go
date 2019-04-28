package http

import (
	"fmt"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
	"net/http"
	"sort"
	"time"
)

func recheckPanel(c *bm.Context) {
	v := new(struct {
		TypeIDStr string `form:"type_id"`
		StartDate int64  `form:"start_date"`
		EndDate   int64  `form:"end_date"`
		UName     string `form:"uname"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var typeIDS []int64
	if typeIDS, err = xstr.SplitInts(v.TypeIDStr); err != nil {
		c.JSON(nil, err)
		return
	}
	recheckViews, err := svc.ArchiveRecheck(c, typeIDS, v.UName, v.StartDate, v.EndDate)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	sort.Slice(recheckViews, func(i, j int) bool {
		return recheckViews[i].Date > recheckViews[j].Date
	})
	c.JSON(recheckViews, nil)
}

func recheckUser(c *bm.Context) {
	v := new(struct {
		TypeIDStr string `form:"type_id"`
		StartDate int64  `form:"start_date"`
		EndDate   int64  `form:"end_date"`
		UName     string `form:"uname"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	var typeIDS []int64
	if typeIDS, err = xstr.SplitInts(v.TypeIDStr); err != nil {
		c.JSON(nil, err)
		return
	}
	recheckViews, err := svc.UserRecheck(c, typeIDS, v.UName, v.StartDate, v.EndDate)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	sort.Slice(recheckViews, func(i, j int) bool {
		return recheckViews[i].Date > recheckViews[j].Date
	})
	c.JSON(recheckViews, nil)
}

func auditCargoCsv(c *bm.Context) {
	v := new(struct {
		StartDate int64  `form:"stime"`
		EndDate   int64  `form:"etime"`
		Uname     string `form:"uname"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	content, err := svc.CsvAuditCargo(c, v.StartDate, v.EndDate, v.Uname)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s~%s-%s", time.Unix(v.StartDate, 0).Format("2006/01/02_15"), time.Unix(v.EndDate, 0).Format("2006/01/02_15"), v.Uname),
	})
}

func auditorCargo(c *bm.Context) {
	v := new(struct {
		StartDate int64  `form:"stime"`
		EndDate   int64  `form:"etime"`
		Uname     string `form:"uname"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	wrappers, _, err := svc.AuditorCargoList(c, v.StartDate, v.EndDate, v.Uname)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	sort.Slice(wrappers, func(i, j int) bool {
		return wrappers[i].Date > wrappers[j].Date
	})
	c.JSON(wrappers, nil)
}

func tagRecheck(c *bm.Context) {
	v := new(struct {
		Uname     string `form:"uname"`
		StartDate int64  `form:"stime"`
		EndDate   int64  `form:"etime"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tagViews, err := svc.TagRecheck(c, v.StartDate, v.EndDate, v.Uname)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	sort.Slice(tagViews, func(i, j int) bool {
		return tagViews[i].Date > tagViews[j].Date
	})
	c.JSON(tagViews, nil)
}

func recheck123(c *bm.Context) {
	v := new(struct {
		StartDate int64  `form:"stime"`
		EndDate   int64  `form:"etime"`
		TypeIDStr string `form:"type_id"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var typeIDS []int64
	if typeIDS, err = xstr.SplitInts(v.TypeIDStr); err != nil {
		c.JSON(nil, err)
		return
	}
	recheckViews, err := svc.Recheck123(c, v.StartDate, v.EndDate, typeIDS)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	sort.Slice(recheckViews, func(i, j int) bool {
		return recheckViews[i].Date > recheckViews[j].Date
	})
	c.JSON(recheckViews, nil)
}

func randomVideo(c *bm.Context) {
	v := new(struct {
		StartDate int64  `form:"stime"`
		EndDate   int64  `form:"etime"`
		TypeIDStr string `form:"type_id"`
		Uname     string `form:"uname"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var typeIDS []int64
	if typeIDS, err = xstr.SplitInts(v.TypeIDStr); err != nil {
		c.JSON(nil, err)
		return
	}
	statViewExts, _, err := svc.RandomVideo(c, v.StartDate, v.EndDate, typeIDS, v.Uname)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	sort.Slice(statViewExts, func(i, j int) bool {
		return statViewExts[i].Date > statViewExts[j].Date
	})
	c.JSON(statViewExts, nil)
}

func fixedVideo(c *bm.Context) {
	v := new(struct {
		StartDate int64  `form:"stime"`
		EndDate   int64  `form:"etime"`
		TypeIDStr string `form:"type_id"`
		Uname     string `form:"uname"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var typeIDS []int64
	if typeIDS, err = xstr.SplitInts(v.TypeIDStr); err != nil {
		c.JSON(nil, err)
		return
	}
	statViewExts, _, err := svc.FixedVideo(c, v.StartDate, v.EndDate, typeIDS, v.Uname)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	sort.Slice(statViewExts, func(i, j int) bool {
		return statViewExts[i].Date > statViewExts[j].Date
	})
	c.JSON(statViewExts, nil)
}

func csvFixedVideo(c *bm.Context) {

	v := new(struct {
		StartDate int64  `form:"stime"`
		EndDate   int64  `form:"etime"`
		Uname     string `form:"uname"`
		TypeIDStr string `form:"type_id"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var typeIDS []int64
	if typeIDS, err = xstr.SplitInts(v.TypeIDStr); err != nil {
		c.JSON(nil, err)
		return
	}
	content, err := svc.CsvFixedVideoAudit(c, v.StartDate, v.EndDate, v.Uname, typeIDS)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s_%s~%s-%s", "(定时发布)视频审核操作数据", time.Unix(v.StartDate, 0).Format("2006-01-02"), time.Unix(v.EndDate, 0).Format("2006-01-02"), v.Uname),
	})
}

func csvRandomVideo(c *bm.Context) {
	v := new(struct {
		StartDate int64  `form:"stime"`
		EndDate   int64  `form:"etime"`
		Uname     string `form:"uname"`
		TypeIDStr string `form:"type_id"`
	})
	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var typeIDS []int64
	if typeIDS, err = xstr.SplitInts(v.TypeIDStr); err != nil {
		c.JSON(nil, err)
		return
	}
	content, err := svc.CsvRandomVideoAudit(c, v.StartDate, v.EndDate, v.Uname, typeIDS)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s_%s~%s-%s", "(非定时)视频审核操作数据", time.Unix(v.StartDate, 0).Format("2006-01-02"), time.Unix(v.EndDate, 0).Format("2006-01-02"), v.Uname),
	})
}

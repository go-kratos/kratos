package http

import (
	"fmt"
	"net/http"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func upIncomeList(c *bm.Context) {
	v := new(struct {
		MIDs      []int64 `form:"mids,split"`
		Type      int     `form:"type"`
		GroupType int     `form:"group_type" default:"1"`
		FromTime  int64   `form:"from_time" validate:"required,min=1"`
		ToTime    int64   `form:"to_time" validate:"required,min=1"`
		MinIncome int64   `form:"min_income"`
		MaxIncome int64   `form:"max_income"`
		From      int     `form:"from" validate:"min=0" default:"0"`
		Limit     int     `form:"limit" validate:"min=1" default:"20"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	data, total, err := incomeSvr.UpIncomeList(c, v.MIDs, v.Type, v.GroupType, v.FromTime, v.ToTime, v.MinIncome, v.MaxIncome, v.From, v.Limit)
	if err != nil {
		log.Error("growup incomeSvr.UpIncomeList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func upIncomeListExport(c *bm.Context) {
	v := new(struct {
		MIDs      []int64 `form:"mids,split"`
		Type      int     `form:"type"`
		GroupType int     `form:"group_type" default:"1"`
		FromTime  int64   `form:"from_time" validate:"min=1,required"`
		ToTime    int64   `form:"to_time" validate:"min=1,required"`
		MinIncome int64   `form:"min_income"`
		MaxIncome int64   `form:"max_income"`
		From      int     `form:"from" validate:"min=0" default:"0"`
		Limit     int     `form:"limit" validate:"min=1" default:"20"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	content, err := incomeSvr.UpIncomeListExport(c, v.MIDs, v.Type, v.GroupType, v.FromTime, v.ToTime, v.MinIncome, v.MaxIncome, v.From, v.Limit)
	if err != nil {
		log.Error("growup incomeSvr.UpIncomeListExport error(%v)", err)
		c.JSON(nil, err)
		return
	}

	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "up_income"),
	})
}

func upIncomeStatis(c *bm.Context) {
	v := new(struct {
		MIDs      []int64 `form:"mids,split"`
		Type      int     `form:"type"`
		GroupType int     `form:"group_type" default:"1"`
		FromTime  int64   `form:"from_time" validate:"required,min=1"`
		ToTime    int64   `form:"to_time" validate:"required,min=1"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	data, err := incomeSvr.UpIncomeStatis(c, v.MIDs, v.Type, v.GroupType, v.FromTime, v.ToTime)
	if err != nil {
		log.Error("growup incomeSvr.UpIncomeStatis error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":   err,
			"status": "fail",
		}))
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"status":  "success",
	}))
}

func archiveStatis(c *bm.Context) {
	v := new(struct {
		CategoryID []int64 `form:"category_id,split"`
		Type       int     `form:"type"`
		GroupType  int     `form:"group_type" default:"1"`
		FromTime   int64   `form:"from_time" validate:"required,min=1"`
		ToTime     int64   `form:"to_time" validate:"required,min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := incomeSvr.ArchiveStatis(c, v.CategoryID, v.Type, v.GroupType, v.FromTime, v.ToTime)
	if err != nil {
		log.Error("growup incomeSvr.ArchiveStatis error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":   err,
			"status": "fail",
		}))
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"status":  "success",
	}))
}

func archiveSection(c *bm.Context) {
	v := new(struct {
		CategoryID []int64 `form:"category_id,split"`
		Type       int     `form:"type"`
		GroupType  int     `form:"group_type" default:"1"`
		FromTime   int64   `form:"from_time" validate:"required,min=1"`
		ToTime     int64   `form:"to_time" validate:"required,min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := incomeSvr.ArchiveSection(c, v.CategoryID, v.Type, v.GroupType, v.FromTime, v.ToTime)
	if err != nil {
		log.Error("growup incomeSvr.ArchiveSection error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":   err,
			"status": "fail",
		}))
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"status":  "success",
	}))
}

func archiveDetail(c *bm.Context) {
	v := new(struct {
		MID       int64 `form:"mid" validate:"required"`
		Type      int   `form:"type"`
		GroupType int   `form:"group_type" default:"1"`
		FromTime  int64 `form:"from_time" validate:"required,min=1"`
		ToTime    int64 `form:"to_time" validate:"required,min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	data, err := incomeSvr.ArchiveDetail(c, v.MID, v.Type, v.GroupType, v.FromTime, v.ToTime)
	if err != nil {
		log.Error("growup incomeSvr.ArchiveDetail error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":   err,
			"status": "fail",
		}))
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"status":  "success",
	}))
}

func archiveTop(c *bm.Context) {
	v := new(struct {
		AIDs      []int64 `form:"aids,split"`
		Type      int     `form:"type"`
		GroupType int     `form:"group_type" default:"1"`
		FromTime  int64   `form:"from_time" validate:"required,min=1"`
		ToTime    int64   `form:"to_time" validate:"required,min=1"`
		From      int     `form:"from" validate:"min=0" default:"0"`
		Limit     int     `form:"limit" validate:"min=1" default:"20"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	data, total, err := incomeSvr.ArchiveTop(c, v.AIDs, v.Type, v.GroupType, v.FromTime, v.ToTime, v.From, v.Limit)
	if err != nil {
		log.Error("growup incomeSvr.ArchiveTop error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":   err,
			"status": "fail",
		}))
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"status":  "success",
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func bgmDetail(c *bm.Context) {
	v := new(struct {
		SID      int64 `form:"sid"`
		FromTime int64 `form:"from_time" validate:"required,min=1"`
		ToTime   int64 `form:"to_time" validate:"required,min=1"`
		From     int   `form:"from" validate:"min=0" default:"0"`
		Limit    int   `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	data, total, err := incomeSvr.BgmDetail(c, v.SID, v.FromTime, v.ToTime, v.From, v.Limit)
	if err != nil {
		log.Error("growup incomeSvr.BgmDetail error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":   err,
			"status": "fail",
		}))
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"status":  "success",
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func archiveBreach(c *bm.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	v := new(struct {
		Type   int     `form:"type"`
		AIDs   []int64 `form:"aids,split" validate:"required"`
		MID    int64   `form:"mid" validate:"required"`
		Reason string  `form:"reason" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}

	err = incomeSvr.ArchiveBreach(c, v.Type, v.AIDs, v.MID, v.Reason, username)
	if err != nil {
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}))
	} else {
		c.JSON(nil, nil)
	}
}

func archiveBlack(c *bm.Context) {
	v := new(struct {
		Type int     `form:"type"`
		AIDs []int64 `form:"aids,split" validate:"required"`
		MID  int64   `form:"mid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	err := incomeSvr.ArchiveBlack(c, v.Type, v.AIDs, v.MID)
	if err != nil {
		log.Error("growup incomeSvr.ArchiveBlack error(%v)", err)
	}
	c.JSON(nil, err)
}

func breachList(c *bm.Context) {
	v := new(struct {
		MIDs     []int64 `form:"mids,split"`
		AIDs     []int64 `form:"aids,split"`
		Type     int     `form:"type"`
		FromTime int64   `form:"from_time" validate:"required,min=1"`
		ToTime   int64   `form:"to_time" validate:"required,min=1"`
		Reason   string  `form:"reason"`
		From     int     `form:"from" validate:"min=0" default:"0"`
		Limit    int     `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	data, total, err := incomeSvr.BreachList(c, v.MIDs, v.AIDs, v.Type, v.FromTime, v.ToTime, v.Reason, v.From, v.Limit)
	if err != nil {
		log.Error("growup incomeSvr.BreachList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func breachStatis(c *bm.Context) {
	v := new(struct {
		MIDs      []int64 `form:"mids,split"`
		AIDs      []int64 `form:"aids,split"`
		Type      int     `form:"type"`
		GroupType int     `form:"group_type"`
		FromTime  int64   `form:"from_time" validate:"required,min=1"`
		ToTime    int64   `form:"to_time" validate:"required,min=1"`
		Reason    string  `form:"reason"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	data, err := incomeSvr.BreachStatis(c, v.MIDs, v.AIDs, v.Type, v.GroupType, v.FromTime, v.ToTime, v.Reason)
	if err != nil {
		log.Error("growup incomeSvr.BreachStatis error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"status":  "success",
	}))
}

func exportBreach(c *bm.Context) {
	v := new(struct {
		MIDs     []int64 `form:"mids,split"`
		AIDs     []int64 `form:"aids,split"`
		Type     int     `form:"type"`
		FromTime int64   `form:"from_time" validate:"required,min=1"`
		ToTime   int64   `form:"to_time" validate:"required,min=1"`
		Reason   string  `form:"reason"`
		From     int     `form:"from" validate:"min=0" default:"0"`
		Limit    int     `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	content, err := incomeSvr.ExportBreach(c, v.MIDs, v.AIDs, v.Type, v.FromTime, v.ToTime, v.Reason, v.From, v.Limit)
	if err != nil {
		log.Error("growup svr.ExportBreach error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "breach_record"),
	})
}

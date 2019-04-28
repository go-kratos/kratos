package http

import (
	"net/http"

	"go-common/app/admin/main/growup/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func addTagInfo(c *bm.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup checkCookie error(%v)", err)
		return
	}
	v := new(model.TagInfo)
	if err = c.Bind(v); err != nil {
		return
	}

	err = svr.AddTagInfo(c, v, username)
	if err != nil {
		log.Error("growup svr.AddTagInfo error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}))
	} else {
		c.JSON(nil, nil)
	}
}

func updateTagInfo(c *bm.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup checkCookie error(%v)", err)
		return
	}
	v := new(model.TagInfo)
	if err = c.Bind(v); err != nil {
		return
	}
	err = svr.UpdateTagInfo(c, v, username)
	if err != nil {
		log.Error("growup svr.UpdateTagInfo error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}))
	} else {
		c.JSON(nil, nil)
	}
}

func modTagState(c *bm.Context) {
	v := new(struct {
		IsDelete int `form:"is_delete"`
		TagID    int `form:"tag_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.ModeTagState(c, v.TagID, v.IsDelete)
	if err != nil {
		log.Error("growup svr.ModTagState error(%v)", err)
	}
	c.JSON(nil, err)
}

func addTagUps(c *bm.Context) {
	v := new(struct {
		TagID    int     `form:"tag_id" validate:"required"`
		MIDs     []int64 `form:"mids,split"`
		IsCommon int     `form:"is_common"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.AddTagUps(c, v.TagID, v.MIDs, v.IsCommon)
	if err != nil {
		log.Error("growup svr.AddTagUPs error(%v)", err)
	}
	c.JSON(nil, err)
}

func releaseUp(c *bm.Context) {
	v := new(struct {
		TagID int   `form:"tag_id" validate:"required"`
		MID   int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.ReleaseUp(c, v.TagID, v.MID)
	if err != nil {
		log.Error("growup svr.ReleaseUps error(%v)", err)
	}
	c.JSON(nil, err)
}

func listTagInfo(c *bm.Context) {
	v := new(struct {
		StartTime  int64   `form:"start_time"`
		EndTime    int64   `form:"end_time"`
		Categories []int64 `form:"categories,split"`
		Business   []int64 `form:"business,split"`
		Tag        string  `form:"tag"`
		Effect     int     `form:"effect"`
		From       int     `form:"from" default:"0" validate:"min=0"`
		Limit      int     `form:"limit" default:"20" validate:"min=1"`
		Sort       string  `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	total, tagInfos, err := svr.QueryTagInfo(c, v.StartTime, v.EndTime, v.Categories, v.Business, v.Tag, v.Effect, v.From, v.Limit, v.Sort)
	if err != nil {
		log.Error("growup svr.queryTagInfo error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    tagInfos,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func listUps(c *bm.Context) {
	v := new(struct {
		TagID int    `form:"tag_id" validate:"required"`
		MID   int64  `form:"mid"`
		From  int    `form:"from" default:"0" validate:"min=0"`
		Limit int    `form:"limit" default:"20" validate:"min=1"`
		Sort  string `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	total, upIncomes, err := svr.ListUps(c, v.TagID, v.MID, v.From, v.Limit)
	if err != nil {
		log.Error("svr.ListUps error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    upIncomes,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func listAvs(c *bm.Context) {
	v := new(struct {
		TagID int    `form:"tag_id" validate:"required"`
		AvID  int64  `form:"av_id"`
		From  int    `form:"from" default:"0" validate:"min=0"`
		Limit int    `form:"limit" default:"20" validate:"min=1"`
		Sort  string `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	total, avIncomes, err := svr.ListAvs(c, v.TagID, v.From, v.Limit, v.AvID)
	if err != nil {
		log.Error("svr.ListAvs error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    avIncomes,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func tagDetails(c *bm.Context) {
	v := new(struct {
		TagID int    `form:"tag_id" validate:"required"`
		From  int    `form:"from" default:"0" validate:"min=0"`
		Limit int    `form:"limit" default:"20" validate:"min=1"`
		Sort  string `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	total, details, err := svr.TagDetails(c, v.TagID, v.From, v.Limit)
	if err != nil {
		log.Error("growup svr.queryDetails error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    details,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func updateActivity(c *bm.Context) {
	v := new(struct {
		TagID      int64 `form:"tag_id" validate:"required"`
		ActivityID int64 `form:"activity_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.UpdateTagActivity(c, v.TagID, v.ActivityID)
	if err != nil {
		log.Error("growup svr.AddActivity error(%v)", err)
	}
	c.JSON(nil, err)
}

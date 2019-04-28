package http

import (
	"net/http"

	"go-common/app/admin/main/growup/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"go-common/library/time"
)

func activityAdd(c *bm.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		log.Error("checkCookie error(%v)", err)
		return
	}
	v := new(model.CActivity)
	if err = c.Bind(v); err != nil {
		return
	}
	if err = svr.AddActivity(c, v, username); err != nil {
		log.Error("svr.AddActivity error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}))
	} else {
		c.JSON(nil, nil)
	}
}

func activityList(c *bm.Context) {
	v := new(struct {
		Name  string `form:"name"`
		From  int    `form:"from" validate:"min=0" default:"0"`
		Limit int    `form:"limit" validate:"min=1" default:"20"`
		Sort  string `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, total, err := svr.ListActivity(c, v.Name, v.From, v.Limit, v.Sort)
	if err != nil {
		log.Error("svr.ListActivity error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code": 0,
		"data": data,
		"paging": map[string]int{
			"size":  v.Limit,
			"total": total,
		},
	}))
}

func activityUpdate(c *bm.Context) {
	var err error
	v := new(model.CActivity)
	if err = c.Bind(v); err != nil {
		return
	}
	if err = svr.UpdateActivity(c, v); err != nil {
		log.Error("svr.UpdateActivity error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}))
	} else {
		c.JSON(nil, nil)
	}
}

func activitySignUp(c *bm.Context) {
	v := new(struct {
		ID    int64 `form:"id"`
		From  int   `form:"from" validate:"min=0" default:"0"`
		Limit int   `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, total, err := svr.ListActivitySignUp(c, v.ID, v.From, v.Limit)
	if err != nil {
		log.Error("svr.ListActivitySignUp error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code": 0,
		"data": data,
		"paging": map[string]int{
			"size":  v.Limit,
			"total": total,
		},
	}))
}

func activityWinners(c *bm.Context) {
	v := new(struct {
		ID    int64 `form:"id"`
		MID   int64 `form:"mid"`
		From  int   `form:"from" validate:"min=0" default:"0"`
		Limit int   `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, total, err := svr.ListActivityWinners(c, v.ID, v.MID, v.From, v.Limit)
	if err != nil {
		log.Error("svr.ListActivityWinners error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code": 0,
		"data": data,
		"paging": map[string]int{
			"size":  v.Limit,
			"total": total,
		},
	}))
}

func activityAward(c *bm.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		log.Error("checkCookie error(%v)", err)
		return
	}
	v := new(struct {
		ID            int64     `form:"id" validate:"required"`
		Name          string    `form:"name" validate:"required"`
		Date          time.Time `form:"date" validate:"required"`
		StatisticsEnd time.Time `form:"statistics_end" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if err = svr.ActivityAward(c, v.ID, v.Name, v.Date, v.StatisticsEnd, username); err != nil {
		log.Error("svr.ActivityAward error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		}))
	} else {
		c.JSON(nil, nil)
	}
}

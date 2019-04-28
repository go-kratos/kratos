package http

import (
	"regexp"
	"strconv"

	"go-common/app/admin/main/laser/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	emailRegex = "^([a-z0-9A-Z]+[-|\\.]?)+[a-z0-9A-Z]@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)?\\.)+[a-zA-Z]{2,}$"
)

func addTask(c *bm.Context) {
	req := c.Request
	v := new(struct {
		MID          int64  `form:"mid" validate:"required"`
		LogDate      int64  `form:"log_date" validate:"required"`
		ContactEmail string `form:"contact_email" validate:"required"`
		Platform     int    `form:"platform" validate:"required"`
		SourceType   int    `form:"source_type" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	if v.LogDate <= 0 || v.MID <= 0 || !checkEmail(v.ContactEmail) || v.Platform <= 0 || v.SourceType <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	userCookie, err := req.Cookie("username")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	username := userCookie.Value

	uidCookie, err := req.Cookie("uid")
	var adminID int64
	if err != nil {
		adminID = 0
	} else {
		adminID, err = strconv.ParseInt(uidCookie.Value, 10, 64)
		if err != nil {
			c.JSON(nil, err)
			return
		}

	}
	err = svc.AddTask(c, v.MID, username, adminID, v.LogDate, v.ContactEmail, v.Platform, v.SourceType)
	if err != nil {
		log.Error("svc.AddTask() error(%v)", err)
	}
	c.JSON(nil, err)
}

func checkEmail(emailAddr string) (match bool) {
	if emailAddr == "" {
		return false
	}
	match, err := regexp.MatchString(emailRegex, emailAddr)
	if err != nil {
		return false
	}
	return match
}

func deleteTask(c *bm.Context) {
	req := c.Request
	v := new(struct {
		TaskID int64 `form:"task_id"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	uidCookie, err := req.Cookie("uid")
	var adminID int64
	if err != nil {
		adminID = 0
	} else {
		adminID, err = strconv.ParseInt(uidCookie.Value, 10, 64)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}

	}
	userCookie, err := req.Cookie("username")
	if err != nil {
		c.JSON(nil, ecode.Unauthorized)
		return
	}
	username := userCookie.Value

	if err = svc.DeleteTask(c, v.TaskID, username, adminID); err != nil {
		log.Error("svc.DeleteTask() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func queryTask(c *bm.Context) {
	v := new(struct {
		Mid          int64  `form:"mid"`
		LogDateStart int64  `form:"log_date_start"`
		LogDateEnd   int64  `form:"log_date_end"`
		SourceType   int    `form:"source_type"`
		Platform     int    `form:"platform"`
		State        int    `form:"state"`
		SortBy       string `form:"sort"`
		PageNo       int    `form:"page_no"`
		PageSize     int    `form:"page_size"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	tasks, count, err := svc.QueryTask(c, v.Mid, v.LogDateStart, v.LogDateEnd, v.SourceType, v.Platform, v.State, v.SortBy, v.PageNo, v.PageSize)
	if err != nil {
		log.Error("svc.QueryTask() error(%v)", err)
		c.JSON(nil, err)
		return
	}
	pager := &model.TaskPager{
		PageSize: v.PageSize,
		PageNo:   v.PageNo,
		Total:    count,
		Items:    tasks,
	}
	c.JSON(pager, nil)
}

func updateTask(c *bm.Context) {
	req := c.Request
	v := new(struct {
		TaskID       int64  `form:"task_id" validate:"required"`
		MID          int64  `form:"mid" validate:"required"`
		LogDate      int64  `form:"log_date" validate:"required"`
		ContactEmail string `form:"contact_email" validate:"required"`
		SourceType   int    `form:"source_type" validate:"required"`
		Platform     int    `form:"platform" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	if v.ContactEmail != "" && !checkEmail(v.ContactEmail) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	uidCookie, err := req.Cookie("uid")
	var adminID int64
	if err != nil {
		adminID = 0
	} else {
		adminID, err = strconv.ParseInt(uidCookie.Value, 10, 64)
		if err != nil {
			c.JSON(nil, err)
			return
		}
	}
	userCookie, err := req.Cookie("username")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	username := userCookie.Value
	err = svc.UpdateTask(c, username, adminID, v.TaskID, v.MID, v.LogDate, v.ContactEmail, v.SourceType, v.Platform)
	if err != nil {
		log.Error("svc.UpdateTask() error(%v)", err)
	}
	c.JSON(nil, err)
}

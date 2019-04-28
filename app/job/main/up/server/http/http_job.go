package http

import (
	"context"
	"go-common/app/job/main/up/model"
	"go-common/app/job/main/up/model/upcrmmodel"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"strings"
	"time"
)

func runJob(c *blademaster.Context) {
	var err error
	var res interface{}
	switch {
	default:
		var arg struct {
			Date string `form:"date"`
			Job  string `form:"job"`
		}
		var err = c.Bind(&arg)
		if err != nil {
			break
		}
		var date time.Time
		if arg.Date == "" {
			date = time.Now()
		} else {
			date, err = time.Parse(upcrmmodel.TimeFmtDate, arg.Date)
			if err != nil {
				log.Error("parse date err")
				break
			}
		}

		switch strings.ToLower(arg.Job) {
		case "task":
			svc.CheckTaskFinish(date)
		case "due":
			svc.CheckDateDueJob(date)
		case "state":
			svc.CheckStateJob(date)
		case "tid":
			svc.UpdateUpTidJob(date)
		case "":
			svc.UpdateUpTidJob(date)
			svc.CheckStateJob(date)
			svc.CheckDateDueJob(date)
			svc.CheckTaskFinish(date)
		}

	}

	if err != nil {
		c.JSON(nil, err)
	} else {
		c.JSON(res, err)
	}
}

func warmUp(c *blademaster.Context) {
	var err error
	var res interface{}
	switch {
	default:
		var arg = &model.WarmUpReq{}
		var err = c.Bind(arg)
		if err != nil {
			break
		}

		go func() {
			res, err = svc.WarmUp(context.Background(), arg)
		}()
	}

	if err != nil {
		c.JSON(nil, err)
	} else {
		c.JSON(res, err)
	}
}

func warmUpMid(c *blademaster.Context) {
	var err error
	var res interface{}
	switch {
	default:
		var arg = &model.WarmUpReq{}
		var err = c.Bind(arg)
		if err != nil {
			break
		}

		go func() {
			res, err = svc.WarmUpMid(context.Background(), arg)
		}()
	}

	if err != nil {
		c.JSON(nil, err)
	} else {
		c.JSON(res, err)
	}
}

func addStaff(c *blademaster.Context) {
	var err error
	var res interface{}
	switch {
	default:
		var arg = &model.AddStaffReq{}
		var err = c.Bind(arg)
		if err != nil {
			break
		}
		res, _ = svc.AddStaff(c, arg)
	}

	if err != nil {
		c.JSON(nil, err)
	} else {
		c.JSON(res, err)
	}
}

func deleteStaff(c *blademaster.Context) {
	var err error
	var res interface{}
	switch {
	default:
		var arg = &model.AddStaffReq{}
		var err = c.Bind(arg)
		if err != nil {
			break
		}
		res, _ = svc.DeleteStaff(c, arg)
	}

	if err != nil {
		c.JSON(nil, err)
	} else {
		c.JSON(res, err)
	}
}

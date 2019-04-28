package http

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func taskList(c *bm.Context) {
	v := new(model.TaskListArg)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.TaskList(c, v))
}

func addTask(c *bm.Context) {
	var (
		uname, _ = c.Get("username")
		v        = new(model.AddTaskArg)
		err      error
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Regex == "" && v.Mids == "" && v.IPs == "" && v.Cids == "" && v.KeyWords == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, err = xstr.SplitInts(v.Mids); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, err = xstr.SplitInts(v.Cids); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.IPs != "" {
		ips := strings.Split(v.IPs, ",")
		for _, ip := range ips {
			tmp := net.ParseIP(ip)
			if tmp == nil {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
	}
	v.Creator = uname.(string)
	c.JSON(nil, dmSvc.AddTask(c, v))
}

func editTaskState(c *bm.Context) {
	v := new(model.EditTasksStateArg)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, dmSvc.EditTaskState(c, v))
}

func reviewTask(c *bm.Context) {
	var (
		reviewer, _ = c.Get("username")
		v           = new(model.ReviewTaskArg)
	)
	if err := c.Bind(v); err != nil {
		return
	}
	v.Reviewer = reviewer.(string)
	c.JSON(nil, dmSvc.ReviewTask(c, v))
}
func taskView(c *bm.Context) {
	v := new(model.TaskViewArg)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(dmSvc.TaskView(c, v))
}

func taskCsv(c *bm.Context) {
	var (
		bs          []byte
		err         error
		contentType = "text/csv"
	)
	v := new(model.TaskCsvArg)
	if err = c.Bind(v); err != nil {
		return
	}
	if bs, err = dmSvc.TaskCsv(c, v.ID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%v.csv", v.ID))
	c.Bytes(http.StatusOK, contentType, bs)
}

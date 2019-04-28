package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strings"
)

func queryTask(c *bm.Context) {
	params := c.Request.Form
	var (
		err        error
		mid        int64
		platform   int
		sourceType int
	)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	platformStr := params.Get("platform")
	if platformStr == "android" {
		platform = 2
	} else {
		if platform, err = strconv.Atoi(platformStr); err != nil || platform <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	sourceTypeStr := params.Get("source_type")
	if sourceType, err = strconv.Atoi(sourceTypeStr); err != nil || sourceType <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	logDate, err := svc.QueryUndoneTaskLogdate(c, mid, platform, sourceType)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	value := make(map[string]interface{})
	value["log_date"] = logDate
	c.JSON(value, nil)
}

func updateTask(c *bm.Context) {
	params := c.Request.Form
	var (
		err       error
		mid       int64
		build     string
		platform  int
		taskState int
		reason    string
	)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	build = params.Get("build")
	platformStr := params.Get("platform")
	taskStateStr := params.Get("task_state")
	reason = params.Get("reason")
	if strings.EqualFold("android", platformStr) {
		platform = 2
	} else if platform, err = strconv.Atoi(platformStr); err != nil || platform <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if taskState, err = strconv.Atoi(taskStateStr); err != nil || taskState <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.UpdateTaskState(c, mid, build, platform, taskState, reason))
}

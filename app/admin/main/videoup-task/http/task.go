package http

import (
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// 任务管理
func delay(c *bm.Context) {
	v := new(struct {
		Taskid int64  `form:"task_id" validate:"required"`
		Reason string `form:"reason"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	uid, _ := getUIDName(c)
	c.JSON(nil, srv.Delay(c, v.Taskid, uid, v.Reason))
}

func taskfree(c *bm.Context) {
	var uid, rows int64
	uid, _ = getUIDName(c)
	rows = srv.Free(c, uid)
	log.Info("释放任务(uid=%d,rows=%d)", uid, rows)
	c.JSON(nil, nil)
}

func next(c *bm.Context) {
	uidS := c.Request.Form.Get("uid")
	uid, err := strconv.Atoi(uidS)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	task, err := srv.Next(c, int64(uid))
	if err != nil {
		log.Error("srv.Next(uid=%d) error(%v)", uid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(task, nil)
}

func listTask(c *bm.Context) {
	v := new(struct {
		UID      int64 `form:"uid" default:"0"`
		IsLeader int8  `form:"isleader"  default:"0"`
		Lt       int8  `form:"listtype"  default:"0"`
		Pn       int   `form:"page"  default:"1"`
		Ps       int   `form:"pagesize"  default:"20"`
	})

	err := c.Bind(v)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	tasks, err := srv.List(c, v.UID, v.Pn, v.Ps, v.Lt, v.IsLeader)
	if err != nil {
		log.Error("srv.List(uid=%d,page=%d, pagesize=%d, listtype=%d, isleader=%d) error(%v)",
			v.UID, v.Pn, v.Ps, v.Lt, v.IsLeader, err)
		c.JSON(nil, err)
		return
	}

	c.JSON(tasks, nil)
}

// info 返回任务信息,复审信息
func info(c *bm.Context) {
	var tid int
	if tid, _ = strconv.Atoi(c.Request.Form.Get("taskid")); tid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	task, err := srv.Info(c, int64(tid))
	if err != nil {
		c.JSON(nil, err)
		return
	}

	form, err := srv.ReviewForm(c, int64(tid))
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSONMap(map[string]interface{}{
		"task":   task,
		"review": form,
	}, nil)
}

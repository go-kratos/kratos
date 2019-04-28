package http

import (
	"net/http"
	"strconv"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func taskTooks(c *bm.Context) {
	req := c.Request
	params := req.Form
	stimeStr := params.Get("stime")
	etimeStr := params.Get("etime")
	if stimeStr == "" {
		stimeStr = time.Now().Format("2006-01-02") + " 00:00:00"
	}
	if etimeStr == "" {
		etimeStr = time.Now().Format("2006-01-02 15:04:05")
	}
	local, _ := time.LoadLocation("Local")
	stime, err := time.ParseInLocation("2006-01-02 15:04:05", stimeStr, local)
	if stime.Unix() < 1 || err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	etime, err := time.ParseInLocation("2006-01-02 15:04:05", etimeStr, local)
	if etime.Unix() < 1 || err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tooks, err := vdaSvc.TaskTooksByHalfHour(c, stime, etime)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(tooks, nil)
}

func next(c *bm.Context) {
	uidS := c.Request.Form.Get("uid")
	uid, err := strconv.Atoi(uidS)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	task, err := vdaSvc.Next(c, int64(uid))
	if err != nil {
		log.Error("vdaSvc.Next(uid=%d) error(%v)", uid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(task, nil)
}

func list(c *bm.Context) {
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

	tasks, err := vdaSvc.List(c, v.UID, v.Pn, v.Ps, v.Lt, v.IsLeader)
	if err != nil {
		log.Error("vdaSvc.List(uid=%d,page=%d, pagesize=%d, listtype=%d, isleader=%d) error(%v)",
			v.UID, v.Pn, v.Ps, v.Lt, v.IsLeader, err)
		c.JSON(nil, err)
		return
	}

	c.JSON(tasks, nil)
}

func info(c *bm.Context) {
	tidS := c.Request.Form.Get("taskid")
	tid, err := strconv.Atoi(tidS)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	task, err := vdaSvc.Info(c, int64(tid))
	if err != nil {
		log.Error("vdaSvc.Info(taskid=%d) error(%v)", tid, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(task, nil)
}

// 权重管理
func addwtconf(c *bm.Context) {
	var err error
	cfg := &archive.WeightConf{}
	if err = c.Bind(cfg); err != nil {
		log.Error("addwtconf error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	uid, uname := getUIDName(c)
	err = vdaSvc.AddWeightConf(c, cfg, uid, uname)
	if err != nil {
		data := map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": err.Error(),
		}
		c.Render(http.StatusOK, render.MapJSON(data))
		return
	}
	c.JSON(nil, nil)
}

func delwtconf(c *bm.Context) {
	var err error
	ids := c.Request.Form.Get("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(vdaSvc.DelWeightConf(c, int64(id)), nil)
}

func listwtconf(c *bm.Context) {
	v := new(archive.Confs)
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(vdaSvc.ListWeightConf(c, v))
}

func maxweight(c *bm.Context) {
	c.JSON(vdaSvc.MaxWeight(c))
}

func listwtlog(c *bm.Context) {
	v := new(struct {
		Taskid int64 `form:"taskid" validate:"required"`
		Pn     int   `form:"page" default:"1"`
		Ps     int   `form:"ps" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	cfg, count, err := vdaSvc.ListWeightLogs(c, v.Taskid, v.Pn)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{})
	data["data"] = cfg
	data["pager"] = map[string]int{
		"current_page": v.Pn,
		"total_items":  int(count),
		"page_size":    20,
	}
	c.JSONMap(data, err)
}

func show(c *bm.Context) {
	c.JSON(vdaSvc.ShowWeightVC(c))
}

func set(c *bm.Context) {
	v := new(archive.WeightVC)
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, vdaSvc.SetWeightVC(c, v))
}

// 登录管理
func on(c *bm.Context) {
	uid, uname := getUIDName(c)
	err := vdaSvc.HandsUp(c, uid, uname)
	if err != nil {
		data := map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": err.Error(),
		}
		c.Render(http.StatusOK, render.MapJSON(data))
		return
	}
	c.JSON(nil, nil)
}

// 踢出
func forceoff(c *bm.Context) {
	uidS := c.Request.Form.Get("uid")
	uid, _ := strconv.ParseInt(uidS, 10, 64)
	if uid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	adminuid, _ := getUIDName(c)
	err := vdaSvc.HandsOff(c, adminuid, uid)
	if err != nil {
		data := map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": err.Error(),
		}
		c.Render(http.StatusOK, render.MapJSON(data))
		return
	}
	c.JSON(nil, nil)
}

func off(c *bm.Context) {
	adminuid, _ := getUIDName(c)
	err := vdaSvc.HandsOff(c, adminuid, 0)
	if err != nil {
		data := map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": err.Error(),
		}
		c.Render(http.StatusOK, render.MapJSON(data))
		return
	}
	c.JSON(nil, nil)
}

func online(c *bm.Context) {
	c.JSON(vdaSvc.Online(c))
}

func inoutlist(c *bm.Context) {
	v := new(struct {
		Unames string `form:"unames" default:""`
		Bt     string `form:"bt" default:""`
		Et     string `form:"et" default:""`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(vdaSvc.InOutList(c, v.Unames, v.Bt, v.Et))
}

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
	c.JSON(nil, vdaSvc.Delay(c, v.Taskid, uid, v.Reason))
}

func taskfree(c *bm.Context) {
	var uid, rows int64
	uid, _ = getUIDName(c)
	rows = vdaSvc.Free(c, uid)
	log.Info("释放任务(uid=%d,rows=%d)", uid, rows)
	c.JSON(nil, nil)
}

func checkgroup() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		uid, _ := getUIDName(ctx)
		role, err := vdaSvc.CheckGroup(ctx, uid)
		if err != nil || role == 0 {
			data := map[string]interface{}{
				"code":    ecode.RequestErr,
				"message": "权限错误",
			}
			ctx.Render(http.StatusOK, render.MapJSON(data))
			ctx.Abort()
			return
		}
	}
}

// 校验任务操作权限
func checkowner() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		tidS := ctx.Request.Form.Get("task_id")
		tid, err := strconv.Atoi(tidS)
		if err != nil {
			ctx.JSON(nil, ecode.RequestErr)
			ctx.Abort()
			return
		}

		uid, _ := getUIDName(ctx)
		if err := vdaSvc.CheckOwner(ctx, int64(tid), uid); err != nil {
			data := map[string]interface{}{
				"code":    ecode.RequestErr,
				"message": err.Error(),
			}
			ctx.Render(http.StatusOK, render.MapJSON(data))
			ctx.Abort()
			return
		}
	}
}

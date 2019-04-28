package http

import (
	"net/http"
	"strconv"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

// 权重管理
func addwtconf(c *bm.Context) {
	var err error
	cfg := &model.WeightConf{}
	if err = c.Bind(cfg); err != nil {
		log.Error("addwtconf error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	uid, uname := getUIDName(c)
	err = srv.AddWeightConf(c, cfg, uid, uname)
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
	c.JSON(srv.DelWeightConf(c, int64(id)), nil)
}

func listwtconf(c *bm.Context) {
	v := new(model.Confs)
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(srv.ListWeightConf(c, v))
}

func maxweight(c *bm.Context) {
	c.JSON(srv.MaxWeight(c))
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

	cfg, count, err := srv.ListWeightLogs(c, v.Taskid, v.Pn)
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
	c.JSON(srv.ShowWeightVC(c))
}

func set(c *bm.Context) {
	v := new(model.WeightVC)
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.SetWeightVC(c, v))
}

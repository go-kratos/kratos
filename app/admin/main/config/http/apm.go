package http

import (
	"go-common/app/admin/main/config/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func apmCopy(c *bm.Context) {
	res := map[string]interface{}{}
	user := user(c)
	v := new(model.ApmCopyReq)
	err := c.Bind(v)
	if err != nil {
		return
	}
	if _, err = svr.AuthApp(c, user, c.Request.Header.Get("Cookie"), v.TreeID); err != nil {
		res["message"] = "服务树权限不足"
		c.JSONMap(res, err)
		return
	}
	cnt := 0
	if err = svr.DB.Model(&model.App{}).Where("tree_id=?", v.TreeID).Count(&cnt).Error; err != nil {
		log.Error("svr.ApmCopy count error(%v)", err)
		res["message"] = "查询该服务失败"
		c.JSONMap(res, err)
		return
	}
	if cnt <= 0 {
		log.Error("svr.ApmCopy count (%v)", cnt)
		res["message"] = "未找到该服务"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, svr.Apm(v.TreeID, v.Name, v.ApmName, user))
}

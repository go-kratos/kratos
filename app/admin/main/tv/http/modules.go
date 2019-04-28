package http

import (
	"strconv"
	"time"

	"go-common/app/admin/main/tv/model"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func modulesAdd(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := new(model.Modules)
	if err = c.Bind(param); err != nil {
		return
	}
	if err = modValid(param.ModCore); err != nil {
		c.JSON(nil, err)
		return
	}
	if err = tvSrv.ModulesAdd(param); err != nil {
		res["message"] = "添加模块列表失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	pageID := strconv.Itoa(int(param.ID))
	if err := tvSrv.SetPublish(c, pageID, model.ModulesPublishNo); err != nil {
		res["message"] = "设置模块发布状态失败!"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func modValid(core model.ModCore) (err error) {
	var t int
	if t, err = strconv.Atoi(core.Type); err != nil {
		log.Error("atoi %s, Err %v", core.Type, err)
		return
	}
	if t < 2 || t > 7 {
		err = ecode.RequestErr
		return
	}
	if sup := tvSrv.TypeSupport(core.SrcType, atoi(core.Source)); !sup {
		err = ecode.RequestErr
	}
	return
}

func modulesEditPost(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := new(model.ModulesAddParam)
	if err = c.Bind(param); err != nil {
		return
	}
	if err = modValid(param.ModCore); err != nil {
		c.JSON(nil, err)
		return
	}
	v := &model.Modules{
		PageID:   param.PageID,
		Flexible: param.Flexible,
		Icon:     param.Icon,
		Title:    param.Title,
		Capacity: param.Capacity,
		More:     param.More,
		Order:    param.Order,
		Moretype: param.Moretype,
		Morepage: param.Morepage,
		ModCore:  param.ModCore,
	}
	if err = tvSrv.ModulesEditPost(param.ID, v); err != nil {
		res["message"] = "编辑数据失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func modulesList(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
		v   []*model.Modules
		p   model.ModPub
	)
	param := new(struct {
		PageID string `form:"page_id" validate:"required"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	if v, err = tvSrv.ModulesList(param.PageID); err != nil {
		res["message"] = "获取模块列表失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if p, err = tvSrv.GetModPub(c, param.PageID); err != nil {
		if err == memcache.ErrNotFound {
			nowTime := time.Now()
			t := nowTime.Format("2006-01-02 15:04:05")
			p = model.ModPub{
				Time:  t,
				State: 0,
			}
		} else {
			res["message"] = "获取模块发布状态失败!" + err.Error()
			c.JSONMap(res, ecode.RequestErr)
			return
		}
	}
	m := &model.ModulesList{
		Items:    v,
		PubState: p.State,
		PubTime:  p.Time,
	}
	c.JSON(m, nil)
}

func modulesEditGet(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
		v   = &model.Modules{}
	)
	param := new(struct {
		ID uint64 `form:"id" validate:"required"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	if v, err = tvSrv.ModulesEditGet(param.ID); err != nil {
		res["message"] = "获取数据失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(v, nil)
}

func modulesPublish(c *bm.Context) {
	var (
		err    error
		res    = map[string]interface{}{}
		tmp    = &model.Param{}
		tmpRes []*model.RegList
	)
	param := new(struct {
		IDs        []int  `form:"ids,split" validate:"required,min=1,dive,gt=0"`
		DeletedIDs []int  `form:"deleted_ids,split" validate:"required,dive,gt=0"`
		PageID     string `form:"page_id" validate:"required"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	if param.PageID != "0" {
		tmp.PageID = param.PageID
		if tmpRes, err = tvSrv.RegList(c, tmp); err != nil || len(tmpRes) != 1 {
			log.Error("search region is publish count(%d) error(%v)", len(tmpRes), err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if tmpRes[0].Valid == 0 {
			log.Error("该分区region(%d)未上线 valid(%d)", param.PageID, tmpRes[0].Valid)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if err := tvSrv.ModulesPublish(c, param.PageID, model.ModulesPublishYes, param.IDs, param.DeletedIDs); err != nil {
		res["message"] = "模块发布失败!" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func supCat(c *bm.Context) {
	if len(tvSrv.SupCats) > 0 {
		c.JSON(tvSrv.SupCats, nil)
		return
	}
	log.Error("SupCats empty")
	c.JSON(nil, ecode.ServiceUnavailable)
}

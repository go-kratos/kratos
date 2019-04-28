package http

import (
	"encoding/json"

	"go-common/app/admin/main/tag/conf"
	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func hotList(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamHotList)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Type != model.HotRegionTag {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	top, view, filter, normal, err := svc.RegionRankList(c, param.Prid, param.Rid, param.Type)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 4)
	data["top"] = top
	data["view"] = view
	data["filter"] = filter
	data["normal"] = normal
	c.JSON(data, nil)
}

func hotArchiveList(c *bm.Context) {
	var (
		err   error
		param = new(model.ParamHotList)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Type != model.HotArchiveTag {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	top, view, filter, checked, tags, err := svc.ArchiveRankList(c, param.Prid, param.Rid, model.HotArchiveTag)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 4)
	data["top"] = top
	data["view"] = view
	data["filter"] = filter
	data["checked"] = checked
	data["mode_tag"] = tags
	c.JSON(data, nil)
}

func hotOperate(c *bm.Context) {
	var (
		err   error
		tag   *model.Tag
		param = new(struct {
			TName string `form:"name"  validate:"required"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if tag, err = svc.OperateHotTag(c, param.TName); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 5)
	data["id"] = tag.ID
	data["name"] = tag.Name
	data["type"] = tag.Type
	data["highLight"] = 0
	data["is_business"] = 1
	c.JSON(data, nil)
}

func updateHotTag(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Data string `form:"data"  validate:"required"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	data := new(model.HotRank)
	if err = json.Unmarshal([]byte(param.Data), &data); err != nil {
		log.Error("json.Unmarshal data failed, error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if data.Type < model.TagHot || data.Type > model.TagSubmission {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(data.Top) > conf.Conf.Tag.MaxSubmitNum {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(data.View) > conf.Conf.Tag.MaxSubmitNum {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.UpdateRank(c, data))
}

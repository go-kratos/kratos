package http

import (
	"encoding/json"
	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func channelAll(c *bm.Context) {
	var (
		err      error
		channels []*model.Channel
	)
	if channels, err = svc.AllChannels(c); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(channels, err)
}

func channeList(c *bm.Context) {
	var (
		err      error
		total    int32
		channels []*model.Channel
		param    = new(model.ParamChanneList)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Order == "" {
		param.Order = model.DefaultOrder
	}
	if param.Sort == "" {
		param.Sort = model.DefaultSort
	}
	if param.Pn < model.DefaultPageNum {
		param.Pn = model.DefaultPageNum
	}
	if param.Ps <= 0 {
		param.Ps = model.DefaultPagesize
	}
	if param.INTShield > model.ChannelStateShieldINT {
		param.INTShield = model.ChannelStateShieldINT
	}
	if channels, total, err = svc.ChanneList(c, param); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["page"] = map[string]interface{}{
		"num":   param.Pn,
		"size":  param.Ps,
		"total": total,
	}
	data["channel"] = channels
	c.JSON(data, nil)
}

func channelInfo(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Tid   int64  `form:"id"`
			TName string `form:"name"`
		})
		channel *model.ChannelInfo
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Tid <= 0 && param.TName == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.TName != "" {
		if param.TName, err = checkName(param.TName, model.TNameMaxLen); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	if channel, err = svc.ChannelInfo(c, param.Tid, param.TName); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(channel, nil)
}

func channelEdit(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Data string `form:"data" validate:"required"`
		})
		channel = new(model.ChannelInfo)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(param.Data), &channel); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if channel.Name == "" && channel.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len([]rune(channel.Content)) > model.DetailContentMax || len([]rune(channel.ShortContent)) > model.ShortContentMax {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if channel.Name != "" {
		if channel.Name, err = checkName(channel.Name, model.TNameMaxLen); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	for _, synonym := range channel.Synonyms {
		if synonym.Alias != "" {
			if synonym.Alias, err = checkName(synonym.Alias, model.TNameMaxLen); err != nil {
				c.JSON(nil, err)
				return
			}
		}
	}
	if int32(len(channel.Synonyms)) > model.ChannelSynonymMax {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	rules, err := svc.RuleChecks(channel.Rules)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	channel.Rules = rules
	_, cname := managerInfo(c)
	c.JSON(nil, svc.EditChannel(c, channel, cname))
}

func channelDelete(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			ID int64 `form:"id" validate:"required,gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	_, cname := managerInfo(c)
	c.JSON(nil, svc.DeleteChannel(c, param.ID, cname))
}

func channelState(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			ID    int64 `form:"id" validate:"required,gt=0"`
			State int32 `form:"state" validate:"gte=0,lt=4"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.State == model.ChanStateStop {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	_, cname := managerInfo(c)
	c.JSON(nil, svc.StateChannel(c, param.ID, param.State, cname))
}

func categoryChannels(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			ID   int64  `form:"id"`
			Name string `form:"name"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.ID <= 0 && param.Name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.Name != "" {
		if param.Name, err = svc.CheckChannelCategory(param.Name); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	offline, online, err := svc.ChannelsByCategory(c, param.ID, param.Name)
	if err != nil {
		c.JSON(nil, err)
	}
	data := map[string]interface{}{
		"online":  online,
		"offline": offline,
	}
	c.JSON(data, nil)
}

func migrateChannel(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Type  int64  `form:"type" validate:"required,gt=0"`
			Tid   int64  `form:"tid"`
			TName string `form:"tname"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.TName == "" && param.Tid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.TName != "" {
		if param.TName, err = checkName(param.TName, model.TNameMaxLen); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	_, cname := managerInfo(c)
	c.JSON(nil, svc.MigrateChannel(c, param.Tid, param.Type, param.TName, cname))
}

func sortChannel(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Type int64   `form:"type" validate:"required,gt=0"`
			Tids []int64 `form:"tids,split" validate:"required,min=1,dive,gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.SortChannels(c, param.Type, param.Tids))
}

func recommandChannel(c *bm.Context) {
	c.JSON(svc.RecommandChannels(c))
}

func sortRecommandChannel(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Tops       []int64 `form:"tops,split"`
			Recommends []int64 `form:"recommends,split"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.SortRecommandChannel(c, param.Tops, param.Recommends))
}

func migrateRecommandChannel(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Tid   int64  `form:"tid"`
			TName string `form:"tname"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.TName == "" && param.Tid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if param.TName != "" {
		if param.TName, err = checkName(param.TName, model.TNameMaxLen); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	_, cname := managerInfo(c)
	c.JSON(nil, svc.MigrateRecommendChannel(c, param.Tid, param.TName, cname))
}

func channelShieldINT(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Tid   int64 `form:"tid" validate:"required,gt=0"`
			State int32 `form:"state" validate:"gte=0,lte=1"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	_, uname := managerInfo(c)
	c.JSON(nil, svc.ChannelShieldINT(c, param.Tid, param.State, uname))
}

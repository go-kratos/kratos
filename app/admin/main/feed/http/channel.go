package http

import (
	"strings"

	"go-common/app/admin/main/feed/model/channel"
	cardmodel "go-common/app/admin/main/feed/model/channel"
	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/model/show"
	"go-common/app/admin/main/feed/util"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func trimSpace(v string) string {
	return strings.TrimSpace(v)
}

func addCardSetup(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &channel.AddCardSetup{}
	if err = c.Bind(req); err != nil {
		return
	}
	uid, name := util.UserInfo(c)
	if name == "" {
		c.JSONMap(map[string]interface{}{"message": "请重新登录"}, ecode.Unauthorized)
		c.Abort()
		return
	}
	req.Value = trimSpace(req.Value)
	if err = chanelSvc.AddCardSetup(req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func cardSetupList(c *bm.Context) {
	var (
		err    error
		cPager *cardmodel.SetupPager
	)
	res := map[string]interface{}{}
	req := &struct {
		ID     int    `form:"id"`
		Type   string `form:"type" validate:"required"`
		Person string `form:"person"`
		Title  string `form:"title"`
		Ps     int    `json:"ps" form:"ps" default:"20"` // 分页大小
		Pn     int    `json:"pn" form:"pn" default:"1"`  // 第几个分页
	}{}
	if err = c.Bind(req); err != nil {
		return
	}
	if cPager, err = chanelSvc.CardSetupList(req.ID, req.Type, req.Person, req.Title, req.Pn, req.Ps); err != nil {
		res["message"] = "卡片获取失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(cPager, nil)
}

func delCardSetup(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &struct {
		ID   int    `form:"id" validate:"required"`
		Type string `form:"type" validate:"required"`
	}{}
	if err = c.Bind(req); err != nil {
		return
	}
	uid, name := util.UserInfo(c)
	if name == "" {
		c.JSONMap(map[string]interface{}{"message": "请重新登录"}, ecode.Unauthorized)
		c.Abort()
		return
	}
	if err = chanelSvc.DelCardSetup(req.ID, req.Type, name, uid); err != nil {
		res["message"] = "卡片删除失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func updateCardSetup(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &cardmodel.UpdateCardSetup{}
	if err = c.Bind(req); err != nil {
		return
	}
	req.Value = trimSpace(req.Value)
	card := &cardmodel.AddCardSetup{
		Type:      req.Type,
		Value:     req.Value,
		Title:     req.Title,
		LongTitle: req.LongTitle,
		Content:   req.Content,
	}
	uid, name := util.UserInfo(c)
	if name == "" {
		c.JSONMap(map[string]interface{}{"message": "请重新登录"}, ecode.Unauthorized)
		c.Abort()
		return
	}
	if err = chanelSvc.UpdateCardSetup(req.ID, card, name, uid); err != nil {
		res["message"] = "卡片更新失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func tabList(c *bm.Context) {
	var (
		err   error
		pager *show.ChannelTabPager
	)
	res := map[string]interface{}{}
	req := &show.ChannelTabLP{}
	if err = c.Bind(req); err != nil {
		return
	}
	if pager, err = chanelSvc.TabList(req); err != nil {
		res["message"] = "列表获取失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(pager, nil)
}

func addTab(c *bm.Context) {
	var (
		err error
		//title string
	)
	res := map[string]interface{}{}
	req := &show.ChannelTabAP{}
	if err = c.Bind(req); err != nil {
		return
	}
	uid, name := util.UserInfo(c)
	if name == "" {
		c.JSONMap(map[string]interface{}{"message": "请重新登录"}, ecode.Unauthorized)
		c.Abort()
		return
	}
	req.Person = name
	req.UID = uid
	if _, err = commonSvc.CardPreview(c, common.CardChannelTab, req.TabID); err != nil {
		return
	}
	if err = chanelSvc.AddTab(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func updateTab(c *bm.Context) {
	var (
		err error
		//title string
	)
	res := map[string]interface{}{}
	req := &show.ChannelTabUP{}
	if err = c.Bind(req); err != nil {
		return
	}
	uid, name := util.UserInfo(c)
	if name == "" {
		c.JSONMap(map[string]interface{}{"message": "请重新登录"}, ecode.Unauthorized)
		c.Abort()
		return
	}
	if req.ID <= 0 {
		c.JSONMap(map[string]interface{}{"message": "ID 参数不合法"}, ecode.RequestErr)
		c.Abort()
		return
	}
	req.Person = name
	req.UID = uid
	if _, err = commonSvc.CardPreview(c, common.CardChannelTab, req.TabID); err != nil {
		return
	}
	if err = chanelSvc.UpdateTab(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func deleteTab(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &struct {
		ID int64 `form:"id" validate:"required"`
	}{}
	if err = c.Bind(req); err != nil {
		return
	}
	uid, name := util.UserInfo(c)
	if name == "" {
		c.JSONMap(map[string]interface{}{"message": "请重新登录"}, ecode.Unauthorized)
		c.Abort()
		return
	}
	if req.ID <= 0 {
		c.JSONMap(map[string]interface{}{"message": "ID 参数不合法"}, ecode.RequestErr)
		c.Abort()
		return
	}
	if err = chanelSvc.DeleteTab(req.ID, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func offlineTab(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &struct {
		ID int64 `form:"id" validate:"required"`
	}{}
	if err = c.Bind(req); err != nil {
		return
	}
	uid, name := util.UserInfo(c)
	if name == "" {
		c.JSONMap(map[string]interface{}{"message": "请重新登录"}, ecode.Unauthorized)
		c.Abort()
		return
	}
	if req.ID <= 0 {
		c.JSONMap(map[string]interface{}{"message": "ID 参数不合法"}, ecode.RequestErr)
		c.Abort()
		return
	}
	if err = chanelSvc.OfflineTab(req.ID, name, uid); err != nil {
		res["message"] = "卡片下线失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

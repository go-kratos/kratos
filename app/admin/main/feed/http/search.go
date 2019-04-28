package http

import (
	"encoding/json"

	searchModel "go-common/app/admin/main/feed/model/search"
	"go-common/app/admin/main/feed/model/show"
	"go-common/app/admin/main/feed/util"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/time"
)

//Black 黑名单
func blackList(c *bm.Context) {
	var (
		err   error
		black []searchModel.Black
	)
	res := map[string]interface{}{}
	if black, err = searchSvc.BlackList(); err != nil {
		res["message"] = "获取失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(black, nil)
}

//addBlack 添加黑名单
func addBlack(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	param := new(searchModel.Black)
	if err = c.Bind(param); err != nil {
		return
	}
	uid, name := managerInfo(c)
	if err = searchSvc.AddBlack(c, param.Searchword, name, uid); err != nil {
		res["message"] = "获取失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//delBlack 删除黑名单
func delBlack(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	param := new(struct {
		ID int `form:"id" validate:"required"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	uid, name := managerInfo(c)
	if err = searchSvc.DelBlack(c, param.ID, name, uid); err != nil {
		res["message"] = "获取失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//openAddDarkword 对外 添加黑马词
func openAddDarkword(c *bm.Context) {
	var (
		err  error
		dark searchModel.OpenDark
	)
	res := map[string]interface{}{}
	param := &struct {
		Data string `form:"data" validate:"required"`
	}{}
	if err = c.Bind(param); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(param.Data), &dark); err != nil {
		res["message"] = "参数有误:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if err = searchSvc.OpenAddDarkword(c, dark); err != nil {
		res["message"] = "添加失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//openBlacklist 对外 黑名单列表
func openBlacklist(c *bm.Context) {
	var (
		err   error
		black []searchModel.Black
	)
	res := map[string]interface{}{}
	if black, err = searchSvc.BlackList(); err != nil {
		res["message"] = "获取失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(black, nil)
}

//OpenHotList 对外 黑名单列表
func openHotList(c *bm.Context) {
	var (
		err error
		hot []searchModel.Intervene
	)
	res := map[string]interface{}{}
	if hot, err = searchSvc.OpenHotList(c); err != nil {
		res["message"] = "获取失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(hot, nil)
}

//openDarkword 对外 获取黑马词
func openDarkword(c *bm.Context) {
	var (
		err  error
		dark []searchModel.Dark
	)
	res := map[string]interface{}{}
	if dark, err = searchSvc.GetDarkPub(c); err != nil {
		res["message"] = "获取失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(dark, nil)
}

//openAddHotword 对外 添加搜索热词
func openAddHotword(c *bm.Context) {
	var (
		err error
		hot searchModel.OpenHot
	)
	res := map[string]interface{}{}
	param := &struct {
		Data string `form:"data" validate:"required"`
	}{}
	if err = c.Bind(param); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(param.Data), &hot); err != nil {
		res["message"] = "参数有误:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if err = searchSvc.OpenAddHotword(c, hot); err != nil {
		res["message"] = "添加失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//publishHotWord publish hot word
func publishHotWord(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	uid, name := managerInfo(c)
	if err = searchSvc.SetHotPub(c, name, uid); err != nil {
		res["message"] = "发布失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//publishDark publish dark word
func publishDarkWord(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	uid, name := managerInfo(c)
	if err = searchSvc.SetDarkPub(c, name, uid); err != nil {
		res["message"] = "发布失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//addInter 添加干预
func addInter(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := searchModel.InterveneAdd{}
	if err = c.Bind(&param); err != nil {
		return
	}
	uid, name := managerInfo(c)
	if err = searchSvc.AddInter(c, param, name, uid); err != nil {
		res["message"] = "添加失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//updateInter 更新干预
func updateInter(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := struct {
		ID         int       `form:"id" validate:"required"`
		Searchword string    `form:"searchword" validate:"required"`
		Rank       int       `form:"position" validate:"required"`
		Tag        string    `form:"tag"`
		Stime      time.Time `form:"stime" validate:"required"`
		Etime      time.Time `form:"etime" validate:"required"`
	}{}
	if err = c.Bind(&param); err != nil {
		return
	}
	inter := searchModel.InterveneAdd{
		Searchword: param.Searchword,
		Rank:       param.Rank,
		Tag:        param.Tag,
		Stime:      param.Stime,
		Etime:      param.Etime,
	}
	uid, name := managerInfo(c)
	if err = searchSvc.UpdateInter(c, inter, param.ID, name, uid); err != nil {
		res["message"] = "更新失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//deleteHot 删除热词
func deleteHot(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := struct {
		ID   int   `form:"id" validate:"required"`
		Type uint8 `form:"type" validate:"required"`
	}{}
	if err = c.Bind(&param); err != nil {
		return
	}
	uid, name := managerInfo(c)
	if err = searchSvc.DeleteHot(c, param.ID, param.Type, name, uid); err != nil {
		res["message"] = "删除失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//deleteDark 删除黑马词
func deleteDark(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := struct {
		ID int `form:"id" validate:"required"`
	}{}
	if err = c.Bind(&param); err != nil {
		return
	}
	uid, name := managerInfo(c)
	if err = searchSvc.DeleteDark(c, param.ID, name, uid); err != nil {
		res["message"] = "删除失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//updateSearch 更新搜索推过来的词
func updateSearch(c *bm.Context) {
	var (
		err error
		res = map[string]interface{}{}
	)
	param := struct {
		ID  int    `form:"id" validate:"required"`
		Tag string `form:"tag"`
	}{}
	if err = c.Bind(&param); err != nil {
		return
	}
	uid, name := managerInfo(c)
	if err = searchSvc.UpdateSearch(c, param.Tag, param.ID, name, uid); err != nil {
		res["message"] = "更新失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

//HotList 搜索热词
func HotList(c *bm.Context) {
	var (
		err    error
		hotout searchModel.HotwordOut
	)
	res := map[string]interface{}{}
	param := struct {
		Date string `form:"date" validate:"required"`
	}{}
	if err = c.Bind(&param); err != nil {
		return
	}
	if hotout, err = searchSvc.HotList(c, param.Date); err != nil {
		res["message"] = "获取热词失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(hotout, nil)
}

//darkList 黑马词
func darkList(c *bm.Context) {
	var (
		err     error
		darkout searchModel.DarkwordOut
	)
	res := map[string]interface{}{}
	param := struct {
		Date string `form:"date" validate:"required"`
	}{}
	if err = c.Bind(&param); err != nil {
		return
	}
	if darkout, err = searchSvc.DarkList(c, param.Date); err != nil {
		res["message"] = "获取黑马词失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(darkout, nil)
}

func searchWebCardList(c *bm.Context) {
	var (
		err   error
		pager *show.SearchWebCardPager
	)
	res := map[string]interface{}{}
	req := &show.SearchWebCardLP{}
	if err = c.Bind(req); err != nil {
		return
	}
	if pager, err = searchSvc.SearchWebCardList(req); err != nil {
		res["message"] = "列表获取失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(pager, nil)
}

func addSearchWebCard(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &show.SearchWebCardAP{}
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
	if err = searchSvc.AddSearchWebCard(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func upSearchWebCard(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &show.SearchWebCardUP{}
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
	if err = searchSvc.UpdateSearchWebCard(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func delSearchWebCard(c *bm.Context) {
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
	if err = searchSvc.DeleteSearchWebCard(req.ID, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func searchWebList(c *bm.Context) {
	var (
		err   error
		pager *show.SearchWebPager
	)
	res := map[string]interface{}{}
	req := &show.SearchWebLP{}
	if err = c.Bind(req); err != nil {
		return
	}
	if pager, err = searchSvc.SearchWebList(req); err != nil {
		res["message"] = "列表获取失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(pager, nil)
}

func openSearchWeb(c *bm.Context) {
	var (
		err   error
		pager []*show.SearchWeb
	)
	res := map[string]interface{}{}
	if pager, err = searchSvc.OpenSearchWebList(); err != nil {
		res["message"] = "列表获取失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(pager, nil)
}

func addSearchWeb(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &show.SearchWebAP{}
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
	if err = searchSvc.AddSearchWeb(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func upSearchWeb(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &show.SearchWebUP{}
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
	if err = searchSvc.UpdateSearchWeb(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func delSearchWeb(c *bm.Context) {
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
	if err = searchSvc.DeleteSearchWeb(req.ID, name, uid); err != nil {
		res["message"] = "卡片删除失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func optSearchWeb(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &struct {
		ID  int64  `form:"id" validate:"required"`
		Opt string `form:"opt" validate:"required"`
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
	if err = searchSvc.OptionSearchWeb(req.ID, req.Opt, name, uid); err != nil {
		res["message"] = "修改失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

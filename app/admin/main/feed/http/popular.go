package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/model/show"
	"go-common/app/admin/main/feed/util"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func eventTopicList(c *bm.Context) {
	var (
		err   error
		pager *show.EventTopicPager
	)
	res := map[string]interface{}{}
	req := &show.EventTopicLP{}
	if err = c.Bind(req); err != nil {
		return
	}
	if pager, err = popularSvc.EventTopicList(req); err != nil {
		res["message"] = "列表获取失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(pager, nil)
}

func addEventTopic(c *bm.Context) {
	var (
		err error
		//title string
	)
	res := map[string]interface{}{}
	req := &show.EventTopicAP{}
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
	if err = popularSvc.AddEventTopic(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func upEventTopic(c *bm.Context) {
	var (
		err error
	)
	res := map[string]interface{}{}
	req := &show.EventTopicUP{}
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
	if err = popularSvc.UpdateEventTopic(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func delEventTopic(c *bm.Context) {
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
	if err = popularSvc.DeleteEventTopic(req.ID, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}
func popularStarsList(c *bm.Context) {
	var (
		err   error
		pager *show.PopularStarsPager
	)
	res := map[string]interface{}{}
	req := &show.PopularStarsLP{}
	if err = c.Bind(req); err != nil {
		return
	}
	if pager, err = popularSvc.PopularStarsList(req); err != nil {
		res["message"] = "列表获取失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(pager, nil)
}

//ValidateUpid .
func ValidateUpid(c context.Context, idStr string) (name string, err error) {
	var (
		id int64
	)
	if id, err = strconv.ParseInt(idStr, 10, 64); err != nil {
		return
	}
	if name, err = commonSvc.CardPreview(c, common.CardUp, id); err != nil {
		return
	}
	return
}

//ValidateAvid .
func ValidateAvid(c context.Context, values string) (err error) {
	type Content struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}
	var contents []*Content
	if err = json.Unmarshal([]byte(values), &contents); err != nil {
		return
	}
	dup := make(map[int64]bool)
	for _, v := range contents {
		if v == nil {
			return fmt.Errorf("不能提交空数据！")
		}
		if dup[v.ID] {
			return fmt.Errorf("重复视频ID (%d)", v.ID)
		}
		dup[v.ID] = true
		if _, err = commonSvc.CardPreview(c, common.CardAv, v.ID); err != nil {
			return
		}
	}
	if len(contents) < 3 {
		return fmt.Errorf("视频组成最少三个")
	}
	if len(contents) > 5 {
		return fmt.Errorf("视频组成不能超过5个")
	}
	return
}

func addPopularStars(c *bm.Context) {
	var (
		err    error
		upName string
	)
	res := map[string]interface{}{}
	req := &show.PopularStarsAP{}
	if err = c.Bind(req); err != nil {
		return
	}
	uid, name := util.UserInfo(c)
	if name == "" {
		c.JSONMap(map[string]interface{}{"message": "请重新登录"}, ecode.Unauthorized)
		c.Abort()
		return
	}
	if err = ValidateAvid(c, req.Content); err != nil {
		c.JSONMap(map[string]interface{}{"message": "视频ID 校验失败：" + err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}
	req.Value = util.TrimStrSpace(req.Value)
	if upName, err = ValidateUpid(c, req.Value); err != nil {
		c.JSONMap(map[string]interface{}{"message": "up主ID 校验失败：" + err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}
	req.Person = name
	req.UID = uid
	req.LongTitle = upName
	if err = popularSvc.AddPopularStars(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func updatePopularStars(c *bm.Context) {
	var (
		err    error
		upName string
	)
	res := map[string]interface{}{}
	req := &show.PopularStarsUP{}
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
	if err = ValidateAvid(c, req.Content); err != nil {
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}
	req.Value = util.TrimStrSpace(req.Value)
	if upName, err = ValidateUpid(c, req.Value); err != nil {
		c.JSONMap(map[string]interface{}{"message": err.Error()}, ecode.RequestErr)
		c.Abort()
		return
	}
	req.LongTitle = upName
	if err = popularSvc.UpdatePopularStars(c, req, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func deletePopularStars(c *bm.Context) {
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
	if err = popularSvc.DeletePopularStars(req.ID, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func rejectPopularStars(c *bm.Context) {
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
	if err = popularSvc.RejectPopularStars(req.ID, name, uid); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

func aiAddPopularStars(c *bm.Context) {
	var (
		err        error
		aids       []byte
		addPopStar []*show.PopularStarsAP
	)
	res := map[string]interface{}{}
	req := &struct {
		Data string `form:"data" validate:"required"`
	}{}
	if err = c.Bind(req); err != nil {
		return
	}
	log.Info("aiAddPopularStars value(%v)", req.Data)
	values := make([]*show.PopularStarsAIAP, 0)
	if err = json.Unmarshal([]byte(req.Data), &values); err != nil {
		log.Error("aiAddPopularStars.Unmarshal value(%v) error(%v)", req.Data, err)
		res["message"] = "数据解析失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	for _, v := range values {
		var (
			AiValues []*show.AiValue
			upName   string
		)
		for _, aid := range v.Aids {
			aiValue := &show.AiValue{
				ID: aid,
			}
			AiValues = append(AiValues, aiValue)
		}
		if aids, err = json.Marshal(AiValues); err != nil {
			log.Error("aiAddPopularStars.Marshal value(%v) error(%v)", v.Aids, err)
			res["message"] = "数据encode失败 " + err.Error()
			c.JSONMap(res, ecode.RequestErr)
		}
		mid := strconv.FormatInt(v.Mid, 10)
		if upName, err = ValidateUpid(c, mid); err != nil {
			c.JSONMap(map[string]interface{}{"message": "up主ID 校验失败：" + err.Error()}, ecode.RequestErr)
			c.Abort()
			return
		}

		tmp := &show.PopularStarsAP{
			LongTitle: upName,
			Value:     strconv.FormatInt(v.Mid, 10),
			Content:   string(aids),
		}
		addPopStar = append(addPopStar, tmp)
	}
	if err = popularSvc.AIAddPopularStars(c, addPopStar); err != nil {
		res["message"] = "卡片创建失败 " + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(nil, nil)
}

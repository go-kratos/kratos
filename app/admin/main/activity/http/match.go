package http

import (
	"fmt"
	"go-common/app/admin/main/activity/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// addMatch 增加赛程
func addMatch(c *bm.Context) {
	arg := new(model.ActMatchs)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addMatch(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// saveMatch 存储赛程
func saveMatch(c *bm.Context) {
	arg := new(model.ActMatchs)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Model(&model.ActMatchs{ID: arg.ID}).Update(arg).Error; err != nil {
		log.Error("saveMatch(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// matchInfo 赛程信息
func matchInfo(c *bm.Context) {
	arg := new(model.ActMatchs)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.First(arg, arg.ID).Error; err != nil {
		log.Error("matcInfo(%d) error(%v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}

// matchList 比赛对象列表
func matchList(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.ActMatchs
	)
	v := new(struct {
		SID    int64 `form:"sid" default:"-1"`
		Status int8  `form:"status" default:"-1"`
		Page   int   `form:"pn" default:"1"`
		Size   int   `form:"ps" default:"20"`
	})
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	db := actSrv.DB
	if v.Page == 0 {
		v.Page = 1
	}
	if v.Size == 0 {
		v.Size = 20
	}
	if v.Status != -1 {
		db = db.Where("status = ?", v.Status)
	}
	if v.SID != -1 {
		db = db.Where("sid = ?", v.SID)
	}
	db = db.Where("status = ?", 0)
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("businessList(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.ActMatchs{}).Count(&count).Error; err != nil {
		log.Error("businessList count error(%v)", err)
		c.JSON(nil, err)
		return
	}

	data := map[string]interface{}{
		"data":  list,
		"pn":    v.Page,
		"ps":    v.Size,
		"total": count,
	}
	c.JSONMap(data, nil)
}

// addMatchObject 比赛对象
func addMatchObject(c *bm.Context) {
	arg := new(model.ActMatchsObject)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addMatch(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// saveMatchObject 更新比赛对象
func saveMatchObject(c *bm.Context) {
	arg := new(model.ActMatchsObject)
	if err := c.Bind(arg); err != nil {
		return
	}
	fmt.Println(arg.AwayName)
	fmt.Println(arg.HomeLogo)
	if err := actSrv.DB.Model(&model.ActMatchsObject{ID: arg.ID}).Update(arg).Error; err != nil {
		log.Error("saveMatch(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// matchObjectInfo 比赛对象信息
func matchObjectInfo(c *bm.Context) {
	arg := new(model.ActMatchsObject)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.First(arg, arg.ID).Error; err != nil {
		log.Error("matcInfo(%d) error(%v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}

// matchObjectList 比赛对象列表
func matchObjectList(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.ActMatchsObject
	)
	v := new(struct {
		SID     int64 `form:"sid" default:"-1"`
		Status  int8  `form:"status" default:"-1"`
		Page    int   `form:"pn" default:"1"`
		Size    int   `form:"ps" default:"20"`
		MatchID int   `form:"match_id"`
	})
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	db := actSrv.DB
	if v.Page == 0 {
		v.Page = 1
	}
	if v.Size == 0 {
		v.Size = 20
	}
	if v.Status != -1 {
		db = db.Where("status = ?", v.Status)
	}
	if v.SID != -1 {
		db = db.Where("sid = ?", v.SID)
	}
	if v.MatchID != 0 {
		db = db.Where("match_id = ?", v.MatchID)
	}
	db = db.Where("status = ?", 0)
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("businessList(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.ActMatchsObject{}).Count(&count).Error; err != nil {
		log.Error("businessList count error(%v)", err)
		c.JSON(nil, err)
		return
	}

	data := map[string]interface{}{
		"data":  list,
		"pn":    v.Page,
		"ps":    v.Size,
		"total": count,
	}
	c.JSONMap(data, nil)
}

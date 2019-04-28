package http

import (
	"go-common/app/admin/main/activity/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// addBws 增加赛程
func addBws(c *bm.Context) {
	arg := new(model.ActBws)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addBws(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// saveBws 存储赛程
func saveBws(c *bm.Context) {
	arg := new(model.ActBws)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Model(&model.ActBws{ID: arg.ID}).Update(arg).Error; err != nil {
		log.Error("saveBws(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// bwsInfo 赛程信息
func bwsInfo(c *bm.Context) {
	arg := new(model.ActBws)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.First(arg, arg.ID).Error; err != nil {
		log.Error("bwsInfo(%d) error(%v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}

// bwsList 比赛对象列表
func bwsList(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.ActBws
	)
	v := new(struct {
		Del  int8 `form:"del" default:"0"`
		Page int  `form:"pn" default:"1"`
		Size int  `form:"ps" default:"20"`
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
	db = db.Where("del = ?", v.Del)
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("bwsList(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.ActBws{}).Count(&count).Error; err != nil {
		log.Error("bwsList count error(%v)", err)
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

// addBwsAchievement 增加赛程
func addBwsAchievement(c *bm.Context) {
	arg := new(model.ActBwsAchievement)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addBws(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// saveBwsAchievement 存储赛程
func saveBwsAchievement(c *bm.Context) {
	arg := new(model.ActBwsAchievement)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Model(&model.ActBwsAchievement{ID: arg.ID}).Update(arg).Error; err != nil {
		log.Error("saveBwsAchievement(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// actBwsAchievement 赛程信息
func bwsAchievement(c *bm.Context) {
	arg := new(model.ActBwsAchievement)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.First(arg, arg.ID).Error; err != nil {
		log.Error("bwsAchievement(%d) error(%v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}

// bwsList 比赛对象列表
func bwsAchievements(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.ActBwsAchievement
	)
	v := new(struct {
		BID  int64 `form:"bid" default:"0"`
		Del  int8  `form:"del" default:"0"`
		Page int   `form:"pn" default:"1"`
		Size int   `form:"ps" default:"20"`
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
	db = db.Where("del = ?", v.Del)
	if v.BID != 0 {
		db = db.Where("bid = ?", v.BID)
	}
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("bwsAchievements(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.ActBwsAchievement{}).Count(&count).Error; err != nil {
		log.Error("bwsAchievements count error(%v)", err)
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

// addBwsField 增加赛程
func addBwsField(c *bm.Context) {
	arg := new(model.ActBwsField)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addBwsFieldws(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// saveBwsField 存储赛程
func saveBwsField(c *bm.Context) {
	arg := new(model.ActBwsField)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Model(&model.ActBwsAchievement{ID: arg.ID}).Update(arg).Error; err != nil {
		log.Error("saveBwsField(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// actBwsAchievement 赛程信息
func bwsField(c *bm.Context) {
	arg := new(model.ActBwsField)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.First(arg, arg.ID).Error; err != nil {
		log.Error("bwsField(%d) error(%v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}

// bwsList 比赛对象列表
func bwsFields(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.ActBwsField
	)
	v := new(struct {
		BID  int64 `form:"bid" default:"0"`
		Del  int8  `form:"del" default:"0"`
		Page int   `form:"pn" default:"1"`
		Size int   `form:"ps" default:"20"`
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
	db = db.Where("del = ?", v.Del)
	if v.BID != 0 {
		db = db.Where("bid = ?", v.BID)
	}
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("bwsFields(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.ActBwsField{}).Count(&count).Error; err != nil {
		log.Error("bwsFields count error(%v)", err)
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

// addBwsPoint
func addBwsPoint(c *bm.Context) {
	arg := new(model.ActBwsPoint)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addBaddBwsFieldws(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// saveBwsPoint 存储赛程
func saveBwsPoint(c *bm.Context) {
	arg := new(model.ActBwsPoint)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Model(&model.ActBwsPoint{ID: arg.ID}).Update(arg).Error; err != nil {
		log.Error("saveBwsPoint(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// actBwsAchievement 赛程信息
func bwsPoint(c *bm.Context) {
	arg := new(model.ActBwsPoint)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.First(arg, arg.ID).Error; err != nil {
		log.Error("bwsPoint(%d) error(%v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}

// bwsList 比赛对象列表
func bwsPoints(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.ActBwsPoint
	)
	v := new(struct {
		FID      int64 `form:"fid" default:"0"`
		BID      int64 `form:"bid" default:"0"`
		LockType int64 `form:"lock_type" default:"lock_type"`
		Del      int8  `form:"del" default:"0"`
		Page     int   `form:"pn" default:"1"`
		Size     int   `form:"ps" default:"20"`
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
	db = db.Where("del = ?", v.Del)
	if v.BID != 0 {
		db = db.Where("bid = ?", v.BID)
	}
	if v.FID != 0 {
		db = db.Where("fid = ?", v.FID)
	}
	if v.LockType != 0 {
		db = db.Where("lock_type = ?", v.LockType)
	}
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("bwsPoints(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.ActBwsPoint{}).Count(&count).Error; err != nil {
		log.Error("bwsPoints count error(%v)", err)
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

// addBwsUserAchievement 保存用户
func addBwsUserAchievement(c *bm.Context) {
	arg := new(model.ActBwsUserAchievement)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addBwsUserAchievement(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// saveBwsUserAchievement 保存用户成就
func saveBwsUserAchievement(c *bm.Context) {
	arg := new(model.ActBwsUserAchievement)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Model(&model.ActBwsUserAchievement{ID: arg.ID}).Update(arg).Error; err != nil {
		log.Error("saveBwsUserAchievement(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// bwsUserAchievement 用户成就信息
func bwsUserAchievement(c *bm.Context) {
	arg := new(model.ActBwsUserAchievement)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.First(arg, arg.ID).Error; err != nil {
		log.Error("bwsUserAchievement(%d) error(%v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}

// bwsUserAchievements 用户成就列表
func bwsUserAchievements(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.ActBwsUserAchievement
	)
	v := new(struct {
		MID  int64 `form:"mid" default:"0"`
		AID  int64 `form:"aid" default:"0"`
		BID  int64 `form:"bid" default:"0"`
		Del  int8  `form:"del" default:"0"`
		Page int   `form:"pn" default:"1"`
		Size int   `form:"ps" default:"20"`
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
	db = db.Where("del = ?", v.Del)
	if v.BID != 0 {
		db = db.Where("bid = ?", v.BID)
	}
	if v.AID != 0 {
		db = db.Where("aid = ?", v.AID)
	}
	if v.MID != 0 {
		db = db.Where("mid = ?", v.MID)
	}
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("bwsUserAchievements(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.ActBwsUserAchievement{}).Count(&count).Error; err != nil {
		log.Error("bwsUserAchievements count error(%v)", err)
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

// addBwsUserPoint 保存用户
func addBwsUserPoint(c *bm.Context) {
	arg := new(model.ActBwsUserPoint)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addBwsUserPoint(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// saveBwsUserAchievement 保存用户成就
func saveBwsUserPoint(c *bm.Context) {
	arg := new(model.ActBwsUserPoint)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Model(&model.ActBwsUserPoint{ID: arg.ID}).Update(arg).Error; err != nil {
		log.Error("saveBwsUserPoint(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// bwsUserPoint 用户点数信息
func bwsUserPoint(c *bm.Context) {
	arg := new(model.ActBwsUserPoint)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.First(arg, arg.ID).Error; err != nil {
		log.Error("bwsUserPoint(%d) error(%v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}

// bwsUserPoints 用户点数列表
func bwsUserPoints(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.ActBwsUserPoint
	)
	v := new(struct {
		MID      int64  `form:"mid" default:"0"`
		KEY      string `form:"key" default:""`
		LockType int64  `form:"lock_type" default:"lock_type"`
		BID      int64  `form:"bid" default:"0"`
		Del      int8   `form:"del" default:"0"`
		Page     int    `form:"pn" default:"1"`
		Size     int    `form:"ps" default:"20"`
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
	db = db.Where("del = ?", v.Del)
	if v.BID != 0 {
		db = db.Where("bid = ?", v.BID)
	}
	if v.KEY != "" {
		db = db.Where("key = ?", v.KEY)
	}
	if v.MID != 0 {
		db = db.Where("mid = ?", v.MID)
	}
	if v.LockType != 0 {
		db = db.Where("lock_type = ?", v.LockType)
	}
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("bwsUserPoints(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.ActBwsUserPoint{}).Count(&count).Error; err != nil {
		log.Error("bwsUserPoints count error(%v)", err)
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

// addBwsUser 添加用户
func addBwsUser(c *bm.Context) {
	arg := new(model.ActBwsUser)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addBwsUserPoint(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// saveBwsUserAchievement 保存用户成就
func saveBwsUser(c *bm.Context) {
	arg := new(model.ActBwsUser)
	if err := c.Bind(arg); err != nil {
		return
	}
	if err := actSrv.DB.Model(&model.ActBwsUser{ID: arg.ID}).Update(arg).Error; err != nil {
		log.Error("saveBwsUserPoint(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

// bwsUserPoint 用户点数信息
func bwsUser(c *bm.Context) {
	arg := new(model.ActBwsUser)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.First(arg, arg.ID).Error; err != nil {
		log.Error("bwsUserPoint(%d) error(%v)", arg.ID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}

// bwsUserPoints 用户点数列表
func bwsUsers(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.ActBwsUser
	)
	v := new(struct {
		MID  int64  `form:"mid" default:"0"`
		KEY  string `form:"key" default:""`
		BID  int64  `form:"bid" default:"0"`
		Del  int8   `form:"del" default:"0"`
		Page int    `form:"pn" default:"1"`
		Size int    `form:"ps" default:"20"`
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
	db = db.Where("del = ?", v.Del)
	if v.BID != 0 {
		db = db.Where("bid = ?", v.BID)
	}
	if v.KEY != "" {
		db = db.Where("key = ?", v.KEY)
	}
	if v.MID != 0 {
		db = db.Where("mid = ?", v.MID)
	}
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("bwsUsers(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.ActBwsUser{}).Count(&count).Error; err != nil {
		log.Error("bwsUsers count error(%v)", err)
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

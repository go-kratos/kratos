package http

import (
	"go-common/app/admin/main/creative/model/operation"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"net/http"
	"time"
)

var (
	typesMap = map[int8]int8{
		0: 0, // 显示
		1: 1, // 隐藏
		2: 2, // 全部
	}
	platformMap = map[int8]int8{
		0:   0,   // web+app
		1:   1,   // app
		2:   2,   // web
		100: 100, // 全平台
	}
)

func viewNotice(c *bm.Context) {
	var (
		err  error
		info = &operation.Operation{}
	)
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Debug().Model(&operation.Operation{}).First(info, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if info.ID == 0 {
		c.JSON(nil, ecode.NothingFound)
	} else {
		layout := "2006-01-02 15:04:05"
		ftime, _ := time.Parse(time.RFC3339, info.Stime)
		info.Stime = ftime.Format(layout)
		ftime, _ = time.Parse(time.RFC3339, info.Etime)
		info.Etime = ftime.Format(layout)
		ftime, _ = time.Parse(time.RFC3339, info.Ctime)
		info.Ctime = ftime.Format(layout)
		ftime, _ = time.Parse(time.RFC3339, info.Mtime)
		info.Mtime = ftime.Format(layout)
		ftime, _ = time.Parse(time.RFC3339, info.Dtime)
		info.Dtime = ftime.Format(layout)
		c.JSON(info, nil)
	}
}

func listNotice(c *bm.Context) {
	var (
		err   error
		ops   []*operation.Operation
		total int
	)
	v := new(struct {
		Type     int8 `form:"type"`
		Platform int8 `form:"platform"`
		Page     int  `form:"page" validate:"min=1"`
		PageSize int  `form:"pagesize" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Page < 1 {
		v.Page = 1
	}
	if v.PageSize < 20 {
		v.PageSize = 20
	}
	if _, ok := typesMap[v.Type]; !ok {
		data := map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": "forbid request with wrong type enum value",
		}
		c.Render(http.StatusOK, render.MapJSON(data))
		return
	}
	if _, ok := platformMap[v.Platform]; !ok {
		data := map[string]interface{}{
			"code":    ecode.RequestErr,
			"message": "forbid request with wrong platform enum value",
		}
		c.Render(http.StatusOK, render.MapJSON(data))
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	where := " dtime = '0000-00-00 00:00:00' AND type = 'play' AND platform = ?"
	if v.Platform == 100 {
		where = " dtime = '0000-00-00 00:00:00' AND type = 'play' "
	}
	if v.Type == 0 {
		where += " AND (stime < ? AND etime > ?) "
	} else if v.Type == 1 {
		where += " AND (stime > ? OR etime < ?) "
	} else if v.Type == 2 {
		where += " AND 1=1 "
	}
	if v.Type == 2 {
		if v.Platform == 100 {
			if err = svc.DB.Debug().Model(&operation.Operation{}).Order("rank ASC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&ops, where).Error; err != nil {
				c.JSON(nil, err)
				return
			}
			svc.DB.Debug().Model(&operation.Operation{}).Where(where).Count(&total)
		} else {
			if err = svc.DB.Debug().Model(&operation.Operation{}).Order("rank ASC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&ops, where, v.Platform).Error; err != nil {
				c.JSON(nil, err)
				return
			}
			svc.DB.Debug().Model(&operation.Operation{}).Where(where, v.Platform).Count(&total)
		}
	} else {
		if v.Platform == 100 {
			if err = svc.DB.Debug().Model(&operation.Operation{}).Order("rank ASC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&ops, where, now, now).Error; err != nil {
				c.JSON(nil, err)
				return
			}
			svc.DB.Debug().Model(&operation.Operation{}).Where(where, now, now).Count(&total)
		} else {
			if err = svc.DB.Debug().Model(&operation.Operation{}).Order("rank ASC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&ops, where, v.Platform, now, now).Error; err != nil {
				c.JSON(nil, err)
				return
			}
			svc.DB.Debug().Model(&operation.Operation{}).Where(where, v.Platform, now, now).Count(&total)
		}
	}
	var opsView []*operation.ViewOperation
	for _, v := range ops {
		var (
			status   string
			timeNow  = time.Now()
			stime, _ = time.Parse(time.RFC3339, v.Stime)
			etime, _ = time.Parse(time.RFC3339, v.Etime)
		)
		if time.Now().Before(etime) && timeNow.After(stime) {
			status = "显示"
		} else {
			status = "隐藏"
		}
		opsView = append(opsView, &operation.ViewOperation{
			ID:       v.ID,
			Type:     v.Type,
			Ads:      v.Ads,
			Platform: v.Platform,
			Rank:     v.Rank,
			Pic:      v.Pic,
			Link:     v.Link,
			Content:  v.Content,
			Username: v.Username,
			Remark:   v.Remark,
			Note:     v.Note,
			AppPic:   v.AppPic,
			Stime:    v.Stime,
			Etime:    v.Etime,
			Ctime:    v.Ctime,
			Mtime:    v.Mtime,
			Dtime:    v.Dtime,
			Status:   status,
		})
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    opsView,
		"pager": map[string]int{
			"page":     v.Page,
			"pagesize": v.PageSize,
			"total":    total,
		},
	}))
}

func addNotice(c *bm.Context) {
	var (
		err error
	)
	v := new(struct {
		Ads      int8   `form:"ads" `
		Rank     int8   `form:"rank"`
		Pic      string `form:"pic" validate:"required"`
		Link     string `form:"link" validate:"required"`
		Content  string `form:"content" validate:"required"`
		Username string `form:"username" validate:"required"`
		Remark   string `form:"remark"`
		Note     string `form:"note"`
		AppPic   string `form:"app_pic" validate:"required"`
		Platform int8   `form:"platform"`
		Stime    string `form:"stime" validate:"required"`
		Etime    string `form:"etime" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	m := &operation.Operation{
		Type:     "play",
		Ads:      v.Ads,
		Rank:     v.Rank,
		Pic:      v.Pic,
		Link:     v.Link,
		Content:  v.Content,
		Username: v.Username,
		Remark:   v.Remark,
		Note:     v.Note,
		AppPic:   v.AppPic,
		Platform: v.Platform,
		Ctime:    time.Now().Format("2006-01-02 15:04:05"),
		Stime:    v.Stime,
		Etime:    v.Etime,
	}
	db := svc.DB.Debug().Model(&operation.Operation{}).Create(m)
	if err = db.Error; err != nil {
		log.Error("creativeSvc.Operation error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(map[string]interface{}{
		"id": db.Value.(*operation.Operation).ID,
	}, nil)
}

func upNotice(c *bm.Context) {
	var (
		op  = &operation.Operation{}
		err error
	)
	v := new(struct {
		ID       int64  `form:"id"`
		Ads      int8   `form:"ads"`
		Platform int8   `form:"platform"`
		Rank     int8   `form:"rank"`
		Pic      string `form:"pic"`
		Link     string `form:"link"`
		Content  string `form:"content"`
		Username string `form:"username"`
		Remark   string `form:"remark"`
		Note     string `form:"note"`
		AppPic   string `form:"app_pic"`
		Stime    string `form:"stime"`
		Etime    string `form:"etime"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Debug().Model(&operation.Operation{}).Find(op, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = svc.DB.Debug().Model(&operation.Operation{ID: v.ID}).Update(map[string]interface {
	}{
		"id":       v.ID,
		"ads":      v.Ads,
		"rank":     v.Rank,
		"pic":      v.Pic,
		"link":     v.Link,
		"content":  v.Content,
		"username": v.Username,
		"remark":   v.Remark,
		"note":     v.Note,
		"app_pic":  v.AppPic,
		"platform": v.Platform,
		"mtime":    time.Now().Format("2006-01-02 15:04:05"),
		"stime":    v.Stime,
		"etime":    v.Etime,
	}).Error; err != nil {
		log.Error("svc.save error(%v)", err)
	}
	c.JSON(nil, err)
}

func delNotice(c *bm.Context) {
	var (
		op  = &operation.Operation{}
		err error
	)
	if err = c.Bind(op); err != nil {
		return
	}
	if op.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Debug().Model(&operation.Operation{}).Find(op, op.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = svc.DB.Debug().Model(&operation.Operation{ID: op.ID}).Update("dtime", time.Now().Format("2006-01-02 15:04:05")).Error; err != nil {
		log.Error("svc.del Notice error(%v)", err)
	}
	c.JSON(nil, err)
}

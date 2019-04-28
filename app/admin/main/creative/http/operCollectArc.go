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

func listCollectArcOper(c *bm.Context) {
	var (
		err   error
		ops   []*operation.Operation
		total int
	)
	v := new(struct {
		Type     int8 `form:"type"`
		Platform int8 `form:"platform"`
		Pn       int  `form:"pn" validate:"min=1"`
		Ps       int  `form:"ps" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Pn < 1 {
		v.Pn = 1
	}
	if v.Ps < 20 {
		v.Ps = 20
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
	where := " dtime = '0000-00-00 00:00:00' AND type = 'collect_arc' AND platform = ?"
	if v.Platform == 100 {
		where = " dtime = '0000-00-00 00:00:00' AND type = 'collect_arc' "
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
			if err = svc.DB.Debug().Model(&operation.Operation{}).Order("rank ASC").Offset((v.Pn-1)*v.Ps).Limit(v.Ps).Find(&ops, where).Error; err != nil {
				c.JSON(nil, err)
				return
			}
			svc.DB.Debug().Model(&operation.Operation{}).Where(where).Count(&total)
		} else {
			if err = svc.DB.Debug().Model(&operation.Operation{}).Order("rank ASC").Offset((v.Pn-1)*v.Ps).Limit(v.Ps).Find(&ops, where, v.Platform).Error; err != nil {
				c.JSON(nil, err)
				return
			}
			svc.DB.Debug().Model(&operation.Operation{}).Where(where, v.Platform).Count(&total)
		}
	} else {
		if v.Platform == 100 {
			if err = svc.DB.Debug().Model(&operation.Operation{}).Order("rank ASC").Offset((v.Pn-1)*v.Ps).Limit(v.Ps).Find(&ops, where, now, now).Error; err != nil {
				c.JSON(nil, err)
				return
			}
			svc.DB.Debug().Model(&operation.Operation{}).Where(where, now, now).Count(&total)
		} else {
			if err = svc.DB.Debug().Model(&operation.Operation{}).Order("rank ASC").Offset((v.Pn-1)*v.Ps).Limit(v.Ps).Find(&ops, where, v.Platform, now, now).Error; err != nil {
				c.JSON(nil, err)
				return
			}
			svc.DB.Debug().Model(&operation.Operation{}).Where(where, v.Platform, now, now).Count(&total)
		}
	}
	var opsView []*operation.ViewOperation
	layout := "2006-01-02 15:04:05"
	for _, v := range ops {
		var (
			status   string
			timeNow  = time.Now()
			stime, _ = time.Parse(time.RFC3339, v.Stime)
			etime, _ = time.Parse(time.RFC3339, v.Etime)
			ctime, _ = time.Parse(time.RFC3339, v.Ctime)
			mtime, _ = time.Parse(time.RFC3339, v.Mtime)
			dtime, _ = time.Parse(time.RFC3339, v.Dtime)
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
			Stime:    stime.Format(layout),
			Etime:    etime.Format(layout),
			Ctime:    ctime.Format(layout),
			Dtime:    dtime.Format(layout),
			Mtime:    mtime.Format(layout),
			Status:   status,
		})
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    opsView,
		"pager": map[string]int{
			"page":     v.Pn,
			"pagesize": v.Ps,
			"total":    total,
		},
	}))
}

func addCollectArcOper(c *bm.Context) {
	var (
		err error
	)
	username, _ := c.Get("username")
	uname, ok := username.(string)
	if !ok || len(uname) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	v := new(struct {
		Ads      int8   `form:"ads" `
		Rank     int8   `form:"rank"`
		Pic      string `form:"pic" validate:"required"`
		Link     string `form:"link" validate:"required"`
		Content  string `form:"content" validate:"required"`
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
		Type:     "collect_arc",
		Ads:      v.Ads,
		Rank:     v.Rank,
		Pic:      v.Pic,
		Link:     v.Link,
		Content:  v.Content,
		Username: uname,
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

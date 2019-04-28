package http

import (
	"encoding/json"
	"fmt"
	"go-common/app/admin/main/creative/model/whitelist"
	accapi "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"net/http"
	"time"
)

var maxMid = int64(2147483647)

func viewWhiteList(c *bm.Context) {
	var (
		err  error
		info = &whitelist.Whitelist{}
	)
	if err = c.Bind(info); err != nil {
		return
	}
	if info.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).First(info, info.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if info.ID == 0 {
		c.JSON(nil, ecode.NothingFound)
	} else {
		pfl, _ := svc.ProfileStat(c, info.MID)
		if pfl.Profile != nil {
			info.Name = pfl.Profile.Name
		}
		info.Fans = pfl.Follower
		info.CurrentLevel = pfl.LevelInfo.Cur
		c.JSON(info, nil)
	}
}

func listWhiteList(c *bm.Context) {
	var (
		err   error
		wls   []*whitelist.Whitelist
		total int
	)
	v := new(struct {
		MID      int64 `form:"mid"`
		AdminMID int64 `form:"admin_mid"`
		State    int8  `form:"state"`
		Type     int8  `form:"type"`
		Page     int   `form:"pn"`
		PageSize int   `form:"ps"`
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
	where := " state = ? AND type = ? "
	if v.MID != 0 {
		where += " AND mid = ?"
		if v.AdminMID != 0 {
			where += " AND admin_id = ?"
			if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Order("ctime DESC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&wls, where, v.State, v.Type, v.MID, v.AdminMID).Error; err != nil {
				c.JSON(nil, err)
			}
			svc.DB.Debug().Model(&whitelist.Whitelist{}).Where(where, v.State, v.Type, v.MID, v.AdminMID).Count(&total)
		} else {
			if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Order("ctime DESC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&wls, where, v.State, v.Type, v.MID).Error; err != nil {
				c.JSON(nil, err)
			}
			svc.DB.Debug().Model(&whitelist.Whitelist{}).Where(where, v.State, v.Type, v.MID).Count(&total)
		}
	} else {
		if v.AdminMID != 0 {
			where += " AND admin_id = ?"
			if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Order("ctime DESC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&wls, where, v.State, v.Type, v.AdminMID).Error; err != nil {
				c.JSON(nil, err)
			}
			svc.DB.Debug().Model(&whitelist.Whitelist{}).Where(where, v.State, v.Type, v.AdminMID).Count(&total)
		} else {
			if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Order("ctime DESC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&wls, where, v.State, v.Type).Error; err != nil {
				c.JSON(nil, err)
			}
			svc.DB.Debug().Model(&whitelist.Whitelist{}).Where(where, v.State, v.Type).Count(&total)
		}
	}
	wls, err = svc.Cards(c, wls)
	if err != nil {
		log.Error("svc.Cards error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    wls,
		"pager": map[string]int{
			"page":     v.Page,
			"pagesize": v.PageSize,
			"total":    total,
		},
	}))
}

func exportWhiteList(c *bm.Context) {
	var (
		err   error
		wls   []*whitelist.Whitelist
		total int
	)
	v := new(struct {
		MID      int64 `form:"mid;" `
		AdminMID int64 `form:"admin_mid;" `
		State    int8  `form:"state"`
		Type     int8  `form:"type"`
		Page     int   `form:"pn"`
		PageSize int   `form:"ps"`
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
	where := " state = ? AND type = ? "
	if v.MID != 0 {
		where += " AND mid = ?"
		if v.AdminMID != 0 {
			where += " AND admin_id = ?"
			if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Order("ctime DESC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&wls, where, v.State, v.Type, v.MID, v.AdminMID).Error; err != nil {
				c.JSON(nil, err)
			}
			svc.DB.Debug().Model(&whitelist.Whitelist{}).Where(where, v.State, v.Type, v.MID, v.AdminMID).Count(&total)
		} else {
			if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Order("ctime DESC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&wls, where, v.State, v.Type, v.MID).Error; err != nil {
				c.JSON(nil, err)
			}
			svc.DB.Debug().Model(&whitelist.Whitelist{}).Where(where, v.State, v.Type, v.MID).Count(&total)
		}
	} else {
		if v.AdminMID != 0 {
			where += " AND admin_id = ?"
			if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Order("ctime DESC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&wls, where, v.State, v.Type, v.AdminMID).Error; err != nil {
				c.JSON(nil, err)
			}
			svc.DB.Debug().Model(&whitelist.Whitelist{}).Where(where, v.State, v.Type, v.AdminMID).Count(&total)
		} else {
			if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Order("ctime DESC").Offset((v.Page-1)*v.PageSize).Limit(v.PageSize).Find(&wls, where, v.State, v.Type).Error; err != nil {
				c.JSON(nil, err)
			}
			svc.DB.Debug().Model(&whitelist.Whitelist{}).Where(where, v.State, v.Type).Count(&total)
		}
	}
	wls, err = svc.Cards(c, wls)
	if err != nil {
		log.Error("svc.Cards error(%v)", err)
		c.JSON(nil, err)
		return
	}
	fWLS, err := formatWhilteList(wls)
	if err != nil {
		log.Error("formatWhilteList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: FormatCSV(fWLS),
		Title:   fmt.Sprintf("%s-%d-%d-%s", time.Now().Format("2006-01-02"), v.Type, v.State, "white_list"),
	})
}

func addWhiteList(c *bm.Context) {
	var (
		err error
		pfl *accapi.ProfileStatReply
	)
	v := new(struct {
		MID      int64  `form:"mid" validate:"required,min=1,gte=1"`
		AdminMID int64  `form:"admin_mid" validate:"required,min=1,gte=1"`
		Comment  string `form:"comment"`
		State    int8   `form:"state"`
		Type     int8   `form:"type"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.MID > maxMid || v.AdminMID > maxMid {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pfl, err = svc.ProfileStat(c, v.MID)
	if err != nil || pfl == nil || pfl.Profile == nil || pfl.Profile.Mid == 0 {
		log.Error("svc.Card zero result error(%v)", err)
		c.JSON(nil, ecode.UserNotExist)
		return
	}
	m := &whitelist.Whitelist{
		State:    v.State,
		Type:     v.Type,
		Ctime:    time.Now().Format("2006-01-02 15:04:05"),
		MID:      v.MID,
		AdminMID: v.AdminMID,
		Comment:  v.Comment,
	}
	db := svc.DB.Debug().Model(&whitelist.Whitelist{}).Create(m)
	if err = db.Error; err != nil {
		log.Error("creativeSvc.whitelist error(%v)", err)
		c.JSON(nil, err)
	}
	c.JSON(map[string]interface{}{
		"id": db.Value.(*whitelist.Whitelist).ID,
	}, nil)
}

func addBatchWhiteList(c *bm.Context) {
	var (
		err error
	)
	v := new(struct {
		Params string `form:"params" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	filterJSON := v.Params
	if len(filterJSON) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	type Filter struct {
		MID      int64  `json:"mid"`
		AdminMID int64  `json:"admin_mid"`
		Comment  string `json:"comment"`
		State    int8   `json:"state"`
		Type     int8   `json:"type"`
	}
	var filtersJSONData []*Filter
	if err = json.Unmarshal([]byte(filterJSON), &filtersJSONData); err != nil {
		err = ecode.RequestErr
		return
	}
	if len(filtersJSONData) == 0 {
		err = ecode.RequestErr
		return
	}
	db := svc.DB.Debug().Model(&whitelist.Whitelist{})
	for _, v := range filtersJSONData {
		m := &whitelist.Whitelist{
			State:    v.State,
			Type:     v.Type,
			Ctime:    time.Now().Format("2006-01-02 15:04:05"),
			MID:      v.MID,
			AdminMID: v.AdminMID,
			Comment:  v.Comment,
		}
		db.Create(m)
		if err = db.Error; err != nil {
			log.Error("creativeSvc.batchWhitelist error(%v)", err)
			c.JSON(nil, err)
			return
		}
	}
	c.JSON(nil, err)
}

func upWhiteList(c *bm.Context) {
	var (
		wl  = &whitelist.Whitelist{}
		pfl *accapi.ProfileStatReply
		err error
	)
	v := new(struct {
		ID       int64  `form:"id"`
		MID      int64  `form:"mid" validate:"required,min=1,gte=1"`
		AdminMID int64  `form:"admin_mid" validate:"required,min=1,gte=1"`
		Comment  string `form:"comment"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 || v.MID > maxMid || v.AdminMID > maxMid {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	pfl, err = svc.ProfileStat(c, v.MID)
	if err != nil || pfl == nil || pfl.Profile == nil || pfl.Profile.Mid == 0 {
		log.Error("svc.Card zero result error(%v)", err)
		c.JSON(nil, ecode.UserNotExist)
		return
	}
	if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Find(wl, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = svc.DB.Debug().Model(&whitelist.Whitelist{ID: v.ID}).Update(map[string]interface {
	}{
		"id":        v.ID,
		"mid":       v.MID,
		"admin_mid": v.AdminMID,
		"comment":   v.Comment,
		"mtime":     time.Now().Format("2006-01-02 15:04:05"),
	}).Error; err != nil {
		log.Error("svc.save error(%v)", err)
	}
	c.JSON(nil, err)
}

func delWhiteList(c *bm.Context) {
	var (
		wl  = &whitelist.Whitelist{}
		err error
	)
	if err = c.Bind(wl); err != nil {
		return
	}
	if wl.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Debug().Model(&whitelist.Whitelist{}).Find(wl, wl.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = svc.DB.Debug().Model(&whitelist.Whitelist{ID: wl.ID}).Update(map[string]interface {
	}{
		"id":    wl.ID,
		"State": 0,
		"mtime": time.Now().Format("2006-01-02 15:04:05"),
	}).Error; err != nil {
		log.Error("svc.del WhiteList error(%v)", err)
	}
	c.JSON(nil, err)
}

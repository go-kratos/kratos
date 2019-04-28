package http

import (
	"encoding/json"
	"go-common/app/admin/main/creative/model/app"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"net/http"
	"sort"
	"time"
)

func viewPortal(c *bm.Context) {
	var (
		err  error
		info = &app.Portal{}
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
	if err = svc.DB.Debug().Model(&app.Portal{}).First(info, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if info.ID == 0 {
		c.JSON(nil, ecode.NothingFound)
	} else {
		layout := "2006-01-02 15:04:05"
		mtime, _ := time.Parse(time.RFC3339, info.MTime)
		info.MTime = mtime.Format(layout)
		ptime, _ := time.Parse(time.RFC3339, info.PTime)
		info.PTime = ptime.Format(layout)
		ctime, _ := time.Parse(time.RFC3339, info.CTime)
		info.CTime = ctime.Format(layout)
		c.JSON(info, nil)
	}
}
func addPortal(c *bm.Context) {
	var (
		err error
	)
	v := new(struct {
		Build    int64  `form:"build"`
		BuildExp string `form:"buildexp"`
		Platform int8   `form:"platform"`
		Compare  int8   `form:"compare"`
		Pos      int16  `form:"pos"`
		Mark     int8   `form:"mark"`
		More     int8   `form:"more"`
		Type     int8   `form:"type"`
		Title    string `form:"title" validate:"required"`
		Icon     string `form:"icon" validate:"required"`
		URL      string `form:"url" validate:"required"`
		Whiteexp string `form:"whiteexp"`
		SubTitle string `form:"subtitle"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if !checkBuild(v.BuildExp) {
		httpCode(c, "buildexp is wrong", ecode.RequestErr)
		return
	}
	if !checkWhite(v.Whiteexp) {
		httpCode(c, "whiteexp is wrong", ecode.RequestErr)
		return
	}
	m := &app.Portal{
		Build:    v.Build,
		BuildExp: v.BuildExp,
		Platform: v.Platform,
		Compare:  v.Compare,
		State:    1,
		Pos:      v.Pos,
		Mark:     v.Mark,
		More:     v.More,
		Type:     v.Type,
		Title:    v.Title,
		Icon:     v.Icon,
		URL:      v.URL,
		CTime:    time.Now().Format("2006-01-02 15:04:05"),
		MTime:    time.Now().Format("2006-01-02 15:04:05"),
		PTime:    time.Now().Format("2006-01-02 15:04:05"),
		SubTitle: v.SubTitle,
		WhiteExp: v.Whiteexp,
	}
	db := svc.DB.Create(m)
	if err = db.Error; err != nil {
		log.Error("creativeSvc.addPortal error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(map[string]interface{}{
		"id": db.Value.(*app.Portal).ID,
	}, nil)
}

func upPortal(c *bm.Context) {
	var (
		err error
		ap  = &app.Portal{}
	)
	v := new(struct {
		ID       int64  `form:"id"`
		Build    int64  `form:"build"`
		BuildExp string `form:"buildexp"`
		Platform int8   `form:"platform"`
		Compare  int8   `form:"compare"`
		Pos      int16  `form:"pos"`
		Mark     int8   `form:"mark"`
		More     int8   `form:"more"`
		Type     int8   `form:"type"`
		Title    string `form:"title" validate:"required"`
		Icon     string `form:"icon" validate:"required"`
		URL      string `form:"url" validate:"required"`
		Whiteexp string `form:"whiteexp"`
		SubTitle string `form:"subtitle"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if !checkBuild(v.BuildExp) {
		httpCode(c, "buildexp is wrong", ecode.RequestErr)
		return
	}
	if !checkWhite(v.Whiteexp) {
		httpCode(c, "whiteexp is wrong", ecode.RequestErr)
		return
	}
	if err = svc.DB.Model(ap).Find(ap, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = svc.DB.Model(&app.Portal{}).Update(map[string]interface{}{
		"id":       v.ID,
		"build":    v.Build,
		"buildexp": v.BuildExp,
		"platform": v.Platform,
		"compare":  v.Compare,
		"state":    1,
		"pos":      v.Pos,
		"mark":     v.Mark,
		"more":     v.More,
		"type":     v.Type,
		"title":    v.Title,
		"icon":     v.Icon,
		"url":      v.URL,
		"mtime":    time.Now().Format("2006-01-02 15:04:05"),
		"subtitle": v.SubTitle,
		"whiteexp": v.Whiteexp,
	}).Error; err != nil {
		log.Error("creativeSvc.addPortal error(%v)", err)
	}
	c.JSON(nil, err)
}

func downPortal(c *bm.Context) {
	var (
		err error
		ap  = &app.Portal{}
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
	if err = svc.DB.Model(ap).Find(ap, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = svc.DB.Model(&app.Portal{ID: v.ID}).Update(map[string]interface{}{
		"id":    v.ID,
		"state": 0,
		"mtime": time.Now().Format("2006-01-02 15:04:05"),
	}).Error; err != nil {
		log.Error("creativeSvc.addPortal error(%v)", err)
	}
	c.JSON(nil, err)
}

// 入口列表查询增加subtitle和whiteexp字段
func portalList(c *bm.Context) {
	var (
		pts   []*app.Portal
		items []*app.Item
		total int64
		err   error
	)
	v := new(struct {
		Pn       int    `form:"pn" validate:"required,min=1"`
		Ps       int    `form:"ps" validate:"required,min=1"`
		State    int8   `form:"state"`
		Type     int8   `form:"type"`
		Platform string `form:"platform"`
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
	if v.Type < 0 {
		v.Type = 0
	}
	db := svc.DB.Model(&app.Portal{}).Where("type=?", v.Type)
	if v.Platform != "" {
		db = db.Where("platform=?", atoi(v.Platform))
	}
	if v.State != 2 {
		db = db.Where("state=?", v.State)
	}
	db.Count(&total)
	if err = db.Order("ctime DESC").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&pts).Error; err != nil {
		log.Error("portalList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if len(pts) > 0 {
		items = make([]*app.Item, 0, len(pts))
		for _, v := range pts {
			i := &app.Item{}
			if v == nil {
				continue
			}
			i.ID = v.ID
			i.Build = v.Build
			i.BuildExp = v.BuildExp
			i.Platform = v.Platform
			i.Compare = v.Compare
			i.State = v.State
			i.Pos = v.Pos
			i.Mark = v.Mark
			i.More = v.More
			i.Type = v.Type
			i.Title = v.Title
			i.Icon = v.Icon
			i.URL = v.URL
			i.SubTitle = v.SubTitle
			if len(v.WhiteExp) > 0 {
				if err = json.Unmarshal([]byte(v.WhiteExp), &i.WhiteExps); err != nil {
					log.Error("json.Unmarshal buildComps failed error(%v)", err)
				}
				sort.Slice(i.WhiteExps, func(m, n int) bool {
					return i.WhiteExps[m].TP <= i.WhiteExps[n].TP
				})
			}
			ct, err := time.Parse(time.RFC3339, v.CTime)
			if err != nil {
				log.Error("ctime time.Parse error(%v)", err)
			}
			mt, err := time.Parse(time.RFC3339, v.MTime)
			if err != nil {
				log.Error("mtime time.Parse error(%v)", err)
			}
			pt, err := time.Parse(time.RFC3339, v.PTime)
			if err != nil {
				log.Error("ptime time.Parse error(%v)", err)
			}
			i.CTime = ct.Unix()
			i.MTime = mt.Unix()
			i.PTime = pt.Unix()
			items = append(items, i)
		}
	} else {
		items = []*app.Item{}
	}
	data := &app.PortalPager{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: items,
		Total: total,
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    data,
	}))
}

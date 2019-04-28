package http

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/history/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// history return the user history for mobile app.
func history(c *bm.Context) {
	var (
		err error
		mid int64
		v   = new(Histroy)
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if tp, ok := business(c); !ok {
		return
	} else if tp > 0 {
		v.TP = tp
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Pn < 1 {
		v.Pn = 1
	}
	if v.Ps > cnf.History.Max || v.Ps <= 0 {
		v.Ps = cnf.History.Max
	}
	list, err := hisSvc.Videos(c, mid, v.Pn, v.Ps, v.TP)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(list, nil)
}

// aids return the user histories.
func aids(c *bm.Context) {
	var (
		mid int64
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	list, err := hisSvc.AVHistories(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(list, nil)
}

func managerHistory(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Mid    int64 `form:"mid" validate:"required,gt=0"`
			OnlyAV bool  `form:"only_av"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	history, err := hisSvc.ManagerHistory(c, v.OnlyAV, v.Mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(history, nil)
}

// clearHistory clear the user histories.
func clearHistory(c *bm.Context) {
	var (
		err error
		mid int64
		v   = new(struct {
			TP int8 `form:"type"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if tp, ok := business(c); !ok {
		return
	} else if tp > 0 {
		v.TP = tp
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var tps []int8
	if v.TP > 0 {
		tps = append(tps, v.TP)
	}
	c.JSON(nil, hisSvc.ClearHistory(c, mid, tps))
}

// delHistory delete the user history by aid.
func delHistory(c *bm.Context) {
	var (
		err error
		mid int64
		v   = new(struct {
			TP   int8    `form:"type"`
			Aids []int64 `form:"aid,split"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if tp, ok := business(c); !ok {
		return
	} else if tp > 0 {
		v.TP = tp
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, hisSvc.DelHistory(c, mid, v.Aids, v.TP))
}

func delete(c *bm.Context) {
	var (
		err error
		mid int64
		v   = new(struct {
			Hid []string `form:"hid,split"`
			Bid []string `form:"bid,split"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var his []*model.History
	if len(v.Bid) == 0 {
		for _, hid := range v.Hid {
			hs := strings.Split(hid, "_")
			if len(hs) == 0 {
				continue
			}
			aid, _ := strconv.ParseInt(hs[0], 10, 0)
			if aid == 0 {
				continue
			}
			var tp int64
			if len(hs) == 2 {
				tp, _ = strconv.ParseInt(hs[1], 10, 0)
			}
			his = append(his, &model.History{
				Aid: aid,
				TP:  int8(tp),
			})
		}
	}
	for _, bid := range v.Bid {
		bs := strings.Split(bid, "_")
		if len(bs) != 2 {
			continue
		}
		aid, _ := strconv.ParseInt(bs[1], 10, 0)
		if aid == 0 {
			continue
		}
		tp, err := model.CheckBusiness(bs[0])
		if err != nil {
			continue
		}
		his = append(his, &model.History{
			Aid: aid,
			TP:  int8(tp),
		})
	}
	if len(his) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, hisSvc.Delete(c, mid, his))
}

// addHistory add history into user redis set.
func addHistory(c *bm.Context) {
	var (
		err error
		mid int64
		h   *model.History
		now = time.Now().Unix()
		v   = new(AddHistory)
	)
	// sid cid tp, dt
	// dt:devece type , sid :season_id,type:video type
	// aid :aid
	// cid:cid,epid
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.TP == model.TypeBangumi || v.TP == model.TypeMovie {
		if v.Sid <= 0 || v.Epid <= 0 {
			if v.Aid <= 0 {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		}
	}
	if v.Aid <= 0 || (v.TP < model.TypeArticle && v.Cid <= 0) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.DT <= 0 {
		if v.Platform == model.PlatformIOS {
			v.DT = model.DeviceIphone
			if v.Device == model.DevicePad {
				v.DT = model.DeviceIpad
			}
		} else if v.Platform == model.PlatformAndroid {
			v.DT = model.DeviceAndroid
		}
	}
	h = &model.History{
		Aid:  v.Aid,
		Unix: now,
		Sid:  v.Sid,
		Epid: v.Epid,
		Cid:  v.Cid,
		TP:   v.TP,
		STP:  v.SubTP,
		DT:   v.DT,
	}
	h.ConvertType()
	c.JSON(nil, hisSvc.AddHistory(c, mid, 0, h))
}

// report report view progress.
func report(c *bm.Context) {
	var (
		err error
		mid int64
		h   *model.History
		now = time.Now().Unix()
		v   = new(HistoryReport)
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if tp, ok := business(c); !ok {
		return
	} else if tp > 0 {
		v.Type = tp
	}
	if v.Progress < 0 {
		v.Progress = model.ProComplete
	}
	if v.Aid <= 0 || (v.Type < model.TypeArticle && v.Cid <= 0) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.SubType <= 0 {
		v.SubType = v.SubTP
	}
	if v.Type == model.TypeBangumi || v.Type == model.TypeMovie {
		if v.Sid <= 0 || v.Epid <= 0 {
			v.Type = model.TypeUGC
		} else {
			v.Type = model.TypePGC
		}
	}
	if v.DT <= 0 {
		if v.Platform == model.PlatformIOS {
			v.DT = model.DeviceIphone
			if v.Device == model.DevicePad {
				v.DT = model.DeviceIpad
			}
		} else if v.Platform == model.PlatformAndroid {
			v.DT = model.DeviceAndroid
			if v.MobileApp == model.MobileAppAndroidTV {
				v.DT = model.DeviceAndroidTV
			}
		}
	}
	h = &model.History{
		Aid:  v.Aid,
		Unix: now,
		Sid:  v.Sid,
		Epid: v.Epid,
		Cid:  v.Cid,
		Pro:  v.Progress,
		TP:   v.Type,
		STP:  v.SubType,
		DT:   v.DT,
	}
	c.JSON(nil, hisSvc.AddHistory(c, mid, v.Realtime, h))
}

// report report view progress.
func innerReport(c *bm.Context) {
	var (
		err error
		h   *model.History
		v   = new(HistoryReport)
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if tp, ok := business(c); !ok {
		return
	} else if tp > 0 {
		v.Type = tp
	}
	if v.Mid == 0 && v.Aid == 0 && v.Type == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Type == model.TypeComic && v.Aid == 0 {
		v.Aid = v.Sid
	}
	if v.PlayTime == 0 {
		v.PlayTime = time.Now().Unix()
	}
	if v.Progress < 0 {
		v.Progress = model.ProComplete
	}
	h = &model.History{
		Aid:  v.Aid,
		Unix: v.PlayTime,
		Sid:  v.Sid,
		Epid: v.Epid,
		Cid:  v.Cid,
		Pro:  v.Progress,
		TP:   v.Type,
		STP:  v.SubType,
		DT:   v.DT,
	}
	c.JSON(nil, hisSvc.AddHistory(c, v.Mid, v.Realtime, h))
}

// reports
func reports(c *bm.Context) {
	var (
		err error
		mid int64
		hs  = make([]*model.History, 0)
		v   = new(struct {
			Type int8   `form:"type"`
			Data string `form:"data"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if tp, ok := business(c); !ok {
		return
	} else if tp > 0 {
		v.Type = tp
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(v.Data), &hs); err != nil {
		log.Error("json.Unmarshal(%s),err:%v.", v.Data, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = hisSvc.AddHistories(c, mid, v.Type, metadata.String(c, metadata.RemoteIP), hs); err != nil {
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(nil, nil)
}

// shadow return the user shadow status.
func shadow(c *bm.Context) {
	var (
		err error
		mid int64
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	status, err := hisSvc.Shadow(c, mid)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(status == model.ShadowOn, nil)
}

// setShadow the user shadow status.
func setShadow(c *bm.Context) {
	var (
		err    error
		mid    int64
		shadow = model.ShadowOff
		v      = new(struct {
			Switch bool `form:"switch"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Switch {
		shadow = model.ShadowOn
	}
	c.JSON(nil, hisSvc.SetShadow(c, mid, shadow))
}

// flush flush users hisotry.
func flush(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Mids  []int64 `form:"mids,split" validate:"required,min=1,dive,gt=0"`
			STime int64   `form:"time"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, hisSvc.FlushHistory(c, v.Mids, v.STime))
}

// position report report view progress.
func position(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Mid int64 `form:"mid"`
			Aid int64 `form:"aid"`
			TP  int8  `form:"type"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if tp, ok := business(c); !ok {
		return
	} else if tp > 0 {
		v.TP = tp
	}
	if v.Mid == 0 && v.Aid == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(hisSvc.Position(c, v.Mid, v.Aid, v.TP))
}

// position report report view progress.
func resource(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Mid int64 `form:"mid"`
			TP  int8  `form:"type"`
			Pn  int   `form:"pn"`
			Ps  int   `form:"ps"`
		})
	)
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, ok := business(c); !ok {
		return
	} else if tp > 0 {
		v.TP = tp
	}
	if v.Mid == 0 || v.TP == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(hisSvc.Histories(c, v.Mid, v.TP, v.Pn, v.Ps))
}

// position report report view progress.
func resources(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Mid  int64   `form:"mid"`
			TP   int8    `form:"type"`
			Aids []int64 `form:"aids,split"`
		})
	)
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tp, ok := business(c); !ok {
		return
	} else if tp > 0 {
		v.TP = tp
	}
	if v.Mid == 0 || v.TP == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(v.Aids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(hisSvc.HistoryType(c, v.Mid, v.TP, v.Aids, metadata.String(c, metadata.RemoteIP)))
}

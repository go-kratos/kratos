package service

import (
	"strconv"
	"time"

	"encoding/json"
	"go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
	binfoc "go-common/library/log/infoc"
	"go-common/library/stat/prom"
)

type displayInfo struct {
	ip       string
	mid      string
	now      string
	client   string
	build    string
	buvid    string
	pagetype string
	pageNo   string
	isRec    string
	showlist json.RawMessage
}

type clickInfo struct {
	mid      string
	client   string
	build    string
	buvid    string
	now      string
	from     string
	itemID   string
	itemType string
	extra    json.RawMessage
}

type aiClickInfo struct {
	mid      string
	client   string
	build    string
	buvid    string
	time     string
	from     string
	itemID   string
	itemType string
	actionID string
	action   string
	extra    json.RawMessage
}

// 用户在列表页停留上报
type showInfo struct {
	ip       string
	time     string
	buvid    string
	mid      string
	client   string
	pageType string
	from     string
	build    string
	extra    string
}

type recItem struct {
	ID          int64  `json:"id"`
	Page        int64  `json:"page"`
	Pos         int64  `json:"pos"`
	View        int64  `json:"view"`
	Fav         int64  `json:"fav"`
	Like        int64  `json:"like"`
	Reply       int64  `json:"reply"`
	Share       int64  `json:"share"`
	AvFeature   string `json:"av_feature,omitempty"`
	UserFeature string `json:"user_feature,omitempty"`
}

// RecommendInfoc .
func (s *Service) RecommendInfoc(mid int64, plat int8, pageType, cid, build int, buvid, ip string, metas []*model.Meta, isRcmd bool, now time.Time, pn int64, sky *model.SkyHorseResp) {
	var isRc = "0"
	if isRcmd {
		isRc = "1"
	}
	skyMap := make(map[int64]string)
	if sky != nil {
		for _, item := range sky.Data {
			skyMap[item.ID] = item.AvFeature
		}
	}
	var list []*recItem
	for i, m := range metas {
		x := &recItem{ID: m.ID, Page: pn, Pos: int64(i + 1), View: m.Stats.View, Fav: m.Stats.Favorite, Like: m.Stats.Like, Reply: m.Stats.Reply, Share: m.Stats.Share}
		if sky != nil && sky.UserFeature != "" {
			x.UserFeature = sky.UserFeature
		}
		x.AvFeature = skyMap[m.ID]
		list = append(list, x)
	}
	var sl = &struct {
		List []*recItem `json:"itemlist"`
	}{
		List: list,
	}
	msg, _ := json.Marshal(sl)
	s.infoc(displayInfo{ip, strconv.FormatInt(mid, 10), strconv.FormatInt(now.Unix(), 10), strconv.Itoa(int(plat)), strconv.Itoa(build), buvid, strconv.Itoa(pageType), strconv.Itoa(cid), isRc, msg})
}

// ViewInfoc .
func (s *Service) ViewInfoc(mid int64, plat int8, build int, itemType, from, buvid string, itemID int64, now time.Time, ua string) {
	var extra = &struct {
		UA string `json:"ua"`
	}{
		UA: ua,
	}
	msg, _ := json.Marshal(extra)
	s.infoc(clickInfo{strconv.FormatInt(mid, 10), strconv.Itoa(int(plat)), strconv.Itoa(build), buvid, strconv.FormatInt(now.Unix(), 10), from, strconv.FormatInt(itemID, 10), itemType, msg})
}

// AIViewInfoc .
func (s *Service) AIViewInfoc(mid int64, plat int8, build int, itemType, from, buvid string, itemID int64, now time.Time, ua string) {
	var extra = &struct {
		UA string `json:"ua"`
	}{
		UA: ua,
	}
	msg, _ := json.Marshal(extra)
	s.infoc(aiClickInfo{
		mid:      strconv.FormatInt(mid, 10),
		client:   model.Client(plat),
		build:    strconv.Itoa(build),
		buvid:    buvid,
		time:     strconv.FormatInt(now.Unix(), 10),
		from:     from,
		itemID:   strconv.FormatInt(itemID, 10),
		itemType: itemType,
		action:   "click",
		actionID: "",
		extra:    msg,
	})
}

// ShowInfoc .
func (s *Service) ShowInfoc(ip string, now time.Time, buvid string, mid int64, client int8, pageType string, from string, build string, ua string, referer string) {
	var extra = &struct {
		UA      string `json:"ua"`
		Referer string `json:"referer"`
	}{
		UA:      ua,
		Referer: referer,
	}
	msg, _ := json.Marshal(extra)
	s.infoc(showInfo{
		ip:       ip,
		time:     strconv.FormatInt(now.Unix(), 10),
		buvid:    buvid,
		mid:      strconv.FormatInt(mid, 10),
		client:   strconv.Itoa(int(client)),
		pageType: pageType,
		from:     from,
		build:    build,
		extra:    string(msg),
	})
}

func (s *Service) infoc(i interface{}) {
	select {
	case s.logCh <- i:
	default:
		log.Warn("infocproc chan full")
	}
}

// writeInfoc
func (s *Service) infocproc() {
	var (
		displayInfoc = binfoc.New(s.c.DisplayInfoc)
		clickInfoc   = binfoc.New(s.c.ClickInfoc)
		aiClickInfoc = binfoc.New(s.c.AIClickInfoc)
		showInfoc    = binfoc.New(s.c.ShowInfoc)
	)
	for {
		i, ok := <-s.logCh
		if !ok {
			log.Warn("infoc proc exit")
			return
		}
		prom.BusinessInfoCount.State("infoc_channel", int64(len(s.logCh)))
		switch l := i.(type) {
		case displayInfo:
			displayInfoc.Info(l.ip, l.now, l.buvid, l.mid, l.client, l.pagetype, l.pageNo, string(l.showlist), l.isRec, l.build)
			log.Info("infocproc displayInfo param(ip:%s,now:%s,buvid:%s,mid:%s,client:%s,pagetype:%s，pageno:%s,showlist:%s,isRec:%s,build:%s)", l.ip, l.now, l.buvid, l.mid, l.client, l.pagetype, l.pageNo, l.showlist, l.isRec, l.build)
		case clickInfo:
			clickInfoc.Info(l.from, l.now, l.buvid, l.mid, l.client, l.itemType, l.itemID, "", l.build)
			log.Info("infocproc clickInfoc param(client:%s,buvid:%s,mid:%s,now:%s,from:%s,build:%s,itemID:%s,itemType:%s)", l.client, l.buvid, l.mid, l.now, l.from, l.build, l.itemID, l.itemType)
		case aiClickInfo:
			aiClickInfoc.Info(l.client, l.buvid, l.mid, l.time, l.from, l.build, l.itemID, l.itemType, l.action, l.actionID, string(l.extra))
			log.Info("infocproc aiclickInfoc param(client:%s,buvid:%s,mid:%s,time:%s,from:%s,build:%s,itemID:%s,itemType:%s,action: %s, actionID: %s, extra: %s)", l.client, l.buvid, l.mid, l.time, l.from, l.build, l.itemID, l.itemType, l.action, l.actionID, string(l.extra))
		case showInfo:
			showInfoc.Info(l.ip, l.time, l.buvid, l.mid, l.client, l.pageType, l.from, l.build, l.extra)
			log.Info("infocproc showInfoc param(%+v)", l)
		}
	}
}

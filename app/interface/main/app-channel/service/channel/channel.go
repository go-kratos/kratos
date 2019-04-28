package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	cardm "go-common/app/interface/main/app-card/model/card"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-channel/model"
	"go-common/app/interface/main/app-channel/model/card"
	"go-common/app/interface/main/app-channel/model/channel"
	"go-common/app/interface/main/app-channel/model/tab"
	tag "go-common/app/interface/main/tag/model"
	locmdl "go-common/app/service/main/location/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	"github.com/dgryski/go-farm"
)

const (
	_initRegionKey = "region_key_%d_%v"
	_initlanguage  = "hans"
	_initVersion   = "region_version"
	_regionRepeat  = "r_%d_%d"
	_maxAtten      = 10 //展示最多10个我的订阅
)

var (
	_tabList = []*channel.TabList{
		&channel.TabList{
			Name:  "推荐",
			URI:   "bilibili://pegasus/channel/feed/%d",
			TabID: "multiple",
		},
		&channel.TabList{
			Name:  "话题",
			URI:   "bilibili://following/topic_detail?id=%d&name=%s",
			TabID: "topic",
		},
	}
)

// Tab channel tab
func (s *Service) Tab(c context.Context, tid, mid int64, tname string, plat int8) (res *channel.Tab, err error) {
	var (
		t *tag.ChannelDetail
	)
	if t, err = s.tg.ChannelDetail(c, mid, tid, tname, s.isOverseas(plat)); err != nil || t == nil {
		log.Error("s.tag.ChannelDetail(%d, %d, %v) error(%v)", mid, tid, tname, err)
		return
	}
	res = &channel.Tab{}
	res.SimilarTagChange(t)
	res.TabList = s.tablist(t)
	return
}

//SubscribeAdd subscribe add
func (s *Service) SubscribeAdd(c context.Context, mid, id int64, now time.Time) (err error) {
	if err = s.tg.SubscribeAdd(c, mid, id, now); err != nil {
		log.Error("s.tg.SubscribeAdd(%d,%d) error(%v)", mid, id, err)
		return
	}
	return
}

//SubscribeCancel subscribe channel
func (s *Service) SubscribeCancel(c context.Context, mid, id int64, now time.Time) (err error) {
	if err = s.tg.SubscribeCancel(c, mid, id, now); err != nil {
		log.Error("s.tg.SubscribeCancel(%d,%d) error(%v)", mid, id, err)
		return
	}
	return
}

// SubscribeUpdate subscribe update
func (s *Service) SubscribeUpdate(c context.Context, mid int64, ids string) (err error) {
	if err = s.tg.SubscribeUpdate(c, mid, ids); err != nil {
		log.Error("s.tg.SubscribeUpdate(%d,%s) error(%v)", mid, ids, err)
		return
	}
	return
}

// List 频道tab页
func (s *Service) List(c context.Context, mid int64, plat int8, build, limit int, ver, mobiApp, device, lang string) (res *channel.List, err error) {
	var (
		rec, atten  []*channel.Channel
		top, bottom []*channel.Region
		max         = 3
	)
	g, _ := errgroup.WithContext(c)
	//获取推荐的三个频道
	g.Go(func() (err error) {
		rec, err = s.Recommend(c, mid, plat)
		if err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	//获取我的订阅
	if mid > 0 {
		g.Go(func() (err error) {
			atten, err = s.Subscribe(c, mid, limit)
			if err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	//获取分区
	g.Go(func() (err error) {
		top, bottom, _, err = s.RegionList(c, plat, build, mobiApp, device, lang)
		if err != nil {
			log.Error("%+v", err)
			err = nil
		}
		return
	})
	g.Wait()
	if tl := len(rec); tl < max {
		if last := max - tl; len(atten) > last {
			rec = append(rec, atten[:last]...)
		} else {
			rec = append(rec, atten...)
		}
	} else {
		rec = rec[:max]
	}
	res = &channel.List{
		RegionTop:    top,
		RegionBottom: bottom,
	}
	if isAudit := s.auditList(mobiApp, plat, build); !isAudit {
		res.RecChannel = rec
		res.AttenChannel = atten
	}
	res.Ver = s.hash(res)
	return
}

// Recommend 推荐
func (s *Service) Recommend(c context.Context, mid int64, plat int8) (res []*channel.Channel, err error) {
	list, err := s.tg.Discover(c, mid, s.isOverseas(plat))
	if err != nil {
		log.Error("%+v", err)
		return
	}
	for _, chann := range list {
		item := &channel.Channel{
			ID:      chann.ID,
			Name:    chann.Name,
			Cover:   chann.Cover,
			IsAtten: chann.Attention,
			Atten:   chann.Sub,
		}
		res = append(res, item)
	}
	return
}

//Subscribe 我订阅的tag（老） standard放前面用户自定义custom放后面
func (s *Service) Subscribe(c context.Context, mid int64, limit int) (res []*channel.Channel, err error) {
	var (
		tinfo []*tag.TagInfo
	)
	list, err := s.tg.Subscribe(c, mid)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	tinfo = list.Standard
	tinfo = append(tinfo, list.Custom...)
	for _, chann := range tinfo {
		item := &channel.Channel{
			ID:      chann.ID,
			Name:    chann.Name,
			Cover:   chann.Cover,
			Atten:   chann.Sub,
			IsAtten: chann.Attention,
			Content: chann.Content,
		}
		res = append(res, item)
	}
	if len(res) > limit && limit > 0 {
		res = res[:limit]
	} else if len(res) == 0 {
		res = []*channel.Channel{}
	}
	return
}

// Discover 发现频道页（推荐走recommend接口，有分类的揍list接口）
func (s *Service) Discover(c context.Context, id, mid int64, plat int8) (res []*channel.Channel, err error) {
	var (
		list []*tag.Channel
	)
	if id > 0 {
		list, err = s.tg.ListByCategory(c, id, mid, s.isOverseas(plat))
		if err != nil {
			log.Error("%+v", err)
			return
		}
	} else {
		list, err = s.tg.Recommend(c, mid, s.isOverseas(plat))
		if err != nil {
			log.Error("%+v", err)
			return
		}
	}
	if len(list) == 0 {
		res = []*channel.Channel{}
		return
	}
	for _, chann := range list {
		item := &channel.Channel{
			ID:      chann.ID,
			Name:    chann.Name,
			Cover:   chann.Cover,
			Atten:   chann.Sub,
			IsAtten: chann.Attention,
			Content: chann.Content,
		}
		res = append(res, item)
	}
	return
}

// Category 频道分类
func (s *Service) Category(c context.Context, plat int8) (res []*channel.Category, err error) {
	category, err := s.tg.Category(c, s.isOverseas(plat))
	if err != nil {
		log.Error("%+v", err)
		return
	}
	res = append(res, &channel.Category{
		ID:   0,
		Name: "推荐",
	})
	for _, cat := range category {
		item := &channel.Category{
			ID:   cat.ID,
			Name: cat.Name,
		}
		res = append(res, item)
	}
	return
}

// RegionList 分区信息
func (s *Service) RegionList(c context.Context, plat int8, build int, mobiApp, device, lang string) (regionTop, regionBottom, regions []*channel.Region, err error) {
	var (
		hantlanguage = "hant"
	)
	if ok := model.IsOverseas(plat); ok && lang != _initlanguage && lang != hantlanguage {
		lang = hantlanguage
	} else if lang == "" {
		lang = _initlanguage
	}
	var (
		rs = s.cachelist[fmt.Sprintf(_initRegionKey, plat, lang)]
		// maxTop = 8
		ridtmp = map[string]struct{}{}
		pids   []string
		auths  map[string]*locmdl.Auth
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	regionTop = []*channel.Region{}
	regionBottom = []*channel.Region{}
	regions = []*channel.Region{}
	for _, rtmp := range rs {
		if rtmp.ReID != 0 { //过滤二级分区
			continue
		}
		if rtmp.Area != "" {
			pids = append(pids, rtmp.Area)
		}
	}
	if len(pids) > 0 {
		auths, _ = s.loc.AuthPIDs(c, strings.Join(pids, ","), ip)
	}
LOOP:
	for _, rtmp := range rs {
		r := &channel.Region{}
		*r = *rtmp
		if r.ReID != 0 { //过滤二级分区
			continue
		}
		var tmpl, limitshow bool
		if limit, ok := s.limitCache[r.ID]; ok {
			for i, l := range s.limitCache[r.ID] {
				if i+1 <= len(limit)-1 {
					if ((l.Condition == "gt" && limit[i+1].Condition == "lt") && (l.Build < limit[i+1].Build)) ||
						((l.Condition == "lt" && limit[i+1].Condition == "gt") && (l.Build > limit[i+1].Build)) {
						if (l.Condition == "gt" && limit[i+1].Condition == "lt") &&
							(build > l.Build && build < limit[i+1].Build) {
							break
						} else if (l.Condition == "lt" && limit[i+1].Condition == "gt") &&
							(build < l.Build && build > limit[i+1].Build) {
							break
						} else {
							tmpl = true
							continue
						}
					}
				}
				if tmpl {
					if i == len(limit)-1 {
						limitshow = true
						break
						// continue LOOP
					}
					tmpl = false
					continue
				}
				if model.InvalidBuild(build, l.Build, l.Condition) {
					limitshow = true
					continue
					// continue LOOP
				} else {
					limitshow = false
					break
				}
			}
		}
		if limitshow {
			continue LOOP
		}
		if r.RID == 65539 {
			if model.IsIOS(plat) {
				r.URI = fmt.Sprintf("%s?from=category", r.URI)
			} else {
				r.URI = fmt.Sprintf("%s?sourceFrom=541", r.URI)
			}
		}
		if auth, ok := auths[r.Area]; ok && auth.Play == locmdl.Forbidden {
			log.Warn("s.invalid area(%v) ip(%v) error(%v)", r.Area, ip, err)
			continue
		}
		if isAudit := s.auditRegion(mobiApp, plat, build, r.RID); isAudit {
			continue
		}
		config, ok := s.configCache[r.ID]
		if !ok {
			continue
		}
		key := fmt.Sprintf(_regionRepeat, r.RID, r.ReID)
		if _, ok := ridtmp[key]; !ok {
			ridtmp[key] = struct{}{}
		} else {
			continue
		}
		for _, conf := range config {
			if conf.ScenesID == 1 /*&& len(regionTop) < maxTop*/ {
				regionTop = append(regionTop, r)
				regions = append(regions, r)
			} else if conf.ScenesID == 0 {
				regionBottom = append(regionBottom, r)
				regions = append(regions, r)
			}
		}
	}
	return
}

func (s *Service) hash(v *channel.List) string {
	bs, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return _initVersion
	}
	return strconv.FormatUint(farm.Hash64(bs), 10)
}

func (s *Service) loadRegionlist() {
	res, err := s.rg.AllList(context.TODO())
	if err != nil {
		log.Error("s.dao.All error(%v)", err)
		return
	}
	tmp := map[string][]*channel.Region{}
	for _, v := range res {
		key := fmt.Sprintf(_initRegionKey, v.Plat, v.Language)
		tmp[key] = append(tmp[key], v)
	}
	if len(tmp) > 0 {
		s.cachelist = tmp
	}
	log.Info("region list cacheproc success")
	limit, err := s.rg.Limit(context.TODO())
	if err != nil {
		log.Error("s.dao.limit error(%v)", err)
		return
	}
	s.limitCache = limit
	log.Info("region limit cacheproc success")
	config, err := s.rg.Config(context.TODO())
	if err != nil {
		log.Error("s.dao.Config error(%v)", err)
		return
	}
	s.configCache = config
	log.Info("region config cacheproc success")
}

// Square 频道广场页
func (s *Service) Square(c context.Context, mid int64, plat int8, build int, loginEvent int32, mobiApp, device, lang, buvid string) (res *channel.Square, err error) {
	res = new(channel.Square)
	var (
		squ     *tag.ChannelSquare
		regions []*channel.Region
		oidNum  = 2
	)
	isAudit := s.auditList(mobiApp, plat, build)
	eg := errgroup.Group{}
	//获取分区
	eg.Go(func() (err error) {
		_, _, regions, err = s.RegionList(c, plat, build, mobiApp, device, lang)
		if err != nil {
			log.Error("%+v", err)
			err = nil
		}
		res.Region = regions
		return
	})
	if !isAudit {
		//获取推荐频道
		eg.Go(func() (err error) {
			var (
				oids            []int64
				tagm            = map[int64]*tag.Tag{}
				chanOids        = map[int64][]*channel.ChanOids{}
				channelCards    = map[int64][]*card.Card{}
				initCardPlatKey = "card_platkey_%d_%d"
			)
			squ, err = s.tg.Square(c, mid, s.c.SquareCount, oidNum, build, loginEvent, plat, buvid, s.isOverseas(plat))
			if err != nil {
				log.Error("%+v", err)
				err = nil
			}
			for _, rec := range squ.Channels {
				cards, ok := s.cardCache[rec.ID]
				if !ok {
					continue
				}
			LOOP:
				for _, c := range cards {
					key := fmt.Sprintf(initCardPlatKey, plat, c.ID)
					cardPlat, ok := s.cardPlatCache[key]
					if !ok {
						continue
					}
					if c.Type != model.GotoAv {
						continue
					}
					for _, l := range cardPlat {
						if model.InvalidBuild(build, l.Build, l.Condition) {
							continue LOOP
						}
					}
					channelCards[c.ChannelID] = append(channelCards[c.ChannelID], c)
				}
			}
			for channelID, recOid := range squ.Oids {
				oids = append(oids, recOid...)
				if cards, ok := channelCards[channelID]; ok {
					for _, c := range cards {
						if c.Type == model.GotoAv {
							chanOids[channelID] = append(chanOids[channelID], &channel.ChanOids{Oid: c.Value, FromType: _fTypeOperation})
							oids = append(oids, c.Value)
						}
					}
				}
				for _, tmpOid := range recOid {
					chanOids[channelID] = append(chanOids[channelID], &channel.ChanOids{Oid: tmpOid, FromType: _fTypeRecommend})
				}
			}
			am, err := s.arc.Archives(c, oids)
			if err != nil {
				return
			}
			for _, rec := range squ.Channels {
				var cardItem []*operate.Card
				tagm[rec.ID] = &tag.Tag{
					ID:      rec.ID,
					Name:    rec.Name,
					Cover:   rec.Cover,
					Content: rec.ShortContent,
					Type:    int8(rec.Type),
					State:   int8(rec.State),
					IsAtten: int8(rec.Attention),
				}
				tagm[rec.ID].Count.Atten = int(rec.Sub)
				for _, oidItem := range chanOids[rec.ID] {
					if len(cardItem) >= 2 {
						break
					}
					if _, ok := am[oidItem.Oid]; !ok {
						continue
					}
					cardItem = append(cardItem, &operate.Card{ID: oidItem.Oid, FromType: oidItem.FromType})
				}
				if len(cardItem) < 2 {
					continue
				}
				var (
					h = cardm.Handle(plat, cdm.CardGt("channel_square"), "channel_square", cdm.ColumnSvrSingle, nil, tagm, nil, nil, nil)
				)
				if h == nil {
					continue
				}
				op := &operate.Card{
					ID:    rec.ID,
					Items: cardItem,
					Plat:  plat,
					Param: strconv.FormatInt(rec.ID, 10),
				}
				h.From(am, op)
				if h.Get() != nil && h.Get().Right {
					res.Square = append(res.Square, h)
				}
			}
			return
		})
	}
	eg.Wait()
	return
}

// Mysub 我订阅的tag（新） standard放前面用户自定义custom放后面
func (s *Service) Mysub(c context.Context, mid int64, limit int) (res *channel.Mysub, err error) {
	var (
		tinfo      []*tag.TagInfo
		subChannel []*channel.Channel
	)
	res = new(channel.Mysub)
	list, err := s.tg.Subscribe(c, mid)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	tinfo = list.Standard
	tinfo = append(tinfo, list.Custom...)
	if len(tinfo) > 0 {
		for _, chann := range tinfo {
			subChannel = append(subChannel, &channel.Channel{
				ID:      chann.ID,
				Name:    chann.Name,
				Cover:   chann.Cover,
				Atten:   chann.Sub,
				IsAtten: chann.Attention,
				Content: chann.Content,
			})
		}
		if len(subChannel) > limit && limit > 0 {
			subChannel = subChannel[:limit]
		}
	}
	res.List = subChannel
	res.DisplayCount = _maxAtten
	return
}

func (s *Service) isOverseas(plat int8) (res int32) {
	if ok := model.IsOverseas(plat); ok {
		res = 1
	} else {
		res = 0
	}
	return
}

func (s *Service) tablist(t *tag.ChannelDetail) (res []*channel.TabList) {
	res = s.defaultTab(t)
	var (
		mpos        []int
		tmpmenus    = map[int]*tab.Menu{}
		menus       = s.menuCache[t.Tag.ID]
		menusTabIDs = map[int64]struct{}{}
	)
	if len(menus) == 0 {
		return
	}
	for _, m := range menus {
		tmpmenus[m.Priority] = m
		mpos = append(mpos, m.Priority)
	}
	for _, pos := range mpos {
		var (
			tmpm *tab.Menu
			ok   bool
		)
		if tmpm, ok = tmpmenus[pos]; !ok || pos == 0 {
			continue
		}
		if _, ok := menusTabIDs[tmpm.TabID]; !ok {
			menusTabIDs[tmpm.TabID] = struct{}{}
		} else {
			continue
		}
		tl := &channel.TabList{}
		tl.TabListChange(tmpm)
		if len(res) < pos {
			res = append(res, tl)
			continue
		}
		res = append(res[:pos-1], append([]*channel.TabList{tl}, res[pos-1:]...)...)
	}
	return
}

func (s *Service) defaultTab(t *tag.ChannelDetail) (res []*channel.TabList) {
	for _, tmp := range _tabList {
		r := &channel.TabList{}
		*r = *tmp
		switch tmp.TabID {
		case "multiple":
			r.URI = fmt.Sprintf(r.URI, t.Tag.ID)
		case "topic":
			r.URI = fmt.Sprintf(r.URI, t.Tag.ID, url.QueryEscape(t.Tag.Name))
		}
		res = append(res, r)
	}
	return
}

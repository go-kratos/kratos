package history

import (
	"context"
	"strconv"

	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/history"
	livemdl "go-common/app/interface/main/app-interface/model/live"
	hismodle "go-common/app/interface/main/history/model"
	artmodle "go-common/app/interface/openplatform/article/model"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

const (
	_tpOld            = -1
	_tpOffline        = 0
	_tpArchive        = 3
	_tpPGC            = 4
	_tpArticle        = 5
	_tpLive           = 6
	_tpCorpus         = 7
	_androidBroadcast = 5305000
)

var (
	gotoDesc = map[int8]string{
		_tpOld:     "av",
		_tpOffline: "av",
		_tpArchive: "av",
		_tpPGC:     "pgc",
		_tpArticle: "article",
		_tpLive:    "live",
		_tpCorpus:  "article",
	}
	badge = map[int8]string{
		1: "番剧",
		2: "电影",
		3: "纪录片",
		4: "国创",
		5: "电视剧",
		7: "综艺",
	}
	busTab = []*history.BusTab{
		{
			Business: "all",
			Name:     "全部",
		},
		{
			Business: "archive",
			Name:     "视频",
		},
		{
			Business: "live",
			Name:     "直播",
		},
		{
			Business: "article",
			Name:     "专栏",
		},
	}
)

// List history list
func (s *Service) List(c context.Context, mid, build int64, pn, ps int, platform string, plat int8) (data []*history.ListRes, err error) {
	res, err := s.historyDao.History(c, mid, pn, ps)
	if len(res) > 50 {
		log.Warn("history lens(%d) mid(%d) pn(%d) ps(%d)", len(res), mid, pn, ps)
	}
	if err != nil {
		log.Error("%+v ", err)
		return
	}
	if len(res) == 0 {
		data = []*history.ListRes{}
		return
	}
	data = s.TogetherHistory(c, res, mid, build, platform, plat)
	return
}

// Live get live for history
func (s *Service) Live(c context.Context, roomIDs []int64) (res []*livemdl.RoomInfo, err error) {
	live, err := s.liveDao.LiveByRIDs(c, roomIDs)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(live) == 0 {
		res = []*livemdl.RoomInfo{}
		return
	}
	for _, lv := range live {
		item := &livemdl.RoomInfo{
			RoomID: lv.RoomID,
			URI:    model.FillURI("live", strconv.FormatInt(lv.RoomID, 10), model.LiveHandler(lv)),
		}
		if lv.Status == 1 {
			item.Status = lv.Status
		}
		res = append(res, item)
	}
	return
}

// LiveList get live list for history
func (s *Service) LiveList(c context.Context, mid, build int64, pn, ps int, platform string, plat int8) (data []*history.ListRes, err error) {
	res, err := s.historyDao.HistoryByTP(c, mid, pn, ps, _tpLive)
	if err != nil {
		log.Error("%+v ", err)
		return
	}
	if len(res) == 0 {
		data = []*history.ListRes{}
		return
	}
	data = s.TogetherHistory(c, res, mid, build, platform, plat)
	return
}

// Cursor for history
func (s *Service) Cursor(c context.Context, mid, build, paramMax int64, ps int, platform string, paramMaxTP, plat int8, businesses []string) (data *history.ListCursor, err error) {
	data = &history.ListCursor{
		List: []*history.ListRes{},
		Tab:  busTab,
	}
	// 国际版不出直播
	if plat == model.PlatAndroidI {
		data.Tab = []*history.BusTab{
			{
				Business: "all",
				Name:     "全部",
			},
			{
				Business: "archive",
				Name:     "视频",
			},
			{
				Business: "article",
				Name:     "专栏",
			},
		}
	}
	curPs := 50
	res, err := s.historyDao.Cursor(c, mid, paramMax, curPs, paramMaxTP, businesses)
	if len(res) > curPs {
		log.Warn("history lens(%d) mid(%d) paramMax(%d) paramMaxTP(%d) curPs(%d)", len(res), mid, paramMax, paramMaxTP, curPs)
	}
	if err != nil {
		log.Error("%+v ", err)
		return
	}
	if len(res) == 0 {
		return
	}
	data.List = s.TogetherHistory(c, res, mid, build, platform, plat)
	if len(data.List) >= ps {
		data.List = data.List[:ps]
	}
	if len(data.List) > 0 {
		data.Cursor = &history.Cursor{
			Max:   data.List[len(data.List)-1].ViewAt,
			MaxTP: data.List[len(data.List)-1].History.Tp,
			Ps:    ps,
		}
	}
	return
}

// TogetherHistory always return 0~50
func (s *Service) TogetherHistory(c context.Context, res []*hismodle.Resource, mid, build int64, platform string, plat int8) (data []*history.ListRes) {
	var (
		aids       []int64
		epids      []int64
		articleIDs []int64
		roomIDs    []int64
		archive    map[int64]*arcmdl.View3
		pgcInfo    map[int64]*history.PGCRes
		article    map[int64]*artmodle.Meta
		live       map[int64]*livemdl.RoomInfo
	)
	i := 0
	for _, his := range res {
		i++
		//由于出现过history吐出数量无限制，限制请求archive的数量逻辑保留
		if i > 80 {
			break
		}
		switch his.TP {
		case _tpOld, _tpOffline, _tpArchive:
			aids = append(aids, his.Oid)
		case _tpPGC:
			aids = append(aids, his.Oid) //用cid拿时长duration
			epids = append(epids, his.Epid)
		case _tpArticle:
			articleIDs = append(articleIDs, his.Oid)
		case _tpLive:
			roomIDs = append(roomIDs, his.Oid)
		case _tpCorpus:
			articleIDs = append(articleIDs, his.Cid)
		default:
			log.Warn("unknow history type(%d) msg(%+v)", his.TP, his)
		}
	}
	eg, ctx := errgroup.WithContext(c)
	if len(aids) > 0 {
		eg.Go(func() (err error) {
			archive, err = s.historyDao.Archive(ctx, aids)
			if err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(epids) > 0 {
		eg.Go(func() (err error) {
			pgcInfo, err = s.historyDao.PGC(ctx, xstr.JoinInts(epids), platform, build, mid)
			if err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(articleIDs) > 0 {
		eg.Go(func() (err error) {
			article, err = s.historyDao.Article(ctx, articleIDs)
			if err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(roomIDs) > 0 {
		eg.Go(func() (err error) {
			live, err = s.liveDao.LiveByRIDs(ctx, roomIDs)
			if err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	eg.Wait()
	for _, his := range res {
		// 国际版不出直播
		if plat == model.PlatAndroidI && his.TP == _tpLive {
			continue
		}
		tmpInfo := &history.ListRes{
			Goto:   gotoDesc[his.TP],
			ViewAt: his.Unix,
		}
		tmpInfo.History.Oid = his.Oid
		tmpInfo.History.Tp = his.TP
		tmpInfo.History.Business = his.Business
		switch his.TP {
		case _tpOld, _tpOffline, _tpArchive:
			arc, ok := archive[his.Oid]
			if !ok {
				continue
			}
			tmpInfo.Title = arc.Title
			tmpInfo.Cover = arc.Pic
			tmpInfo.Mid = arc.Author.Mid
			tmpInfo.Name = arc.Author.Name
			tmpInfo.Progress = his.Pro
			tmpInfo.Videos = arc.Videos
			for _, p := range arc.Pages {
				if p.Cid == his.Cid {
					tmpInfo.Duration = p.Duration
					tmpInfo.History.Cid = p.Cid
					tmpInfo.History.Page = p.Page
					tmpInfo.History.Part = p.Part
					break
				}
			}
			tmpInfo.URI = model.FillURI(tmpInfo.Goto, strconv.FormatInt(his.Oid, 10), model.AvHandler(arc.Archive3))
		case _tpPGC:
			pgc, okPGC := pgcInfo[his.Epid]
			arc, okArc := archive[his.Oid]
			if !okPGC || !okArc {
				continue
			}
			tmpInfo.Title = pgc.Title
			tmpInfo.ShowTitle = pgc.ShowTitle
			tmpInfo.Cover = pgc.Cover
			tmpInfo.Badge = badge[his.STP]
			tmpInfo.Progress = his.Pro
			tmpInfo.URI = model.FillURI(tmpInfo.Goto, strconv.FormatInt(his.Sid, 10), nil)
			for _, p := range arc.Pages {
				if p.Cid == his.Cid {
					tmpInfo.Duration = p.Duration
					break
				}
			}
		case _tpArticle:
			art, ok := article[his.Oid]
			if !ok {
				continue
			}
			tmpInfo.Title = art.Title
			tmpInfo.Covers = art.ImageURLs
			tmpInfo.Mid = art.Author.Mid
			tmpInfo.Name = art.Author.Name
			tmpInfo.Badge = "专栏"
			tmpInfo.URI = model.FillURI(tmpInfo.Goto, strconv.FormatInt(his.Oid, 10), nil)
		case _tpLive:
			lv, ok := live[his.Oid]
			if !ok {
				continue
			}
			tmpInfo.Title = lv.Title
			tmpInfo.Cover = lv.UserCover
			if lv.UserCover == "" {
				tmpInfo.Cover = lv.Cover
			}
			tmpInfo.Mid = lv.Mid
			tmpInfo.Name = lv.Name
			tmpInfo.TagName = lv.TagName
			if lv.Status == 1 { //1是直播中，0、2是未开播
				tmpInfo.LiveStatus = 1
			}
			if model.IsAndroid(plat) && build < _androidBroadcast {
				lv = nil
			}
			tmpInfo.URI = model.FillURI(tmpInfo.Goto, strconv.FormatInt(his.Oid, 10), model.LiveHandler(lv))
		case _tpCorpus:
			art, ok := article[his.Cid]
			if !ok {
				continue
			}
			tmpInfo.Title = art.Title
			tmpInfo.Covers = art.ImageURLs
			tmpInfo.Mid = art.Author.Mid
			tmpInfo.Name = art.Author.Name
			tmpInfo.Badge = "专栏"
			tmpInfo.URI = model.FillURI("article", strconv.FormatInt(his.Cid, 10), nil)
		default:
			continue
		}
		data = append(data, tmpInfo)
	}
	return
}

// Del for history
func (s *Service) Del(c context.Context, mid int64, hisRes []*hismodle.Resource) (err error) {
	err = s.historyDao.Del(c, mid, hisRes)
	if err != nil {
		log.Error("%+v ", err)
		return
	}
	return
}

// Clear for history
func (s *Service) Clear(c context.Context, mid int64, businesses []string) (err error) {
	err = s.historyDao.Clear(c, mid, businesses)
	if err != nil {
		log.Error("%+v ", err)
		return
	}
	return
}

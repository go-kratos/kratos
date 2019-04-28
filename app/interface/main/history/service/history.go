package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go-common/app/interface/main/history/model"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var (
	_countLimit  = 50
	_emptyVideos = []*model.Video{}
	_emptyHis    = []*model.History{}
	_emptyHismap = make(map[int64]*model.History)
)

func (s *Service) key(aid int64, typ int8) string {
	if typ < model.TypeArticle {
		return strconv.FormatInt(aid, 10)
	}
	return fmt.Sprintf("%d_%d", aid, typ)
}

// AddHistories batch update history.
// +wd:ignore
func (s *Service) AddHistories(c context.Context, mid int64, typ int8, ip string, hs []*model.History) (err error) {
	var (
		ok       bool
		his, h   *model.History
		hm2      map[string]*model.History
		hmc, res map[int64]*model.History
		expire   = time.Now().Unix() - 60*60*24*90
	)
	if len(hs) > _countLimit {
		return ecode.TargetNumberLimit
	}
	s.serviceAdds(mid, hs)
	if hm2, err = s.historyDao.Map(c, mid); err != nil {
		return
	}
	if typ < model.TypeArticle {
		if hmc, err = s.historyDao.CacheMap(c, mid); err != nil {
			return err
		}
	}
	res = make(map[int64]*model.History)
	for _, his = range hs {
		if his.Unix < expire {
			continue
		}
		if len(hm2) > 0 {
			if h, ok = hm2[s.key(his.Aid, his.TP)]; ok && his.Unix < h.Unix {
				continue
			}
		}
		if len(hmc) > 0 {
			if h, ok = hmc[his.Aid]; ok && his.Unix < h.Unix {
				continue
			}
		}
		// TODO comment && merge
		res[his.Aid] = his
	}
	if err = s.historyDao.AddMap(c, mid, res); err == nil {
		return
	}
	if typ < model.TypeArticle {
		s.historyDao.AddCacheMap(c, mid, res)
	}
	return nil
}

// AddHistory add hisotry progress into hbase.
func (s *Service) AddHistory(c context.Context, mid, rtime int64, h *model.History) (err error) {
	if h.TP < model.TypeUnknown || h.TP > model.TypeComic {
		err = ecode.RequestErr
		return
	}
	if h.Aid == 0 {
		return ecode.RequestErr
	}
	if h.TP == model.TypeBangumi || h.TP == model.TypeMovie || h.TP == model.TypePGC {
		msg := playPro{
			Type:     h.TP,
			SubType:  h.STP,
			Mid:      mid,
			Sid:      h.Sid,
			Epid:     h.Epid,
			Cid:      h.Cid,
			Progress: h.Pro,
			IP:       metadata.String(c, metadata.RemoteIP),
			Ts:       h.Unix,
			RealTime: rtime,
		}
		s.addPlayPro(&msg)
	}
	// NOTE if login user to save history
	if mid == 0 {
		return
	}
	return s.addHistory(c, mid, h)
}

// Progress get view progress from cache/hbase.
func (s *Service) Progress(c context.Context, mid int64, aids []int64) (res map[int64]*model.History, err error) {
	if mid == 0 {
		res = _emptyHismap
		return
	}
	if s.migration(mid) {
		res, err = s.servicePosition(c, mid, model.BusinessByTP(model.TypeUGC), aids)
		if err == nil {
			return
		}
	}
	if res, _, err = s.historyDao.Cache(c, mid, aids); err != nil {
		return
	} else if len(res) == 0 {
		res = _emptyHismap
	}
	return
}

// Position get view progress from cache/hbase.
func (s *Service) Position(c context.Context, mid int64, aid int64, typ int8) (res *model.History, err error) {
	if mid == 0 {
		err = ecode.NothingFound
		return
	}
	if s.migration(mid) {
		var hm map[int64]*model.History
		hm, err = s.servicePosition(c, mid, model.BusinessByTP(typ), []int64{aid})
		if err == nil && hm != nil {
			if res = hm[aid]; res == nil {
				err = ecode.NothingFound
			}
			return
		}
	}
	if typ < model.TypeArticle {
		var hm map[int64]*model.History
		hm, _, err = s.historyDao.Cache(c, mid, []int64{aid})
		if err != nil {
			return
		}
		if len(hm) > 0 {
			if h, ok := hm[aid]; ok {
				res = h
			}
		}
		if res == nil {
			err = ecode.NothingFound
		}
		return
	}
	var mhis map[string]*model.History
	mhis, err = s.historyDao.Map(c, mid)
	if err != nil {
		return
	}
	if len(mhis) > 0 {
		key := fmt.Sprintf("%d_%d", aid, typ)
		if h, ok := mhis[key]; ok {
			res = h
		}
	}
	if res == nil {
		err = ecode.NothingFound
	}
	return
}

// addHistory add new history into set.
func (s *Service) addHistory(c context.Context, mid int64, h *model.History) (err error) {
	var cmd int64
	// note: the type is video to increase experience of user .
	if h.TP < model.TypeArticle {
		s.historyDao.PushFirstQueue(c, mid, h.Aid, h.Unix)
	}
	if cmd, err = s.Shadow(c, mid); err != nil {
		return
	}
	if cmd == model.ShadowOn {
		return
	}
	h.Mid = mid
	s.serviceAdd(h)
	// note: the type is video to redis`cache .
	if h.TP >= model.TypeArticle {
		s.addProPub(h)
		if !s.conf.History.Pub {
			err = s.historyDao.Add(c, mid, h)
		}
		return
	}
	// NOTE first view
	if h.Pro < 30 && h.Pro != -1 {
		h.Pro = 0
	}
	// NOTE after 30s
	if err = s.historyDao.AddCache(c, mid, h); err != nil {
		return
	}
	s.addMerge(mid, h.Unix)
	return
}

// ClearHistory clear user's historys.
func (s *Service) ClearHistory(c context.Context, mid int64, tps []int8) (err error) {
	s.serviceClear(mid, tps)
	if len(tps) == 0 {
		s.historyDao.ClearCache(c, mid)
		err = s.historyDao.Clear(c, mid)
		s.userActionLog(mid, model.HistoryClear)
		return
	}
	tpsMap := make(map[int8]bool)
	for _, tp := range tps {
		tpsMap[tp] = true
	}
	var (
		histories map[string]*model.History
		dels      []*model.History
		chis      map[int64]*model.History
		aids      []int64
	)
	if histories, err = s.historyDao.Map(c, mid); err != nil {
		return
	}
	chis, _ = s.historyDao.CacheMap(c, mid)
	for _, v := range chis {
		histories[s.key(v.Aid, v.TP)] = v
	}
	logMap := make(map[int8]struct{})
	for _, h := range histories {
		oldType := h.TP
		h.ConvertType()
		if tpsMap[h.TP] {
			if (h.TP == model.TypeUGC) || (h.TP == model.TypePGC) {
				aids = append(aids, h.Aid)
			}
			h.TP = oldType
			dels = append(dels, h)
			logMap[oldType] = struct{}{}
		}
	}
	if len(dels) == 0 {
		return
	}
	if len(aids) > 0 {
		s.historyDao.DelCache(c, mid, aids)
	}
	if err = s.historyDao.Delete(c, mid, dels); err != nil {
		return
	}
	for k := range logMap {
		s.userActionLog(mid, fmt.Sprintf(model.HistoryClearTyp, model.BusinessByTP(k)))
	}
	return
}

// DelHistory  delete  user's history  del archive .
// +wd:ignore
func (s *Service) DelHistory(ctx context.Context, mid int64, aids []int64, typ int8) (err error) {
	if err = s.serviceDels(ctx, mid, aids, typ); err != nil {
		return
	}
	if err = s.historyDao.DelAids(ctx, mid, aids); err != nil {
		return
	}
	if typ >= model.TypeArticle {
		return
	}
	return s.historyDao.DelCache(ctx, mid, aids)
}

// Videos get videos of user view history.
// +wd:ignore
func (s *Service) Videos(c context.Context, mid int64, pn, ps int, typ int8) (res []*model.Video, err error) {
	var (
		arc         *arcmdl.View3
		ok          bool
		arcs        map[int64]*arcmdl.View3
		video       *model.Video
		history     []*model.History
		his         *model.History
		aids, epids []int64
		aidFavs     map[int64]bool
		epban       map[int64]*model.Bangumi
	)
	res = _emptyVideos
	mOK := s.migration(mid)
	if mOK {
		businesses := []string{model.BusinessByTP(model.TypeUGC), model.BusinessByTP(model.TypePGC)}
		history, epids, err = s.servicePnPsCursor(c, mid, businesses, pn, ps)
	}
	if !mOK || err != nil {
		if history, epids, err = s.histories(c, mid, pn, ps, true); err != nil {
			return
		}
	}
	if len(history) == 0 {
		return
	}
	for _, his = range history {
		aids = append(aids, his.Aid)
	}
	if typ >= model.TypeArticle {
		// TODO Article info .
		return
	}
	// bangumi info
	if len(epids) > 0 {
		epban = s.bangumi(c, mid, epids)
	}
	// archive info
	arcAids := &arcmdl.ArgAids2{Aids: aids}
	if arcs, err = s.arcRPC.Views3(c, arcAids); err != nil {
		return
	} else if len(arcs) == 0 {
		return
	}
	// favorite info
	aidFavs = s.favoriteds(c, mid, aids)
	res = make([]*model.Video, 0, len(aids))
	for _, his = range history {
		if arc, ok = arcs[his.Aid]; !ok || arc.Archive3 == nil {
			continue
		}
		// NOTE all no pay
		arc.Rights.Movie = 0
		video = &model.Video{
			Archive3: arc.Archive3,
			ViewAt:   his.Unix,
			DT:       his.DT,
			STP:      his.STP,
			TP:       his.TP,
			Progress: his.Pro,
			Count:    len(arc.Pages),
		}
		if aidFavs != nil {
			video.Favorite = aidFavs[his.Aid]
		}
		for n, p := range arc.Pages {
			if p.Cid == his.Cid {
				p.Page = int32(n + 1)
				video.Page = p
				break
			}
		}
		if video.TP == model.TypeBangumi || video.TP == model.TypeMovie || video.TP == model.TypePGC {
			if epban != nil {
				video.BangumiInfo = epban[his.Epid]
			}
			video.Count = 0
		}
		res = append(res, video)
	}
	return
}

// AVHistories return the user all av  history.
// +wd:ignore
func (s *Service) AVHistories(c context.Context, mid int64) (hs []*model.History, err error) {
	if s.migration(mid) {
		businesses := []string{model.BusinessByTP(model.TypeUGC), model.BusinessByTP(model.TypePGC)}
		hs, _, err = s.servicePnPsCursor(c, mid, businesses, 1, s.conf.History.Max)
		if err == nil {
			return
		}
	}
	hs, _, err = s.histories(c, mid, 1, s.conf.History.Max, true)
	return
}

// Histories return the user all av  history.
func (s *Service) Histories(c context.Context, mid int64, typ int8, pn, ps int) (res []*model.Resource, err error) {
	var hs []*model.History
	mOK := s.migration(mid)
	if mOK {
		var businesses []string
		if typ > 0 {
			businesses = []string{model.BusinessByTP(typ)}
		}
		hs, _, err = s.servicePnPsCursor(c, mid, businesses, pn, ps)
	}
	if !mOK || err != nil {
		if typ >= model.TypeArticle {
			hs, err = s.platformHistories(c, mid, typ, pn, ps)
		} else {
			hs, _, err = s.histories(c, mid, pn, ps, false)
		}
	}
	if err != nil {
		return
	}
	if len(hs) == 0 {
		return
	}
	for _, h := range hs {
		h.ConvertType()
		r := &model.Resource{
			Mid:  h.Mid,
			Oid:  h.Aid,
			Sid:  h.Sid,
			Epid: h.Epid,
			Cid:  h.Cid,
			DT:   h.DT,
			Pro:  h.Pro,
			Unix: h.Unix,
			TP:   h.TP,
			STP:  h.STP,
		}
		res = append(res, r)
	}
	return
}

// histories get aids of user view history
func (s *Service) histories(c context.Context, mid int64, pn, ps int, onlyAV bool) (his []*model.History, epids []int64, err error) {
	var (
		size       int
		ok         bool
		e, h       *model.History
		mhis       map[string]*model.History
		chis, ehis map[int64]*model.History
		dhis       []*model.History
		start      = (pn - 1) * ps
		end        = start + ps - 1
		total      = s.conf.History.Total
	)
	if mhis, err = s.historyDao.Map(c, mid); err != nil {
		err = nil
		mhis = make(map[string]*model.History)
	}
	chis, _ = s.historyDao.CacheMap(c, mid)
	for _, v := range chis {
		mhis[s.key(v.Aid, v.TP)] = v
	}
	if len(mhis) == 0 {
		his = _emptyHis
		return
	}
	ehis = make(map[int64]*model.History, len(mhis))
	for _, h = range mhis {
		if onlyAV && h.TP >= model.TypeArticle {
			continue
		}
		if (h.TP == model.TypeBangumi || h.TP == model.TypeMovie || h.TP == model.TypePGC) && h.Sid != 0 {
			if e, ok = ehis[h.Sid]; !ok || h.Unix > e.Unix {
				ehis[h.Sid] = h
				if e != nil {
					dhis = append(dhis, e)
				}
			}
		} else {
			his = append(his, h)
		}
	}
	for _, h = range ehis {
		if h.Epid != 0 {
			epids = append(epids, h.Epid)
		}
		his = append(his, h)
	}
	sort.Sort(model.Histories(his))
	if size = len(his); size > total {
		dhis = append(dhis, his[total:]...)
		s.delChan.Do(c, func(ctx context.Context) {
			s.historyDao.Delete(ctx, mid, dhis)
		})
		his = his[:total]
		size = total
	}
	switch {
	case size > start && size > end:
		his = his[start : end+1]
	case size > start && size <= end:
		his = his[start:]
	default:
		his = _emptyHis
	}
	return
}

// platformHistories get aids of user view history
func (s *Service) platformHistories(c context.Context, mid int64, typ int8, pn, ps int) (his []*model.History, err error) {
	var (
		size  int
		h     *model.History
		mhis  map[string]*model.History
		start = (pn - 1) * ps
		end   = start + ps - 1
		total = s.conf.History.Total
	)
	if mhis, err = s.historyDao.Map(c, mid); err != nil {
		return
	}
	if len(mhis) == 0 {
		his = _emptyHis
		return
	}
	for _, h = range mhis {
		if typ != h.TP {
			continue
		}
		his = append(his, h)
	}
	sort.Sort(model.Histories(his))
	if size = len(his); size > total {
		his = his[:total]
		size = total
	}
	switch {
	case size > start && size > end:
		his = his[start : end+1]
	case size > start && size <= end:
		his = his[start:]
	default:
		his = _emptyHis
	}
	return
}

// HistoryType get aids of user view history
// +wd:ignore
func (s *Service) HistoryType(c context.Context, mid int64, typ int8, oids []int64, ip string) (his []*model.History, err error) {
	var (
		mhis map[string]*model.History
	)
	if s.migration(mid) {
		his, err = s.serviceHistoryType(c, mid, model.BusinessByTP(typ), oids)
		if err == nil {
			return
		}
	}
	if mhis, err = s.historyDao.Map(c, mid); err != nil {
		return
	}
	if len(mhis) == 0 {
		his = _emptyHis
		return
	}
	for _, oid := range oids {
		key := fmt.Sprintf("%d_%d", oid, typ)
		if h, ok := mhis[key]; ok && h != nil {
			his = append(his, h)
		}
	}
	return
}

// HistoryCursor return the user all av  history.
func (s *Service) HistoryCursor(c context.Context, mid, max, viewAt int64, ps int, tp int8, tps []int8, ip string) (res []*model.Resource, err error) {
	var hs []*model.History
	if s.migration(mid) {
		var businesses []string
		for _, b := range tps {
			businesses = append(businesses, model.BusinessByTP(b))
		}
		res, err = s.serviceHistoryCursor(c, mid, max, businesses, model.BusinessByTP(tp), viewAt, ps)
		if err == nil {
			return
		}
	}
	hs, err = s.historyCursor(c, mid, max, viewAt, ps, tp, tps, ip)
	if err != nil {
		return
	}
	if len(hs) == 0 {
		return
	}
	for _, h := range hs {
		r := &model.Resource{
			TP:       h.TP,
			STP:      h.STP,
			Mid:      h.Mid,
			Oid:      h.Aid,
			Sid:      h.Sid,
			Epid:     h.Epid,
			Cid:      h.Cid,
			Business: h.Business,
			DT:       h.DT,
			Pro:      h.Pro,
			Unix:     h.Unix,
		}
		res = append(res, r)
	}
	return
}

// historyCursor get aids of user view history.
func (s *Service) historyCursor(c context.Context, mid, max, viewAt int64, ps int, tp int8, tps []int8, ip string) (his []*model.History, err error) {
	var (
		ok         bool
		e, h       *model.History
		mhis       map[string]*model.History
		chis, ehis map[int64]*model.History
		dhis       []*model.History
		total      = s.conf.History.Total
		tpMap      = make(map[int8]bool)
	)
	for _, tp := range tps {
		tpMap[tp] = true
	}
	if mhis, err = s.historyDao.Map(c, mid); err != nil {
		err = nil
		mhis = make(map[string]*model.History)
	}
	for k, v := range mhis {
		v.ConvertType()
		if (len(tps) > 0) && !tpMap[v.TP] {
			delete(mhis, k)
		}
	}
	chis, _ = s.historyDao.CacheMap(c, mid)
	for k, v := range chis {
		v.ConvertType()
		if (len(tps) > 0) && !tpMap[v.TP] {
			delete(chis, k)
			continue
		}
		mhis[s.key(v.Aid, v.TP)] = v
	}
	if len(mhis) == 0 {
		his = _emptyHis
		return
	}
	ehis = make(map[int64]*model.History, len(mhis))
	for _, h = range mhis {
		if (h.TP == model.TypePGC) && h.Sid != 0 {
			if e, ok = ehis[h.Sid]; !ok || h.Unix > e.Unix {
				ehis[h.Sid] = h
				if e != nil {
					dhis = append(dhis, e)
				}
			}
		} else {
			his = append(his, h)
		}
	}
	for _, h = range ehis {
		his = append(his, h)
	}
	sort.Sort(model.Histories(his))
	if len(his) > total {
		dhis = append(dhis, his[total:]...)
		s.delChan.Do(c, func(ctx context.Context) {
			s.historyDao.Delete(ctx, mid, dhis)
		})
		his = his[:total]
	}
	if viewAt != 0 || max != 0 && tp != 0 {
		for index, h := range his {
			if viewAt != 0 && h.Unix <= viewAt || h.Aid == max && h.TP == tp {
				index++
				if index+ps <= len(his) {
					return his[index : index+ps], nil
				}
				return his[index:], nil
			}
		}
		return _emptyHis, nil
	}
	if len(his) >= ps {
		return his[:ps], nil
	}
	return
}

// SetShadow set the user switch.
// +wd:ignore
func (s *Service) SetShadow(c context.Context, mid, value int64) (err error) {
	s.serviceHide(mid, value == model.ShadowOn)
	if err = s.historyDao.SetInfoShadow(c, mid, value); err != nil {
		return
	}
	return s.historyDao.SetShadowCache(c, mid, value)
}

// Delete .
func (s *Service) Delete(ctx context.Context, mid int64, his []*model.History) (err error) {
	if err = s.serviceDel(ctx, mid, his); err != nil {
		return
	}
	if err = s.historyDao.Delete(ctx, mid, his); err != nil {
		return
	}
	var aids []int64
	for _, h := range his {
		if h.TP < model.TypeArticle {
			aids = append(aids, h.Aid)
		}
	}
	if len(aids) == 0 {
		return
	}
	return s.historyDao.DelCache(ctx, mid, aids)
}

// Shadow return the user switch by mid.
// +wd:ignore
func (s *Service) Shadow(c context.Context, mid int64) (value int64, err error) {
	var (
		ok    bool
		cache = true
	)
	if s.migration(mid) {
		value, err = s.serviceHideState(c, mid)
		if err == nil {
			return
		}
	}
	if value, err = s.historyDao.ShadowCache(c, mid); err != nil {
		err = nil
		cache = false
	} else if value != model.ShadowUnknown {
		return
	}
	if value, err = s.historyDao.InfoShadow(c, mid); err != nil {
		ok = true
	}
	if cache {
		s.cache.Do(c, func(ctx context.Context) {
			s.historyDao.SetShadowCache(ctx, mid, value)
			if ok && value == model.ShadowOn {
				s.historyDao.SetInfoShadow(ctx, mid, value)
			}
		})
	}
	return
}

// FlushHistory flush to hbase from cache.
func (s *Service) FlushHistory(c context.Context, mids []int64, stime int64) (err error) {
	var (
		aids, miss []int64
		res        map[int64]*model.History
		limit      = s.conf.History.Cache
	)
	for _, mid := range mids {
		if aids, err = s.historyDao.IndexCacheByTime(c, mid, stime); err != nil {
			log.Error("s.historyDao.IndexCacheByTime(%d,%v) error(%v)", mid, stime, err)
			err = nil
			continue
		}
		if len(aids) == 0 {
			continue
		}
		if res, miss, err = s.historyDao.Cache(c, mid, aids); err != nil {
			log.Error("historyDao.Cache(%d,%v) miss:%v error(%v)", mid, aids, miss, err)
			err = nil
			continue
		}
		// * typ < model.TypeArticle all can .
		if err = s.historyDao.AddMap(c, mid, res); err != nil {
			log.Error("historyDao.AddMap(%d,%+v) error(%v)", mid, res, err)
			err = nil
		}
		if err = s.historyDao.TrimCache(c, mid, limit); err != nil {
			log.Error("historyDao.TrimCache(%d,%d) error(%v)", mid, limit, err)
			err = nil
		}
	}
	return
}

func (s *Service) merge(hmap map[int64]int64) {
	var (
		size   = int64(s.conf.History.Size)
		merges = make(map[int64][]*model.Merge, size)
	)
	for k, v := range hmap {
		merges[k%size] = append(merges[k%size], &model.Merge{Mid: k, Now: v})
	}
	for k, v := range merges {
		s.historyDao.Merge(context.TODO(), int64(k), v)
	}
}

func (s *Service) bangumi(c context.Context, mid int64, epids []int64) (bangumiMap map[int64]*model.Bangumi) {
	var n = 50
	bangumiMap = make(map[int64]*model.Bangumi, len(epids))
	for len(epids) > 0 {
		if n > len(epids) {
			n = len(epids)
		}
		epban, _ := s.historyDao.Bangumis(c, mid, epids[:n])
		epids = epids[n:]
		for k, v := range epban {
			bangumiMap[k] = v
		}
	}
	return
}

// ManagerHistory ManagerHistory.
// +wd:ignore
func (s *Service) ManagerHistory(c context.Context, onlyAV bool, mid int64) (his []*model.History, err error) {
	var (
		mhis       map[string]*model.History
		chis, ehis map[int64]*model.History
	)
	if mhis, err = s.historyDao.Map(c, mid); err != nil {
		err = nil
		mhis = make(map[string]*model.History)
	}
	chis, _ = s.historyDao.CacheMap(c, mid)
	for _, v := range chis {
		mhis[s.key(v.Aid, v.TP)] = v
	}
	if len(mhis) == 0 {
		his = _emptyHis
		return
	}
	ehis = make(map[int64]*model.History, len(mhis))
	for _, h := range mhis {
		if onlyAV && h.TP >= model.TypeArticle {
			continue
		}
		if (h.TP == model.TypeBangumi || h.TP == model.TypeMovie || h.TP == model.TypePGC) && h.Sid != 0 {
			if e, ok := ehis[h.Sid]; !ok || h.Unix > e.Unix {
				ehis[h.Sid] = h
			}
		} else {
			his = append(his, h)
		}
	}
	for _, h := range ehis {
		his = append(his, h)
	}
	sort.Sort(model.Histories(his))
	return
}

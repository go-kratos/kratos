package like

import (
	"bytes"
	"context"
	"strconv"
	"sync"

	dao "go-common/app/interface/main/activity/dao/like"
	likemdl "go-common/app/interface/main/activity/model/like"
	tagmdl "go-common/app/interface/main/tag/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	thpmdl "go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

const (
	_aidBulkSize     = 100
	_tagBlkSize      = 50
	_tagArcType      = 3
	_tagLikePoint    = 100
	_orderTypeCtime  = "ctime"
	_orderTypeRandom = "random"
	_specialLikeRate = 1000
	_businessLike    = "archive"
)

var (
	_emptyLikeList = make([]*likemdl.Like, 0)
	_emptyArcs     = make([]*api.Arc, 0)
)

// UpdateActSourceList update act arc list.
func (s *Service) updateActSourceList(c context.Context, sid int64, typ string) (err error) {
	var (
		likes []*likemdl.Item
	)
	if likes, err = s.dao.LikeList(c, sid); err != nil {
		log.Error("UpdateActSourceList s.dao.LikeList(%d) error(%v)", sid, err)
		return
	}
	s.cache.Do(c, func(c context.Context) {
		if typ == _typeAll {
			s.updateActCacheList(c, sid, likes)
		}
		if typ == _typeRegion {
			s.updateActRegionList(c, sid, likes)
		}
	})
	return
}

func (s *Service) updateActCacheList(c context.Context, sid int64, likes []*likemdl.Item) (err error) {
	var (
		aids []int64
		tags map[int64][]*tagmdl.Tag
		arcs map[int64]*api.Arc
	)
	likeMap := make(map[int64]*likemdl.Item, len(likes))
	for _, v := range likes {
		if v.Wid > 0 {
			aids = append(aids, v.Wid)
			likeMap[v.Wid] = v
		}
	}
	if len(aids) == 0 {
		return
	}
	if tags, err = s.arcTags(c, aids); err != nil {
		return
	}
	if arcs, err = s.archives(c, aids); err != nil {
		return
	}
	arcTagMap := make(map[int64][]*likemdl.Item, len(s.dialectTags))
	tagLikePtTmp := make(map[int64]int32, len(s.dialectTags))
	for aid, arcTag := range tags {
		for _, tag := range arcTag {
			if _, ok := s.dialectTags[tag.ID]; ok {
				arcTagMap[tag.ID] = append(arcTagMap[tag.ID], likeMap[aid])
				if arc, ok := arcs[aid]; ok && arc.IsNormal() {
					tagLikePtTmp[tag.ID] += arc.Stat.Like
				}
			}
		}
	}
	tagPtMap := make(map[int64]int32, len(s.dialectTags))
	for tagID, v := range arcTagMap {
		s.dao.SetLikeTagCache(c, sid, tagID, v)
		if like, ok := tagLikePtTmp[tagID]; ok {
			tagPt := int32(len(v)*_tagLikePoint) + like
			tagPtMap[tagID] = tagPt
		}
	}
	s.dao.SetTagLikeCountsCache(c, sid, tagPtMap)
	regionMap := make(map[int16][]*likemdl.Item, len(s.dialectRegions))
	for _, arc := range arcs {
		if region, ok := s.arcType[int16(arc.TypeID)]; ok {
			if _, ok := s.dialectRegions[region.Pid]; ok {
				regionMap[region.Pid] = append(regionMap[region.Pid], likeMap[arc.Aid])
			}
		}
	}
	for rid, v := range regionMap {
		s.dao.SetLikeRegionCache(c, sid, rid, v)
	}
	return
}

func (s *Service) updateActRegionList(c context.Context, sid int64, likes []*likemdl.Item) (err error) {
	var (
		aids []int64
		arcs map[int64]*api.Arc
	)
	likeMap := make(map[int64]*likemdl.Item, len(likes))
	for _, v := range likes {
		if v.Wid > 0 {
			aids = append(aids, v.Wid)
			likeMap[v.Wid] = v
		}
	}
	if len(aids) == 0 {
		return
	}
	if arcs, err = s.archives(c, aids); err != nil {
		return
	}
	regionMap := make(map[int16][]*likemdl.Item)
	for _, arc := range arcs {
		if region, ok := s.arcType[int16(arc.TypeID)]; ok {
			regionMap[region.Pid] = append(regionMap[region.Pid], likeMap[arc.Aid])
		}
	}
	for rid, v := range regionMap {
		s.dao.SetLikeRegionCache(c, sid, rid, v)
	}
	return
}

func (s *Service) archives(c context.Context, aids []int64) (archives map[int64]*api.Arc, err error) {
	var (
		mutex         = sync.Mutex{}
		aidsLen       = len(aids)
		group, errCtx = errgroup.WithContext(c)
	)
	archives = make(map[int64]*api.Arc, aidsLen)
	for i := 0; i < aidsLen; i += _aidBulkSize {
		var partAids []int64
		if i+_aidBulkSize > aidsLen {
			partAids = aids[i:]
		} else {
			partAids = aids[i : i+_aidBulkSize]
		}
		group.Go(func() (err error) {
			var arcs map[int64]*api.Arc
			arg := &archive.ArgAids2{Aids: partAids}
			if arcs, err = s.arcRPC.Archives3(errCtx, arg); err != nil {
				log.Error("s.arcRPC.Archives(%v) error(%v)", partAids, err)
				return
			}
			mutex.Lock()
			for _, v := range arcs {
				archives[v.Aid] = v
			}
			mutex.Unlock()
			return
		})
	}
	err = group.Wait()
	return
}

func (s *Service) arcTags(c context.Context, aids []int64) (tags map[int64][]*tagmdl.Tag, err error) {
	var (
		tagErr error
		mutex  = sync.Mutex{}
	)
	group, errCtx := errgroup.WithContext(c)
	aidsLen := len(aids)
	tags = make(map[int64][]*tagmdl.Tag, aidsLen)
	for i := 0; i < aidsLen; i += _tagBlkSize {
		var partAids []int64
		if i+_tagBlkSize > aidsLen {
			partAids = aids[i:]
		} else {
			partAids = aids[i : i+_tagBlkSize]
		}
		group.Go(func() (err error) {
			var tmpRes map[int64][]*tagmdl.Tag
			arg := &tagmdl.ArgResTags{Oids: partAids, Type: _tagArcType}
			if tmpRes, tagErr = s.tagRPC.ResTags(errCtx, arg); tagErr != nil {
				dao.PromError("ResTags接口错误", "s.tag.ResTag(%+v) error(%v)", arg, tagErr)
				return
			}
			mutex.Lock()
			for aid, tmpTags := range tmpRes {
				tags[aid] = tmpTags
			}
			mutex.Unlock()
			return nil
		})
	}
	group.Wait()
	return
}

// TagArcList tag arc list.
func (s *Service) TagArcList(c context.Context, sid, tagID int64, pn, ps int, typ, ip string) (list []*likemdl.Like, cnt int, err error) {
	var (
		likes      []*likemdl.Item
		start, end int
		aids       []int64
		archives   map[int64]*api.Arc
	)
	if sid != s.c.Rule.DialectSid {
		err = ecode.RequestErr
		return
	}
	if _, ok := s.dialectTags[tagID]; !ok {
		err = ecode.RequestErr
		return
	}
	if cnt, err = s.dao.LikeTagCnt(c, sid, tagID); err != nil {
		log.Error("TagArcList s.dao.LikeTagCnt sid(%d) tagID(%d) error(%v)", sid, tagID, err)
		return
	}
	if start, end, err = s.fmtStartEnd(pn, ps, cnt, typ); err != nil {
		err = nil
		list = _emptyLikeList
		return
	}
	if likes, err = s.dao.LikeTagCache(c, sid, tagID, start, end); err != nil {
		log.Error("TagArcList s.dao.LikeTagCache sid(%d) tagID(%d) start(%d) end(%d) error(%+v)", sid, tagID, start, end, err)
		return
	}
	for _, v := range likes {
		if v.Wid > 0 {
			aids = append(aids, v.Wid)
		}
	}
	if len(aids) == 0 {
		list = _emptyLikeList
		return
	}
	if archives, err = s.arcRPC.Archives3(c, &archive.ArgAids2{Aids: aids, RealIP: ip}); err != nil {
		log.Error("TagArcList s.arcRPC.Archives3 aids(%v) error(%+v)", aids, err)
		return
	}
	for _, v := range likes {
		if arc, ok := archives[v.Wid]; ok && arc.IsNormal() {
			list = append(list, &likemdl.Like{Item: v, Archive: arc})
		}
	}
	l := len(list)
	if l == 0 {
		list = _emptyLikeList
		return
	}
	if typ == _orderTypeRandom {
		s.shuffle(l, func(i, j int) {
			list[i], list[j] = list[j], list[i]
		})
	}
	return
}

// RegionArcList region arc list.
func (s *Service) RegionArcList(c context.Context, sid int64, rid int16, pn, ps int, typ, ip string) (list []*likemdl.Like, cnt int, err error) {
	var (
		likes      []*likemdl.Item
		start, end int
		aids       []int64
		archives   map[int64]*api.Arc
	)
	if sid != s.c.Rule.DialectSid {
		err = ecode.RequestErr
		return
	}
	if _, ok := s.dialectRegions[rid]; !ok {
		err = ecode.RequestErr
		return
	}
	if cnt, err = s.dao.LikeRegionCnt(c, sid, rid); err != nil {
		log.Error("RegionArcList s.dao.LikeRegionCnt sid(%d) rid(%d) error(%v)", sid, rid, err)
		return
	}
	if start, end, err = s.fmtStartEnd(pn, ps, cnt, typ); err != nil {
		err = nil
		list = _emptyLikeList
		return
	}
	if likes, err = s.dao.LikeRegionCache(c, sid, rid, start, end); err != nil {
		log.Error("RegionArcList s.dao.LikeRegionCache sid(%d) rid(%d) start(%d) end(%d) error(%+v)", sid, rid, start, end, err)
		return
	}
	for _, v := range likes {
		if v.Wid > 0 {
			aids = append(aids, v.Wid)
		}
	}
	if len(aids) == 0 {
		list = _emptyLikeList
		return
	}
	if archives, err = s.arcRPC.Archives3(c, &archive.ArgAids2{Aids: aids, RealIP: ip}); err != nil {
		log.Error("RegionArcList s.arcRPC.Archives3 aids(%v) error(%+v)", aids, err)
		return
	}
	for _, v := range likes {
		if arc, ok := archives[v.Wid]; ok && arc.IsNormal() {
			list = append(list, &likemdl.Like{Item: v, Archive: arc})
		}
	}
	l := len(list)
	if l == 0 {
		list = _emptyLikeList
		return
	}
	if typ == _orderTypeRandom {
		s.shuffle(l, func(i, j int) {
			list[i], list[j] = list[j], list[i]
		})
	}
	return
}

// TagLikeCounts .
func (s *Service) TagLikeCounts(c context.Context, sid int64) (data map[int64]int32, err error) {
	if sid != s.c.Rule.DialectSid {
		err = ecode.RequestErr
		return
	}
	return s.dao.TagLikeCountsCache(c, sid, s.c.Rule.DialectTags)
}

func (s *Service) fmtStartEnd(pn, ps, cnt int, typ string) (start, end int, err error) {
	if typ == _orderTypeCtime {
		start = (pn - 1) * ps
		end = start + ps - 1
		if start > cnt {
			err = ecode.NothingFound
			return
		}
		if end > cnt {
			end = cnt
		}
	} else {
		if ps >= cnt-1 {
			start = 0
		} else {
			start = s.r.Intn(cnt - ps - 1)
		}
		end = start + ps - 1
	}
	return
}

func (s *Service) shuffle(l int, swap func(i, j int)) {
	for i := l - 1; i > 0; i-- {
		j := s.r.Intn(i + 1)
		swap(i, j)
	}
}

// LikeInitialize initialize like cache data .
func (s *Service) LikeInitialize(c context.Context, lid int64) (err error) {
	if lid < 0 {
		lid = 0
	}
	var likesItem []*likemdl.Item
	for {
		if likesItem, err = s.dao.LikeListMoreLid(c, lid); err != nil {
			log.Error("dao.LikeInitialize(%d) error(%+v)", lid, err)
			break
		}
		if len(likesItem) == 0 {
			log.Info("LikeInitialize end success")
			break
		}
		for _, val := range likesItem {
			item := val
			if lid < item.ID {
				lid = item.ID
			}
			id := item.ID
			//the likes offline is stored with empty data
			if item.State != 1 {
				item = &likemdl.Item{}
			}
			s.cache.Do(c, func(c context.Context) {
				s.dao.AddCacheLike(c, id, item)
			})
		}
	}
	s.cache.Do(c, func(c context.Context) {
		s.LikeMaxIDInitialize(c)
	})
	return
}

// LikeMaxIDInitialize likes max id initialize
func (s *Service) LikeMaxIDInitialize(c context.Context) (err error) {
	var likeItem *likemdl.Item
	if likeItem, err = s.dao.LikeMaxID(c); err != nil {
		log.Error("s.dao.LikeMaxID() error(%+v)", err)
		return
	}
	if likeItem.ID >= 0 {
		if err = s.dao.AddCacheLikeMaxID(c, likeItem.ID); err != nil {
			log.Error("s.dao.AddCacheLikeMaxID(%d),error(%v)", likeItem.ID, err)
		}
	}
	return
}

// LikeUp update likes cache and like maxID cache
func (s *Service) LikeUp(c context.Context, lid int64) (err error) {
	var (
		likeItem  *likemdl.Item
		likeMaxID int64
	)
	group, ctx := errgroup.WithContext(c)
	group.Go(func() (e error) {
		if likeItem, e = s.dao.RawLike(ctx, lid); e != nil {
			log.Error("LikeUp:s.dao.RawLike(%d) error(%+v)", lid, e)
		}
		return
	})
	group.Go(func() (e error) {
		if likeMaxID, e = s.dao.CacheLikeMaxID(ctx); e != nil {
			log.Error("LikeUp:s.dao.CacheLikeMaxID() error(%v)", e)
		}
		return
	})
	if err = group.Wait(); err != nil {
		log.Error("LikeUp: group.Wait() error(%v)", err)
		return
	}
	if likeMaxID < lid {
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddCacheLikeMaxID(c, lid)
		})
	}
	if likeItem.ID == 0 {
		likeItem = &likemdl.Item{}
	}
	s.cache.Do(c, func(c context.Context) {
		s.dao.AddCacheLike(c, lid, likeItem)
	})
	return
}

// AddLikeCtimeCache add cache .
func (s *Service) AddLikeCtimeCache(c context.Context, lid int64) (err error) {
	var (
		likeItem *likemdl.Item
		cItems   = make([]*likemdl.Item, 0, 1)
	)
	if likeItem, err = s.dao.RawLike(c, lid); err != nil {
		log.Error("LikeUp:s.dao.RawLike(%d) error(%+v)", lid, err)
		return
	}
	if likeItem.ID > 0 {
		eg, errCtx := errgroup.WithContext(c)
		cItems = append(cItems, likeItem)
		eg.Go(func() (e error) {
			e = s.dao.LikeListCtime(errCtx, likeItem.Sid, cItems)
			return
		})
		eg.Go(func() (e error) {
			// 初始化排行榜数据
			e = s.dao.SetRedisCache(c, likeItem.Sid, lid, 0, likeItem.Type)
			return
		})
		if err = eg.Wait(); err != nil {
			log.Error("AddLikeCtimeCache eg.Wait() error(%+v)", err)
		}
	}
	return
}

// DelLikeCtimeCache delete ctime cache.
func (s *Service) DelLikeCtimeCache(c context.Context, lid, sid int64, likeType int) (err error) {
	var (
		cItems = make([]*likemdl.Item, 0, 1)
	)
	likeItem := &likemdl.Item{
		ID:   lid,
		Sid:  sid,
		Type: likeType,
	}
	cItems = append(cItems, likeItem)
	if err = s.dao.DelLikeListCtime(c, likeItem.Sid, cItems); err != nil {
		log.Error("s.dao.DelLikeListCtime(%v) error (%v)", likeItem, err)
	}
	return
}

// SubjectStat get subject stat .
func (s *Service) SubjectStat(c context.Context, sid int64) (score *likemdl.SubjectScore, err error) {
	if sid == s.c.Rule.S8Sid {
		var arcScore, artScore int64
		group, errCtx := errgroup.WithContext(c)
		group.Go(func() error {
			var (
				stat   *likemdl.SubjectStat
				arcErr error
			)
			if stat, arcErr = s.dao.CacheSubjectStat(errCtx, s.c.Rule.S8ArcSid); arcErr != nil {
				log.Error("s.dao.CacheSubjectStat sid(%d) error(%v)", sid, arcErr)
			}
			if stat == nil {
				stat = new(likemdl.SubjectStat)
			}
			arcScore = stat.Count*_specialLikeRate + stat.Like
			return nil
		})
		group.Go(func() error {
			var (
				stat   *likemdl.SubjectStat
				artErr error
			)
			if stat, artErr = s.dao.CacheSubjectStat(errCtx, s.c.Rule.S8ArtSid); artErr != nil {
				log.Error("s.dao.CacheSubjectStat sid(%d) error(%v)", sid, artErr)
			}
			if stat == nil {
				stat = new(likemdl.SubjectStat)
			}
			artScore = stat.Count*_specialLikeRate + stat.Like
			return nil
		})
		group.Wait()
		score = &likemdl.SubjectScore{Score: arcScore + artScore}
	} else {
		var stat *likemdl.SubjectStat
		if stat, err = s.dao.CacheSubjectStat(c, sid); err != nil {
			log.Error("s.dao.CacheSubjectStat sid(%d) error(%v)", sid, err)
			err = nil
		}
		if stat == nil {
			stat = new(likemdl.SubjectStat)
		}
		if sid == s.c.Rule.KingStorySid {
			score = &likemdl.SubjectScore{Score: stat.View + stat.Fav + stat.Coin + stat.Like}
		} else {
			score = &likemdl.SubjectScore{Score: stat.Count*_specialLikeRate + stat.Like}
		}
	}
	return
}

// SetSubjectStat set subject stat .
func (s *Service) SetSubjectStat(c context.Context, stat *likemdl.SubjectStat) (err error) {
	return s.dao.AddCacheSubjectStat(c, stat.Sid, stat)
}

// ViewRank get view rank arcs.
func (s *Service) ViewRank(c context.Context, sid int64, pn, ps int) (list []*api.Arc, count int, err error) {
	var (
		aidsCache       string
		aids, pieceAids []int64
		arcs            map[int64]*api.Arc
	)
	if aidsCache, err = s.dao.CacheViewRank(c, sid); err != nil {
		log.Error("ViewRank s.dao.CacheViewRank(%d) error(%v)", sid, err)
		return
	}
	if aids, err = xstr.SplitInts(aidsCache); err != nil {
		log.Error("ViewRank xstr.SplitInts(%d,%s) error(%v)", sid, aidsCache, err)
		return
	}
	count = len(aids)
	start := (pn - 1) * ps
	end := start + ps - 1
	if count < start {
		list = _emptyArcs
		return
	}
	if count > end {
		pieceAids = aids[start : end+1]
	} else {
		pieceAids = aids[start:]
	}
	if arcs, err = s.arcRPC.Archives3(c, &archive.ArgAids2{Aids: pieceAids}); err != nil {
		log.Error("ViewRank s.arcRPC.Archives3(%v) error(%v)", aids, err)
		return
	}
	for _, aid := range pieceAids {
		if arc, ok := arcs[aid]; ok && arc.IsNormal() {
			list = append(list, arc)
		}
	}
	if len(list) == 0 {
		list = _emptyArcs
	}
	return
}

// SetViewRank set view rank arcs.
func (s *Service) SetViewRank(c context.Context, sid int64, aids []int64) (err error) {
	aidsStr := xstr.JoinInts(aids)
	if err = s.dao.AddCacheViewRank(c, sid, aidsStr); err != nil {
		log.Error("SetViewRank s.dao.AddCacheViewRank(%d,%s) error(%v)", sid, aidsStr, err)
	}
	return
}

// ObjectGroup group like data.
func (s *Service) ObjectGroup(c context.Context, sid int64, ck string) (data map[int64][]*likemdl.GroupItem, err error) {
	var sids []int64
	if sids, err = s.dao.SourceItemData(c, sid); err != nil {
		log.Error("ObjectGroup SourceItemData(%d) error(%+v)", sid, err)
		return
	}
	if len(sids) == 0 {
		log.Warn("ObjectGroup sid(%d) len(sids) == 0", sid)
		err = ecode.NothingFound
		return
	}
	data = make(map[int64][]*likemdl.GroupItem, len(sids))
	group, errCtx := errgroup.WithContext(c)
	mutex := sync.Mutex{}
	for _, v := range sids {
		groupSid := v
		group.Go(func() error {
			item, e := s.dao.GroupItemData(errCtx, groupSid, ck)
			if e != nil {
				log.Error("ObjectGroup s.dao.GroupItemData(%d) error(%+v)", groupSid, e)
			} else {
				mutex.Lock()
				data[groupSid] = item
				mutex.Unlock()
			}
			return nil
		})
	}
	group.Wait()
	return
}

// SetLikeContent .
func (s *Service) SetLikeContent(c context.Context, lid int64) (err error) {
	var (
		conts map[int64]*likemdl.LikeContent
	)
	if conts, err = s.dao.RawLikeContent(c, []int64{lid}); err != nil {
		log.Error("s.dao.RawLikeContent(%d) error(%+v)", lid, err)
		return
	}
	if _, ok := conts[lid]; !ok {
		conts = make(map[int64]*likemdl.LikeContent, 1)
		conts[lid] = &likemdl.LikeContent{}
	}
	if err = s.dao.AddCacheLikeContent(c, conts); err != nil {
		log.Error("s.dao.AddCacheLikeContent(%d) error(%+v)", lid, err)
	}
	return
}

// AddLikeActCache .
func (s *Service) AddLikeActCache(c context.Context, sid, lid, score int64) (err error) {
	var (
		likeItem *likemdl.Item
	)
	if likeItem, err = s.dao.Like(c, lid); err != nil {
		log.Error("AddLikeActCache:s.dao.Like(%d) error(%+v)", lid, err)
		return
	}
	if likeItem.ID == 0 {
		return
	}
	if err = s.dao.SetRedisCache(c, sid, lid, score, likeItem.Type); err != nil {
		log.Error("AddLikeActCache:s.dao.SetRedisCache(%d,%d,%d) error(%+v)", sid, lid, score, err)
	}
	return
}

// LikeActCache .
func (s *Service) LikeActCache(c context.Context, sid, lid int64) (res int64, err error) {
	return s.dao.LikeActZscore(c, sid, lid)
}

// BatchInsertLikeExtend batch insert like_extend table.
func (s *Service) BatchInsertLikeExtend(c context.Context, extends []*likemdl.Extend) (res int64, err error) {
	var (
		buf  bytes.Buffer
		cnt  int
		rows int64
	)
	for _, v := range extends {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(v.Lid, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.Like, 10))
		buf.WriteString("),")
		cnt++
		if cnt%500 == 0 {
			buf.Truncate(buf.Len() - 1)
			if rows, err = s.dao.AddExtend(c, buf.String()); err != nil {
				log.Error("dao.dealAddExtend() error(%+v)", err)
				return
			}
			res += rows
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
		if rows, err = s.dao.AddExtend(c, buf.String()); err != nil {
			log.Error("dao.dealAddExtend() error(%+v)", err)
			return
		}
		res += rows
	}
	return
}

// arcTag get archive and tags.
func (s *Service) arcTag(c context.Context, list []*likemdl.List, order string, mid int64) (err error) {
	var (
		arcsReply   *api.ArcsReply
		lt          = len(list)
		wids        = make([]int64, 0, lt)
		tagRes      map[int64][]string
		hasLikeList map[int64]int8
	)
	for _, v := range list {
		if v.Wid > 0 {
			wids = append(wids, v.Wid)
		}
	}
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (e error) {
		arcsReply, e = s.arcClient.Arcs(errCtx, &api.ArcsRequest{Aids: wids})
		return
	})
	eg.Go(func() (e error) {
		tagRes, e = s.dao.MultiTags(errCtx, wids)
		return
	})
	if mid != 0 && (order == dao.EsOrderLikes || order == dao.ActOrderCtime) {
		eg.Go(func() (e error) {
			hasLikeList, e = s.thumbup.HasLike(errCtx, &thpmdl.ArgHasLike{Business: _businessLike, MessageIDs: wids, Mid: mid})
			return
		})
	}
	if err = eg.Wait(); err != nil {
		log.Error("arcTag:eg.Wait() error(%+v)", err)
		return
	}
	for _, v := range list {
		if v.Wid == 0 {
			continue
		}
		obj := new(likemdl.ArgTag)
		if _, ok := arcsReply.Arcs[v.Wid]; ok {
			obj.Archive = arcsReply.Arcs[v.Wid]
		}
		if _, ok := tagRes[v.Wid]; ok {
			obj.Tags = tagRes[v.Wid]
		}
		v.Object = obj
		if _, ok := hasLikeList[v.Wid]; ok {
			v.HasLikes = hasLikeList[v.Wid]
		}
	}
	return
}

// LikeOidsInfo .
func (s *Service) LikeOidsInfo(c context.Context, sType int, oids []int64) (res map[int64]*likemdl.Item, err error) {
	if res, err = s.dao.OidInfoFromES(c, oids, sType); err != nil {
		log.Error("s.dao.OidInfoFromES(%v,%d) error(%v)", oids, sType, err)
	}
	return
}

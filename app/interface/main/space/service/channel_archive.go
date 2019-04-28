package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/interface/main/space/conf"
	"go-common/app/interface/main/space/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

const _aidBulkSize = 50

// AddChannelArc add channel archive.
func (s *Service) AddChannelArc(c context.Context, mid, cid int64, aids []int64) (fakeAids []int64, err error) {
	var (
		lastID          int64
		orderNum        int
		chAids, addAids []int64
		arcs            map[int64]*arcmdl.Arc
		videos          []*model.ChannelArc
		videoMap        map[int64]int64
		remainVideos    []*model.ChannelArcSort
		ts              = time.Now()
	)
	fakeAids = make([]int64, 0)
	if _, _, err = s.channel(c, mid, cid); err != nil {
		log.Error("s.dao.Channel(%d,%d) error(%v)", mid, cid, err)
		return
	}
	if videos, err = s.dao.ChannelVideos(c, mid, cid, false); err != nil {
		log.Error("s.dao.channelVideos(%d,%d) error(%v)", mid, cid, err)
		return
	} else if orderNum = len(videos); orderNum > 0 {
		if len(aids)+orderNum > conf.Conf.Rule.MaxChArcLimit {
			err = ecode.ChMaxArcCount
			return
		}
		videoMap = make(map[int64]int64)
		for _, video := range videos {
			chAids = append(chAids, video.Aid)
			videoMap[video.Aid] = video.Aid
		}
	}
	for _, aid := range aids {
		if _, ok := videoMap[aid]; ok {
			fakeAids = append(fakeAids, aid)
		} else {
			addAids = append(addAids, aid)
		}
	}
	if len(addAids) == 0 {
		err = ecode.ChAidsExist
		return
	}
	if err = s.arcsCheck(c, mid, chAids); err != nil {
		return
	}
	if arcs, err = s.archives(c, addAids); err != nil {
		log.Error("s.arc.Archive3(%v) error(%v)", addAids, err)
		return
	}
	for _, aid := range addAids {
		if arc, ok := arcs[aid]; !ok || !arc.IsNormal() || arc.Author.Mid != mid {
			fakeAids = append(fakeAids, aid)
			continue
		}
		orderNum++
		remainVideos = append(remainVideos, &model.ChannelArcSort{Aid: aid, OrderNum: orderNum})
	}
	if len(remainVideos) == 0 {
		err = ecode.ChAidsExist
		return
	}
	if lastID, err = s.dao.AddChannelArc(c, mid, cid, ts, remainVideos); err != nil {
		log.Error("s.dao.AddChannelArc(mid:%d,cid:%d) error(%v)", mid, cid, err)
		return
	} else if lastID > 0 {
		var arcs []*model.ChannelArc
		for _, v := range remainVideos {
			arc := &model.ChannelArc{ID: lastID, Mid: mid, Cid: cid, Aid: v.Aid, OrderNum: v.OrderNum, Mtime: xtime.Time(ts.Unix())}
			arcs = append(arcs, arc)
		}
		s.dao.AddChannelArcCache(context.Background(), mid, cid, arcs)
	}
	return
}

func (s *Service) arcsCheck(c context.Context, mid int64, aids []int64) (err error) {
	var arcs map[int64]*arcmdl.Arc
	if arcs, err = s.archives(c, aids); err != nil {
		log.Error("s.archives error(%v)", err)
		return
	}
	for _, aid := range aids {
		if arc, ok := arcs[aid]; !ok || !arc.IsNormal() || arc.Author.Mid != mid {
			err = ecode.ChFakeAid
			return
		}
	}
	return
}

// DelChannelArc delete channel archive.
func (s *Service) DelChannelArc(c context.Context, mid, cid, aid int64) (err error) {
	var (
		affected int64
		orderNum int
		videos   []*model.ChannelArc
	)
	if videos, err = s.dao.ChannelVideos(c, mid, cid, false); err != nil {
		log.Error("s.dao.Channel(%d,%d) error(%v)", mid, cid, err)
		return
	} else if len(videos) == 0 {
		err = ecode.ChNoArcs
		return
	} else {
		check := false
		for _, video := range videos {
			if aid == video.Aid {
				check = true
				orderNum = video.OrderNum
			}
		}
		if !check {
			err = ecode.ChNoArc
			return
		}
	}
	if affected, err = s.dao.DelChannelArc(c, mid, cid, aid, orderNum); err != nil {
		log.Error("s.dao.DelChannelArc(%d,%d) error(%v)", mid, aid, err)
		return
	} else if affected > 0 {
		s.dao.DelChannelArcCache(c, mid, cid, aid)
		s.setChannelArcSortCache(c, mid, cid)
	}
	return
}

// SortChannelArc sort channel archive.
func (s *Service) SortChannelArc(c context.Context, mid, cid, aid int64, orderNum int) (err error) {
	var (
		videos                                 []*model.ChannelArc
		bfSortBegin, bfSortEnd, chSort, afSort []*model.ChannelArcSort
		affected                               int64
		aidIndex, aidOn                        int
		aidCheck                               bool
		ts                                     = time.Now()
	)
	if videos, err = s.dao.ChannelVideos(c, mid, cid, false); err != nil {
		log.Error("s.dao.ChannelVideos(%d,%d) error(%v)", mid, cid, err)
		return
	} else if len(videos) == 0 {
		err = ecode.ChNoArcs
		return
	} else {
		videoLen := len(videos)
		if orderNum > videoLen {
			err = ecode.RequestErr
			return
		}
		for index, video := range videos {
			if aid == video.Aid {
				aidCheck = true
				aidIndex = index
				aidOn = video.OrderNum
				break
			}
		}
		if !aidCheck {
			err = ecode.RequestErr
			return
		}
		if orderNum > aidOn {
			chSort = append(chSort, &model.ChannelArcSort{Aid: aid, OrderNum: orderNum})
			for i, v := range videos {
				if i < videoLen-orderNum {
					bfSortBegin = append(bfSortBegin, &model.ChannelArcSort{Aid: v.Aid, OrderNum: v.OrderNum})
				} else if i >= videoLen-orderNum && i < aidIndex {
					chSort = append(chSort, &model.ChannelArcSort{Aid: v.Aid, OrderNum: v.OrderNum - 1})
				} else if i > aidIndex {
					bfSortEnd = append(bfSortEnd, &model.ChannelArcSort{Aid: v.Aid, OrderNum: v.OrderNum})
				}
			}
		} else if orderNum < aidOn {
			for i, v := range videos {
				if i < aidIndex {
					bfSortBegin = append(bfSortBegin, &model.ChannelArcSort{Aid: v.Aid, OrderNum: v.OrderNum})
				} else if i > aidIndex && i <= videoLen-orderNum {
					chSort = append(chSort, &model.ChannelArcSort{Aid: v.Aid, OrderNum: v.OrderNum + 1})
				} else if i > videoLen-orderNum {
					bfSortEnd = append(bfSortEnd, &model.ChannelArcSort{Aid: v.Aid, OrderNum: v.OrderNum})
				}
			}
			chSort = append(chSort, &model.ChannelArcSort{Aid: aid, OrderNum: orderNum})
		} else {
			return
		}
		afSort = append(afSort, bfSortBegin...)
		afSort = append(afSort, chSort...)
		afSort = append(afSort, bfSortEnd...)
	}
	if affected, err = s.dao.EditChannelArc(c, mid, cid, ts, chSort); err != nil {
		log.Error("s.dao.s.dao.EditChannelArc(%d,%d,%d,%d) error(%v)", mid, cid, aid, orderNum, err)
		return
	} else if affected > 0 {
		s.dao.SetChannelArcSortCache(c, mid, cid, afSort)
	}
	return
}

// ChannelVideos get channel and channel video info.
func (s *Service) ChannelVideos(c context.Context, mid, cid int64, pn, ps int, isGuest, order bool) (res *model.ChannelDetail, err error) {
	var (
		channel *model.Channel
		start   = (pn - 1) * ps
		end     = start + ps - 1
	)
	if channel, err = s.Channel(c, mid, cid); err != nil {
		return
	}
	res = &model.ChannelDetail{Channel: channel}
	res.Archives, err = s.channelArc(c, mid, cid, start, end, isGuest, order)
	return
}

func (s *Service) channelVideos(c context.Context, mid, cid int64, start, end int, order bool) (res []*model.ChannelArc, err error) {
	var (
		videos   []*model.ChannelArc
		addCache = true
	)
	if res, err = s.dao.ChannelArcsCache(c, mid, cid, start, end, order); err != nil {
		addCache = false
	} else if len(res) > 0 {
		return
	}
	if videos, err = s.dao.ChannelVideos(c, mid, cid, order); err != nil {
		log.Error("s.dao.ChannelVideos(%d,%d) error(%v)", mid, cid, err)
		return
	} else if len(videos) > 0 {
		if addCache {
			s.cache.Do(c, func(c context.Context) {
				s.dao.SetChannelArcsCache(c, mid, cid, videos)
				s.setChannelArcSortCache(c, mid, cid)
			})
		}
		length := len(videos)
		if length < start {
			res = make([]*model.ChannelArc, 0)
			return
		}
		if length > end {
			res = videos[start : end+1]
		} else {
			res = videos[start:]
		}
	}
	return
}

// CheckChannelVideo check useless channel video.
func (s *Service) CheckChannelVideo(c context.Context, mid, cid int64) (err error) {
	var (
		videos []*model.ChannelArc
		aids   []int64
	)
	if videos, err = s.dao.ChannelVideos(c, mid, cid, false); err != nil {
		log.Error("s.dao.channelVideos(%d,%d) error(%v)", mid, cid, err)
		return
	}
	for _, v := range videos {
		aids = append(aids, v.Aid)
	}
	err = s.arcsCheck(c, mid, aids)
	return
}

func (s *Service) channelArc(c context.Context, mid, cid int64, start, end int, isGuest, order bool) (res []*arcmdl.Arc, err error) {
	var (
		videoAids []*model.ChannelArc
		archives  map[int64]*arcmdl.Arc
		aids      []int64
	)
	if videoAids, err = s.channelVideos(c, mid, cid, start, end, order); err != nil {
		log.Error("s.dao.ChannelVideos(%d,%d) error(%v)", mid, cid, err)
		return
	} else if len(videoAids) == 0 {
		res = _emptyChArc
		return
	}
	for _, video := range videoAids {
		aids = append(aids, video.Aid)
	}
	if archives, err = s.archives(c, aids); err != nil {
		log.Error("s.arc.Archives3(%v) error(%v)", aids, err)
		return
	}
	for _, video := range videoAids {
		if arc, ok := archives[video.Aid]; ok {
			if arc.IsNormal() {
				if arc.Access >= 10000 {
					arc.Stat.View = -1
				}
				res = append(res, arc)
			} else {
				res = append(res, &arcmdl.Arc{Aid: video.Aid, Title: arc.Title, Pic: arc.Pic, Stat: arc.Stat, PubDate: arc.PubDate, State: arc.State})
			}
		}
	}
	return
}

func (s *Service) setChannelArcSortCache(c context.Context, mid, cid int64) (err error) {
	var (
		videos []*model.ChannelArc
		sorts  []*model.ChannelArcSort
	)
	if videos, err = s.dao.ChannelVideos(c, mid, cid, false); err != nil {
		log.Error("s.dao.ChannelVideos(%d,%d) error(%v)", mid, cid, err)
		return
	} else if len(videos) == 0 {
		return
	}
	for _, v := range videos {
		sort := &model.ChannelArcSort{Aid: v.Aid, OrderNum: v.OrderNum}
		sorts = append(sorts, sort)
	}
	return s.dao.SetChannelArcSortCache(c, mid, cid, sorts)
}

func (s *Service) archives(c context.Context, aids []int64) (archives map[int64]*arcmdl.Arc, err error) {
	var (
		mutex         = sync.Mutex{}
		aidsLen       = len(aids)
		group, errCtx = errgroup.WithContext(c)
	)
	archives = make(map[int64]*arcmdl.Arc, aidsLen)
	for i := 0; i < aidsLen; i += _aidBulkSize {
		var partAids []int64
		if i+_aidBulkSize > aidsLen {
			partAids = aids[i:]
		} else {
			partAids = aids[i : i+_aidBulkSize]
		}
		group.Go(func() (err error) {
			var arcs *arcmdl.ArcsReply
			arg := &arcmdl.ArcsRequest{Aids: partAids}
			if arcs, err = s.arcClient.Arcs(errCtx, arg); err != nil {
				log.Error("s.arcClient.Arcs(%v) error(%v)", partAids, err)
				return
			}
			mutex.Lock()
			for _, v := range arcs.Arcs {
				archives[v.Aid] = v
			}
			mutex.Unlock()
			return
		})
	}
	err = group.Wait()
	return
}

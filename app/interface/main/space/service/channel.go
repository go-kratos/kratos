package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/interface/main/space/conf"
	"go-common/app/interface/main/space/model"
	arcmdl "go-common/app/service/main/archive/api"
	filmdl "go-common/app/service/main/filter/model/rpc"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

var (
	_emptyChArc        = make([]*arcmdl.Arc, 0)
	_emptyChList       = make([]*model.Channel, 0)
	_emptyChDetailList = make([]*model.ChannelDetail, 0)
	_nameErrorLevel    = int8(20)
	_introWarnLevel    = int8(20)
	_introErrorLevel   = int8(30)
)

// ChannelList get channel list.
func (s *Service) ChannelList(c context.Context, mid int64, isGuest bool) (channels []*model.Channel, err error) {
	var (
		channelExtra map[int64]*model.ChannelExtra
		cids         []int64
		addCache     = true
	)
	if channels, err = s.dao.ChannelListCache(c, mid); err != nil {
		addCache = false
	} else if len(channels) > 0 {
		return
	}
	if channels, err = s.dao.ChannelList(c, mid); err != nil {
		log.Error("s.dao.ChannelList(%d) error(%v)", mid, err)
		return
	}
	if len(channels) == 0 {
		channels = _emptyChList
		return
	}
	for _, channel := range channels {
		cids = append(cids, channel.Cid)
	}
	if channelExtra, err = s.channelExtra(c, mid, cids); err != nil {
		err = nil
		return
	}
	for _, channel := range channels {
		if _, ok := channelExtra[channel.Cid]; ok {
			channel.Count = channelExtra[channel.Cid].Count
			channel.Cover = channelExtra[channel.Cid].Cover
		}
	}
	if addCache {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetChannelListCache(c, mid, channels)
		})
	}
	return
}

// Channel get channel info.
func (s *Service) Channel(c context.Context, mid, cid int64) (channel *model.Channel, err error) {
	var (
		extra    *model.ChannelExtra
		arcReply *arcmdl.ArcReply
		addCache bool
	)
	if channel, addCache, err = s.channel(c, mid, cid); err != nil {
		log.Error("s.channel(%d,%d) error(%v)", mid, cid, err)
		return
	}
	if extra, err = s.dao.ChannelExtra(c, mid, cid); err != nil {
		log.Error("s.dao.ChannelExtra(%d,%d) error(%v)", mid, cid, err)
		err = nil
	} else if extra != nil {
		channel.Count = extra.Count
		if extra.Aid > 0 {
			if arcReply, err = s.arcClient.Arc(c, &arcmdl.ArcRequest{Aid: extra.Aid}); err != nil {
				log.Error("s.arcClient.Arc(%d) error(%v)", extra.Aid, err)
				err = nil
			} else {
				channel.Cover = arcReply.Arc.Pic
			}
		}
	}
	if addCache {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetChannelCache(c, mid, cid, channel)
		})
	}
	return
}

func (s *Service) channel(c context.Context, mid, cid int64) (res *model.Channel, addCache bool, err error) {
	addCache = true
	if res, err = s.dao.ChannelCache(c, mid, cid); err != nil {
		addCache = false
	} else if res != nil {
		return
	}
	if res, err = s.dao.Channel(c, mid, cid); err != nil {
		log.Error("s.dao.Channel(%d,%d) error(%v)", mid, cid, err)
	} else if res == nil {
		err = ecode.NothingFound
	}
	return
}

// ChannelIndex get channel index info.
func (s *Service) ChannelIndex(c context.Context, mid int64, isGuest bool) (channelDetails []*model.ChannelDetail, err error) {
	var (
		channels []*model.Channel
		detail   *model.ChannelDetail
	)
	if channels, err = s.ChannelList(c, mid, isGuest); err != nil {
		log.Error("s.Channel(%d) error(%v)", mid, err)
		return
	} else if len(channels) == 0 {
		channelDetails = _emptyChDetailList
		return
	}
	group, errCtx := errgroup.WithContext(c)
	mutex := sync.Mutex{}
	for _, channel := range channels {
		cid := channel.Cid
		group.Go(func() (err error) {
			if detail, err = s.ChannelVideos(errCtx, mid, cid, 1, conf.Conf.Rule.ChIndexCnt, isGuest, false); err != nil {
				log.Error("s.ChannelVideos(%d,%d) error(%v)", mid, cid, err)
				err = nil
			} else if detail != nil {
				mutex.Lock()
				channelDetails = append(channelDetails, detail)
				mutex.Unlock()
			}
			return
		})
	}
	group.Wait()
	if len(channelDetails) == 0 {
		channelDetails = _emptyChDetailList
	}
	return
}

// AddChannel add channel.
func (s *Service) AddChannel(c context.Context, mid int64, name, intro string) (cid int64, err error) {
	var (
		afIntro string
		ts      = time.Now()
	)
	if _, err = s.realName(c, mid); err != nil {
		return
	}
	if err = s.channelCheck(c, mid, 0, name, true, true); err != nil {
		log.Error("s.channelCheck(%d,%s) error(%v)", mid, name, err)
		return
	}
	if afIntro, err = s.channelFilter(c, name, intro); err != nil {
		log.Error("s.channelFilter(%s,%s) error(%v)", name, intro, err)
		return
	}
	if cid, err = s.dao.AddChannel(c, mid, name, afIntro, ts); err != nil {
		log.Error("s.dao.AddChannel(%d,%s,%s) error(%v)", mid, name, intro, err)
		return
	} else if cid > 0 {
		s.cache.Do(c, func(c context.Context) {
			ch := &model.Channel{Cid: cid, Mid: mid, Name: name, Intro: intro, Mtime: xtime.Time(ts.Unix())}
			s.dao.SetChannelCache(c, mid, cid, ch)
		})
	}
	return
}

// EditChannel edit channel.
func (s *Service) EditChannel(c context.Context, mid, cid int64, name, intro string) (err error) {
	var (
		affected int64
		afIntro  string
		ts       = time.Now()
	)
	if _, err = s.realName(c, mid); err != nil {
		return
	}
	if err = s.channelCheck(c, mid, cid, name, true, false); err != nil {
		log.Error("s.channelCheck(%d,%d,%s) error(%v)", mid, cid, name, err)
		return
	}
	if afIntro, err = s.channelFilter(c, name, intro); err != nil {
		log.Error("s.channelFilter(%s,%s) error(%v)", name, intro, err)
		return
	}
	//if channel,err := s.Channel(c,mid,cid,ip)
	if affected, err = s.dao.EditChannel(c, mid, cid, name, afIntro, ts); err != nil {
		log.Error("s.dao.EditChannel(%d,%s,%s) error(%v)", mid, name, intro, err)
		return
	} else if affected > 0 {
		s.cache.Do(c, func(c context.Context) {
			ch := &model.Channel{Cid: cid, Mid: mid, Name: name, Intro: intro, Mtime: xtime.Time(ts.Unix())}
			s.dao.SetChannelCache(c, mid, cid, ch)
		})
	}
	return
}

// DelChannel del channel.
func (s *Service) DelChannel(c context.Context, mid, cid int64) (err error) {
	var affected int64
	if affected, err = s.dao.DelChannel(c, mid, cid); err != nil {
		log.Error("s.dao.DelChannel(%d,%d) error(%v)", mid, cid, err)
		return
	} else if affected > 0 {
		s.dao.DelChannelCache(c, mid, cid)
		s.dao.DelChannelArcsCache(c, mid, cid)
	}
	return
}

func (s *Service) channelExtra(c context.Context, mid int64, cids []int64) (extra map[int64]*model.ChannelExtra, err error) {
	if len(cids) == 0 {
		return
	}
	var (
		arcsReply *arcmdl.ArcsReply
		aids      = make([]int64, 0, len(cids))
	)
	extra = make(map[int64]*model.ChannelExtra, len(cids))
	for _, cid := range cids {
		var data *model.ChannelExtra
		if data, err = s.dao.ChannelExtra(c, mid, cid); err != nil {
			log.Error("s.dao.ChannelExtra(%d,%d) error(%v)", mid, cid, err)
			continue
		} else if data != nil {
			extra[cid] = &model.ChannelExtra{Aid: data.Aid, Cid: data.Cid, Count: data.Count}
			if data.Aid > 0 {
				aids = append(aids, data.Aid)
			}
		}
	}
	if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
		log.Error("s.arcClient.Arcs(%v) error (%v)", aids, err)
		return
	}
	for _, cid := range cids {
		if _, ok := extra[cid]; ok {
			if arc, ok := arcsReply.Arcs[extra[cid].Aid]; ok {
				extra[cid].Cover = arc.Pic
			}
		}
	}
	return
}

func (s *Service) channelCheck(c context.Context, mid, cid int64, name string, nameCheck, countCheck bool) (err error) {
	var (
		channels []*model.Channel
		dbCheck  = false
	)
	if channels, err = s.dao.ChannelListCache(c, mid); err != nil {
		err = nil
		dbCheck = true
	} else if len(channels) == 0 {
		dbCheck = true
	}
	if dbCheck {
		if channels, err = s.dao.ChannelList(c, mid); err != nil {
			log.Error("s.dao.ChannelList(%d) error(%v)", mid, err)
			return
		}
	}
	if cnt := len(channels); cnt > 0 {
		if countCheck && cnt > conf.Conf.Rule.MaxChLimit {
			err = ecode.ChMaxCount
			return
		}
		if nameCheck {
			for _, channel := range channels {
				if name == channel.Name && cid != channel.Cid {
					err = ecode.ChNameExist
					return
				}
			}
		}
	}
	return
}

func (s *Service) channelFilter(c context.Context, name, intro string) (afterIntro string, err error) {
	var (
		filterRes map[string]*filmdl.FilterRes
		arg       = &filmdl.ArgMfilter{Area: "common", Message: map[string]string{"name": name, "intro": intro}}
	)
	afterIntro = intro
	if filterRes, err = s.filter.MFilter(c, arg); err != nil {
		log.Error("s.filter.MFilter(%v) error(%v)", arg, err)
		return
	}
	for k, v := range filterRes {
		if k == "name" && v.Level >= _nameErrorLevel {
			err = ecode.ChNameBanned
			return
		}
		if k == "intro" {
			if v.Level == _introWarnLevel {
				afterIntro = v.Result
			} else if v.Level >= _introErrorLevel {
				err = ecode.ChIntroBanned
				return
			}
		}
	}
	return
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/interface/main/dm2/model"
	account "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	archive "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// .
var (
	_dmFlagFmt = `{"rec_flag":%d,"rec_text":"%s","rec_switch":%d}`
)

// View dm view
func (s *Service) View(c context.Context, mid, aid, oid int64, tp int32, plat int32) (res *model.ViewDm, err error) {
	var (
		sub *model.Subject
		eg  errgroup.Group
	)
	// hot cache
	if hit, ok := s.localViewCache[keyLocalView(aid, oid)]; ok {
		return hit, nil
	}
	if sub, err = s.subject(c, model.SubTypeVideo, oid); err != nil {
		return
	}
	res = &model.ViewDm{
		Closed: sub.State == model.SubStateClosed,
		Flag:   json.RawMessage([]byte(fmt.Sprintf(_dmFlagFmt, s.conf.DmFlag.RecFlag, s.conf.DmFlag.RecText, s.conf.DmFlag.RecSwitch))),
	}
	// mask
	eg.Go(func() (err error) {
		var (
			mask     *model.Mask
			maskPlat int8
		)
		switch plat {
		case model.PlatWeb:
			maskPlat = model.MaskPlatWeb
		case model.PlatAndroid, model.PlatIPhone, model.PlatIPad, model.PlatPadHd:
			maskPlat = model.MaskPlatMbl
		default:
			return
		}
		if mask, err = s.MaskListWithSub(c, oid, maskPlat, sub); err != nil {
			log.Error("View.MaskListWithSub(oid:%v) error(%v)", oid, err)
			return
		}
		res.ViewDmMask = mask
		return
	})
	// subtitle
	eg.Go(func() (err error) {
		var (
			subtitle *model.ViewSubtitle
		)
		if subtitle, err = s.viewSubtitles(c, aid, oid, tp); err != nil {
			log.Error("View.viewSubtitles(aid:%v,oid:%v) error(%v)", aid, oid, err)
			return
		}
		res.Subtitle = subtitle
		return
	})
	// special dm
	// TODO special dm
	// dm seg rule
	eg.Go(func() (err error) {
		var (
			dmSeg *model.ViewDmSeg
		)
		if dmSeg, err = s.viewDmSeg(c, aid, oid); err != nil {
			log.Error("View.viewDmSeg(aid:%v,oid:%v) error(%v)", aid, oid, err)
			return
		}
		res.ViewDmSeg = dmSeg
		return
	})
	// ignore error
	eg.Wait()
	return
}

func (s *Service) viewSubtitles(c context.Context, aid, oid int64, tp int32) (viewSubtitle *model.ViewSubtitle, err error) {
	var (
		videoSubtitles  []*model.VideoSubtitle
		subtitles       []*model.ViewVideoSubtitle
		reply           *account.InfoReply
		subtitleSubject *model.SubtitleSubjectReply
	)
	if videoSubtitles, err = s.getVideoSubtitles(c, oid, tp); err != nil {
		log.Error("View.getVideoSubtitles(oid:%v) error(%v)", oid, err)
		return
	}
	for _, videoSubtitle := range videoSubtitles {
		subtitle := &model.ViewVideoSubtitle{
			ID:          videoSubtitle.ID,
			Lan:         videoSubtitle.Lan,
			LanDoc:      videoSubtitle.LanDoc,
			SubtitleURL: videoSubtitle.SubtitleURL,
		}
		if videoSubtitle.AuthorMid > 0 {
			if reply, _ = s.accountRPC.Info3(c, &account.MidReq{Mid: videoSubtitle.AuthorMid}); reply != nil {
				subtitle.Author = &model.ViewAuthor{
					Mid:  reply.GetInfo().GetMid(),
					Name: reply.GetInfo().GetName(),
					Sex:  reply.GetInfo().GetSex(),
					Face: reply.GetInfo().GetFace(),
					Sign: reply.GetInfo().GetSign(),
					Rank: reply.GetInfo().GetRank(),
				}
			}
		}
		subtitles = append(subtitles, subtitle)
	}
	if subtitleSubject, err = s.SubtitleSubject(c, aid); err != nil {
		log.Error("View.subtitleSubject(aid:%v) error(%v)", aid, err)
		return
	}
	viewSubtitle = &model.ViewSubtitle{
		Subtitles: subtitles,
	}
	if subtitleSubject != nil {
		viewSubtitle.Lan = subtitleSubject.Lan
		viewSubtitle.LanDoc = subtitleSubject.LanDoc
	}
	return
}

func (s *Service) viewDmSeg(c context.Context, aid, oid int64) (dmSeg *model.ViewDmSeg, err error) {
	var (
		duration int64
		cnt      int64
	)
	if duration, err = s.videoDuration(c, aid, oid); err != nil {
		return
	}
	cnt = duration / model.DefaultPageSize
	if duration%model.DefaultPageSize > 0 {
		cnt++
	}
	dmSeg = &model.ViewDmSeg{
		PageSize: model.DefaultPageSize,
		Total:    cnt,
	}
	return
}

func (s *Service) viewProc() {
	if len(s.conf.Localcache.ViewAids) <= 0 {
		return
	}
	ticker := time.NewTicker(time.Duration(s.conf.Localcache.ViewExpire))
	defer ticker.Stop()
	for range ticker.C {
		s.cacheView(s.conf.Localcache.ViewAids)
	}
}

func keyLocalView(aid, oid int64) string {
	return fmt.Sprintf("dm_view_%d_%d", aid, oid)
}

func (s *Service) cacheView(aids []int64) {
	var (
		sub      *model.Subject
		pages    []*api.Page
		err      error
		cacheMap = make(map[string]*model.ViewDm)
	)
	for _, aid := range aids {
		pages, err = s.arcRPC.Page3(context.Background(), &archive.ArgAid2{
			Aid: aid,
		})
		if err != nil {
			log.Error("localCacheView.Page3(aid:%v) error(%v)", aid, err)
			continue
		}
		for _, page := range pages {
			if sub, err = s.subject(context.Background(), model.SubTypeVideo, page.Cid); err != nil {
				continue
			}
			res := &model.ViewDm{
				Closed: sub.State == model.SubStateClosed,
				Flag:   json.RawMessage([]byte(fmt.Sprintf(_dmFlagFmt, s.conf.DmFlag.RecFlag, s.conf.DmFlag.RecText, s.conf.DmFlag.RecSwitch))),
			}
			// ignore error
			if res.Subtitle, err = s.viewSubtitles(context.Background(), aid, page.Cid, model.SubTypeVideo); err != nil {
				log.Error("View.viewSubtitles(aid:%v,oid:%v) error(%v)", aid, page.Cid, err)
				err = nil
			}
			if res.ViewDmSeg, err = s.viewDmSeg(context.Background(), aid, page.Cid); err != nil {
				log.Error("View.viewDmSeg(aid:%v,oid:%v) error(%v)", aid, page.Cid, err)
				err = nil
			}
			cacheMap[keyLocalView(aid, page.Cid)] = res
		}
	}
	s.localViewCache = cacheMap
}

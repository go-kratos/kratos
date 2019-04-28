package service

import (
	"context"
	"time"

	"go-common/app/job/main/dm2/model"
	"go-common/app/service/main/archive/api"
	archiveMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

var (
	_maskJobDay = []int32{3, 7}
)

// maskProc .
func (s *Service) maskProc() {
	var (
		err    error
		c      = context.Background()
		ticker *time.Ticker
	)
	if s.conf.MaskCate == nil {
		return
	}
	ticker = time.NewTicker(time.Duration(s.conf.MaskCate.Interval))
	for range ticker.C {
		if err = s.maskSchedule(c); err != nil {
			log.Error("maskProc.error(%v)", err)
			continue
		}
	}
}

func (s *Service) maskSchedule(c context.Context) (err error) {
	var (
		ok                               bool
		now                              = time.Now()
		expire                           = now.Add(time.Duration(s.conf.MaskCate.Interval))
		expireStr                        = expire.Format(time.RFC3339)
		oldExpireStr, oldExpireGetSetStr string
		oldExpire                        time.Time
	)
	if ok, err = s.dao.SetnxMaskJob(c, expireStr); err != nil {
		return
	}
	// redis中不存在
	if ok {
		if err = s.maskJob(c); err != nil {
			s.dao.DelMaskJob(c)
			log.Error("maskJob,error(%v)", err)
			return
		}
		return
	}
	// redis中已经存在
	// 判断是否过期了
	if oldExpireStr, err = s.dao.GetMaskJob(c); err != nil {
		return
	}
	if oldExpire, err = time.Parse(time.RFC3339, oldExpireStr); err != nil {
		return
	}
	if oldExpire.Sub(now) > 0 {
		return
	}
	if oldExpireGetSetStr, err = s.dao.GetSetMaskJob(c, expireStr); err != nil {
		return
	}
	if oldExpireGetSetStr != oldExpireStr {
		return
	}
	if err = s.maskJob(c); err != nil {
		s.dao.DelMaskJob(c)
		log.Error("maskJob,error(%v)", err)
		return
	}
	return
}

// 执行任务
func (s *Service) maskJob(c context.Context) (err error) {
	for _, tid := range s.conf.MaskCate.Tids {
		if err = s.maskOneCate(c, tid); err != nil {
			log.Error("maskOneCate(tid:%v),error(%v)", tid, err)
			return
		}
	}
	return
}

func (s *Service) maskOneCate(c context.Context, tid int64) (err error) {
	var (
		err1 error
		resp *model.RankRecentResp
		aids []int64
	)
	for _, day := range _maskJobDay {
		if resp, err = s.dao.RankList(c, tid, day); err != nil {
			log.Error("RankList(tid:%v,day:%v),error(%v)", tid, day, err)
			return
		}
		for idx, recentRegion := range resp.List {
			if idx >= s.conf.MaskCate.Limit {
				break
			}
			aids = append(aids, recentRegion.Aid)
			for _, other := range recentRegion.Others {
				aids = append(aids, other.Aid)
			}
		}
	}
	for _, aid := range aids {
		if err1 = s.maskOneArchive(c, aid); err1 != nil {
			log.Error("maskOneArchive.err aid:%v,error(%v)", aid, err1)
			continue
		}
		log.Info("maskOneArchive.ok aid:%v", aid)
	}
	return
}

func (s *Service) maskOneArchive(c context.Context, aid int64) (err error) {
	var (
		pages []*api.Page
	)
	if pages, err = s.arcRPC.Page3(c, &archiveMdl.ArgAid2{Aid: aid}); err != nil {
		log.Error("s.arcRPC.Page3(aid:%v),error(%v)", aid, err)
		return
	}
	for _, page := range pages {
		if err = s.maskOneVideo(c, page.Cid); err != nil {
			log.Error("maskOneVideo(oid:%v),error(%v)", page.Cid, err)
			return
		}
	}
	return
}

// runGenMask send to gen mask url
func (s *Service) maskOneVideo(c context.Context, oid int64) (err error) {
	var (
		subject  *model.Subject
		archive3 *api.Arc
		err1     error
		duration int64
		typeID   int32
	)
	if subject, err = s.subject(c, model.SubTypeVideo, oid); err != nil {
		log.Error("s.subject(oid:%v),error(%v)", oid, err)
		return
	}
	if subject.AttrVal(model.AttrSubMaskOpen) == model.AttrYes {
		return
	}
	if archive3, err1 = s.arcRPC.Archive3(c, &archiveMdl.ArgAid2{Aid: subject.Pid}); err1 == nil && archive3 != nil {
		duration = archive3.Duration
		typeID = archive3.TypeID
	}
	if err = s.dao.GenerateMask(c, oid, subject.Mid, model.MaskPlatAll, model.MaskPriorityLow, subject.Pid, duration, typeID); err != nil {
		log.Error("GenerateMask(oid:%v),error(%v)", oid, err)
		return
	}
	subject.AttrSet(model.AttrYes, model.AttrSubMaskOpen)
	if _, err = s.dao.UpdateSubAttr(c, subject.Type, subject.Oid, subject.Attr); err != nil {
		log.Error("UpdateSubAttr(oid:%v,attr:%v),error(%v)", oid, subject.Attr, err)
		return
	}
	return
}

func (s *Service) maskMidProc() {
	var (
		c    = context.Background()
		mids []int64
		err  error
	)
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	for range ticker.C {
		if mids, err = s.dao.MaskMids(c); err != nil {
			continue
		}
		s.maskMid = mids
		log.Info("update mask mid(%v)", s.maskMid)
	}
}

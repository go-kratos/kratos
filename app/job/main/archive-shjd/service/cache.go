package service

import (
	"context"
	"strconv"

	"go-common/app/job/main/archive-shjd/model"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// UpdateCache is
func (s *Service) UpdateCache(old *model.Archive, nw *model.Archive, action string) (err error) {
	defer func() {
		if err == nil {
			s.notifyPub.Send(context.Background(), strconv.FormatInt(nw.AID, 10), &model.Notify{Table: _tableArchive, Nw: nw, Old: old, Action: action})
			return
		}
		// retry
		item := &model.RetryItem{
			Old:    old,
			Nw:     nw,
			Tp:     model.TypeForUpdateArchive,
			Action: action,
		}
		if err1 := s.PushItem(context.TODO(), item); err1 != nil {
			log.Error("s.PushItem(%+v) error(%+v)", item, err1)
			return
		}
	}()
	args := &arcmdl.ArgCache2{}
	args.Aid = nw.AID
	args.Tp = arcmdl.CacheUpdate
	if old == nil {
		// insert
		if nw.State >= 0 {
			args.Tp = arcmdl.CacheAdd
		}
	} else {
		if nw.State >= 0 {
			args.Tp = arcmdl.CacheAdd
		} else {
			args.Tp = arcmdl.CacheDelete
		}
		if nw.Mid != old.Mid {
			args.OldMid = old.Mid
		}
		if old.TypeID != nw.TypeID {
			fieldAgs := &arcmdl.ArgFieldCache2{Aid: nw.AID, TypeID: nw.TypeID, OldTypeID: old.TypeID}
			for cluster, arc := range s.arcRPCs {
				if err = arc.ArcFieldCache2(context.TODO(), fieldAgs); err != nil {
					log.Error("s.arcRPC.ArcFieldCache2(%s, %+v) error(%+v)", cluster, fieldAgs, err)
					return
				}
			}
		}
	}
	for cluster, arc := range s.arcRPCs {
		if err = arc.ArcCache2(context.TODO(), args); err != nil {
			log.Error("s.arcRPC.ArcCache2(%s,%+v) error(%v)", cluster, args, err)
			return
		}
	}
	return
}

// UpdateVideoCache is
func (s *Service) UpdateVideoCache(aid, cid int64) (err error) {
	defer func() {
		if err == nil {
			return
		}
		// retry
		item := &model.RetryItem{
			AID: aid,
			CID: cid,
			Tp:  model.TypeForUpdateVideo,
		}
		if err1 := s.PushItem(context.TODO(), item); err1 != nil {
			log.Error("s.PushItem(%+v) error(%+v)", item, err1)
			return
		}
	}()
	for cluster, arc := range s.arcRPCs {
		if err = arc.UpVideo2(context.TODO(), &arcmdl.ArgVideo2{Aid: aid, Cid: cid}); err != nil {
			if ecode.Cause(err).Equal(ecode.NothingFound) {
				err = nil
				return
			}
			err = errors.Wrapf(err, "s.arcRPC.UpVideo2 cluster(%s)", cluster)
		}
	}
	return
}

// DelteVideoCache del video cache
func (s *Service) DelteVideoCache(aid, cid int64) (err error) {
	defer func() {
		if err == nil {
			return
		}
		// retry
		item := &model.RetryItem{
			AID: aid,
			CID: cid,
			Tp:  model.TypeForDelVideo,
		}
		if err1 := s.PushItem(context.TODO(), item); err1 != nil {
			log.Error("s.PushItem(%+v) error(%+v)", item, err1)
			return
		}
	}()
	for cluster, arc := range s.arcRPCs {
		if err = arc.DelVideo2(context.TODO(), &arcmdl.ArgVideo2{Aid: aid, Cid: cid}); err != nil {
			if ecode.Cause(err).Equal(ecode.NothingFound) {
				err = nil
				return
			}
			err = errors.Wrapf(err, "s.arcRPC.VdelVideo2 cluster(%s)", cluster)
		}
	}
	return
}

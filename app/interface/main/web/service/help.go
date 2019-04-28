package service

import (
	"context"

	"go-common/app/interface/main/web/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const _firstPage = 1

var _emptyHelpList = make([]*model.HelpList, 0)

// HelpList get help menu list
func (s *Service) HelpList(c context.Context, pTypeID string) (res []*model.HelpList, err error) {
	if res, err = s.dao.HlCache(c, pTypeID); err != nil || len(res) == 0 {
		if res, err = s.dao.HelpList(context.Background(), pTypeID); err != nil {
			log.Error("s.do.HelpList(%s) error(%v)", pTypeID, err)
			return
		}
		if len(res) > 0 {
			s.cache.Do(c, func(c context.Context) {
				s.dao.SetHlCache(c, pTypeID, res)
			})
		}
	}
	return
}

// HelpDetail get help detail
func (s *Service) HelpDetail(c context.Context, fID, qTypeID string, keyFlag, pn, ps int) (resD []*model.HelpDeatil, resL []*model.HelpList, total int, err error) {
	remoteIP := metadata.String(c, metadata.RemoteIP)
	if resD, total, err = s.dao.DetailCache(c, qTypeID, keyFlag, pn, ps); err != nil || len(resD) == 0 {
		if resD, total, err = s.dao.HelpDetail(context.Background(), qTypeID, keyFlag, pn, ps, remoteIP); err != nil {
			log.Error("s.do.HelpDetail(%s,%d,%d,%d) error(%v)", qTypeID, keyFlag, pn, ps, err)
		}
		if pn == _firstPage && len(resD) > 0 {
			s.cache.Do(c, func(c context.Context) {
				s.dao.SetDetailCache(c, qTypeID, keyFlag, pn, ps, total, resD)
			})
		}
	}
	if fID == "" {
		resL = _emptyHelpList
	} else {
		if resL, err = s.HelpList(c, fID); err != nil {
			log.Error("s.HelpList(%s) error(%v)", fID, err)
		}
	}
	return
}

// HelpSearch get help search
func (s *Service) HelpSearch(c context.Context, pTypeID, keyWords string, keyFlag, pn, ps int) (res []*model.HelpDeatil, total int, err error) {

	if res, total, err = s.dao.HelpSearch(context.Background(), pTypeID, keyWords, keyFlag, pn, ps); err != nil {
		log.Error("s.do.HelpDetail(%s,%d,%d,%d) error(%v)", keyWords, keyFlag, pn, ps, err)
	}
	return
}

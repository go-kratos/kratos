package archive

import (
	"context"

	"go-common/app/interface/main/creative/model/archive"
	pubSvc "go-common/app/interface/main/creative/service"
	"go-common/library/ecode"
	"go-common/library/log"
)

// HistoryList get the history of aid
func (s *Service) HistoryList(c context.Context, mid, aid int64, ip string) (historys []*archive.ArcHistory, err error) {
	if historys, err = s.arc.HistoryList(c, mid, aid, ip); err != nil {
		log.Error("s.arc.HistoryList(%d,%d) err(%v)", mid, aid, err)
		return
	}
	for key, history := range historys {
		if history.Mid > 0 && history.Mid != mid {
			err = ecode.ArchiveOwnerErr
			return
		}
		historys[key].Cover = pubSvc.CoverURL(history.Cover)
	}
	return
}

// HistoryView get the history of hid
func (s *Service) HistoryView(c context.Context, mid, hid int64, ip string) (history *archive.ArcHistory, err error) {
	if history, err = s.arc.HistoryView(c, mid, hid, ip); err != nil {
		log.Error("s.arc.HistoryView(%d,%d) err(%v)", mid, hid, err)
		return
	}
	if history.Mid > 0 && history.Mid != mid {
		err = ecode.ArchiveOwnerErr
		history = nil
		return
	}
	history.Cover = pubSvc.CoverURL(history.Cover)
	return
}

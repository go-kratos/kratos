package service

import (
	"context"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
	"time"
)

//EditHistory 根据稿件的某条编辑历史id，获取当时完整的稿件、分p视频编辑历史
func (s *Service) EditHistory(c context.Context, hid int64) (h *archive.EditHistory, err error) {
	arcHistory, err := s.arc.HistoryByID(c, hid)
	if err != nil {
		log.Error("EditHistory s.arc.HistoryByID(hid(%d)) error(%v)", hid, err)
		return nil, err
	}

	vHistory, err := s.arc.VideoHistoryByHID(c, hid)
	if err != nil {
		log.Error("EditHistory s.arc.VideoHistoryByHID(hid(%d)) error(%v)", hid, err)
		return nil, err
	}

	h = &archive.EditHistory{
		ArcHistory: arcHistory,
		VHistory:   vHistory,
	}
	return
}

//AllEditHistory 根据aid获取 其所有的用户编辑历史
func (s *Service) AllEditHistory(c context.Context, aid int64) (hs []*archive.EditHistory, err error) {
	stime := time.Now().Add(time.Hour * 720 * -1)
	arcHistory, err := s.arc.HistoryByAID(c, aid, stime)
	if err != nil {
		log.Error("AllEditHistory s.arc.HistoryByAID(aid(%d)) error(%v)", aid, err)
		hs = []*archive.EditHistory{}
		return
	}

	var (
		videoHistory []*archive.VideoHistory
		prev         *archive.EditHistory
		total        int
	)

	total = len(arcHistory)
	hs = make([]*archive.EditHistory, total)
	for i := total - 1; i >= 0; i-- {
		h := arcHistory[i]
		videoHistory, err = s.arc.VideoHistoryByHID(c, h.ID)
		if err != nil {
			log.Error("AllEditHistory s.arc.VideoHistoryByHID(hid(%d), aid(%d)) error(%v)", h.ID, aid, err)
			return
		}
		one := &archive.EditHistory{
			ArcHistory: h,
			VHistory:   videoHistory,
		}

		//only show diff between next edit archive
		show, diff := one.Diff(prev)
		hs[i] = show
		if diff {
			prev = one
		}
	}

	return
}

package service

import (
	"context"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddTransferJob add job
func (s *Service) AddTransferJob(c context.Context, from, to, mid int64, offset float64, state int8) (err error) {
	_, err = s.dao.InsertTransferJob(c, from, to, mid, offset, 0)
	if err != nil {
		log.Error("s.dao.InsertTransferJob(%d, %d %d %f) error(%v)", from, to, mid, offset, err)
	}
	return
}

// TransferList transfer list
func (s *Service) TransferList(c context.Context, cid, state, pn, ps int64) (res []*model.TransList, total int64, err error) {
	var (
		aids []int64
		cids []int64
	)
	if res, total, err = s.dao.TransferList(c, cid, state, pn, ps); err != nil {
		log.Error("s.dao.TransferList(%d, %d) error(%v)", cid, state, err)
		return
	}
	if len(res) <= 0 {
		return
	}
	for _, idx := range res {
		cids = append(cids, idx.From)
	}
	subs, err := s.dao.Subjects(c, model.SubTypeVideo, cids)
	if err != nil {
		return
	}
	for _, idx := range subs {
		aids = append(aids, idx.Pid)
	}
	avm, err := s.dao.ArchiveVideos(c, aids)
	if err != nil {
		log.Error("s.dao.ArchiveInfo(aid:%v) error(%v)", aids, err)
		return
	}
	for _, idx := range res {
		sub, ok := subs[idx.From]
		if !ok {
			continue
		}
		info, ok := avm[sub.Pid] // get archive info
		if !ok {
			continue
		}
		idx.Title = info.Archive.Title
	}
	return
}

// ReTransferJob retransfer job
func (s *Service) ReTransferJob(c context.Context, id, mid int64) (err error) {
	job, err := s.dao.CheckTransferID(c, id)
	if err != nil {
		log.Error("dao.CheckTransferID(%d) err(%v)", id, err)
		return
	}
	if job.State != model.TransferJobStatFailed || job.MID != mid {
		err = ecode.RequestErr
		return
	}
	_, err = s.dao.SetTransferState(c, id, model.TransferJobStatInit)
	if err != nil {
		log.Error("dao.SetTransferState(id:%d,state:%d) err(%v)", id, model.TransferJobStatInit, err)
	}
	return
}

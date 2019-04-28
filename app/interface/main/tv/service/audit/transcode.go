package audit

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
)

const (
	_typeUGC = "ugc"
	_typePGC = "pgc"
)

type cidExistFunc = func(context.Context, int64) ([]int64, error)
type cidTransFunc = func(context.Context, []int64, int64) error
type reqTrans struct {
	CID       int64
	Action    int64
	CheckFunc cidExistFunc
	TransFunc cidTransFunc
}

// Transcode update the video/ep's transcoded status
func (s *Service) Transcode(req *model.ReqTransode) (err error) {
	var ctx = context.TODO()
	if req.ContType == _typePGC {
		err = s.transPGC(ctx, req.CID, req.Action)
	} else if req.ContType == _typeUGC {
		err = s.transUGC(ctx, req.CID, req.Action)
	} else {
		err = ecode.TvDangbeiWrongType
	}
	return
}

func commonTrans(ctx context.Context, req reqTrans) (err error) {
	var ids []int64
	if ids, err = req.CheckFunc(ctx, req.CID); err != nil {
		return
	}
	if len(ids) == 0 {
		return ecode.NothingFound
	}
	err = req.TransFunc(ctx, ids, req.Action)
	return
}

func (s *Service) transPGC(ctx context.Context, cid int64, action int64) (err error) {
	return commonTrans(ctx, reqTrans{
		CID:       cid,
		Action:    action,
		CheckFunc: s.auditDao.PgcCID,
		TransFunc: s.auditDao.PgcTranscode,
	})
}

func (s *Service) transUGC(ctx context.Context, cid int64, action int64) (err error) {
	return commonTrans(ctx, reqTrans{
		CID:       cid,
		Action:    action,
		CheckFunc: s.auditDao.UgcCID,
		TransFunc: s.auditDao.UgcTranscode,
	})
}

// ApplyPGC saves the pgc transcode apply time
func (s *Service) ApplyPGC(ctx context.Context, req *model.ReqApply) (err error) {
	return commonTrans(ctx, reqTrans{
		CID:       req.CID,
		Action:    req.ApplyTime,
		CheckFunc: s.auditDao.PgcCID,
		TransFunc: s.auditDao.ApplyPGC,
	})
}

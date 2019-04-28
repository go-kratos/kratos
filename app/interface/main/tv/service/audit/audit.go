package audit

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// HandleAudits treats the slice of IDList
func (s *Service) HandleAudits(ctx context.Context, vals []*model.IDList) (err error) {
	var tx *sql.Tx
	if tx, err = s.auditDao.BeginTran(ctx); err != nil {
		log.Error("audit HandleAudits BeginTran Err %v", err)
		return
	}
	for _, v := range vals {
		if err = s.handleAudit(ctx, v, tx); err != nil {
			log.Error("HandleAudits audit (%v), err %v", v, err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

// handleAudit checks the prefix to dispatch the task to ugc/pgc season/ep
func (s *Service) handleAudit(ctx context.Context, val *model.IDList, tx *sql.Tx) (err error) {
	var (
		op = new(model.AuditOp)
	)
	if err = op.FromIDList(val); err != nil {
		log.Error("audit handle Type error %v", val)
		return
	}
	switch op.ContentType {
	case model.UgcArc:
		if arcCMS, err := s.cmsDao.LoadArcMeta(ctx, op.KID); err == nil && arcCMS != nil && arcCMS.NotDeleted() {
			return auditCore(ctx, op, tx, s.auditDao.UpdateArc)
		}
	case model.UgcVideo:
		if videoCMS, err := s.cmsDao.LoadVideoMeta(ctx, op.KID); err == nil && videoCMS != nil && videoCMS.NotDeleted() {
			return auditCore(ctx, op, tx, s.auditDao.UpdateVideo)
		}
	case model.PgcSn:
		if snCMS, err := s.cmsDao.SnAuth(ctx, op.KID); err == nil && snCMS != nil && snCMS.NotDeleted() {
			return auditCore(ctx, op, tx, s.auditDao.UpdateSea)
		}
	case model.PgcEp:
		if epCMS, err := s.cmsDao.EpAuth(ctx, op.KID); err == nil && epCMS != nil && epCMS.NotDeleted() {
			return auditCore(ctx, op, tx, s.auditDao.UpdateCont)
		}
	default:
		log.Error("audit handle Content Type Error %s", op.ToMsg())
		return ecode.TvDangbeiWrongType
	}
	return ecode.NothingFound
}

type doAudit func(ctx context.Context, v *model.AuditOp, tx *sql.Tx) (err error)

func auditCore(c context.Context, v *model.AuditOp, tx *sql.Tx, updateFunc doAudit) (err error) {
	if err = updateFunc(c, v, tx); err != nil {
		log.Error("%s fail(%v)", v.ToMsg(), err)
	}
	return
}

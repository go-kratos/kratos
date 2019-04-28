package service

import (
	"context"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AuditSubtitle audit subtitle by uper pr assitant
func (s *Service) AuditSubtitle(c context.Context, oid int64, subtitleID int64, mid int64, pass bool, rejectComment string) (err error) {
	var (
		draft   *model.Subtitle
		subject *model.Subject
	)
	if draft, err = s.getSubtitle(c, oid, subtitleID); err != nil {
		log.Error("s.getSubtitle(oid:%v,subtitleID:%v),error(%v)", oid, subtitleID, err)
		return
	}
	if draft == nil {
		err = ecode.NothingFound
		return
	}
	// up主，协管有权限
	if subject, err = s.subject(c, draft.Type, draft.Oid); err != nil {
		return
	}
	// 非up主，无权限
	if subject.Mid != mid {
		err = ecode.SubtitlePermissionDenied
		return
	}
	if draft.Status != model.SubtitleStatusToAudit && draft.Status != model.SubtitleStatusPublish {
		err = ecode.SubtitleUnValid
		return
	}
	draft.RejectComment = rejectComment
	if !pass {
		if draft.Status == model.SubtitleStatusPublish {
			if err = s.backPubSubtitle(c, draft); err != nil {
				return
			}
			return
		}
		if err = s.auditReject(c, draft); err != nil {
			log.Error("params(draft:%+v).error(%v)", draft, err)
			return
		}
	} else {
		if err = s.auditPass(c, draft); err != nil {
			log.Error("params(draft:%+v).error(%v)", draft, err)
			return
		}
	}
	return
}

// auditReject subtitle submit
func (s *Service) auditReject(c context.Context, subtitle *model.Subtitle) (err error) {
	subtitle.Status = model.SubtitleStatusAuditBack
	subtitle.PubTime = time.Now().Unix()
	if err = s.dao.UpdateSubtitle(c, subtitle); err != nil {
		log.Error("params(%+v).error(%v)", subtitle, err)
		return
	}
	s.dao.DelSubtitleDraftCache(context.Background(), subtitle.Oid, subtitle.Type, subtitle.Mid, subtitle.Lan)
	s.dao.DelSubtitleCache(context.Background(), subtitle.Oid, subtitle.ID)
	return
}

func (s *Service) auditPass(c context.Context, subtitle *model.Subtitle) (err error) {
	var (
		tx          *sql.Tx
		subtitlePub *model.SubtitlePub
	)
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Error("params(subtitle:%+v).err(%v)", subtitle, err)
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("params(subtitle:%+v).err(%v)", subtitle, err)
			return
		}
	}()
	subtitle.Status = model.SubtitleStatusPublish
	subtitle.PubTime = time.Now().Unix()
	if tx, err = s.dao.BeginBiliDMTrans(c); err != nil {
		log.Error("error(%v)", err)
		return
	}
	if err = s.dao.TxUpdateSubtitle(tx, subtitle); err != nil {
		log.Error("params(%+v).error(%v)", subtitle, err)
		return
	}
	subtitlePub = &model.SubtitlePub{
		Oid:        subtitle.Oid,
		Type:       subtitle.Type,
		Lan:        subtitle.Lan,
		SubtitleID: subtitle.ID,
	}
	if err = s.dao.TxAddSubtitlePub(tx, subtitlePub); err != nil {
		log.Error("params(%+v).error(%v)", subtitlePub, err)
		return
	}
	if err = s.dao.DelSubtitleDraftCache(c, subtitle.Oid, subtitle.Type, subtitle.Mid, subtitle.Lan); err != nil {
		log.Error("DelSubtitleDraftCache.params(subtitle:%+v).err(%v)", subtitle, err)
		return
	}
	if err = s.dao.DelSubtitleCache(c, subtitle.Oid, subtitle.ID); err != nil {
		log.Error("DelSubtitleCache.params(subtitle:%+v).err(%v)", subtitle, err)
		return
	}
	if err = s.dao.DelVideoSubtitleCache(c, subtitle.Oid, subtitle.Type); err != nil {
		log.Error("DelVideoSubtitleCache.params(subtitle:%+v).err(%v)", subtitle, err)
		return
	}
	return
}

func (s *Service) backPubSubtitle(c context.Context, subtitle *model.Subtitle) (err error) {
	var (
		tx          *sql.Tx
		subtitleNew *model.Subtitle
		subtitlePub *model.SubtitlePub
	)
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			return
		}
	}()
	subtitle.Status = model.SubtitleStatusAuditBack
	subtitle.PubTime = time.Now().Unix()
	if tx, err = s.dao.BeginBiliDMTrans(c); err != nil {
		log.Error("error(%v)", err)
		return
	}
	if err = s.dao.TxUpdateSubtitle(tx, subtitle); err != nil {
		log.Error("params(%+v) error(%v)", subtitle, err)
		return
	}
	if subtitleNew, err = s.dao.TxGetSubtitleOne(tx, subtitle.Oid, subtitle.Type, subtitle.Lan); err != nil {
		log.Error("params(%+v) error(%v)", subtitle, err)
		return
	}
	subtitlePub = &model.SubtitlePub{
		Oid:  subtitle.Oid,
		Type: subtitle.Type,
		Lan:  subtitle.Lan,
	}
	if subtitleNew == nil {
		subtitlePub.IsDelete = true
	} else {
		subtitlePub.SubtitleID = subtitleNew.ID
	}
	if err = s.dao.TxAddSubtitlePub(tx, subtitlePub); err != nil {
		log.Error("params(%+v) error(%v)", subtitlePub, err)
		return
	}
	if err = s.dao.DelSubtitleCache(context.Background(), subtitle.Oid, subtitle.ID); err != nil {
		log.Error("params(oid:%v,subtitleID:%v) error(%v)", subtitle.Oid, subtitle.ID, err)
		return
	}
	if err = s.dao.DelVideoSubtitleCache(context.Background(), subtitle.Oid, subtitle.Type); err != nil {
		log.Error("params(oid:%v,subtitleID:%v) error(%v)", subtitle.Oid, subtitle.ID, err)
		return
	}
	return
}

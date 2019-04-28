package service

import (
	"context"
	"time"

	"go-common/app/admin/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// SubtitleStatusList .
func (s *Service) SubtitleStatusList(c context.Context) (res map[uint8]string, err error) {
	return model.StatusContent, nil
}

// CheckHasDraft .
func (s *Service) CheckHasDraft(c context.Context, subtitle *model.Subtitle) (ok bool, err error) {
	var (
		draftCount int64
	)
	if draftCount, err = s.dao.CountSubtitleDraft(c, subtitle.Oid, subtitle.Mid, subtitle.Lan, subtitle.Type); err != nil {
		log.Error("CheckHasDraft,params(subtitle:%+v),error(%v)", subtitle, err)
		return
	}
	if draftCount > 0 {
		ok = true
	}
	return
}

// RebuildSubtitle .
// need transtaion
// 1、更新自身状态
// 2、重新查询发布的字幕id，插入到发布表
// 3、删除缓存
func (s *Service) RebuildSubtitle(c context.Context, subtitle *model.Subtitle) (err error) {
	var (
		tx                *sql.Tx
		subtitlePublishID int64
		subtitlePub       *model.SubtitlePub
	)
	switch subtitle.Status {
	case model.SubtitleStatusDraft, model.SubtitleStatusToAudit:
		subtitle.PubTime = 0
	default:
		subtitle.PubTime = time.Now().Unix()
	}
	if tx, err = s.dao.BeginBiliDMTrans(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		if err = tx.Commit(); err != nil {
			return
		}
	}()
	if err = s.dao.TxUpdateSubtitle(tx, subtitle); err != nil {
		log.Error("RebuildSubtitle.TxUpdateSubtitle(subtitle:%+v),error(%v)", subtitle, err)
		return
	}
	if subtitlePublishID, err = s.dao.TxGetSubtitleID(tx, subtitle.Oid, subtitle.Type, subtitle.Lan); err != nil {
		log.Error("RebuildSubtitle.TxGetSubtitleID(params:%+v),error(%v)", subtitle, err)
		return
	}
	subtitlePub = &model.SubtitlePub{
		Oid:        subtitle.Oid,
		Type:       subtitle.Type,
		Lan:        subtitle.Lan,
		SubtitleID: subtitlePublishID,
	}
	if subtitlePublishID <= 0 {
		subtitlePub.IsDelete = true
	}
	if err = s.dao.TxUpdateSubtitlePub(tx, subtitlePub); err != nil {
		log.Error("RebuildSubtitle.TxUpdateSubtitlePub(subtitlePub:%+v),error(%v)", subtitlePub, err)
		return
	}
	return
}

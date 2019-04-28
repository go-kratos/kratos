package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) genSubtitleID(c context.Context) (subtitleID int64, err error) {
	subtitleID, err = s.seqRPC.ID(c, s.seqSubtitleArg)
	if err != nil {
		return
	}
	return
}

// SaveSubtitleDraft save subtitle
func (s *Service) SaveSubtitleDraft(c context.Context, aid, oid int64, tp int32, mid int64, lan string, submit, sign bool, originSubtitleID int64, data []byte) (detectErrs []*model.SubtitleDetectError, err error) {
	var (
		subject  *model.Subject
		draft    *model.Subtitle
		authorID int64
	)
	if subject, err = s.subject(c, tp, oid); err != nil {
		log.Error("params(tp:%v, oid:%v).error(%v)", tp, oid, err)
		return
	}
	if err = s.checkSubtitleLan(c, lan); err != nil {
		return
	}
	if err = s.checkAidOid(c, aid, oid); err != nil {
		return
	}
	if err = s.checkSubtitlePermission(c, aid, oid, tp, mid, subject); err != nil {
		return
	}
	// TODO remove error
	if detectErrs, err = s.checkSubtitleData(c, aid, oid, data); err != nil {
		return
	}
	if len(detectErrs) > 0 {
		return
	}
	if err = s.checkSubtitleLocked(c, submit, oid, tp, lan, mid); err != nil {
		return
	}
	if authorID, err = s.checkSubtitleAuthor(c, oid, originSubtitleID, lan, mid); err != nil {
		return
	}
	status := model.SubtitleStatusDraft
	if submit {
		status = model.SubtitleStatusCheckToAudit
		if mid == subject.Mid {
			status = model.SubtitleStatusCheckPublish
		}
	}
	if draft, err = s.buildSubtitleDraft(c, oid, tp, mid, authorID, lan, data, status, sign); err != nil {
		return
	}
	if err = s.addSubtitleDraft(c, draft); err != nil {
		return
	}
	if status == model.SubtitleStatusCheckToAudit || status == model.SubtitleStatusCheckPublish {
		s.dao.SendSubtitleCheck(c, draft.CheckSum, &model.SubtitleCheckMsg{
			Oid:        oid,
			SubtitleID: draft.ID,
		})
	}
	return
}

func (s *Service) addSubtitleDraft(c context.Context, draft *model.Subtitle) (err error) {
	if draft.ID > 0 {
		if err = s.updateSubtitle(c, draft); err != nil {
			return
		}
	} else {
		if draft.ID, err = s.genSubtitleID(c); err != nil {
			return
		}
		if _, err = s.addSubtitle(c, draft); err != nil {
			return
		}
	}
	return
}

// buildSubtitleDraft when save draft or save to submit
func (s *Service) buildSubtitleDraft(c context.Context, oid int64, tp int32, mid, authorID int64, lan string, data []byte, status model.SubtitleStatus, sign bool) (draft *model.Subtitle, err error) {
	var (
		subtitleURL string
		checkSum    string
		subject     *model.Subject
		lanCode     int64
	)
	if lanCode = s.subtitleLans.GetByLan(lan); lanCode <= 0 {
		err = ecode.SubtitleIllegalLanguage
		return
	}
	if draft, err = s.getSubtitlDraft(c, oid, tp, mid, uint8(lanCode)); err != nil {
		log.Error("params(oid:%v,tp:%v,mid:%v,lanCode:%v).error(%v)", oid, tp, mid, lanCode, err)
		return
	}
	if draft == nil {
		if subject, err = s.subject(c, tp, oid); err != nil {
			log.Error("params(oid:%v,tp:%v).error(%v)", oid, tp, err)
			return
		}
		draft = &model.Subtitle{
			Oid:      oid,
			Type:     tp,
			Mid:      mid,
			Aid:      subject.Pid,
			Lan:      uint8(lanCode),
			AuthorID: mid,
			UpMid:    subject.Mid,
			PubTime:  0,
			IsSign:   sign,
			Status:   model.SubtitleStatusDraft,
		}
	}
	if draft.Status != model.SubtitleStatusDraft && draft.Status != model.SubtitleStatusToAudit && draft.Status != model.SubtitleStatusCheckToAudit {
		err = ecode.SubtitlePermissionDenied
		return
	}
	sha := sha1.Sum(data)
	if checkSum = hex.EncodeToString(sha[:]); checkSum != draft.CheckSum {
		if subtitleURL, err = s.dao.UploadBfs(c, "", data); err != nil {
			log.Error("UploadBfs.error(%v)", err)
			return
		}
		draft.SubtitleURL = subtitleURL
		draft.CheckSum = checkSum
	}
	draft.Status = status
	draft.IsSign = sign
	if status == model.SubtitleStatusCheckPublish {
		draft.PubTime = time.Now().Unix()
	}
	draft.AuthorID = authorID
	draft.RejectComment = ""
	return
}

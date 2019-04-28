package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// txUpArcAttrs update archive attrs
func (s *Service) txUpArcAttrs(tx *sql.Tx, a *archive.Archive, attrs map[uint]int32, redirectURL string) (cont []string, err error) {
	const template = "[%s]从[%s]设为[%s]"
	var yesOrNo = map[int32]string{archive.AttrYes: "是", archive.AttrNo: "否"}
	for bit, attr := range attrs {
		if bit == archive.AttrBitJumpURL && (attr == archive.AttrNo || redirectURL == "") {
			redirectURL = ""
			attr = archive.AttrNo
		}
		var arcRows int64
		if arcRows, err = s.arc.TxUpArcAttr(tx, a.Aid, bit, attr); err != nil {
			log.Error("s.arc.TxUpArcAttr(%d,%d,%d) error(%v)", a.Aid, bit, attr, err)
			return
		}
		if bit == archive.AttrBitJumpURL {
			var additRows int64
			if additRows, err = s.arc.TxUpAdditRedirect(tx, a.Aid, redirectURL); err != nil {
				log.Error("s.arc.TxUpAddit(%d,%s) error(%v)", a.Aid, redirectURL, err)
				return
			}
			if arcRows != 0 || additRows != 0 {
				log.Info("aid(%d) update addit redirectURL(%s)", a.Aid, redirectURL)
				var desc = redirectURL
				if desc == "" {
					desc = "空"
				}
				cont = append(cont, fmt.Sprintf("[%s]设为[%s]", "跳转地址", desc))
			}
		}
		a.AttrSet(attr, bit)
		if arcRows <= 0 {
			continue
		}
		cont = append(cont, fmt.Sprintf(template, archive.BitDesc(bit), yesOrNo[^attr&1], yesOrNo[attr]))
		log.Info("aid(%d) update archive bit(%d) bitdesc(%s) attrs(%d)", a.Aid, bit, archive.BitDesc(bit), attr)
	}
	return
}

// txUpArcForbidAttrs update archive forbid attrs
func (s *Service) txUpArcForbidAttrs(c context.Context, tx *sql.Tx, a *archive.Archive, forbidAttrs map[string]map[uint]int32) (conts []string, err error) {
	const template = "[%s]从[%s]设为[%s]"
	var (
		yesOrNo = map[int32]string{archive.AttrYes: "是", archive.AttrNo: "否"}
		forbid  *archive.ForbidAttr
		change  bool
	)
	// forbid
	if forbid, err = s.arc.Forbid(c, a.Aid); err != nil {
		log.Error("s.arc.Forbid(%d) error(%v)", a.Aid, err)
		return
	}
	for name, attrs := range forbidAttrs {
		for bit, attr := range attrs {
			if change = forbid.SetAttr(name, attr, bit); change {
				conts = append(conts, fmt.Sprintf(template, archive.ForbidBitDesc(name, bit), yesOrNo[^attr&1], yesOrNo[attr]))
				log.Info("aid(%d) update archive forbid name(%s) bit(%d) bitdesc(%s) attrs(%d)", a.Aid, name, bit, archive.ForbidBitDesc(name, bit), attr)
			}
		}
	}
	forbid.Convert()
	if _, err = s.arc.TxUpForbid(tx, forbid); err != nil {
		log.Error("s.arc.TxUpForbid(%+v) error(%v)", forbid, err)
	}
	return
}

//txAddFirstPass 添加稿件第一次过审的记录
func (s *Service) txAddFirstPass(c context.Context, tx *sql.Tx, aid int64, state int8) (firstPass bool, err error) {
	if !archive.NormalState(state) || s.hadPassed(c, aid) {
		return
	}

	if err = s.arc.AddFirstPass(tx, aid); err != nil {
		log.Error("txUpArcState s.arc.AddFirstPass error(%v) aid(%d)", err, aid)
		return
	}

	firstPass = true
	return
}

//txUpArcState 更新稿件属性时，联动添加第一次过审记录
func (s *Service) txUpArcState(c context.Context, tx *sql.Tx, aid int64, state int8) (firstPass bool, err error) {
	if _, err = s.arc.TxUpArcState(tx, aid, state); err != nil {
		log.Error("txUpArcState s.arc.TxUpArcState error(%v) aid(%d) state(%d)", err, aid, state)
		return
	}

	if firstPass, err = s.txAddFirstPass(c, tx, aid, state); err != nil {
		log.Error("txUpArcState s.txAddFirstPass error(%v) aid(%d) state(%d)", err, aid, state)
		return
	}
	return
}

// txUpArcMainState update archive states
func (s *Service) txUpArcMainState(c context.Context, tx *sql.Tx, aid, forward int64, typeID, access int16, state, round int8, reason string) (racs int16, err error) {
	log.Info("aid(%d) get archive state(%d)", aid, state)
	if ok := s.isAccess(c, aid); ok && access == archive.AccessDefault {
		access = archive.AccessMember
	}
	if _, err = s.txUpArcState(c, tx, aid, state); err != nil {
		log.Error("txUpArcMainState s.txUpArcState error(%v) aid(%d) state(%d)", err, aid, state)
		return
	}
	log.Info("aid(%d) update archive state(%d)", aid, state)
	if _, err = s.arc.TxUpArcAccess(tx, aid, access); err != nil {
		log.Error("s.arc.TxUpArcAccess(%d,%d) error(%v)", aid, access, err)
		return
	}
	racs = access
	log.Info("aid(%d) update archive access(%d)", aid, access)
	if typeID != 0 {
		if _, err = s.arc.TxUpArcTypeID(tx, aid, typeID); err != nil {
			log.Error("s.arc.TxUpArcTypeID(%d,%d) error(%v)", aid, typeID, err)
			return
		}
		log.Info("aid(%d) update archive type_id(%d)", aid, typeID)
	}
	if _, err = s.arc.TxUpArcRound(tx, aid, round); err != nil {
		log.Error("s.arc.TxUpArcRound(%d,%d) error(%v)", aid, round, err)
		return
	}
	log.Info("aid(%d) update archive round(%d)", aid, round)
	if _, err = s.arc.TxUpArcReason(tx, aid, forward, reason); err != nil {
		log.Error("s.arc.TxUpArcReason(%d,%d,%s) error(%v)", aid, forward, reason, err)
		return
	}
	log.Info("aid(%d) update archive reason(%s) forward_id(%d)", aid, reason, forward)
	return
}

// txUpArcDelay update archive delay
func (s *Service) txUpArcDelay(c context.Context, tx *sql.Tx, aid, mid int64, state int8, isDelay bool, dTime xtime.Time) (rs int8, cont string, err error) {
	rs = state
	var delay *archive.Delay
	if delay, _ = s.arc.Delay(c, aid, archive.DelayTypeForUser); delay == nil && dTime <= 0 {
		return
	}
	if !isDelay || archive.NotAllowDelay(state) {
		if _, err = s.arc.TxDelDelay(tx, aid, archive.DelayTypeForUser); err != nil {
			log.Error("s.arc.TxDelDelay(%d) error(%v)", aid, err)
			return
		}
		cont = archive.Operformat(archive.OperTypeDelay, dTime.Time().Format("2006-01-02 15:04:05"), "无", archive.OperStyleOne)
		log.Info("aid(%d) err delay time delete archive_delay", aid)
		return
	}
	if dTime <= 0 && delay != nil {
		dTime = delay.DTime
	}
	if _, err = s.arc.TxUpDelay(tx, mid, aid, state, archive.DelayTypeForUser, dTime); err != nil {
		log.Error("s.arc.TxUpDelay(%d，%d,%d,%d,%d) error(%v)", aid, mid, state, archive.DelayTypeForUser, dTime, err)
		return
	}
	if archive.NormalState(state) {
		rs = archive.StateForbidUserDelay
	}
	if delay != nil && dTime != delay.DTime {
		cont = archive.Operformat(archive.OperTypeDelay, delay.DTime.Time().Format("2006-01-02 15:04:05"), dTime.Time().Format("2006-01-02 15:04:05"), archive.OperStyleOne)
	} else if delay == nil {
		cont = archive.Operformat(archive.OperTypeDelay, "无", dTime.Time().Format("2006-01-02 15:04:05"), archive.OperStyleOne)
	}
	log.Info("aid(%d) second_round submit update archive_delay mid(%d) state(%d) type(%d) dtime(%v)", aid, mid, state, archive.DelayTypeForUser, dTime)
	return
}

// TxUpArchiveAttr update archive attr by aid.
func (s *Service) TxUpArchiveAttr(c context.Context, tx *sql.Tx, a *archive.Archive, aid, uid int64, attrs map[uint]int32, forbidAttrs map[string]map[uint]int32, redirectURL string) (conts []string, err error) {
	var cont []string
	log.Info("aid(%d) begin tran change attribute", aid)
	if cont, err = s.txUpArcAttrs(tx, a, attrs, redirectURL); err != nil {
		log.Error("s.txUpArcAttrs(%d) error(%v)", aid, err)
		return
	}
	conts = append(conts, cont...)
	if _, err = s.txUpArcForbidAttrs(c, tx, a, forbidAttrs); err != nil {
		log.Error("s.txUpArcForbidAttrs(%d) error(%v)", aid, err)
		return
	}
	log.Info("aid(%d) end tran change attribute", aid)
	return
}

// txUpVideoStatus update video status by vid and cid.
func (s *Service) txUpVideoStatus(tx *sql.Tx, vid, cid int64, status int16) (err error) {
	//update archive_video
	if _, err = s.arc.TxUpVideoStatus(tx, vid, status); err != nil {
		log.Error("s.arc.TxUpVideoStatus vid(%d) status(%d) error(%v)", vid, status, err)
		return
	}
	//update video
	if _, err = s.arc.TxUpStatus(tx, cid, status); err != nil {
		log.Error("s.arc.TxUpStatus cid(%d) status(%d)  error(%v)", cid, status, err)
		return
	}
	//update archive_video_relation to 0
	if _, err = s.arc.TxUpRelationState(tx, vid, archive.VideoStatusOpen); err != nil {
		log.Error("s.arc.TxUpRelationState cid(%d) status(%d)  error(%v)", cid, archive.VideoStatusOpen, err)
		return
	}
	log.Info("vid(%d) cid(%d) update video status(%d)", vid, cid, status)
	return
}

// txAddVideo insert video get vid.
func (s *Service) txAddVideo(tx *sql.Tx, v *archive.Video) (vid int64, err error) {
	if vid, err = s.arc.TxAddVideo(tx, v); err != nil {
		log.Error("s.arc.TxAddVideo video(%+v) error(%v)", v, err)
		return
	}
	v.ID = vid
	if _, err = s.arc.TxAddRelation(tx, v); err != nil {
		log.Error("s.arc.TxAddRelation video(%+v) error(%v)", v, err)
		return
	}
	log.Info("aid(%d) update video vid(%d) cid(%d) index(%d) title(%s) desc(%s) filename(%s) srctype(%s)", v.Aid, vid, v.Cid, v.Index, v.Title, v.Desc, v.Filename, v.SrcType)
	return
}

// txUpVideo update video title and desc by vid.
func (s *Service) txUpVideo(tx *sql.Tx, vid int64, title, desc string) (err error) {
	if _, err = s.arc.TxUpVideo(tx, vid, title, desc); err != nil {
		log.Error("s.arc.TxUpVideo vid(%d) title(%s) desc(%s) error(%v)", vid, title, desc, err)
		return
	}
	if _, err = s.arc.TxUpRelation(tx, vid, title, desc); err != nil {
		log.Error("s.arc.TxUpRelation vid(%d) title(%s) desc(%s) error(%v)", vid, title, desc, err)
		return
	}
	log.Info("vid(%d) update video title(%s) desc(%s)", vid, title, desc)
	return
}

// txUpVideoLink update video weblink by vid and cid.
func (s *Service) txUpVideoLink(tx *sql.Tx, vid, cid int64, webLink string) (err error) {
	if _, err = s.arc.TxUpVideoLink(tx, vid, webLink); err != nil {
		log.Error("s.arc.TxUpVideoLink(%d,%s)", vid, webLink)
		return
	}
	if _, err = s.arc.TxUpWebLink(tx, cid, webLink); err != nil {
		log.Error("s.arc.TxUpWebLink(%d,%s)", cid, webLink)
		return
	}
	log.Info("vid(%d) cid(%d) update webLink(%s)", vid, cid, webLink)
	return
}

// txDelVideo delete video by vid.
func (s *Service) txDelVideo(tx *sql.Tx, vp *archive.VideoParam) (err error) {
	if _, err = s.arc.TxUpVideoStatus(tx, vp.ID, archive.VideoStatusDelete); err != nil {
		log.Error("s.arc.TxUpVideoStatus(%d,%d)", vp.ID, archive.VideoStatusDelete)
		return
	}
	if _, err = s.arc.TxUpRelationState(tx, vp.ID, archive.VideoStatusDelete); err != nil {
		log.Error("s.arc.TxUpRelationState(%d,%d)", vp.ID, archive.VideoStatusDelete)
		return
	}
	log.Info("del video cid(%d) vid(%d)", vp.Cid, vp.ID)
	return
}

// txUpVideoIndex update video index by vid.
func (s *Service) txUpVideoIndex(tx *sql.Tx, vid int64, index int) (err error) {
	if _, err = s.arc.TxUpVideoIndex(tx, vid, index); err != nil {
		log.Error("s.arc.TxUpVideoIndex(%d,%d)", vid, index)
		return
	}
	if _, err = s.arc.TxUpRelationOrder(tx, vid, index); err != nil {
		log.Error("s.arc.TxUpRelationOrder(%d,%d)", vid, index)
		return
	}
	log.Info("vid(%d) update index(%d)", vid, index)
	return
}

func (s *Service) txUpVideoAttr(c context.Context, tx *sql.Tx, vid, cid int64, attrs map[uint]int32) (conts []string, attribute int32, err error) {
	var v *archive.Video
	if v, err = s.arc.NewVideoByID(c, vid); err != nil || v == nil {
		log.Error("s.arc.VideoByID(%d) error(%v)", vid, err)
		return
	}
	const template = "[%s]从[%s]设为[%s]"
	var yesOrNo = map[int32]string{archive.AttrYes: "是", archive.AttrNo: "否"}
	for bit, attr := range attrs {
		var rows int64
		if rows, err = s.arc.TxUpVideoAttr(tx, vid, bit, attr); err != nil {
			log.Error("s.arc.TxUpVideoAttr id(%d) bit(%d) attr(%d)  error(%v)", vid, bit, attr, err)
			return
		}
		if _, err = s.arc.TxUpAttr(tx, cid, bit, attr); err != nil {
			log.Error("s.arc.TxUpAttr cid(%d) bit(%d) attr(%d)  error(%v)", cid, bit, attr, err)
			return
		}
		v.AttrSet(attr, bit)
		if rows <= 0 {
			continue
		}
		conts = append(conts, fmt.Sprintf(template, archive.BitDesc(bit), yesOrNo[^attr&1], yesOrNo[attr]))
		log.Info("vid(%d) update video bit(%d) bitdesc(%s) attrs(%d)", vid, bit, archive.BitDesc(bit), attr)
	}
	attribute = v.Attribute
	return
}

// txUpVideoAudit update video audit by videoParam.
func (s *Service) txUpVideoAudit(tx *sql.Tx, vp *archive.VideoParam) (err error) {
	if err = s.txUpVideoStatus(tx, vp.ID, vp.Cid, vp.Status); err != nil {
		log.Error("s.arc.TxUpVideoStatus id(%d) cid(%d) status(%d) error(%v)", vp.ID, vp.Cid, vp.Status, err)
		return
	}
	log.Info("aid(%d) vid(%d) update video status(%d)", vp.Aid, vp.ID, vp.Status)
	if _, err = s.arc.TxAddAudit(tx, vp.Aid, vp.ID, vp.TagID, vp.Oname, vp.Note, vp.Reason); err != nil {
		log.Error("s.arc.TxAddAudit(%d,%d,%d,%s,%s,%s)", vp.Aid, vp.ID, vp.TagID, vp.Oname, vp.Note, vp.Reason)
		return
	}
	log.Info("aid(%d) vid(%d) audit log tag(%d) oname(%s) note(%s) reason(%s)", vp.Aid, vp.ID, vp.TagID, vp.Oname, vp.Note, vp.Reason)
	return
}

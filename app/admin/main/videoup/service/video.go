package service

import (
	"context"
	"runtime"
	"strings"
	"time"

	"fmt"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// VideoAudit upload first_round info.
func (s *Service) VideoAudit(c context.Context, vp *archive.VideoParam, attrs map[uint]int32) (err error) {
	var (
		tx      *sql.Tx
		qaVideo *archive.QAVideo
		qaErr   error
	)
	defer func() {
		if r := recover(); r != nil {
			if tx != nil {
				tx.Rollback()
			}

			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			log.Error("wocao jingran recover le error(%v) panic stack(%s)", r, string(buf[:n]))
		}
	}()

	//只有一审任务会去新增质检任务
	if vp.TaskID > 0 {
		qaVideo, qaErr = s.fetchQAVideo(c, vp)
		if qaErr != nil {
			log.Error("VideoAudit s.fetchQAVideo error(%v) aid(%d) cid(%d)", qaErr, vp.Aid, vp.Cid)
		}
	}

	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	if err = s.txUpVideo(tx, vp.ID, vp.Title, vp.Desc); err != nil {
		tx.Rollback()
		log.Error("s.txUpVideo(%d,%s,%s) error(%v)", vp.ID, vp.Title, vp.Desc, err)
		return
	}
	log.Info("aid(%d) update video  vid(%d) Title(%s) Desc(%s)  status(%d) ", vp.Aid, vp.ID, vp.Title, vp.Desc, vp.Status)
	var (
		operConts   []string
		attr        int32
		attrChanged bool
	)
	if operConts, attr, err = s.txUpVideoAttr(c, tx, vp.ID, vp.Cid, attrs); err != nil {
		tx.Rollback()
		return
	}
	if len(operConts) > 0 {
		attrChanged = true
	}
	if err = s.txUpVideoAudit(tx, vp); err != nil {
		tx.Rollback()
		log.Error("s.txUpVideoAudit(vp(%+v)) error(%v)", vp, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	oper := &archive.VideoOper{AID: vp.Aid, UID: vp.UID, VID: vp.ID, Attribute: attr, Status: vp.Status, Remark: vp.Note}
	operConts = append(operConts, s.diffVideoOper(vp)...)
	oper.Content = strings.Join(operConts, "，")
	s.addVideoOper(c, oper)
	// databus
	s.busFirstRoundByPlayurl(c, vp, attrChanged)
	go s.busSecondRoundUpCredit(vp.Aid, vp.Cid, vp.Mid, vp.UID, int8(vp.Status), 0, vp.ReasonID, vp.Reason)
	// log  set attr
	vp.Attribute = attr
	s.sendVideoLog(c, vp, oper.Content)
	if qaVideo != nil && qaErr == nil {
		qaVideo.Attribute = attr
		s.addQAVideo(c, qaVideo)
	}
	return
}

func (s *Service) busFirstRoundByPlayurl(c context.Context, vp *archive.VideoParam, attrChanged bool) (err error) {
	arc, err := s.arc.Archive(c, vp.Aid)
	if err != nil {
		log.Error("s.arc.Archive() (%d) error(%v)", vp.Aid, err)
		return
	}
	playurl := vp.Playurl
	if len(playurl) == 0 {
		playurl, err = s.arc.VideoPlayurl(c, vp.Cid)
		if err != nil {
			log.Error("s.busfirstroundbyplayurl get playurl by vid(%d) error(%v)", vp.ID, err)
			return
		}
		if len(playurl) == 0 {
			err = fmt.Errorf("video(%d) not exist", vp.ID)
			log.Error("s.busfirstroundbyplayurl get playurl by vid(%d) error(%v)", vp.ID, err)
			return
		}
	}
	log.Info("begin to bus first round by playurl(%s) vp(%v) err(%v)", playurl, vp, err)

	//adminChange=true发送报备邮件
	adminChange := false
	if vp.Status == archive.VideoStatusRecycle || vp.Status == archive.VideoStatusLock || attrChanged {
		adminChange = true
	}
	//get fans
	fans := int64(0)
	if pfl, _ := s.profile(c, arc.Mid); pfl != nil {
		fans = pfl.Follower
	}
	if strings.Contains(playurl, "/ugc/") {
		err = s.busUGCFirstRound(vp.Aid, vp.Filename, vp.EncodePurpose, vp.VideoDesign, vp.Status, vp.Encoding, vp.RegionID, arc.TypeID, fans, adminChange)
	} else {
		err = s.busFirstRound(vp.Aid, vp.Filename, vp.EncodePurpose, vp.VideoDesign, vp.Status, vp.Encoding, vp.RegionID, arc.TypeID, fans, adminChange)
	}
	return
}

// BatchVideo  batch async video audit.
func (s *Service) BatchVideo(c context.Context, vps []*archive.VideoParam, action string) (err error) {
	var mp = &archive.MultSyncParam{}
	var ok bool
	for _, vp := range vps {
		mp.Action = action
		mp.VideoParam = vp
		if ok, err = s.busCache.PushMultSync(c, mp); err != nil {
			log.Error("s.busCache.PushMultSync(vp(%+v)) error(%v)", vp, err)
			return
		}
		if !ok {
			log.Warn("s.busCache.PushMultSync(vp(%+v))", vp)
			continue
		}
	}
	return
}

func (s *Service) dealVideo(c context.Context, vp *archive.VideoParam) (err error) {
	var tx *sql.Tx
	var attrChanged bool
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if err = s.txUpVideoAudit(tx, vp); err != nil {
		tx.Rollback()
		log.Error("s.txUpVideoAudit(vp(%+v)) error(%v)", vp, err)
		return
	}
	//todo 批量支持家长模式 级联通过类tag逻辑
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	oper := &archive.VideoOper{AID: vp.Aid, UID: vp.UID, VID: vp.ID, Attribute: vp.Attribute, Status: vp.Status, Remark: vp.Note}
	operConts := s.diffVideoOper(vp)
	oper.Content = strings.Join(operConts, "，")
	s.addVideoOper(c, oper)
	s.busFirstRoundByPlayurl(c, vp, attrChanged)
	s.sendVideoLog(c, vp, oper.Content)
	return
}

// UpVideo update archive_video title && desc
func (s *Service) UpVideo(c context.Context, vp *archive.VideoParam) (err error) {
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if vp.ID != 0 {
		if err = s.txUpVideo(tx, vp.ID, vp.Title, vp.Desc); err != nil {
			tx.Rollback()
			log.Error("s.txUpVideo(%d,%s,%s) error(%v)", vp.ID, vp.Title, vp.Desc, err)
			return
		}
		vp.Note = "审核人员修改分P标题和简介!"
	} else {
		var (
			v     *archive.Video
			ctime = xtime.Time(time.Now().Unix())
		)
		v = &archive.Video{Aid: vp.Aid, Title: vp.Title, Desc: vp.Desc, Filename: vp.Filename, SrcType: vp.SrcType, Cid: vp.Cid,
			Duration: vp.Duration, Filesize: vp.Filesize, Resolutions: vp.Resolutions, Index: vp.Index, Status: vp.Status, Playurl: vp.Playurl,
			Attribute: vp.Attribute, FailCode: vp.FailCode, XcodeState: vp.XcodeState, WebLink: vp.WebLink, CTime: ctime, MTime: ctime}
		if vp.ID, err = s.txAddVideo(tx, v); err != nil {
			tx.Rollback()
			log.Error("s.arc.txAddVideo(%+v) error(%v)", v, err)
			return
		}
		vp.Note = "审核人员添加分P!"
	}
	if _, err = s.arc.TxAddAudit(tx, vp.Aid, vp.ID, vp.TagID, vp.Oname, vp.Note, vp.Reason); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxAddAudit(%d,%d,%d,%s,%s,%s)", vp.Aid, vp.ID, vp.TagID, vp.Oname, vp.Note, vp.Reason)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.busSecondRound(vp.Aid, 0, false, false, false, false, false, false, "", nil)
	s.sendVideoLog(c, vp, vp.Note)
	return
}

// UpWebLink update archive_video weblink
func (s *Service) UpWebLink(c context.Context, vp *archive.VideoParam) (err error) {
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if err = s.txUpVideoLink(tx, vp.ID, vp.Cid, vp.WebLink); err != nil {
		tx.Rollback()
		log.Error("s.txUpVideoLink(%d,%d,%s)", vp.ID, vp.Cid, vp.WebLink)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.busSecondRound(vp.Aid, 0, false, false, false, false, false, false, "", nil)
	s.sendVideoLog(c, vp, "更新weblink")
	return
}

// DelVideo delete archive_video.
func (s *Service) DelVideo(c context.Context, vp *archive.VideoParam) (err error) {
	var vs []*archive.Video
	if vs, _ = s.arc.NewVideosByAid(c, vp.Aid); len(vs) <= 0 {
		err = ecode.ArchiveNotExist
		return
	}
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if err = s.txDelVideo(tx, vp); err != nil {
		tx.Rollback()
		log.Error("s.txDelVideo(%d)", vp.ID)
		return
	}
	indexOrder := 1
	for _, v := range vs {
		if v.ID == vp.ID {
			continue
		}
		if err = s.txUpVideoIndex(tx, v.ID, indexOrder); err != nil {
			tx.Rollback()
			log.Error("s.txUpVideoIndex(%d,%d)", v.ID, indexOrder)
			return
		}
		indexOrder++
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.busSecondRound(vp.Aid, 0, false, false, false, false, false, false, "", nil)
	s.sendVideoLog(c, vp, "删除视频")
	return
}

// ChangeIndex change archive_video index.
func (s *Service) ChangeIndex(c context.Context, lo *archive.IndexParam) (err error) {
	if vs, _ := s.arc.NewVideosByAid(c, lo.Aid); len(vs) <= 0 {
		err = ecode.ArchiveNotExist
		return
	}
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	for _, v := range lo.ListOrder {
		if err = s.txUpVideoIndex(tx, v.ID, v.Index); err != nil {
			tx.Rollback()
			log.Error("s.arc.TxUpVideoIndex(%d,%d)", v.ID, v.Index)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.busSecondRound(lo.Aid, 0, false, false, false, false, false, false, "", nil)
	return
}

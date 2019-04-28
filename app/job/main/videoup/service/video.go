package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/message"
	"go-common/app/job/main/videoup/model/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

// bvc2VuConsumer is bvc 2 videoup message consumer.
func (s *Service) bvc2VuConsumer() {
	defer s.wg.Done()
	var (
		msgs = s.bvc2VuSub.Messages()
		err  error
		c    = context.TODO()
	)
	for {
		func() {
			msg, ok := <-msgs
			if !ok {
				log.Error("s,bvc2VuSub.Message closed")
				return
			}
			defer s.Rescue(string(msg.Value))
			msg.Commit()
			s.bvc2VuMo++
			m := &message.BvcVideo{}
			if err = json.Unmarshal(msg.Value, m); err != nil {
				log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
				return
			}
			if time.Now().Unix()-m.Timestamp > s.c.BvcConsumeTimeout {
				log.Info("bvcMessage consume delayed! key(%s) value(%s) partition(%d) offset(%d) route(%s) commit start", msg.Key, msg.Value, msg.Partition, msg.Offset, m.Route)
				s.bvc2VuDelayMo++
			}
			log.Info("bvcMessage key(%s) value(%s) partition(%d) offset(%d) route(%s) commit start", msg.Key, msg.Value, msg.Partition, msg.Offset, m.Route)
			s.promDatabus.Incr(m.Route)
			switch m.Route {
			case message.RouteUploadInfo:
				s.promVideoS.Incr("xcode_sd")
				s.uploadInfo(c, m)
			case message.RouteXcodeSDFail:
				err = s.xcodeSDFail(c, m)
			case message.RouteXcodeSdFinish:
				s.promVideoE.Incr("xcode_sd")
				s.promVideoS.Incr("first_round")
				err = s.xcodeSDFinish(c, m)
			case message.RouteXcodeHDFail:
				log.Info("bvcMessage key(%s) value(%s) partition(%d) offset(%d) route(%s) before enter func", msg.Key, msg.Value, msg.Partition, msg.Offset, m.Route)
				err = s.xcodeHDFail(c, m)
			case message.RouteXcodeHDFinish:
				s.promVideoE.Incr("xcode_hd")
				s.promVideoS.Incr("dispatch")
				err = s.xcodeHDFinish(c, m)
			case message.RouteDispatchRunning:
				err = s.dispatchRunning(c, m)
			case message.RouteDispatchFinish:
				s.promVideoE.Incr("dispatch")
				err = s.dispatchFinish(c, m)
			case message.RouteVideoshotpv:
				err = s.videoshotPv(c, m)
			default:
				log.Warn("bvc2VuConsumer unknown route(%s) message(%s)", m.Route, msg.Value)
			}
			if err == nil {
				log.Info("bvcMessage key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
			} else {
				log.Error("bvcMessage key(%s) value(%s) partition(%d) offset(%d) no commit error(%v)", msg.Key, msg.Value, msg.Partition, msg.Offset, err)
			}
		}()
	}
}

func (s *Service) videoshotPv(c context.Context, m *message.BvcVideo) (err error) {
	var count = len(m.ImgURLs)
	if count == 0 {
		return
	}
	var v *archive.Video
	if v, err = s.arc.NewVideo(c, m.Filename); err != nil {
		return
	}
	if v == nil {
		log.Warn("filename(%s) videoshotpv video not exist", m.Filename)
		return
	}
	s.arc.AddVideoShot(c, v.Cid, count)
	// double write...
	//s.videoshotAdd(v.Cid, m.BinURL, m.ImgURLs)
	return
}

func (s *Service) uploadInfo(c context.Context, m *message.BvcVideo) (err error) {
	s.redis.AddFilename(c, m.Filename)
	log.Info("filename(%s) upload success", m.Filename)
	return
}

func (s *Service) xcodeSDFail(c context.Context, m *message.BvcVideo) (err error) {
	var (
		v *archive.Video
		a *archive.Archive
	)
	if v, a, err = s.archiveVideo(c, m.Filename); err != nil {
		log.Error("s.archiveVideo(%s) error(%v)", m.Filename, err)
		return
	}
	if v.XcodeState >= archive.VideoXcodeSDFail {
		// NOTE: xcodeFail=1, xcodeState must uploadInfo=0
		log.Warn("archive(%d) video(%s) already(%d)", a.Aid, m.Filename, v.XcodeState)
		return // NOTE: is or not return???
	}
	v.Status = archive.VideoStatusXcodeFail
	v.XcodeState = archive.VideoXcodeSDFail
	v.FailCode = archive.XcodeFailCodes[m.FailInfo]
	// begin transcation
	var (
		tx     *xsql.Tx
		change bool
	)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran archive(%d) filename(%s) error(%v)", a.Aid, m.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) begin sd_fail transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = s.tranVideo(c, tx, a, v); err != nil {
		tx.Rollback()
		log.Error("s.tranVideo(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) sd_fail tranVideo fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if change, err = s.tranArchive(c, tx, a, v, nil); err != nil {
		tx.Rollback()
		log.Error("s.tranArchive(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) sd_fail tranArchive fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) end sd_fail transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if change {
		s.sendMsg(c, a, v)
	}
	return
}

func (s *Service) xcodeSDFinish(c context.Context, m *message.BvcVideo) (err error) {
	var (
		v *archive.Video
		a *archive.Archive
	)
	if v, a, err = s.archiveVideo(c, m.Filename); err != nil {
		log.Error("s.archiveVideo(%s) error(%v)", m.Filename, err)
		return
	}
	// if v.Playurl == m.PlayURL && v.Duration == m.Duration { // NOTE: check playurl&duration or xcode_state???
	if v.XcodeState >= archive.VideoXcodeSDFinish {
		// NOTE: sdFinish=2, xcideState must uploadInfo=0||xcodesdfail=1
		log.Warn("archive(%d) video(%s) already(%d)", a.Aid, m.Filename, v.XcodeState)
		return // NOTE: is or not return???
	}
	// if video already deleted, no dispatch no update
	if v.Status == archive.VideoStatusDelete {
		log.Info("xcodeSDFinish archive(%d) video(%s) video already deleted", a.Aid, m.Filename)
		v.Status = archive.VideoStatusDelete
	} else {
		v.Status = archive.VideoStatusWait // NOTE: default -1
	}
	// if archive already deleted, video state should be?
	if a.State == archive.StateForbidUpDelete {
		log.Info("xcodeSDFinish archive(%d) video(%s) archive already deleted", a.Aid, m.Filename)
		v.Status = archive.VideoStatusDelete
	}
	v.XcodeState = archive.VideoXcodeSDFinish
	v.Playurl = m.PlayURL
	v.Duration = m.Duration
	// begin transcation
	var (
		tx *xsql.Tx
	)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran archive(%d) filename(%s) error(%v)", a.Aid, m.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) begin sd_finish transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = s.tranVideo(c, tx, a, v); err != nil {
		tx.Rollback()
		log.Error("s.tranVideo(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	if v.Status != archive.VideoStatusDelete {
		log.Info("archive(%d) filename(%s) sd_finish tranVideo fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
		if _, err = s.tranArchive(c, tx, a, v, nil); err != nil {
			tx.Rollback()
			log.Error("s.tranArchive(%d, %s) error(%v)", a.Aid, v.Filename, err)
			return
		}
	}
	if err = s.tranArcCover(c, tx, a, v); err != nil {
		tx.Rollback()
		log.Error("s.tranArcCover(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) sd_finish tranArchive fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if _, err = s.arc.TranVideoOper(c, tx, a.Aid, v.ID, v.Status, v.Attribute); err != nil {
		tx.Rollback()
		log.Error("s.arc.TranVideoOper(%d, %d, %d) error(%v)", a.Aid, v.ID, v.Status, err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	if s.canDo(a.Mid) {
		s.syncRetry(context.TODO(), a.Aid, 0, redis.ActionForVideocovers, a.Cover, a.Cover)
	}

	log.Info("archive(%d) filename(%s) end sd_finish transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	return
}

func (s *Service) xcodeHDFail(c context.Context, m *message.BvcVideo) (err error) {
	log.Info("xcode hd fail filename (%s)", m.Filename)
	var (
		v *archive.Video
		a *archive.Archive
	)
	if v, a, err = s.archiveVideo(c, m.Filename); err != nil {
		log.Error("s.archiveVideo(%s) error(%v)", m.Filename, err)
		return
	}
	if v.XcodeState >= archive.VideoXcodeHDFail {
		// NOTE: hdfail=3, xcodeState must uploadInfo=0||xcodesdfail=1||sdfinish=2
		log.Warn("archive(%d) video(%s) already(%d)", a.Aid, m.Filename, v.XcodeState)
		return // NOTE: is or not return???
	}
	v.Status = archive.VideoStatusXcodeFail
	v.XcodeState = archive.VideoXcodeHDFail
	v.FailCode = archive.XcodeFailCodes[m.FailInfo]
	// begin transcation
	var (
		tx     *xsql.Tx
		change bool
	)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran archive(%d) filename(%s) error(%v)", a.Aid, m.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) begin hd_fail transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = s.tranVideo(c, tx, a, v); err != nil {
		tx.Rollback()
		log.Error("s.tranVideo(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) hd_fail tranVideo fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if change, err = s.tranArchive(c, tx, a, v, nil); err != nil {
		tx.Rollback()
		log.Error("s.tranArchive(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) hd_fail tranArchive fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) end hd_fail transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if change {
		s.sendMsg(c, a, v)
	}
	return
}

func (s *Service) xcodeHDFinish(c context.Context, m *message.BvcVideo) (err error) {
	var (
		v *archive.Video
		a *archive.Archive
	)
	if v, a, err = s.archiveVideo(c, m.Filename); err != nil {
		log.Error("s.archiveVideo(%s) error(%v)", m.Filename, err)
		return
	}
	if v.XcodeState >= archive.VideoXcodeHDFinish {
		// NOTE: hdFinish=3, xcodeState must uploadInfo=0||xcodesdfail=1||sdfinish=2||hdfail=3
		log.Warn("archive(%d) video(%s) already(%d)", a.Aid, m.Filename, v.XcodeState)
		return // NOTE: is or not return???
	}
	// make sure filename not exist in redis, otherwise videoup can not submit!!!
	s.redis.DelFilename(c, m.Filename)
	log.Info("filename(%s) del_filename from redis success", m.Filename)
	// start deal hd finish
	v.XcodeState = archive.VideoXcodeHDFinish
	v.Resolutions = m.Resolutions
	v.Filesize = m.Filesize
	v.Duration = m.Duration
	v.Dimensions = fmt.Sprintf("%d,%d,%d", m.Width, m.Height, m.Rotate)
	// begin transcation
	var (
		tx *xsql.Tx
	)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran archive(%d) filename(%s) error(%v)", a.Aid, m.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) begin hd_finish transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = s.tranVideo(c, tx, a, v); err != nil {
		tx.Rollback()
		log.Error("s.tranVideo(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) hd_finish tranVideo fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	// only hd5???
	var attr int32
	if strings.Contains(m.Resolutions, "hdflv2") || strings.Contains(m.Resolutions, "112") {
		attr = archive.AttrYes
	} else {
		attr = archive.AttrNo
	}
	if _, err = s.arc.TxUpAttrBit(tx, a.Aid, attr, archive.AttrBitHasHD5); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpAttrBit(%d, %d, hd5) error(%v)", a.Aid, attr, err)
		return
	}
	log.Info("archive(%d) filename(%s) hd_finish attrBitHD5 fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) end hd_finish transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	return
}

func (s *Service) dispatchRunning(c context.Context, m *message.BvcVideo) (err error) {
	var (
		v *archive.Video
		a *archive.Archive
	)
	if v, a, err = s.archiveVideo(c, m.Filename); err != nil {
		log.Error("s.archiveVideo(%s) error(%v)", m.Filename, err)
		return
	}
	if v.XcodeState >= archive.VideoDispatchRunning {
		// NOTE: dispathRun=4, xcodeState must uploadInfo=0||xcodefail=1||sdfinish=2||xcodehdfail=3||hdFinish=4
		log.Warn("archive(%d) video(%s) already(%d)", a.Aid, m.Filename, v.XcodeState)
		return // NOTE: is or not return???
	}
	v.XcodeState = archive.VideoDispatchRunning
	// begin transcation
	var (
		tx *xsql.Tx
	)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran archive(%d) filename(%s) error(%v)", a.Aid, m.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) begin dispatch_run transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = s.tranVideo(c, tx, a, v); err != nil {
		tx.Rollback()
		log.Error("s.tranVideo(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) dispatch_run tranVideo fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if _, err = s.tranArchive(c, tx, a, v, nil); err != nil {
		tx.Rollback()
		log.Error("s.tranArchive(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) dispatch_run tranArchive fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) end dispatch_run transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	return
}

func (s *Service) dispatchFinish(c context.Context, m *message.BvcVideo) (err error) {
	var (
		v *archive.Video
		a *archive.Archive
	)
	if v, a, err = s.archiveVideo(c, m.Filename); err != nil {
		log.Error("s.archiveVideo(%s) error(%v)", m.Filename, err)
		return
	}
	if v.XcodeState >= archive.VideoDispatchFinish {
		// NOTE: dispathFinish=5, xcodeState must uploadInfo=0||xcodefail=1||sdfinish=2||hdFinish=3||dispathRun=4
		log.Warn("archive(%d) video(%s) already(%d)", a.Aid, m.Filename, v.XcodeState)
		return // NOTE: is or not return???
	}
	v.XcodeState = archive.VideoDispatchFinish
	// begin transcation
	var (
		tx               *xsql.Tx
		sChange, rChange bool
	)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran archive(%d) filename(%s) error(%v)", a.Aid, m.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) begin dispatch_finish transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if err = s.tranVideo(c, tx, a, v); err != nil {
		tx.Rollback()
		log.Error("s.tranVideo(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) dispatch_finish tranVideo fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if sChange, err = s.tranArchive(c, tx, a, v, nil); err != nil {
		tx.Rollback()
		log.Error("s.tranArchive(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) dispatch_finish tranArchive fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	var round int8
	if round, err = s.tranRound(c, tx, a); err != nil {
		tx.Rollback()
		return
	}
	rChange = round != a.Round
	log.Info("archive(%d) filename(%s) dispatch_finish tranRound fininsh old_round(%d) new_round(%d)", a.Aid, v.Filename, a.Round, round)
	a.Round = round
	if sChange || rChange {
		if err = s.tranArchiveOper(tx, a); err != nil {
			tx.Rollback()
			return
		}
	}
	log.Info("archive(%d) filename(%s) dispatch_finish round_opr fininsh round(%d)", a.Aid, v.Filename, a.Round)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) end dispatch_finish transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if sChange {
		//稿件二审in/out量监控  每个aid 只统计一次 （1,自动过审2,手动审核需分开处理）（多P也只统计一次） start,end
		var had bool
		if archive.NormalState(a.State) {
			//monitor  second_round 自动开放 in/out diff
			s.promVideoS.Incr("second_round")
			s.promVideoE.Incr("second_round")

			//auto open
			s.syncBVC(c, a)
			s.sendAuditMsg(c, message.RouteAutoOpen, a.Aid)
			if is, _ := s.IsUpperFirstPass(c, a.Mid, a.Aid); is {
				go s.sendNewUpperMsg(c, a.Mid, a.Aid)
			}
		} else if had, _ = s.redis.SetMonitorCache(c, a.Aid); had {
			s.promVideoS.Incr("second_round")
		}
		s.sendMsg(c, a, v)
		if a.State == archive.StateForbidFixed {
			s.addClickToRedis(c, a.Aid)
		}
	}
	return
}

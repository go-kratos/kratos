package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/message"
	"go-common/app/admin/main/videoup/model/up"
	"go-common/library/log"
)

func (s *Service) busFirstRound(aid int64, fn, encodePurpose string, videoDesign *message.VideoDesign, status int16, encoding int8, regionID, typeID int16, fans int64, adminChange bool) (err error) {
	var msg = &message.Videoup{
		Route:          message.RouteFirstRound,
		Fans:           fans,
		Timestamp:      time.Now().Unix(),
		Aid:            aid,
		Filename:       fn,
		Status:         status,
		Xcode:          encoding,
		EncodeRegionID: regionID,
		VideoDesign:    videoDesign,
		AdminChange:    adminChange,
		EncodeTypeID:   typeID,
	}
	if len(encodePurpose) != 0 {
		msg.EncodePurpose = encodePurpose
	}
	log.Info("filename(%s) start to send firstRound(%+v) to databus", fn, msg)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busUGCFirstRound(aid int64, fn, encodePurpose string, videoDesign *message.VideoDesign, status int16, encoding int8, regionID, typeID int16, fans int64, adminChange bool) (err error) {
	var msg = &message.Videoup{
		Route:          message.RouteUGCFirstRound,
		Fans:           fans,
		Timestamp:      time.Now().Unix(),
		Aid:            aid,
		Filename:       fn,
		Status:         status,
		Xcode:          encoding,
		EncodeRegionID: regionID,
		VideoDesign:    videoDesign,
		AdminChange:    adminChange,
		EncodeTypeID:   typeID,
	}
	if len(encodePurpose) != 0 {
		msg.EncodePurpose = encodePurpose
	}
	log.Info("filename(%s) start to send UGCFirstRound(%+v) to databus", fn, msg)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busArchiveForceSync(aid int64) (err error) {
	var msg = &message.Videoup{
		Route:     message.RouteForceSync,
		Timestamp: time.Now().Unix(),
		Aid:       aid,
	}
	log.Info("aid(%d) send busArchiveForceSync to databus ", aid)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busModifyArchive(aid int64, aChange, vChange bool) (err error) {
	var msg = &message.Videoup{
		Route:       message.RouteModifyArchive,
		Timestamp:   time.Now().Unix(),
		Aid:         aid,
		EditArchive: aChange,
		EditVideo:   vChange,
	}
	log.Info("aid(%d) send modifyArchive to databus by EditArchive(%v) EditVideo(%v)", aid, aChange, vChange)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busSecondRound(aid, missionID int64, notify bool, email, changeTypeID, changeCopyright, changeTitle, ChangeCover bool, fromList string, ap *archive.ArcParam) (err error) {
	sendEmail := true
	if ap != nil {
		sendEmail = !ap.NoEmail
	}
	var msg = &message.Videoup{
		Route:           message.RouteSecondRound,
		Aid:             aid,
		Notify:          notify,
		MissionID:       missionID,
		Timestamp:       time.Now().Unix(),
		AdminChange:     email,
		ChangeTypeID:    changeTypeID,
		ChangeCopyright: changeCopyright,
		ChangeTitle:     changeTitle,
		ChangeCover:     ChangeCover,
		FromList:        fromList,
		SendEmail:       sendEmail,
	}
	log.Info("aid(%d) start to send secondRound msg(%+v) to databus", aid, msg)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busSecondRoundUpCredit(aid, cid, mid, uid int64, state, round int8, reasonID int64, reason string) (err error) {
	if mid == 0 || aid == 0 {
		return
	}
	var msg = &up.CreditLog{
		Type:         round,
		Optyte:       state,
		Reason:       reasonID,
		BusinessType: up.CreditBusinessTypeArchive,
		MID:          mid,
		OID:          aid,
		UID:          uid,
		Content:      reason,
		Ctime:        time.Now().Unix(),
		Extra:        map[string]interface{}{"cid": cid},
	}
	var c = context.TODO()
	log.Info("aid(%d) start to send busSecondRoundUpCredit msg(%+v) to databus", aid, msg)
	if err = s.upCreditPub.Send(c, string(msg.OID)+string(msg.UID), msg); err != nil {
		log.Error("aid(%d)  s.upCreditPub.Send(%+v) error(%v)", msg.OID, msg, err)
	}
	return
}

//func (s *Service) busDeleteVideo(aid int64, filename string) (err error) {
//	var msg = &message.Videoup{
//		Route:     message.RouteDeleteVideo,
//		Timestamp: time.Now().Unix(),
//		Aid:       aid,
//		Filename:  filename,
//	}
//	log.Info("aid(%d) filename(%s) start to send deleteVideo to databus", aid, filename)
//	if err = s.busSendMsg(msg); err != nil {
//		s.msgCh <- msg
//	}
//	return
//}

func (s *Service) busSendMsg(msg *message.Videoup) (err error) {
	var c = context.TODO()
	switch msg.Route {
	case message.RouteFirstRound, message.RouteUGCFirstRound, message.RouteDeleteVideo:
		if err = s.videoupPub.Send(c, msg.Filename, msg); err != nil {
			log.Error("filename(%s) %s s.videoupPub.Send(%+v) error(%v)", msg.Filename, msg.Route, msg, err)
		}
	case message.RouteSecondRound, message.RouteModifyArchive, message.RouteForceSync:
		if err = s.videoupPub.Send(c, strconv.FormatInt(msg.Aid, 10), msg); err != nil {
			log.Error("aid(%d) %s s.videoupPub.Send(%+v) error(%v)", msg.Aid, msg.Route, msg, err)
		}
	default:
		log.Warn("databuserr can't process the type (%s)", msg.Route)
	}
	return
}

// databus err proc
func (s *Service) msgproc() {
	// NOTE: chan
	s.wg.Add(1)
	go func() {
		var (
			c   = context.TODO()
			msg *message.Videoup
			ok  bool
			err error
		)
		defer s.wg.Done()
		for {
			if msg, ok = <-s.msgCh; !ok {
				log.Info("msgproc s.msgCh stop")
				return
			}
			log.Info("aid(%d) filename(%s) get msg(%+v) from s.msgCh", msg.Aid, msg.Filename, msg)
			if err = s.busSendMsg(msg); err != nil {
				s.busCache.PushMsgCache(c, msg)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	// NOTE: from redis list when chan error
	s.wg.Add(1)
	go func() {
		var (
			c   = context.TODO()
			msg *message.Videoup
			err error
		)
		defer s.wg.Done()
		for {
			if s.closed {
				log.Info("second msgproc service is close")
				return
			}
			if msg, err = s.busCache.PopMsgCache(c); err != nil {
				log.Error("msgproc s.busCache.PopMsgCache() error(%v)", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if msg == nil {
				select {
				case <-time.After(3 * time.Minute):
					continue
				case <-s.stop:
					return
				}
			}
			log.Info("aid(%d) filename(%s) get msg(%+v) from redis", msg.Aid, msg.Filename, msg)
			if err = s.busSendMsg(msg); err != nil {
				s.busCache.PushMsgCache(c, msg)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

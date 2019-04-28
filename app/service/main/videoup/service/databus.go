package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/app/service/main/videoup/model/dede"
	"go-common/app/service/main/videoup/model/message"
	"go-common/app/service/main/videoup/model/prom"
	"go-common/library/log"
)

func (s *Service) busSyncCid(pad *dede.PadInfo) (err error) {
	for fn, cid := range pad.Fnm {
		var msg = &message.Videoup{
			Route:     message.RouteSyncCid,
			Timestamp: time.Now().Unix(),
			Filename:  fn,
			Cid:       cid,
			DMIndex:   "vupload_" + strconv.FormatInt(cid, 10),
		}
		log.Info("aid(%d) filename(%s) start to send cid(%d) to databus", pad.Aid, fn, cid)
		if err = s.busSendMsg(msg); err != nil {
			s.msgCh <- msg
		}
	}
	return
}

func (s *Service) busAddArchive(aid int64, upFrom int8) (err error) {
	var msg = &message.Videoup{
		Route:     message.RouteAddArchive,
		UpFrom:    upFrom,
		Timestamp: time.Now().Unix(),
		Aid:       aid,
	}
	log.Info("aid(%d) filename() start to send addArchive to databus", aid)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busSyncArchive(aid int64) (err error) {
	var msg = &message.Videoup{
		Route:     message.RouteForceSync,
		Timestamp: time.Now().Unix(),
		Aid:       aid,
	}
	log.Info("aid(%d) filename() start to send force_sync to databus", aid)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busModifyArchive(aid int64, aChange, vChange, tagChange, addvideos, typeIDChange bool) (err error) {
	var msg = &message.Videoup{
		Route:        message.RouteModifyArchive,
		Timestamp:    time.Now().Unix(),
		Aid:          aid,
		EditArchive:  aChange,
		EditVideo:    vChange,
		TagChange:    tagChange,
		AddVideos:    addvideos,
		ChangeTypeID: typeIDChange,
	}
	log.Info("aid(%d) filename() start to send modifyArchive to databus by EditArchive(%v) EditVideo(%v) tagchange(%v) addvideos(%v) typeidchange(%v)", aid, aChange, vChange, tagChange, addvideos, typeIDChange)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busUGCSubmit(pad *dede.PadInfo) (err error) {
	for filename, cid := range pad.Fnm {
		var msg = &message.Videoup{
			Route:     message.RouteUGCSubmit,
			Timestamp: time.Now().Unix(),
			Filename:  filename,
			Cid:       cid,
			DMIndex:   "vupload_" + strconv.FormatInt(cid, 10),
			Submit:    1,
		}
		log.Info("aid(%d) filename(%s) start to send ugc submit to databus", pad.Aid, filename)
		if err = s.busSendMsg(msg); err != nil {
			s.msgCh <- msg
		}
	}
	return
}

func (s *Service) busPGCSubmit(nvs []*archive.Video) (err error) {
	for _, v := range nvs {
		// add pub agent
		s.agent.AgentMsg(context.TODO(), message.RoutePGCSubmit, v.Filename)
		// databus
		var msg = &message.Videoup{
			Route:     message.RoutePGCSubmit,
			Timestamp: time.Now().Unix(),
			Filename:  v.Filename,
			Submit:    1,
		}
		log.Info("aid(%d) filename(%s) start to send pgc submit to databus", v.Aid, v.Filename)
		if err = s.busSendMsg(msg); err != nil {
			s.msgCh <- msg
		}
	}
	return
}

func (s *Service) busDRMSubmit(nvs []*archive.Video) (err error) {
	for _, v := range nvs {
		var msg = &message.Videoup{
			Route:     message.RouteDRMSubmit,
			Timestamp: time.Now().Unix(),
			Filename:  v.Filename,
			Submit:    1,
		}
		log.Info("aid(%d) filename(%s) start to send drm submit to databus", v.Aid, v.Filename)
		if err = s.busSendMsg(msg); err != nil {
			s.msgCh <- msg
		}
	}
	return
}

func (s *Service) busDeleteVideo(dvs []*archive.Video) (err error) {
	for _, v := range dvs {
		var msg = &message.Videoup{
			Route:     message.RouteDeleteVideo,
			Timestamp: time.Now().Unix(),
			Aid:       v.Aid,
			Filename:  v.Filename,
		}
		log.Info("aid(%d) filename(%s) start to send deleteVideo to databus", v.Aid, v.Filename)
		if err = s.busSendMsg(msg); err != nil {
			s.msgCh <- msg
		}
	}
	return
}

func (s *Service) busModifyVideo(evs []*archive.Video) (err error) {
	for _, v := range evs {
		var msg = &message.Videoup{
			Route:     message.RouteModifyVideo,
			Timestamp: time.Now().Unix(),
			Aid:       v.Aid,
			Filename:  v.Filename,
		}
		log.Info("aid(%d) filename(%s) start to send modifyVideo to databus", v.Aid, v.Filename)
		if err = s.busSendMsg(msg); err != nil {
			s.msgCh <- msg
		}
	}
	return
}

func (s *Service) busDeleteArchive(aid int64) (err error) {
	var msg = &message.Videoup{
		Route:     message.RouteDeleteArchive,
		Aid:       aid,
		Timestamp: time.Now().Unix(),
	}
	log.Info("aid(%d) filename() start to send deleteArchive to databus", aid)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busSecondRound(aid, missionID int64) (err error) {
	var msg = &message.Videoup{
		Route:     message.RouteSecondRound,
		Aid:       aid,
		MissionID: missionID,
		Timestamp: time.Now().Unix(),
	}
	log.Info("aid(%d) start to send secondRound msg(%+v) to databus", aid, msg)
	if err = s.busSendMsg(msg); err != nil {
		s.msgCh <- msg
	}
	return
}

func (s *Service) busSendMsg(msg *message.Videoup) (err error) {
	var c = context.TODO()
	switch msg.Route {
	case message.RouteSyncCid, message.RouteDeleteVideo, message.RouteModifyVideo, message.RouteUGCSubmit:
		if err = s.videoupPub.Send(c, msg.Filename, msg); err != nil {
			log.Error("aid(%d) filename(%s) %s s.videoupPub.Send(%+v) error(%v)", msg.Aid, msg.Filename, msg.Route, msg, err)
		}
	case message.RouteModifyArchive, message.RouteAddArchive, message.RouteDeleteArchive, message.RouteSecondRound, message.RouteForceSync:
		if err = s.videoupPub.Send(c, strconv.FormatInt(msg.Aid, 10), msg); err != nil {
			log.Error("aid(%d) filename(%s) %s s.videoupPub.Send(%+v) error(%v)", msg.Aid, msg.Filename, msg.Route, msg, err)
		}
	case message.RoutePGCSubmit, message.RouteDRMSubmit:
		if err = s.videoupPGCPub.Send(c, msg.Filename, msg); err != nil {
			log.Error("aid(%d) filename(%s) %s s.videoupPGCPub.Send(%+v) error(%v)", msg.Aid, msg.Filename, msg.Route, msg, err)
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
			s.promErr.Incr(prom.RouteDatabusTry)
			log.Info("aid(%d) filename(%s) get msg(%+v) from redis", msg.Aid, msg.Filename, msg)
			if err = s.busSendMsg(msg); err != nil {
				s.busCache.PushMsgCache(c, msg)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

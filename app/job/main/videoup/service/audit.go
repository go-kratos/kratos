package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/message"
	"go-common/app/job/main/videoup/model/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// videoupConsumer is videoup message consumer.
func (s *Service) videoupConsumer() {
	defer s.wg.Done()
	var (
		msgs = s.videoupSub.Messages()
		err  error
		c    = context.TODO()
	)
	for {
		func() {
			var (
				msg    *databus.Message
				ok     bool
				offset int64
				dbus   *archive.Databus
			)
			msg, ok = <-msgs
			if !ok {
				for part, lastOffset := range s.videoupSubIdempotent {
					if _, err = s.arc.UpDBus(c, s.c.VideoupSub.Group, s.c.VideoupSub.Topic, part, lastOffset); err != nil {
						continue
					}
				}
				log.Error("s.videoupSub.Messages closed")
				return
			}
			defer s.Rescue(string(msg.Value))
			msg.Commit()
			s.videoupMo++
			m := &message.Videoup{}
			if err = json.Unmarshal(msg.Value, m); err != nil {
				log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
				return
			}
			offset, ok = s.videoupSubIdempotent[msg.Partition]
			if !ok {
				if dbus, err = s.arc.DBus(c, s.c.VideoupSub.Group, msg.Topic, msg.Partition); err != nil {
					return
				}
				if dbus == nil {
					if _, err = s.arc.AddDBus(c, s.c.VideoupSub.Group, msg.Topic, msg.Partition, msg.Offset); err != nil {
						return
					}
					offset = msg.Offset
				} else {
					offset = dbus.Offset
				}

			}
			// if last offset > current offset -> continue
			if offset > msg.Offset {
				log.Error("key(%s) value(%s) partition(%s) offset(%d) is too early", msg.Key, msg.Value, msg.Partition, msg.Offset)
				return
			}
			s.videoupSubIdempotent[msg.Partition] = msg.Offset
			s.promDatabus.Incr(m.Route)
			log.Info("videoupMessage key(%s) value(%s) partition(%d) offset(%d) route(%s) commit start", msg.Key, msg.Value, msg.Partition, msg.Offset, m.Route)
			switch m.Route {
			case message.RouteSyncCid:
				err = s.syncCid(c, m)
			case message.RouteFirstRound:
				err = s.firstRound(c, m)
			case message.RouteUGCFirstRound:
				err = s.firstRound(c, m)
			case message.RouteSecondRound:
				err = s.secondRound(c, m)
			case message.RouteAddArchive:
				err = s.addArchive(c, m)
			case message.RouteModifyArchive:
				err = s.modifyArchive(c, m)
			case message.RouteDeleteArchive:
				err = s.deleteArchive(c, m)
			case message.RouteDeleteVideo:
				err = s.deleteVideo(c, m)
			case message.RouteModifyVideo:
				log.Info("databus modifyVideo(%+v)", m)
			default:
				log.Warn("videoupConsumer unknown message route(%s)", m.Route)
			}
			if err == nil {
				log.Info("videoupMessage key(%s) value(%s) partition(%d) offset(%d) end", msg.Key, msg.Value, msg.Partition, msg.Offset)
			} else {
				log.Error("videoupMessage key(%s) value(%s) partition(%d) offset(%d) error(%v)", msg.Key, msg.Value, msg.Partition, msg.Offset, err)
			}
		}()
	}
}

func (s *Service) syncCid(c context.Context, m *message.Videoup) (err error) {
	// make sure filename not exist in redis, otherwise videoup can not submit!!!
	s.redis.DelFilename(c, m.Filename)
	log.Info("filename(%s) del_filename from redis success", m.Filename)
	return
}

func (s *Service) firstRound(c context.Context, m *message.Videoup) (err error) {
	var (
		v *archive.Video
		a *archive.Archive
	)
	if m.Xcode == 1 {
		s.promVideoS.Incr("xcode_hd")
	}
	s.promVideoE.Incr("first_round")
	if v, a, err = s.archiveVideoByAid(c, m.Filename, m.Aid); err != nil {
		log.Error("s.archiveVideoByAid(%s, %d) error(%v)", m.Filename, m.Aid, err)
		return
	}
	if a.State == archive.StateForbidUpDelete {
		log.Warn("archive(%d) filename(%s) state(%d) is deleted", a.Aid, v.Filename, a.State)
		return
	}
	// make sure filename not exist in redis, otherwise videoup can not submit!!!
	s.redis.DelFilename(c, m.Filename)
	log.Info("filename(%s) del_filename from redis success", m.Filename)
	var (
		tx               *xsql.Tx
		sChange, rChange bool
		round            int8
		ad               = archive.AuditParam{IsAudit: true}
	)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran archive(%d) filename(%s) error(%v)", a.Aid, m.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) begin first_round transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if sChange, err = s.tranArchive(c, tx, a, v, &ad); err != nil {
		tx.Rollback()
		log.Error("s.tranArchive(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) first_round tranArchive fininsh a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)
	if round, err = s.tranRound(c, tx, a); err != nil {
		tx.Rollback()
		log.Error("s.tranRound error(%v)", err)
		return
	}
	rChange = round != a.Round
	log.Info("archive(%d) firstRound upRound(%d)", a.Aid, a.Round)
	a.Round = round
	if sChange || rChange {
		if err = s.tranArchiveOper(tx, a); err != nil {
			tx.Rollback()
			return
		}
	}
	log.Info("archive(%d) filename(%s) first_round tranRound fininsh a_round(%d)", a.Aid, v.Filename, a.Round)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(%d, %s) error(%v)", a.Aid, v.Filename, err)
		return
	}
	log.Info("archive(%d) filename(%s) end first_round transcation a_state(%d) v_status(%d)", a.Aid, v.Filename, a.State, v.Status)

	s.sendPostFirstRound(c, message.RoutePostFirstRound, a.Aid, v.Filename, m.AdminChange)
	if sChange {
		if a.IsForbid() {
			//todo fixbug 导致打回的再自动开放视频 16s
			s.syncBVC(c, a)
			s.sendAuditMsg(c, message.RouteFirstRoundForbid, a.Aid)
		}
		s.sendMsg(c, a, v)
		if archive.NormalState(a.State) {
			// add for open state after dispatch
			s.sendAuditMsg(c, message.RouteAutoOpen, a.Aid)
			if is, _ := s.IsUpperFirstPass(c, a.Mid, a.Aid); is {
				go s.sendNewUpperMsg(c, a.Mid, a.Aid)
			}
			log.Info("firstRound aid(%d) pubbed. Do sync.", a.Aid)
		}
	}
	if s.canDo(a.Mid) {
		s.syncRetry(context.TODO(), a.Aid, 0, redis.ActionForVideocovers, a.Cover, a.Cover)
	}
	return
}

func (s *Service) secondRound(c context.Context, m *message.Videoup) (err error) {
	var (
		a     *archive.Archive
		addit *archive.Addit
	)
	if a, err = s.arc.Archive(c, m.Aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", m.Aid, err)
		return
	}
	if addit, err = s.arc.Addit(c, m.Aid); err != nil {
		//为了不影响主流程，这里只是记个错误日志
		//!!!注意!!! addit可能是nil，用之前需要判断
		log.Error("s.arc.Addit(%d) error(%v)", m.Aid, err)
		err = nil
	}
	if a == nil {
		log.Warn("archive(%d) is not exist", m.Aid)
		return
	}
	//check first admin audit
	var had bool
	if had, _ = s.redis.SetMonitorCache(c, a.Aid); had {
		s.redis.DelMonitorCache(c, a.Aid)
		s.promVideoE.Incr("second_round")
	}
	// start archive
	var tx *xsql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran error(%v)", err)
		return
	}
	if err = s.tranSumDuration(c, tx, a); err != nil {
		tx.Rollback()
		log.Error("s.tranSumDuration error(%v)", err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	if a.State != archive.StateForbidFixed {
		s.syncBVC(c, a)
	}

	if archive.NormalState(a.State) {
		s.changeMission(c, a, m.MissionID)
	} else {
		s.unBindMission(c, a, m.MissionID)
	}
	if m.IsSendNotify {
		s.sendMsg(c, a, nil)
	}
	//TODO 将其它发送系统通知的逻辑都放到sendMsg2Upper()里
	go s.sendMsg2Upper(c, m, a, addit)
	//回查阶段修改了稿件信息需要发送通知
	if a.State >= archive.StateOpen && ((s.c.ChangeDebug && s.c.ChangeMid == a.Mid) || !s.c.ChangeDebug) && (m.ChangeTypeID || m.ChangeTitle || m.ChangeCopyright || m.ChangeCover) {
		log.Info("archive(%d) secondRound sendChangeMsg(%+v)", a.Aid, m)
		s.sendChangeMsg(c, a, nil, m)
	}

	if s.canDo(a.Mid) {
		s.syncRetry(context.TODO(), a.Aid, 0, redis.ActionForVideocovers, a.Cover, a.Cover)
	}

	return
}

// sendMsg2Upper 发送系统通知给UP主
func (s *Service) sendMsg2Upper(c context.Context, m *message.Videoup, a *archive.Archive, addit *archive.Addit) {
	if archive.NormalState(a.State) && m.MissionID != 0 && addit != nil && addit.MissionID == 0 { //m.MissionID是审核人员提交修改之前的活动id
		s.sendMissionMsg(c, a)
	}
	is, err := s.IsUpperFirstPass(c, a.Mid, a.Aid)
	if err != nil {
		log.Error("s.IsUpperFirstPass(%d,%d) error(%v)", a.Mid, a.Aid, err)
		err = nil
	} else if is {
		// UP主第一次稿件过审发送
		s.sendNewUpperMsg(c, a.Mid, a.Aid)
	}
}

func (s *Service) addArchive(c context.Context, m *message.Videoup) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, m.Aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", m.Aid, err)
		return
	}
	if a == nil {
		log.Warn("archive(%d) is not exist", m.Aid)
		return
	}
	if s.canDo(a.Mid) {
		s.syncRetry(context.TODO(), a.Aid, 0, redis.ActionForVideocovers, a.Cover, a.Cover)
	}

	var (
		tx             *xsql.Tx
		round          int8
		state, _, _, _ = s.archiveState(c, a, nil, nil)
	)
	if archive.NormalState(state) { // NOTE: add archive must audit maybe code mode.
		state = archive.StateForbidWait
	}
	if state == a.State {
		log.Warn("archive(%d) add_archive newState(%d)==oldState(%d)", a.Aid, state, a.State)
		return
	}
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran archive(%d) error(%v)", a.Aid, err)
		return
	}
	log.Info("archive(%d) begin add_archive transcation a_state(%d)", a.Aid, a.State)
	// archive
	if _, err = s.arc.TxUpState(tx, a.Aid, state); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpState(%d, %d) error(%v)", a.Aid, state, err)
		return
	}
	a.State = state
	log.Info("archive(%d) add_archive upState(%d)", a.Aid, a.State)
	if round, err = s.tranRound(c, tx, a); err != nil {
		tx.Rollback()
		log.Error("s.tranRound error(%v)", err)
		return
	}
	a.Round = round
	log.Info("archive(%d) add_archive upRound(%d)", a.Aid, a.Round)
	if err = s.tranArchiveOper(tx, a); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(%d) error(%v)", a.Aid, err)
		return
	}
	log.Info("archive(%d) end add_archive transcation a_state(%d)", a.Aid, a.State)
	return
}

func (s *Service) modifyArchive(c context.Context, m *message.Videoup) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, m.Aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", m.Aid, err)
		return
	}
	// start archive
	if a.NotAllowUp() {
		log.Warn("archive(%d) modify_archive state(%d) not allow update", a.Aid, a.State)
		return
	}
	if s.canDo(a.Mid) {
		s.syncRetry(context.TODO(), a.Aid, 0, redis.ActionForVideocovers, a.Cover, a.Cover)
	}

	// TODO comment when online
	if !m.EditArchive && !m.EditVideo {
		// to sync pubbed arc if not edited or just edit order.
		if archive.NormalState(a.State) {
			s.sendAuditMsg(c, message.RouteForceSync, a.Aid)
			log.Info("modifyArchive aid(%d) not modified or changed order. Do sync.", a.Aid)
		}
		return
	}
	var (
		tx             *xsql.Tx
		state, _, _, _ = s.archiveState(c, a, nil, nil)
		round          int8
	)
	if archive.NormalState(state) || state == archive.StateForbidUserDelay {
		if state != archive.StateForbidFixed {
			state = archive.StateForbidWait
		}
	}
	if state == a.State {
		return
	}
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran archive(%d) error(%v)", a.Aid, err)
		return
	}
	log.Info("archive(%d) begin modify_archive transcation a_state(%d)", a.Aid, a.State)
	// archive
	if _, err = s.arc.TxUpState(tx, a.Aid, state); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpState(%d, %d) error(%v)", a.Aid, state, err)
		return
	}
	log.Info("archive(%d) modify_archive a.State(%d) upState(%d)", a.Aid, a.State, state)
	a.State = state
	if round, err = s.tranRound(c, tx, a); err != nil {
		tx.Rollback()
		return
	}
	a.Round = round
	log.Info("archive(%d) modify_archive upRound(%d)", a.Aid, a.Round)
	if err = s.tranArchiveOper(tx, a); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(%d) error(%v)", a.Aid, err)
		return
	}
	log.Info("archive(%d) end modify_archive transcation a_state(%d)", a.Aid, a.State)
	if a.State == archive.StateForbidFixed {
		s.addClickToRedis(c, a.Aid)
	}

	return
}

func (s *Service) deleteArchive(c context.Context, m *message.Videoup) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, m.Aid); err != nil {
		log.Error("s.arc.Archive(%d) error(%v)", m.Aid, err)
		return
	}

	log.Info("archive delete_archive aid(%d)", m.Aid)
	s.syncBVC(c, a)
	return
}

func (s *Service) deleteVideo(c context.Context, m *message.Videoup) (err error) {
	var (
		v   *archive.Video
		tx  *xsql.Tx
		sum int64
	)
	if v, err = s.arc.NewVideo(c, m.Filename); err != nil || v == nil {
		return
	}
	log.Info("archive delete_video filename(%s) cid(%d) begin", m.Filename, v.Cid)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran error(%v)", err)
		return
	}
	if sum, err = s.arc.NewSumDuration(c, v.Aid); err != nil {
		tx.Rollback()
		log.Error("s.arc.SumDuration(%d) filename(%s) error(%v)", v.Aid, v.Filename, err)
		return
	}
	if _, err = s.arc.TxUpArcDuration(tx, v.Aid, sum); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpArcDuration(%d, %d) filename(%s) error(%v)", v.Aid, sum, v.Filename, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	log.Info("archive(%d) filename(%s) upArcDuration(%d)", v.Aid, v.Filename, sum)
	log.Info("archive delete_video filename(%s) cid(%d) end", m.Filename, v.Cid)
	return
}

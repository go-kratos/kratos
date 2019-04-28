package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/message"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) delayproc() {
	defer s.wg.Done()
	for {
		if s.closed {
			return
		}
		s.pubDelay()
		time.Sleep(1 * time.Minute)
	}
}

func (s *Service) pubDelay() {
	var (
		c      = context.TODO()
		delays []*archive.Delay
		err    error
		now    = time.Now()
	)
	if delays, err = s.arc.NowDelays(c, now); err != nil {
		log.Error("s.arc.NowDelays err(%v)", err)
		return
	}
	if len(delays) == 0 {
		log.Info("%s nothing need delay pub", now.Format("2006-01-02 15:04:05"))
		return
	}
	var delayIds = make([]int64, 0, len(delays))
	for _, v := range delays {
		log.Info("archive(%d) type(%d) begin delay pub", v.Aid, v.Type)
		var (
			a  *archive.Archive
			tx *sql.Tx
		)
		delayIds = append(delayIds, v.ID)
		if a, err = s.arc.Archive(c, v.Aid); err != nil || a == nil {
			log.Error("s.arc.Archive(%d) error(%v)", v.Aid, err)
			continue
		}
		//定时发布稿件 到临界点必须 state=-40
		if v.Type == archive.DelayTypeForUser && a.State != archive.StateForbidUserDelay {
			log.Info("archive(%d) type(%d) state(%d) delay failed", v.Aid, v.Type, a.State)
			now := now.Unix()
			s.msg.Send(c, "1_7_9", "您的稿件定时发布失败", fmt.Sprintf(`您的稿件《%s》（av%d）未能成功发布。原因：未能及时通过审核 您可以在稿件编辑页面重新编辑定时发布的时间，为此带来的麻烦，深感抱歉。#{点击进入编辑>>}{"http://member.bilibili.com/v/video/submit.html?type=edit&aid=%d"}`,
				a.Title, a.Aid, a.Aid), a.Mid, now)
			s.msg.Send(c, "113_1_1", "您的视频定时发布失败", fmt.Sprintf(`您的视频《%s》（av%d）定时发布失败。原因：未能及时通过审核。您可在稿件编辑页重新编辑定时发布时间。为此带来的麻烦，深感抱歉。#{点击进行编辑>>}{"https://member.bilibili.com/v/video/submit.html?type=edit&aid=%d"}`,
				a.Title, a.Aid, a.Aid), a.Mid, now)
			continue
		}
		if v.Type == archive.DelayTypeForUser && a.State == archive.StateForbidUserDelay {
			a.State = archive.StateOpen
			var addit *archive.Addit
			if addit, err = s.arc.Addit(c, a.Aid); err != nil && addit != nil && addit.MissionID > 0 {
				a.State = archive.StateOpen
			} else if !s.isWhite(a.Mid) && !s.isBlack(a.Mid) {
				if pfl, _ := s.profile(c, a.Mid); pfl != nil && pfl.Follower < int64(s.fansCache) && s.isRoundType(a.TypeID) {
					a.State = archive.StateOrange // NOTE: auto open must
				}
			}
		} else {
			log.Info("archive(%d) type(%d) state(%d) delay failed", v.Aid, v.Type, v.State)
			continue
		}
		if tx, err = s.arc.BeginTran(c); err != nil {
			log.Error("s.arc.BeginTran error(%v)", err)
			continue
		}
		firstPass := false
		if firstPass, err = s.txUpArcState(c, tx, a.Aid, a.State); err != nil {
			tx.Rollback()
			continue
		}
		log.Info("archive(%d) type(%d) delay upState(%d)", v.Aid, v.Type, a.State)
		if firstPass {
			a.PTime = xtime.Time(now.Unix())
			if _, err = s.arc.TxUpPTime(tx, a.Aid, now); err != nil {
				tx.Rollback()
				continue
			}
			log.Info("archive(%d) type(%d) delay upPTime(%d)", v.Aid, v.Type, a.PTime)
		}
		if archive.NormalState(a.State) {
			var round int8
			if round, err = s.tranRound(c, tx, a); err != nil {
				tx.Rollback()
				continue
			}
			a.Round = round
			if err = s.tranArchiveOper(tx, a); err != nil {
				tx.Rollback()
				continue
			}
		} else {
			if err = s.tranArchiveOper(tx, a); err != nil {
				tx.Rollback()
				return
			}
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.commit error(%v)", err)
			continue
		}
		if archive.NormalState(a.State) {
			s.syncBVC(c, a)
			s.changeMission(c, a, 0)
			now := now.Unix()
			s.msg.Send(c, "1_7_10", "您的稿件已经发布", fmt.Sprintf(`您的稿件《%s》（av%d）已经发布成功，#{点击查看>>}{"http://www.bilibili.com/video/av%d/"} `,
				a.Title, a.Aid, a.Aid), a.Mid, now)
			s.msg.Send(c, "113_1_1", "您的视频定时发布成功", fmt.Sprintf(`您的视频《%s》（av%d）已经发布成功，#{点击查看>>}{"http://www.bilibili.com/video/av%d/"} `,
				a.Title, a.Aid, a.Aid), a.Mid, now)
			s.sendAuditMsg(c, message.RouteDelayOpen, a.Aid)
			if is, _ := s.IsUpperFirstPass(c, a.Mid, a.Aid); is {
				go s.sendNewUpperMsg(c, a.Mid, a.Aid)
			}
		}
		log.Info("archive(%d) type(%d) end delay pub", v.Aid, v.Type)
	}
	s.arc.DelDelayByIds(c, delayIds)
}

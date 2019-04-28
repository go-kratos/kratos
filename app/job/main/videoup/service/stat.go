package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/message"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// statConsumer is stat message consumer.
func (s *Service) statConsumer() {
	defer s.wg.Done()
	const magicClick = 3000
	var (
		msgs = s.statSub.Messages()
		err  error
		c    = context.TODO()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.statSub.Messages closed")
			return
		}
		msg.Commit()
		s.statMo++
		m := &message.StatMsg{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("stat %s json.Unmarshal(%s) error(%v)", msg, msg.Value, err)
			continue
		}
		if m.Type != "archive" || m.ID < 8050956 || m.Count < s.thrMin || m.Count > s.thrMax+magicClick {
			continue
		}
		s.checkClick(c, m)
		log.Info("statConsumer key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
	}
}

func (s *Service) checkClick(c context.Context, m *message.StatMsg) (err error) {
	var a *archive.Archive
	if a, err = s.arc.Archive(c, m.ID); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v) or a==nil", m.ID, err)
		return
	}
	if a.Round != archive.RoundReviewFirstWaitTrigger {
		log.Warn("archive(%d) round(%d) not 31 wait trigger", a.Aid, a.Round)
		return
	}
	var (
		thr      = s.thrTpsCache[s.sfTpsCache[a.TypeID]]
		click, _ = s.redis.ArcClick(c, a.Aid)
	)
	if m.Count-click < thr || thr == 0 {
		log.Warn("archive(%d) typeThr(%d) nowClick(%d) beforeClick(%d) round(%d)", a.Aid, thr, m.Count, click, a.Round)
		return
	}
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran(%d) error(%v)", m.ID, err)
		return
	}
	var round = archive.RoundTriggerClick
	if _, err = s.arc.TxUpRound(tx, a.Aid, round); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxUpRound(%d, %d) error(%d)", a.Aid, round, err)
		return
	}
	a.Round = round
	if err = s.tranArchiveOper(tx, a); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	log.Info("archive(%d) typeThr(%d) nowClick(%d) upRound success", a.Aid, thr, m.Count)
	return
}

func (s *Service) addClickToRedis(c context.Context, aid int64) (err error) {
	click, err := s.arc.Stat(c, aid)
	if err != nil {
		log.Error("s.arc.Stat(%d) error(%v)", aid, err)
		return
	}
	if err = s.redis.AddArcClick(c, aid, click); err != nil {
		log.Error("s.redis.AddArcClick(%d) error(%v)", click, err)
	}
	return
}

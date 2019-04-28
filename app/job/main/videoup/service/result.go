package service

import (
	"encoding/json"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/message"
	"go-common/library/log"
)

const (
	_archive = "archive"
)

// arcResultConsumer consume archive result databus message
func (s *Service) arcResultConsumer() {
	defer s.wg.Done()
	var (
		msgs = s.arcResultSub.Messages()
		err  error
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.arcResultSub.Messages closed")
			return
		}
		msg.Commit()
		s.arcResultMo++
		m := &message.ArcResult{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		newArc := &archive.Result{}
		if err = json.Unmarshal(m.New, newArc); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", m.New, err)
			continue
		}
		log.Info("arcResultConsumer Topic(%s) partition(%d) offset(%d)  commit start", msg.Topic, msg.Partition, msg.Offset)
		if m.Table == _archive {
			log.Info("arcResultConsumer aid(%d) SendBblog msg(%v)", newArc.Aid, newArc)
			s.sendBblog(newArc)
		}
	}
}

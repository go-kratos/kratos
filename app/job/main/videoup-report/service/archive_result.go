package service

import (
	"encoding/json"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/log"
)

func (s *Service) arcResultConsume() {
	defer s.waiter.Done()
	var (
		err  error
		msgs = s.arcResultSub.Messages()
	)
	for {
		msg, open := <-msgs
		if !open {
			log.Info("arcResultConsume s.arcResultSub.Messages is closed")
			return
		}
		msg.Commit()
		if msg == nil {
			continue
		}
		log.Info("arcResultConsume consume key(%s) offset(%d) value(%s)", msg.Key, msg.Offset, string(msg.Value))

		m := &archive.Message{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("arcResultConsume json.Unmarshal error(%v)", err)
			continue
		}
		if m.Table != _archive {
			continue
		}

		nw := &archive.Archive{}
		if err = json.Unmarshal(m.New, nw); err != nil {
			log.Error("arcResultConsume json.Unmarshal error(%v) msg.new(%s)", err, string(m.New))
			continue
		}
		nw.ID = nw.AID

		var old *archive.Archive
		if m.Action == _insertAct {
			old = nil
		} else if m.Action == _updateAct {
			old = &archive.Archive{}
			if err = json.Unmarshal(m.Old, old); err != nil {
				log.Error("arcResultConsume json.Unmarshal error(%v) msg.old(%s)", err, string(m.Old))
				continue
			}
			old.ID = old.AID
		}

		go s.arcStateChange(nw, old, false)
	}
}

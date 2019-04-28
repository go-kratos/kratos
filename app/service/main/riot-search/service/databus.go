package service

import (
	"encoding/json"

	"go-common/app/service/main/riot-search/model"
	"go-common/library/log"
)

func (s *Service) watcherproc() {
	defer func() {
		log.Error("watcherproc quit")
	}()
	msgs := s.databus.Messages()
	states := model.PubStates
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.event.Messages closed")
			return
		}
		msg.Commit()
		log.Info("key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
		var m model.ArchiveMessage
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}

		if m.Table != "archive" {
			continue
		}
		if m.New == nil {
			log.Error("dirty data from databus value(%v)", string(msg.Value))
			continue
		}
		switch m.Action {
		case "insert":
			if states.Legal(m.New.State) {
				s.pool.JobQueue <- func() {
					s.dao.Insert(m.New.AID, m.New.Title, false)
					log.Info("riot: insert data into index id(%d) content(%s)", m.New.AID, m.New.Title)
				}
			}
		case "update":
			if m.Old == nil {
				log.Error("dirty data from databus value(%v)", msg.Value)
				continue
			} else if states.Legal(m.New.State) && m.New.Title != m.Old.Title {
				s.pool.JobQueue <- func() {
					s.dao.Insert(m.New.AID, m.New.Title, false)
					log.Info("riot: update data into index id(%d) content(%s)", m.New.AID, m.New.Title)
				}
			} else if !states.Legal(m.New.State) && states.Legal(m.Old.State) {
				s.pool.JobQueue <- func() {
					s.dao.Remove(m.New.AID, false)
					log.Info("riot: remove data id(%d) state(%d)", m.New.AID, m.New.State)
				}
			} else {
				log.Info("ignore action(%s) value(%s)", m.Action, msg.Value)
			}
		default:
		}
	}
}

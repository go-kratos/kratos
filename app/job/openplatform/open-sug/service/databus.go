package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/job/openplatform/open-sug/conf"
	"go-common/app/job/openplatform/open-sug/model"
	"go-common/library/log"
)

func (s *Service) pgcConsumePROC() {
	var (
		msgs = s.pgcSub.Messages()
		err  error
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("pgc databus Consumer exit")
			return
		}
		s.pgcMsgCnt++
		msg.Commit()
		log.Info("message commit key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
		m := &model.Message{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		func(m *model.Message) {
			s.wg.Add(1)
			//defer s.wg.Done()
			go s.subproc()

			if m.Table == "t_chn_season2" || m.Table == "t_jp_season2" {
				s.seasonUpdate(m.Action, m.New, m.Old)
				s.seasonMsgCnt++
			}
		}(m)
		log.Info("pgcConsumeproc table:%s key:%s partition:%d offset:%d", m.Table, msg.Key, msg.Partition, msg.Offset)
	}
}

func (s *Service) seasonUpdate(action string, n, o json.RawMessage) {
	newSeason := new(model.Season)
	if err := json.Unmarshal(n, newSeason); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", n, err)
		return
	}
	oldSeason := new(model.Season)
	if err := json.Unmarshal(o, oldSeason); err != nil {
		log.Error("json.Unmarshal(%s) (%s)error(%v)", string(o), o, err)
		return
	}
	if action == "insert" || action == "update" {
		if newSeason.FieldDiff(oldSeason) {
			s.dao.Index(
				context.TODO(),
				s.envIndex(s.c.ElasticSearch.Season.Index),
				s.c.ElasticSearch.Season.Type,
				strconv.Itoa(newSeason.ID),
				newSeason.EsFormat(),
			)
		}
	}
}

func (s *Service) envIndex(index string) string {
	return fmt.Sprintf("%s_%s", s.c.Env, index)
}

func (s *Service) existsOrCreate(c conf.EsIndex) {
	if s.dao.IndexExists(context.TODO(), s.envIndex(c.Index)) {
		log.Info("索引(%s)已存在", s.envIndex(c.Index))
		return
	}
	if !s.dao.CreateIndex(context.TODO(), s.envIndex(c.Index), c.Mapping) {
		panic("创建索引失败")
	}
}

func (s *Service) subproc() {
	defer s.wg.Done()
}

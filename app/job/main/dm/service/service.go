package service

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"

	"go-common/app/job/main/dm/conf"
	"go-common/app/job/main/dm/dao"
	"go-common/app/job/main/dm/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
)

var (
	errSubNotExist = errors.New("subject not exist")
)

// Service rpc service.
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// databus sub
	dmMetaCsmr *databus.Databus
	// cache
	cache *fanout.Fanout
}

// New new rpc service.
func New(c *conf.Config) *Service {
	s := &Service{
		c:          c,
		dao:        dao.New(c),
		dmMetaCsmr: databus.New(c.Databus.DMMetaCsmr),
		cache:      fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
	}
	// 消费DMMeta-T消息
	go s.dmMetaCsmproc()
	return s
}

// Ping check if service is ok.
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}

func (s *Service) dmMetaCsmproc() {
	var (
		err        error
		c          = context.TODO()
		regexIndex = regexp.MustCompile("dm_index_[0-9]+")
	)
	for {
		msg, ok := <-s.dmMetaCsmr.Messages()
		if !ok {
			log.Error("dmmeta binlog consumer exit")
			return
		}
		m := &model.BinlogMsg{}
		if err = json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if regexIndex.MatchString(m.Table) {
			if err = s.trackDMMeta(c, m); err != nil {
				log.Error("s.trackDMMeta(%s) error(%v)", m, err)
				continue
			}
		}
		if err = msg.Commit(); err != nil {
			log.Error("commit offset(%v) error(%v)", msg, err)
		}
	}
}

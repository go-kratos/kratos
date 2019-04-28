package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"sort"
	"strings"
	"time"

	"go-common/app/service/main/history/model"
	"go-common/library/log"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline"
)

func (s *Service) serviceConsumeproc() {
	var (
		err  error
		msgs = s.serviceHisSub.Messages()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.serviceConsumeproc closed")
			return
		}
		if s.c.Job.IgnoreMsg {
			err = msg.Commit()
			log.Info("serviceConsumeproc key:%s partition:%d offset:%d err: %+v, ts:%v ignore", msg.Key, msg.Partition, msg.Offset, err, msg.Timestamp)
			continue
		}
		ms := make([]*model.Merge, 0, 32)
		if err = json.Unmarshal(msg.Value, &ms); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		for _, x := range ms {
			key := fmt.Sprintf("%d-%d-%d", x.Mid, x.Bid, x.Kid)
			s.merge.SyncAdd(context.Background(), key, x)
		}
		err := msg.Commit()
		log.Info("serviceConsumeproc key:%s partition:%d offset:%d err: %+v, len(%v)", msg.Key, msg.Partition, msg.Offset, err, len(ms))
	}
}

func (s *Service) serviceFlush(merges []*model.Merge) {
	// 相同的mid聚合在一起
	sort.Slice(merges, func(i, j int) bool { return merges[i].Mid < merges[j].Mid })
	var ms []*model.Merge
	for _, m := range merges {
		if (len(ms) < s.c.Job.ServiceBatch) || (ms[len(ms)-1].Mid == m.Mid) {
			ms = append(ms, m)
			continue
		}
		s.FlushCache(context.Background(), ms)
		ms = []*model.Merge{m}
	}
	if len(ms) > 0 {
		s.FlushCache(context.Background(), ms)
	}
}

// FlushCache  数据从缓存写入到DB中
func (s *Service) FlushCache(c context.Context, merges []*model.Merge) (err error) {
	var histories []*model.History
	if histories, err = s.dao.HistoriesCache(c, merges); err != nil {
		log.Error("historyDao.Cache(%+v) error(%v)", merges, err)
		return
	}
	prom.BusinessInfoCount.Add("histories-db", int64(len(histories)))
	if err = s.limit.WaitN(context.Background(), len(histories)); err != nil {
		log.Error("s.limit.WaitN(%v) err: %+v", len(histories), err)
	}
	for {
		if err = s.dao.AddHistories(c, histories); err != nil {
			prom.BusinessInfoCount.Add("retry", int64(len(histories)))
			time.Sleep(time.Duration(s.c.Job.RetryTime))
			continue
		}
		break
	}
	s.cache.Do(c, func(c context.Context) {
		for _, merge := range merges {
			limit := s.c.Job.CacheLen
			s.dao.TrimCache(context.Background(), merge.Business, merge.Mid, limit)
		}
	})
	return
}

func (s *Service) initMerge() {
	s.merge = pipeline.NewPipeline(s.c.Merge)
	s.merge.Split = func(a string) int {
		midStr := strings.Split(a, "-")[0]
		return int(crc32.ChecksumIEEE([]byte(midStr)))
	}
	s.merge.Do = func(c context.Context, ch int, values map[string][]interface{}) {
		var merges []*model.Merge
		for _, vs := range values {
			var t int64
			var m *model.Merge
			for _, v := range vs {
				prom.BusinessInfoCount.Incr("dbus-msg")
				if v.(*model.Merge).Time >= t {
					m = v.(*model.Merge)
				}
			}
			if m.Mid%1000 == 0 {
				log.Info("debug: merge mid:%v, ch:%v, value:%+v", m.Mid, ch, m)
			}
			merges = append(merges, m)
		}
		prom.BusinessInfoCount.Add(fmt.Sprintf("ch-%v", ch), int64(len(merges)))
		s.serviceFlush(merges)
	}
	s.merge.Start()
}

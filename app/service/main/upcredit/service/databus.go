package service

import (
	"context"
	"encoding/json"
	"go-common/app/service/main/upcredit/model/canal"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/library/log"
	"strings"
)

const (
	tableUpQualityInfo = "up_quality_info_"
)

func (s *Service) arcCreditLogConsume() {
	defer func() {
		s.wg.Done()
		log.Error("arcCreditLogConsume stop!")
	}()
	var (
		msgs = s.creditLogSub.Messages()
		err  error
	)
	log.Info("arcCreditLogConsume start")
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.arcSub.Messages closed")
			return
		}
		msg.Commit()
		//s.arcMo++
		m := &upcrmmodel.ArgCreditLogAdd{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		var c = context.TODO()
		s.LogCredit(c, m)
		log.Info("arcCreditLogConsume key(%s) value(%s) partition(%d) offset(%d) commit, mid=%d, bustype=%d, optype=%d", msg.Key, msg.Value, msg.Partition, msg.Offset, m.Mid, m.BusinessType, m.OpType)
	}
}

func (s *Service) arcBusinessBinLogCanalConsume() {
	defer func() {
		s.wg.Done()
		log.Error("business CanalConsume stop!")
	}()
	var (
		msgs = s.businessBinLogSub.Messages()
		err  error
	)
	log.Info("business CanalConsume start")
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.businessBinLogSub.Messages closed")
			return
		}
		m := &canal.Msg{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			msg.Commit()
			continue
		}
		if strings.HasPrefix(m.Table, tableUpQualityInfo) {
			s.onBusinessDatabus(m)
		}
		msg.Commit()
		log.Info("businessBinLog consume key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, string(msg.Value), msg.Partition, msg.Offset)
	}
}

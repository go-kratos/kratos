package service

import (
	"context"
	"encoding/json"

	"go-common/app/service/main/assist/model/message"
	"go-common/library/log"
)

func (s *Service) relationConsumer() {
	defer s.wg.Done()
	var (
		msgs = s.relationSub.Messages()
		err  error
		c    = context.TODO()
	)
	for {
		msg, ok := <-msgs
		log.Info("relationConsumer msg info |Key(%+v)|Topic(%+v)|Partition(%+v)", msg.Key, msg.Topic, msg.Partition)
		if !ok {
			log.Error("s.relationSub.Message closed")
			return
		}
		msg.Commit()
		rm := &message.Relation{}
		if err = json.Unmarshal(msg.Value, rm); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			err = nil
			continue
		}
		// 不关注表取mod的分配表名，只关注单向的关系链状态流转
		breakOnDirection := (rm.Old.Attr == 2 && rm.New.Attr == 0) || (rm.Old.Attr == 6 && rm.New.Attr == 0)
		if rm.Action == "update" && breakOnDirection {
			log.Info("relationConsumer rm.Action(%s) OldAttr(%d),NewAttr(%d) mid(%d) fid(%d) table(%s)", rm.Action, rm.Old.Attr, rm.New.Attr, rm.New.MID, rm.New.FID, rm.Table)
			if err = s.DelAssist(c, rm.New.FID, rm.New.MID); err != nil {
				log.Error("DelAssist(%v), mid(%d),fid(%d) error(%v)", string(msg.Value), rm.New.MID, rm.New.FID, err)
				err = nil
				continue
			}
		}
		log.Info("relationConsumer key(%s) partition(%d) offset(%d) commit", msg.Key, msg.Partition, msg.Offset)
	}
}

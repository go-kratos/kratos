package service

import (
	"context"
	"encoding/json"
	"runtime/debug"

	"go-common/app/job/main/ugcpay/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_tableOrderUser     = "order_user"
	_tableAsset         = "asset"
	_tableAssetRelation = "asset_relation"
)

func (s *Service) binlogproc() (err error) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("binlogproc panic(%+v) :\n %s", x, debug.Stack())
			go s.binlogproc()
		}
	}()
	var (
		c = context.Background()
	)
	for res := range s.binlogMQ.Messages() {
		if err != nil {
			log.Error("%+v", err)
			err = nil
		}
		msg := &model.Message{}
		if err = json.Unmarshal(res.Value, msg); err != nil {
			err = errors.WithStack(err)
			continue
		}
		switch msg.Table {
		case _tableOrderUser:
			ms := &model.BinlogOrderUser{}
			if err = json.Unmarshal(msg.New, ms); err != nil {
				err = errors.Wrapf(err, "%s", msg.New)
				continue
			}
			log.Info("Delete order_user cache : %+v", ms)
			if err = s.dao.DelCacheOrderUser(c, ms.OrderID); err != nil {
				continue
			}
		case _tableAsset:
			ms := &model.BinlogAsset{}
			if err = json.Unmarshal(msg.New, ms); err != nil {
				err = errors.Wrapf(err, "%s", msg.New)
				continue
			}
			log.Info("Delete asset cache : %+v", ms)
			if err = s.dao.DelCacheAsset(c, ms.OID, ms.OType, ms.Currency); err != nil {
				continue
			}
		case _tableAssetRelation:
			ms := &model.BinlogAssetRelation{}
			if err = json.Unmarshal(msg.New, ms); err != nil {
				err = errors.Wrapf(err, "%s", msg.New)
				continue
			}
			log.Info("Delete asset_relation cache : %+v", ms)
			if err = s.dao.DelCacheAssetRelationState(c, ms.OID, ms.OType, ms.MID); err != nil {
				continue
			}
		}
		if err = res.Commit(); err != nil {
			err = errors.Wrapf(err, "binlogproc commit")
			continue
		}
		log.Info("binlogproc consume key:%v, topic: %v, part:%v, offset:%v, message %s,", res.Key, res.Topic, res.Partition, res.Offset, res.Value)
	}
	return
}

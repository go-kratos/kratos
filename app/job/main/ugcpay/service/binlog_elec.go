package service

import (
	"context"
	"encoding/json"
	"runtime/debug"
	"strconv"
	"time"

	"go-common/app/job/main/ugcpay/model"
	// "go-common/app/service/main/ugcpay-rank/internal/service/rank"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_tableOldElecOrder       = "elec_pay_order"
	_tableOldElecMessage     = "elec_message"
	_tableOldElecUserSetting = "elec_user_setting"
)

var (
	_ctx = context.Background()
)

func (s *Service) elecbinlogproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("binlogproc panic(%+v) :\n %s", x, debug.Stack())
			go s.binlogproc()
		}
	}()
	log.Info("Start binlogproc")
	var (
		err error
	)
	for res := range s.elecBinlogMQ.Messages() {
		if err != nil {
			log.Error("binlogproc consume key:%v, topic: %v, part:%v, offset:%v, message %s, err: %+v", res.Key, res.Topic, res.Partition, res.Offset, res.Value, err)
			err = nil
		}
		if time.Since(time.Unix(res.Timestamp, 0)) >= time.Hour*24 {
			log.Error("binlogproc consume expired msg, key:%v, topic: %v, part:%v, offset:%v, message %s", res.Key, res.Topic, res.Partition, res.Offset, res.Value)
			continue
		}
		msg := &model.Message{}
		if err = json.Unmarshal(res.Value, msg); err != nil {
			err = errors.WithStack(err)
			continue
		}
		switch msg.Table {
		case _tableOldElecOrder:
			data := &model.DBOldElecPayOrder{}
			if err = json.Unmarshal(msg.New, data); err != nil {
				err = errors.Wrapf(err, "%s", msg.New)
				continue
			}
			if err = s.handleElecOrder(_ctx, data); err != nil {
				continue
			}
		case _tableOldElecMessage:
			data := &model.DBOldElecMessage{}
			if err = json.Unmarshal(msg.New, data); err != nil {
				err = errors.Wrapf(err, "%s", msg.New)
				continue
			}
			if err = s.handleOldElecMessage(_ctx, data); err != nil {
				continue
			}
		case _tableOldElecUserSetting:
			data := &model.DBOldElecUserSetting{}
			if err = json.Unmarshal(msg.New, data); err != nil {
				err = errors.Wrapf(err, "%s", msg.New)
				continue
			}
			if err = s.handleElecUserSetting(_ctx, data); err != nil {
				continue
			}
		default:
			log.Error("binlogproc unknown table: %s", msg.Table)
		}
		if err = res.Commit(); err != nil {
			err = errors.Wrapf(err, "binlogproc commit")
			continue
		}
		log.Info("binlogproc consume msg, key:%v, topic: %v, part:%v, offset:%v, message %s", res.Key, res.Topic, res.Partition, res.Offset, res.Value)
	}
	log.Info("End binlogproc")
}

func (s *Service) handleOldElecMessage(ctx context.Context, msg *model.DBOldElecMessage) (err error) {
	log.Info("handleOldElecMessage message: %+v", msg)
	var (
		ver     int64
		verTime time.Time
		avID    int64
	)
	if verTime, err = time.Parse("2006-01", msg.DateVer); err != nil {
		err = errors.WithStack(err)
		return
	}
	ver = monthlyBillVer(verTime)

	if msg.AVID != "" {
		if avID, err = strconv.ParseInt(msg.AVID, 10, 64); err != nil {
			log.Error("%+v", errors.WithStack(err))
			avID = 0
			err = nil
		}
	}

	switch msg.Type {
	// 用户对up主留言
	case 1:
		dbMSG := &model.DBElecMessage{
			ID:      msg.ID,
			Ver:     ver,
			AVID:    avID,
			UPMID:   msg.RefMID,
			PayMID:  msg.MID,
			Message: msg.Message,
			Replied: msg.State == 1,
			Hidden:  msg.State == 2,
			CTime:   msg.ParseCTime(),
			MTime:   msg.ParseMTime(),
		}
		if err = s.dao.UpsertElecMessage(ctx, dbMSG); err != nil {
			return
		}
		if err = s.dao.RankElecUpdateMessage(ctx, dbMSG.AVID, dbMSG.UPMID, dbMSG.PayMID, dbMSG.Ver, dbMSG.Message, dbMSG.Hidden); err != nil {
			return
		}
	// up主回复用户
	case 2:
		dbReply := &model.DBElecReply{
			ID:     msg.ID,
			MSGID:  msg.RefID,
			Reply:  msg.Message,
			Hidden: msg.State == 2,
			CTime:  msg.ParseCTime(),
			MTime:  msg.ParseMTime(),
		}
		if err = s.dao.UpsertElecReply(ctx, dbReply); err != nil {
			return
		}
	default:
		log.Error("old_ele_message unknown type: %+v", msg)
	}
	return
}

func (s *Service) handleElecOrder(ctx context.Context, order *model.DBOldElecPayOrder) (err error) {
	log.Info("handleElecOrder order: %+v", order)

	if !order.IsPaid() {
		return
	}

	var ok bool
	ok, err = s.dao.AddCacheOrderID(ctx, order.OrderID)
	if err != nil {
		return
	}
	// 重复消费
	if !ok {
		log.Info("handleElecOrder order: %+v, has consumed before", order)
		err = nil
		return
	}

	if order.IsHiddnRank() {
		log.Info("handleElecOrder order: %+v which app_id == 19", order)
		return
	}

	tradeInfo, err := s.dao.RawOldElecTradeInfo(ctx, order.OrderID)
	if err != nil {
		return
	}
	avID := int64(0)
	if tradeInfo != nil {
		// log.Info("RawOldElecTradeInfo data not found, order: %+v", order)
		if avID, err = strconv.ParseInt(tradeInfo.AVID, 10, 64); err != nil {
			log.Error("handleElecOrder cant convert avID from: %s, err: %+v", tradeInfo.AVID, err)
			avID = 0
		}
	}

	var (
		ver    = monthlyBillVer(order.ParseMTime())
		hidden = false
	)
	// 更新DB
	tx, err := s.dao.BeginTranRank(ctx)
	if err != nil {
		return
	}
	rollbackFN := func() {
		if theErr := s.dao.DelCacheOrderID(ctx, order.OrderID); theErr != nil {
			log.Error("%+v", theErr)
		}
		tx.Rollback()
	}
	if avID != 0 {
		if err = s.dao.TXUpsertElecAVRank(ctx, tx, 0, avID, order.UPMID, order.PayMID, order.ElecNum, hidden); err != nil {
			rollbackFN()
			return
		}
		if err = s.dao.TXUpsertElecAVRank(ctx, tx, ver, avID, order.UPMID, order.PayMID, order.ElecNum, hidden); err != nil {
			rollbackFN()
			return
		}
	}

	if err = s.dao.TXUpsertElecUPRank(ctx, tx, 0, order.UPMID, order.PayMID, order.ElecNum, hidden); err != nil {
		rollbackFN()
		return
	}
	if err = s.dao.TXUpsertElecUPRank(ctx, tx, ver, order.UPMID, order.PayMID, order.ElecNum, hidden); err != nil {
		rollbackFN()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}

	err = s.dao.RankElecUpdateOrder(ctx, avID, order.UPMID, order.PayMID, ver, order.ElecNum)
	return
}

func (s *Service) handleElecUserSetting(ctx context.Context, setting *model.DBOldElecUserSetting) (err error) {
	if setting.Status > 0 {
		log.Info("handleElecUserSetting add setting: %+v", setting)
		err = s.dao.ElecAddSetting(ctx, model.DefaultUserSetting, setting.MID, setting.BitValue())
	} else {
		log.Info("handleElecUserSetting delete setting: %+v", setting)
		err = s.dao.ElecDeleteSetting(ctx, model.DefaultUserSetting, setting.MID, setting.BitValue())
	}
	if err != nil {
		return
	}
	// 清理缓存
	if err = s.dao.DelCacheUserSetting(ctx, setting.MID); err != nil {
		log.Error("DelCacheUserSetting: %d, err: %+v", setting.MID, err)
		err = nil
	}
	return
}

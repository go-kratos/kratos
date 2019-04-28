package v1

import (
	"context"
	"encoding/json"

	v1pb "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/app/service/live/dao-anchor/dao"
	"go-common/app/service/live/dao-anchor/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

//ValidDaysStatistics  有效天数相关数据(同步落地)
func (s *ConsumerService) ValidDaysStatistics(ctx context.Context, msg *databus.Message) (err error) {
	defer msg.Commit()
	return s.internalValidDaysStatistics(ctx, msg)
}

func (s *ConsumerService) internalValidDaysStatistics(ctx context.Context, msg *databus.Message) (err error) {
	msgContent := new(model.TopicCommonMsg)
	if err = json.Unmarshal(msg.Value, &msgContent); err != nil {
		log.Error("ValidDaysStatistics_error:msg=%v;err=%v;msgInfo=%v", msg, err, msgContent)
		return
	}
	log.Info("ValidDaysStatistics:msgInfo=%v", string(msg.Value))
	cycle := msgContent.Cycle
	value := msgContent.Value
	roomId := msgContent.RoomId
	attrId := dao.ATTRID_VALID_LIVE_DAYS
	typ := msgContent.Type
	reply, err := getValidDaysAttrInfo(cycle, typ)
	if err != nil {
		log.Error("ValidDaysStatistics_check_attr_sub_id_error")
		return
	}
	attrSubId := reply.AttrSubId

	//contentType := reply.ContentType
	//redisKeyInfo := s.dao.SRoomRecordCurrent(contentType, roomId)
	//s.dao.Incr(ctx, redisKeyInfo.RedisKey, 1, redisKeyInfo.TimeOut)
	//有效开播天数为离线计算数据，直接入库
	req := &v1pb.RoomAttrCreateReq{
		RoomId:    roomId,
		AttrId:    int64(attrId),
		AttrSubId: attrSubId,
		AttrValue: value,
	}
	resp, err := s.dao.RoomAttrCreate(ctx, req)
	if err != nil || resp.AffectedRows <= 0 {
		log.Error("ValidDaysStatistics: msg=%v; err=%v; req=%v", msg, err, req)
	}
	return
}

func getValidDaysAttrInfo(cycle int64, typ int64) (resp *attrInfo, err error) {

	days := cycle / 60 / 60 / 24
	resp = &attrInfo{}
	switch days {
	case 7:
		if typ == dao.VALID_LIVE_DAYS_TYPE_1 {
			resp.AttrSubId = dao.ATTRSUBID_VALID_LIVE_DAYS_TYPE_1_DAY_7
			resp.ContentType = dao.VALID_LIVE_DAYS_TYPE_1_DAY_7
		} else if typ == dao.VALID_LIVE_DAYS_TYPE_2 {
			resp.AttrSubId = dao.ATTRSUBID_VALID_LIVE_DAYS_TYPE_2_DAY_7
			resp.ContentType = dao.VALID_LIVE_DAYS_TYPE_2_DAY_7
		}
		break
	case 14:
		if typ == dao.VALID_LIVE_DAYS_TYPE_1 {
			resp.AttrSubId = dao.ATTRSUBID_VALID_LIVE_DAYS_TYPE_1_DAY_14
			resp.ContentType = dao.VALID_LIVE_DAYS_TYPE_1_DAY_14
		}
	case 30:
		if typ == dao.VALID_LIVE_DAYS_TYPE_2 {
			resp.AttrSubId = dao.ATTRSUBID_VALID_LIVE_DAYS_TYPE_2_DAY_30
			resp.ContentType = dao.VALID_LIVE_DAYS_TYPE_2_DAY_30
		}
		break
	default:
		err = ecode.DaoAnchorCheckAttrSubIdERROR
		break
	}
	return
}

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

//PopularityStatistics  离线人气相关统计，同步落地
func (s *ConsumerService) PopularityStatistics(ctx context.Context, msg *databus.Message) (err error) {
	defer msg.Commit()
	return s.internalPopularityStatistics(ctx, msg)
}

func (s *ConsumerService) internalPopularityStatistics(ctx context.Context, msg *databus.Message) (err error) {
	msgContent := new(model.TopicCommonMsg)
	if err = json.Unmarshal(msg.Value, &msgContent); err != nil {
		log.Error("PopularityStatistics_error:msg=%v;err=%v;msgInfo=%v", msg, err, msgContent)
		return
	}
	log.Info("PopularityStatistics:msgInfo=%v", string(msg.Value))
	cycle := msgContent.Cycle
	value := msgContent.Value
	roomId := msgContent.RoomId
	attrId := dao.ATTRID_POPULARITY
	//expireTime := time.Now().Unix()
	//人气过渡时间，先消费databus,若为1分钟，则为实时人气，插入extend表
	if cycle == 60 {
		req := &v1pb.RoomExtendUpdateReq{
			RoomId:          roomId,
			PopularityCount: value,
			Fields:		[]string{"popularity_count"},
		}
		s.dao.RoomExtendUpdate(ctx, req)
		return
	}
	reply, err := getPopularityAttrInfo(cycle)
	if err != nil {
		log.Error("PopularityStatistics_check_attr_sub_id_error")
		return
	}
	attrSubId := reply.AttrSubId
	//contentType := reply.ContentType
	//redisKeyInfo := s.dao.SRoomRecordCurrent(contentType, roomId)
	//s.dao.Incr(ctx, redisKeyInfo.RedisKey, 1, redisKeyInfo.TimeOut)
	if attrSubId == dao.ATTRSUBID_POPULARITY_MAX_TO_ARG_7 ||
		attrSubId == dao.ATTRSUBID_POPULARITY_MAX_TO_ARG_30 {
		//@todo 增加过期标记
		req := &v1pb.RoomAttrCreateReq{
			RoomId:    roomId,
			AttrId:    int64(attrId),
			AttrSubId: attrSubId,
			AttrValue: value,
		}
		resp, err := s.dao.RoomAttrCreate(ctx, req)
		if err != nil || resp.AffectedRows <= 0 {
			log.Error("PopularityStatistics: msg=%v; err=%v; req=%v", msg, err, req)
		}
	}
	return
}

func getPopularityAttrInfo(cycle int64) (resp *attrInfo, err error) {

	days := cycle / 60 / 60 / 24
	resp = &attrInfo{}
	switch days {
	case 7:
		resp.AttrSubId = dao.ATTRSUBID_POPULARITY_MAX_TO_ARG_7
		resp.ContentType = dao.POPULARITY_MAX_TO_ARG_7
		break
	case 30:
		resp.AttrSubId = dao.ATTRSUBID_POPULARITY_MAX_TO_ARG_30
		resp.ContentType = dao.POPULARITY_MAX_TO_ARG_30
		break
	default:
		err = ecode.DaoAnchorCheckAttrSubIdERROR
		break
	}
	return
}

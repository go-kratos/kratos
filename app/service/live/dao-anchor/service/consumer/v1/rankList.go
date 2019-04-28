package v1

import (
	"context"
	"encoding/json"

	v1pb "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/app/service/live/dao-anchor/dao"
	"go-common/app/service/live/dao-anchor/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

//LiveRoomTag  房间实时标签，同步落地
func (s *ConsumerService) LiveRankList(ctx context.Context, msg *databus.Message) (err error) {
	defer msg.Commit()
	return s.internalLiveRankList(ctx, msg)
}

func (s *ConsumerService) internalLiveRankList(ctx context.Context, msg *databus.Message) (err error) {
	msgContent := new(model.LiveRankListMsg)
	if err = json.Unmarshal(msg.Value, &msgContent); err != nil {
		log.Error("LiveRankList_error:msg=%v;err=%v;msgInfo=%v", msg, err, msgContent)
		return
	}
	log.Info("LiveRankList:msgInfo=%v", string(msg.Value))
	//expireTime := msgContent.ExpireTime // todo @wangyao
	attrId := dao.ATTRID_RANK_LIST
	attrSubId := msgContent.RankId
	rankList := msgContent.RankList
	if rankList == nil || len(rankList) <= 0 {
		log.Error("LiveRankList_get_rank_list_error:msg=%v", msgContent)
		return
	}
	deleteReq := &v1pb.DeleteAttrReq{
		AttrId:    int64(attrId),
		AttrSubId: attrSubId,
	}
	_, err = s.dao.DeleteAttr(ctx, deleteReq)
	if err != nil {
		log.Error("LiveRankList_DeleteAttr_Error: msg=%v; err=%v; deleteReq=%v", msg, err, deleteReq)
		return
	}
	for roomId, value := range rankList {
		req := &v1pb.RoomAttrSetExReq{
			RoomId:    roomId,
			AttrId:    int64(attrId),
			AttrSubId: attrSubId,
			AttrValue: value,
		}
		_, err := s.dao.RoomAttrSetEx(ctx, req)
		if err != nil {
			log.Error("LiveRankList: msg=%v; err=%v; req=%v", string(msg.Value), err, req)
		}
	}
	return
}

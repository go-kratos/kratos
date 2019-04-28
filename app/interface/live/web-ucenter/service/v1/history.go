package v1

import (
	"context"

	"github.com/pkg/errors"

	historypb "go-common/app/interface/live/web-ucenter/api/http/v1"
	"go-common/app/interface/live/web-ucenter/conf"
	"go-common/app/interface/live/web-ucenter/dao"
	historydao "go-common/app/interface/live/web-ucenter/dao/history"
	"go-common/app/service/live/room/api/liverpc/v2"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *historydao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: historydao.New(c),
	}
	return s
}

// GetHistoryByUid 获取直播历史记录
func (s *Service) GetHistoryByUid(ctx context.Context, req *historypb.GetHistoryReq) (resp *historypb.GetHistoryResp, err error) {
	uid, ok := metadata.Value(ctx, metadata.Mid).(int64)
	if !ok {
		err = errors.Wrap(err, "未取到uid")
		return
	}
	mainHistoryInfo, err := s.dao.GetMainHistory(ctx, int32(uid))
	if err != nil {
		err = errors.Wrap(err, "Call GetMainHistory err")
		return
	}
	if mainHistoryInfo == nil {
		return
	}
	RoomIds := make([]int64, 0)
	for _, v := range mainHistoryInfo {
		RoomIds = append(RoomIds, v.RoomId)
	}
	reply, err := dao.RoomAPI.V2Room.GetByIds(ctx, &v2.RoomGetByIdsReq{Ids: RoomIds})
	if err != nil {
		err = errors.Wrap(err, "Call GetByIds err")
		return
	}
	if reply.GetCode() != 0 {
		err = ecode.Int(int(reply.GetCode()))
		return
	}
	roomInfos := reply.Data
	resp = &historypb.GetHistoryResp{}
	for _, RoomId := range RoomIds {
		list := &historypb.GetHistoryResp_List{}
		room, ok := roomInfos[RoomId]
		if !ok {
			log.Warn("[GetHistoryByUid] req(%v), uid(%d), failed to get room(%d) info from (%v)", req, uid, RoomId, roomInfos)
			continue
		}
		list.Roomid = RoomId
		list.Uid = int32(room.Uid)
		list.Uname = room.Uname
		list.Title = room.Title
		list.Face = room.Face
		list.LiveStatus = int32(room.LiveStatus)
		list.FansNum = int32(room.Attentions)
		list.AreaV2Id = int32(room.AreaV2Id)
		list.AreaV2Name = room.AreaV2Name
		list.LiveStatus = int32(room.LiveStatus)
		list.UserCover = room.UserCover
		list.AreaV2ParentId = int32(room.AreaV2ParentId)
		list.AreaV2ParentName = room.AreaV2ParentName
		list.Tags = room.Tags
		resp.List = append(resp.List, list)
	}
	resp.Title = "哔哩哔哩直播 - 观看历史"
	resp.Count = int32(len(roomInfos))

	return
}

// DelHistory 删除直播历史记录
func (s *Service) DelHistory(ctx context.Context, req *historypb.DelHistoryReq) (resp *historypb.DelHistoryResp, err error) {
	uid, ok := metadata.Value(ctx, metadata.Mid).(int64)
	if !ok {
		err = errors.Wrap(err, "未取到uid")
		return
	}
	reply, err := s.dao.DelHistory(ctx, uid)
	resp = &historypb.DelHistoryResp{}
	if err != nil || reply != 0 {
		err = ecode.Int(int(reply))
		return
	}
	return
}

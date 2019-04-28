package service

import (
	"context"
	daoAnchorV1 "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/library/log"
)


func (s *Service) getOnlineListByAttrs(ctx context.Context, attrs []*daoAnchorV1.AttrReq) (roomList map[int64]*daoAnchorV1.AttrResp, err error){
	roomList = make(map[int64]*daoAnchorV1.AttrResp)
	RoomOnlineListByAttrsResp, err := s.daoAnchor.RoomOnlineListByAttrs(ctx, &daoAnchorV1.RoomOnlineListByAttrsReq{
		Attrs: attrs,
	})
	if err != nil {
		log.Error("[getOnlineListByAttrs]rpc_error:%+v", err)
		return
	}

	if RoomOnlineListByAttrsResp == nil || len(RoomOnlineListByAttrsResp.Attrs) <= 0 {
		log.Error("[getOnlineListByAttrs]RoomList_empty")
		return
	}

	roomList = RoomOnlineListByAttrsResp.Attrs
	return
}
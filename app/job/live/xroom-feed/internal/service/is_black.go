package service

import (
	"context"
	"time"

	roomServeice "go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

func (s *Service) isBlackRoomID(roomID int64) bool {
	bl := s.indexBlackList.Load()

	backList, ok := bl.(map[int64]bool)
	if !ok {
		log.Warn("[isBlackRoomID] cache err: %+v", bl)
		backList = s.getBlackList()
	}
	if _, ok := backList[roomID]; !ok {
		return false
	}
	return true
}

func (s *Service) getBlackList() (res map[int64]bool) {
	cCtx := metadata.NewContext(context.TODO(), metadata.MD{metadata.Color: env.Color})
	ctx, _ := context.WithTimeout(cCtx, time.Duration(200*time.Millisecond))

	res = make(map[int64]bool)
	req := &roomServeice.RoomMngIsBlackReq{}
	resp, err := s.roomService.V1RoomMng.IsBlack(ctx, req)
	if err != nil {
		log.Error("[getBlackList] rpc V1RoomMng.IsBlack err:%v", err)
		return
	}

	if resp.Code != 0 {
		log.Error("[getBlackList] rpc V1RoomMng.IsBlack err code:%d", resp.Code)
		return
	}
	for roomID := range resp.Data {
		res[roomID] = true
	}
	return
}

func (s *Service) loadBlackList() {
	m := s.getBlackList()
	s.indexBlackList.Store(m)
}

func (s *Service) blackListProc() {
	time.Sleep(time.Second * 60)

	s.loadBlackList()
}

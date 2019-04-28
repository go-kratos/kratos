package service

import (
	"context"
	"sort"
	"strings"

	"go-common/app/service/main/broadcast/model"
)

var (
	_emptyTops = make([]*model.Top, 0)
)

// OnlineTop get the top online.
func (s *Service) OnlineTop(c context.Context, business string, n int) (tops []*model.Top, err error) {
	for roomKey, cnt := range s.roomCount {
		if strings.HasPrefix(roomKey, business) {
			_, roomID, err := model.DecodeRoomKey(roomKey)
			if err != nil {
				continue
			}
			top := &model.Top{
				RoomID: roomID,
				Count:  cnt,
			}
			tops = append(tops, top)
		}
	}
	sort.Slice(tops, func(i, j int) bool {
		return tops[i].Count > tops[j].Count
	})
	if len(tops) > n {
		tops = tops[:n]
	}
	if len(tops) == 0 {
		tops = _emptyTops
	}
	return
}

// OnlineRoom get rooms online.
func (s *Service) OnlineRoom(c context.Context, business string, rooms []string) (res map[string]int32, err error) {
	res = make(map[string]int32, len(rooms))
	for _, roomID := range rooms {
		res[roomID] = s.roomCount[model.EncodeRoomKey(business, roomID)]
	}
	return
}

// OnlineTotal get all online.
func (s *Service) OnlineTotal(c context.Context) (int64, int64) {
	return s.totalIPs, s.totalConns
}

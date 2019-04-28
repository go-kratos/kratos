package live

import (
	"strconv"

	"go-common/app/interface/main/app-view/model"
)

type Live struct {
	Mid    int64  `json:"mid"`
	RoomID int64  `json:"roomid"`
	Title  string `json:"title"`
	URI    string `json:"uri,omitempty"`
}

type RoomInfo struct {
	UID           string `json:"uid"`
	RoomID        string `json:"roomid"`
	Title         string `json:"title"`
	BroadcastType int    `json:"broadcast_type"`
}

func (l *Live) LiveChange(r *RoomInfo) {
	var (
		err error
	)
	if l.Mid, err = strconv.ParseInt(r.UID, 10, 64); err != nil {
		return
	}
	if l.RoomID, err = strconv.ParseInt(r.RoomID, 10, 64); err != nil {
		return
	}
	l.Title = r.Title
	l.URI = model.FillURI(model.GotoLive, strconv.FormatInt(l.RoomID, 10), model.LiveRoomHandler(r.BroadcastType))
}

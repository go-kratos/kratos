package server

import (
	"encoding/json"
	"fmt"

	"go-common/library/log"
)

func encodeRoomID(aid, cid int64) string {
	return fmt.Sprintf("video://%d/%d", aid, cid)
}

// NoAuthParam .
type NoAuthParam struct {
	Key    string `json:"key,omitempty"`
	Aid    int64  `json:"aid,omitempty"`
	RoomID int64  `json:"roomid,omitempty"`
	UserID int64  `json:"uid,omitempty"`
	From   int64  `json:"from,omitempty"`
}

// NoAuth .
func (s *Server) NoAuth(ver int16, token []byte, ip string) (userID int64, roomID, key string, rpt *Report, err error) {
	param := NoAuthParam{}
	if err = json.Unmarshal(token, &param); err != nil {
		log.Error("json.Unmarshal(%d, %s) error(%v)", ver, token, err)
		return
	}
	if param.Key != "" {
		key = param.Key
	} else {
		key = s.NextKey()
	}
	userID = param.UserID
	roomID = encodeRoomID(param.Aid, param.RoomID)
	rpt = &Report{
		From: param.From,
		Aid:  param.Aid,
		Cid:  param.RoomID,
		Mid:  param.UserID,
		Key:  key,
		IP:   ip,
	}
	return
}

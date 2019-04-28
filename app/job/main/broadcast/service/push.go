package service

import (
	"encoding/json"
	"fmt"

	pb "go-common/app/interface/main/broadcast/api/grpc/v1"
	"go-common/app/service/main/broadcast/model"
	"go-common/library/log"
)

const (
	_pushMsg          = "push"
	_broadcastMsg     = "broadcast"
	_broadcastRoomMsg = "broadcast_room"
)

type databusMsg struct {
	Type        string          `json:"type,omitempty"`
	Operation   int32           `json:"operation,omitempty"`
	Server      string          `json:"server,omitempty"`
	Keys        []string        `json:"keys,omitempty"`
	Room        string          `json:"room,omitempty"`
	Speed       int32           `json:"speed,omitempty"`
	Platform    string          `json:"platform,omitempty"`
	ContentType int32           `json:"content_type,omitempty"`
	Message     json.RawMessage `json:"message,omitempty"`
}

func (s *Service) pushMsg(msg []byte) (err error) {
	m := &databusMsg{}
	if err = json.Unmarshal(msg, m); err != nil {
		log.Error("json.Unmarshal(%s) error(%s)", msg, err)
		return
	}
	log.Info("push message:%s", msg)
	switch m.Type {
	case _pushMsg:
		if len(m.Keys) == 0 {
			log.Error("service push keys is invalid:%+v", m)
			return
		}
		s.pushKeys(m.Operation, m.Server, m.Keys, m.Message, m.ContentType)
	case _broadcastMsg:
		s.broadcast(m.Operation, m.Message, m.Speed, m.Platform, m.ContentType)
	case _broadcastRoomMsg:
		if m.Room == "" {
			log.Error("service broadcast room is invalid:%+v", m)
			return
		}
		if err = s.room(m.Room).Push(m.Operation, m.Message, m.ContentType); err != nil {
			log.Error("room.Push(%s) roomId:%s error(%v)", m.Message, m.Room, err)
		}
		// NOTE: 弹幕兼容老协议推送，等web播放器接完可以去掉了
		if m.Operation == 1000 {
			s.broadcastRoom(m.Room, 5, []byte(fmt.Sprintf(`{"cmd":"DM","info":%s}`, m.Message)))
		}
	default:
		log.Error("unknown message type:%s", m.Type)
	}
	return
}

// pushKeys push a message to a batch of subkeys.
func (s *Service) pushKeys(operation int32, serverID string, subKeys []string, body []byte, contentType int32) {
	p := &model.Proto{
		Ver:         1,
		Operation:   operation,
		ContentType: contentType,
		Body:        body,
	}
	var args = pb.PushMsgReq{
		Keys:    subKeys,
		ProtoOp: operation,
		Proto:   p,
	}
	if c, ok := s.cometServers[serverID]; ok {
		if err := c.Push(&args); err != nil {
			log.Error("c.Push(%v) serverId:%s error(%v)", args, serverID, err)
		}
	}
}

// broadcast broadcast a message to all.
func (s *Service) broadcast(operation int32, body []byte, speed int32, platform string, contentType int32) {
	p := &model.Proto{
		Ver:         1,
		Operation:   operation,
		ContentType: contentType,
		Body:        body,
	}
	comets := s.cometServers
	speed /= int32(len(comets))
	var args = pb.BroadcastReq{
		ProtoOp:  operation,
		Proto:    p,
		Speed:    speed,
		Platform: platform,
	}
	for serverID, c := range comets {
		if err := c.Broadcast(&args); err != nil {
			log.Error("c.Broadcast(%v) serverId:%s error(%v)", args, serverID, err)
		}
	}
}

// broadcastRoomRawBytes broadcast aggregation messages to room.
func (s *Service) broadcastRoomRawBytes(roomID string, body []byte) {
	args := pb.BroadcastRoomReq{
		RoomID: roomID,
		Proto: &model.Proto{
			Ver:       1,
			Operation: model.OpRaw,
			Body:      body,
		},
	}
	comets := s.cometServers
	for serverID, c := range comets {
		if err := c.BroadcastRoom(&args); err != nil {
			log.Error("c.BroadcastRoom(%v) roomID:%s serverId:%s  error(%v)", args, roomID, serverID, err)
		}
	}
}

// broadcastRoom broadcast messages to room.
func (s *Service) broadcastRoom(roomID string, op int32, body []byte) {
	args := pb.BroadcastRoomReq{
		RoomID: roomID,
		Proto: &model.Proto{
			Ver:       1,
			Operation: op,
			Body:      body,
		},
	}
	comets := s.cometServers
	for serverID, c := range comets {
		if err := c.BroadcastRoom(&args); err != nil {
			log.Error("c.BroadcastRoom(%v) roomID:%s serverId:%s  error(%v)", args, roomID, serverID, err)
		}
	}
}

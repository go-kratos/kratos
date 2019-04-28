package dao

import (
	"context"
	"encoding/json"

	"go-common/library/log"
)

const (
	_pushMsg          = "push"
	_broadcastMsg     = "broadcast"
	_broadcastRoomMsg = "broadcast_room"
)

type pushMsg struct {
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

// PushMsg push a message to databus.
func (d *Dao) PushMsg(c context.Context, op int32, server, msg string, keys []string, contentType int32) (err error) {
	pushMsg := &pushMsg{
		Type:        _pushMsg,
		Operation:   op,
		Server:      server,
		Keys:        keys,
		Message:     []byte(msg),
		ContentType: contentType,
	}
	if err = d.pushBus.Send(c, keys[0], pushMsg); err != nil {
		log.Error("PushMsg.send(server:%v,pushMsg:%v).error(%v)", server, pushMsg, err)
	}
	return
}

// BroadcastRoomMsg push a message to databus.
func (d *Dao) BroadcastRoomMsg(c context.Context, op int32, room, msg string, contentType int32) (err error) {
	pushMsg := &pushMsg{
		Type:        _broadcastRoomMsg,
		Operation:   op,
		Room:        room,
		Message:     []byte(msg),
		ContentType: contentType,
	}
	if err = d.pushBus.Send(c, room, pushMsg); err != nil {
		log.Error("BroadcastRoomMsg.send(room:%v,pushMsg:%v).error(%v)", room, pushMsg, err)
	}
	return
}

// BroadcastMsg push a message to databus.
func (d *Dao) BroadcastMsg(c context.Context, op, speed int32, msg, platform string, contentType int32) (err error) {
	pushMsg := &pushMsg{
		Operation:   op,
		Type:        _broadcastMsg,
		Speed:       speed,
		Message:     []byte(msg),
		Platform:    platform,
		ContentType: contentType,
	}
	if err = d.pushBus.Send(c, _broadcastMsg, pushMsg); err != nil {
		log.Error("BroadcastMsg.send(_broadcastMsg:%v,pushMsg:%v).error(%v)", _broadcastMsg, pushMsg, err)
	}
	return
}

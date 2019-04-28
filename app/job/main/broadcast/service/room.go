package service

import (
	"time"

	"go-common/app/service/main/broadcast/libs/bytes"
	"go-common/app/service/main/broadcast/model"
	"go-common/library/log"
)

// RoomOptions room options.
type RoomOptions struct {
	BatchNum   int
	SignalTime time.Duration
}

// Room room.
type Room struct {
	s     *Service
	id    string
	proto chan *model.Proto
}

var (
	roomReadyProto = new(model.Proto)
)

// NewRoom new a room struct, store channel room info.
func NewRoom(s *Service, id string, options RoomOptions) (r *Room) {
	r = &Room{
		s:     s,
		id:    id,
		proto: make(chan *model.Proto, options.BatchNum*2),
	}
	go r.pushproc(options.BatchNum, options.SignalTime)
	return
}

// Push push msg to the room, if chan full discard it.
func (r *Room) Push(op int32, msg []byte, contentType int32) (err error) {
	var p = &model.Proto{
		Ver:         1,
		Operation:   op,
		ContentType: contentType,
		Body:        msg,
	}
	select {
	case r.proto <- p:
	default:
		err = ErrRoomFull
	}
	return
}

// pushproc merge proto and push msgs in batch.
func (r *Room) pushproc(batch int, sigTime time.Duration) {
	var (
		n    int
		last time.Time
		p    *model.Proto
		buf  = bytes.NewWriterSize(int(model.MaxBodySize))
	)
	log.Info("start room:%s goroutine", r.id)
	td := time.AfterFunc(sigTime, func() {
		select {
		case r.proto <- roomReadyProto:
		default:
		}
	})
	defer td.Stop()
	for {
		if p = <-r.proto; p == nil {
			break // exit
		} else if p != roomReadyProto {
			// merge buffer ignore error, always nil
			p.WriteTo(buf)
			if n++; n == 1 {
				last = time.Now()
				td.Reset(sigTime)
				continue
			} else if n < batch {
				if sigTime > time.Since(last) {
					continue
				}
			}
		} else {
			if n == 0 {
				break
			}
		}
		r.s.broadcastRoomRawBytes(r.id, buf.Buffer())
		// TODO use reset buffer
		// after push to room channel, renew a buffer, let old buffer gc
		buf = bytes.NewWriterSize(buf.Size())
		n = 0
		if r.s.conf.Room.Idle != 0 {
			td.Reset(time.Duration(r.s.conf.Room.Idle))
		} else {
			td.Reset(time.Minute)
		}
	}
	r.s.roomsMutex.Lock()
	delete(r.s.rooms, r.id)
	r.s.roomsMutex.Unlock()
	log.Info("room:%s goroutine exit", r.id)
}

package server

import (
	"sync"

	"go-common/app/service/main/broadcast/model"
)

// Room is a room.
type Room struct {
	ID        string
	rLock     sync.RWMutex
	next      *Channel
	drop      bool
	Online    int32 // dirty read is ok
	AllOnline int32
}

// NewRoom new a room struct, store channel room info.
func NewRoom(id string) (r *Room) {
	r = new(Room)
	r.ID = id
	r.drop = false
	r.next = nil
	r.AllOnline = roomOnline(id)
	return
}

// Put put channel into the room.
func (r *Room) Put(ch *Channel) (err error) {
	r.rLock.Lock()
	if !r.drop {
		if r.next != nil {
			r.next.Prev = ch
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch // insert to header
		r.Online++
	} else {
		err = ErrRoomDroped
	}
	r.rLock.Unlock()
	return
}

// Del delete channel from the room.
func (r *Room) Del(ch *Channel) bool {
	r.rLock.Lock()
	if ch.Next != nil {
		// if not footer
		ch.Next.Prev = ch.Prev
	}
	if ch.Prev != nil {
		// if not header
		ch.Prev.Next = ch.Next
	} else {
		r.next = ch.Next
	}
	r.Online--
	r.drop = (r.Online == 0)
	r.rLock.Unlock()
	return r.drop
}

// Push push msg to the room, if chan full discard it.
func (r *Room) Push(p *model.Proto) {
	r.rLock.RLock()
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Push(p)
	}
	r.rLock.RUnlock()
}

// Close close the room.
func (r *Room) Close() {
	r.rLock.RLock()
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Close()
	}
	r.rLock.RUnlock()
}

// OnlineNum the room all online.
func (r *Room) OnlineNum() int32 {
	if r.AllOnline > 0 {
		return r.AllOnline
	}
	return r.Online
}

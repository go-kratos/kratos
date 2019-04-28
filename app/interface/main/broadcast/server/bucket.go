package server

import (
	"sync"
	"sync/atomic"

	"go-common/app/interface/main/broadcast/conf"
	"go-common/app/service/main/broadcast/model"
)

// ProtoRoom room proto.
type ProtoRoom struct {
	RoomID string
	Proto  *model.Proto
}

// Bucket is a channel holder.
type Bucket struct {
	c     *conf.Bucket
	cLock sync.RWMutex        // protect the channels for chs
	chs   map[string]*Channel // map sub key to a channel
	// room
	rooms       map[string]*Room // bucket room channels
	routines    []chan ProtoRoom
	routinesNum uint64

	ipCnts     map[string]int32
	roomIPCnts map[string]int32
}

// NewBucket new a bucket struct. store the key with im channel.
func NewBucket(c *conf.Bucket) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, c.Channel)
	b.ipCnts = make(map[string]int32)
	b.roomIPCnts = make(map[string]int32)
	b.c = c
	b.rooms = make(map[string]*Room, c.Room)
	b.routines = make([]chan ProtoRoom, c.RoutineAmount)
	for i := uint64(0); i < c.RoutineAmount; i++ {
		c := make(chan ProtoRoom, c.RoutineSize)
		b.routines[i] = c
		go b.roomproc(c)
	}
	return
}

// ChannelCount channel count in the bucket
func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

// RoomCount room count in the bucket
func (b *Bucket) RoomCount() int {
	return len(b.rooms)
}

// RoomsCount get all room id where online number > 0.
func (b *Bucket) RoomsCount() (res map[string]int32) {
	var (
		roomID string
		room   *Room
	)
	b.cLock.RLock()
	res = make(map[string]int32)
	for roomID, room = range b.rooms {
		if room.Online > 0 {
			res[roomID] = room.Online
		}
	}
	b.cLock.RUnlock()
	return
}

// ChangeRoom change ro room
func (b *Bucket) ChangeRoom(nrid string, ch *Channel) (err error) {
	var (
		nroom *Room
		ok    bool
		oroom = ch.Room
	)
	// change to no room
	if nrid == model.NoRoom {
		if oroom != nil && oroom.Del(ch) {
			b.DelRoom(oroom)
		}
		ch.Room = nil
		return
	}
	b.cLock.Lock()
	if nroom, ok = b.rooms[nrid]; !ok {
		nroom = NewRoom(nrid)
		b.rooms[nrid] = nroom
	}
	b.cLock.Unlock()
	if err = nroom.Put(ch); err != nil {
		return
	}
	ch.Room = nroom
	if oroom != nil && oroom.Del(ch) {
		b.DelRoom(oroom)
	}
	return
}

// Put put a channel according with sub key.
func (b *Bucket) Put(rid string, ch *Channel) (err error) {
	var (
		room *Room
		ok   bool
	)
	b.cLock.Lock()
	// close old channel
	if dch := b.chs[ch.Key]; dch != nil {
		dch.Close()
	}
	b.chs[ch.Key] = ch
	if rid != model.NoRoom && rid != "" {
		if room, ok = b.rooms[rid]; !ok {
			room = NewRoom(rid)
			b.rooms[rid] = room
		}
		ch.Room = room
		b.roomIPCnts[ch.IP]++
	}
	b.ipCnts[ch.IP]++
	b.cLock.Unlock()
	if room != nil {
		err = room.Put(ch)
	}
	return
}

// Del delete the channel by sub key.
func (b *Bucket) Del(dch *Channel) {
	var (
		ok   bool
		ch   *Channel
		room *Room
	)
	b.cLock.Lock()
	if ch, ok = b.chs[dch.Key]; ok {
		if room = ch.Room; room != nil {
			if b.roomIPCnts[ch.IP] > 1 {
				b.roomIPCnts[ch.IP]--
			} else {
				delete(b.roomIPCnts, ch.IP)
			}
		}
		if ch == dch {
			delete(b.chs, ch.Key)
		}
		// ip counter
		if b.ipCnts[ch.IP] > 1 {
			b.ipCnts[ch.IP]--
		} else {
			delete(b.ipCnts, ch.IP)
		}
	}
	b.cLock.Unlock()
	if room != nil && room.Del(ch) {
		// if empty room, must delete from bucket
		b.DelRoom(room)
	}
}

// Channel get a channel by sub key.
func (b *Bucket) Channel(key string) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[key]
	b.cLock.RUnlock()
	return
}

// Broadcast push msgs to all channels in the bucket.
func (b *Bucket) Broadcast(p *model.Proto, op int32, platform string) {
	var ch *Channel
	b.cLock.RLock()
	for _, ch = range b.chs {
		if !ch.NeedPush(op, platform) {
			continue
		}
		ch.Push(p)
	}
	b.cLock.RUnlock()
}

// Room get a room by roomid.
func (b *Bucket) Room(rid string) (room *Room) {
	b.cLock.RLock()
	room = b.rooms[rid]
	b.cLock.RUnlock()
	return
}

// DelRoom delete a room by roomid.
func (b *Bucket) DelRoom(room *Room) {
	b.cLock.Lock()
	delete(b.rooms, room.ID)
	b.cLock.Unlock()
	room.Close()
}

// BroadcastRoom broadcast a message to specified room
func (b *Bucket) BroadcastRoom(arg ProtoRoom) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.c.RoutineAmount
	b.routines[num] <- arg
}

// Rooms get all room id where online number > 0.
func (b *Bucket) Rooms() (res map[string]struct{}) {
	var (
		roomID string
		room   *Room
	)
	res = make(map[string]struct{})
	b.cLock.RLock()
	for roomID, room = range b.rooms {
		if room.Online > 0 {
			res[roomID] = struct{}{}
		}
	}
	b.cLock.RUnlock()
	return
}

// IPCount get ip count.
func (b *Bucket) IPCount() (res map[string]struct{}) {
	var (
		ip string
	)
	b.cLock.RLock()
	res = make(map[string]struct{}, len(b.ipCnts))
	for ip = range b.ipCnts {
		res[ip] = struct{}{}
	}
	b.cLock.RUnlock()
	return
}

// RoomIPCount get ip count.
func (b *Bucket) RoomIPCount() (res map[string]struct{}) {
	var (
		ip string
	)
	b.cLock.RLock()
	res = make(map[string]struct{}, len(b.roomIPCnts))
	for ip = range b.roomIPCnts {
		res[ip] = struct{}{}
	}
	b.cLock.RUnlock()
	return
}

// UpRoomsCount update all room count
func (b *Bucket) UpRoomsCount(roomCountMap map[string]int32) {
	var (
		roomID string
		room   *Room
	)
	b.cLock.RLock()
	for roomID, room = range b.rooms {
		room.AllOnline = roomCountMap[roomID]
	}
	b.cLock.RUnlock()
}

// roomproc
func (b *Bucket) roomproc(c chan ProtoRoom) {
	for {
		arg := <-c
		if room := b.Room(arg.RoomID); room != nil {
			room.Push(arg.Proto)
		}
	}
}

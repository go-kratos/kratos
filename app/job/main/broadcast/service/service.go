package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go-common/app/job/main/broadcast/conf"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

const (
	broadcastAppID = "push.interface.broadcast"
)

var (
	// ErrComet commet error.
	ErrComet = errors.New("comet rpc is not available")
	// ErrCometFull comet chan full.
	ErrCometFull = errors.New("comet proto chan full")
	// ErrRoomFull room chan full.
	ErrRoomFull = errors.New("room proto chan full")
)

// Service is a service.
type Service struct {
	conf         *conf.Config
	consumer     *databus.Databus
	cometServers map[string]*Comet
	rooms        map[string]*Room
	roomsMutex   sync.RWMutex
	options      RoomOptions
}

// New new a service and return.
func New(c *conf.Config) *Service {
	if c.Room.Refresh <= 0 {
		c.Room.Refresh = xtime.Duration(time.Second)
	}
	s := &Service{
		conf:         c,
		consumer:     databus.New(c.Databus),
		cometServers: make(map[string]*Comet),
		rooms:        make(map[string]*Room, 1024),
		roomsMutex:   sync.RWMutex{},
		options: RoomOptions{
			BatchNum:   c.Room.Batch,
			SignalTime: time.Duration(c.Room.Signal),
		},
	}
	dis := discovery.New(c.Discovery)
	s.watchComet(dis.Build(broadcastAppID))
	go s.consume()
	return s
}

func (s *Service) consume() {
	msgs := s.consumer.Messages()
	for {
		msg, ok := <-msgs
		if !ok {
			log.Warn("[job] consumer has been closed")
			return
		}
		if msg.Topic != s.conf.Databus.Topic {
			log.Error("unknown message:%v", msg)
			continue
		}
		s.pushMsg(msg.Value)
		msg.Commit()
	}
}

// Close close the resources.
func (s *Service) Close() error {
	if err := s.consumer.Close(); err != nil {
		return err
	}
	for _, c := range s.cometServers {
		if err := c.Close(); err != nil {
			log.Error("c.Close() error(%v)", err)
		}
	}
	return nil
}

func (s *Service) watchComet(resolver naming.Resolver) {
	event := resolver.Watch()
	select {
	case _, ok := <-event:
		if !ok {
			panic("watchComet init failed")
		}
		if ins, ok := resolver.Fetch(context.Background()); ok {
			if err := s.newAddress(ins); err != nil {
				panic(err)
			}
			log.Info("watchComet init newAddress:%+v", ins)
		}
	case <-time.After(10 * time.Second):
		log.Error("watchComet init instances timeout")
	}
	go func() {
		for {
			if _, ok := <-event; !ok {
				log.Info("watchComet exit")
				return
			}
			ins, ok := resolver.Fetch(context.Background())
			if ok {
				if err := s.newAddress(ins); err != nil {
					log.Error("watchComet newAddress(%+v) error(%+v)", ins, err)
					continue
				}
				log.Info("watchComet change newAddress:%+v", ins)
			}
		}
	}()
}

func (s *Service) newAddress(insMap map[string][]*naming.Instance) error {
	ins := insMap[env.Zone]
	if len(ins) == 0 {
		return fmt.Errorf("watchComet instance is empty")
	}
	comets := map[string]*Comet{}
	options := CometOptions{
		RoutineSize: s.conf.Routine.Size,
		RoutineChan: s.conf.Routine.Chan,
	}
	for _, data := range ins {
		if old, ok := s.cometServers[data.Hostname]; ok {
			comets[data.Hostname] = old
			continue
		}
		c, err := NewComet(data, s.conf, options)
		if err != nil {
			log.Error("watchComet NewComet(%+v) error(%v)", data, err)
			return err
		}
		comets[data.Hostname] = c
		log.Info("watchComet AddComet grpc:%+v", data)
	}
	for key, old := range s.cometServers {
		if _, ok := comets[key]; !ok {
			old.cancel()
			log.Info("watchComet DelComet:%s", key)
		}
	}
	s.cometServers = comets
	return nil
}

func (s *Service) room(roomID string) *Room {
	s.roomsMutex.RLock()
	room, ok := s.rooms[roomID]
	s.roomsMutex.RUnlock()
	if !ok {
		s.roomsMutex.Lock()
		if room, ok = s.rooms[roomID]; !ok {
			room = NewRoom(s, roomID, s.options)
			s.rooms[roomID] = room
		}
		s.roomsMutex.Unlock()
		log.Info("new a room:%s active:%d", roomID, len(s.rooms))
	}
	return room
}

package service

import (
	"context"
	"fmt"
	"net/url"
	"sync/atomic"
	"time"

	client "go-common/app/interface/main/broadcast/api/grpc/v1"
	"go-common/app/job/main/broadcast/conf"
	"go-common/library/log"
	"go-common/library/naming"
)

// CometOptions comet options.
type CometOptions struct {
	RoutineSize uint64
	RoutineChan uint64
}

// Comet is a broadcast comet.
type Comet struct {
	serverID        string
	broadcastClient client.ZergClient
	pushChan        []chan *client.PushMsgReq
	roomChan        []chan *client.BroadcastRoomReq
	broadcastChan   chan *client.BroadcastReq
	pushChanNum     uint64
	roomChanNum     uint64
	options         CometOptions
	ctx             context.Context
	cancel          context.CancelFunc
}

// Push push a user message.
func (c *Comet) Push(arg *client.PushMsgReq) (err error) {
	idx := atomic.AddUint64(&c.pushChanNum, 1) % c.options.RoutineSize
	c.pushChan[idx] <- arg
	return
}

// BroadcastRoom broadcast a room message.
func (c *Comet) BroadcastRoom(arg *client.BroadcastRoomReq) (err error) {
	idx := atomic.AddUint64(&c.roomChanNum, 1) % c.options.RoutineSize
	c.roomChan[idx] <- arg
	return
}

// Broadcast broadcast a message.
func (c *Comet) Broadcast(arg *client.BroadcastReq) (err error) {
	c.broadcastChan <- arg
	return
}

// process
func (c *Comet) process(pushChan chan *client.PushMsgReq, roomChan chan *client.BroadcastRoomReq, broadcastChan chan *client.BroadcastReq) {
	var err error
	for {
		select {
		case broadcastArg := <-broadcastChan:
			_, err = c.broadcastClient.Broadcast(context.Background(), &client.BroadcastReq{
				Proto:    broadcastArg.Proto,
				ProtoOp:  broadcastArg.ProtoOp,
				Speed:    broadcastArg.Speed,
				Platform: broadcastArg.Platform,
			})
			if err != nil {
				log.Error("c.broadcastClient.Broadcast(%s, %v, reply) serverId:%d error(%v)", broadcastArg, c.serverID, err)
			}
		case roomArg := <-roomChan:
			_, err = c.broadcastClient.BroadcastRoom(context.Background(), &client.BroadcastRoomReq{
				RoomID: roomArg.RoomID,
				Proto:  roomArg.Proto,
			})
			if err != nil {
				log.Error("c.broadcastClient.BroadcastRoom(%s, %v, reply) serverId:%d error(%v)", roomArg, c.serverID, err)
			}
		case pushArg := <-pushChan:
			_, err = c.broadcastClient.PushMsg(context.Background(), &client.PushMsgReq{
				Keys:    pushArg.Keys,
				Proto:   pushArg.Proto,
				ProtoOp: pushArg.ProtoOp,
			})
			if err != nil {
				log.Error("c.broadcastClient.PushMsg(%s, %v, reply) serverId:%d error(%v)", pushArg, c.serverID, err)
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// Close close the resouces.
func (c *Comet) Close() (err error) {
	finish := make(chan bool)
	go func() {
		for {
			n := len(c.broadcastChan)
			for _, ch := range c.pushChan {
				n += len(ch)
			}
			for _, ch := range c.roomChan {
				n += len(ch)
			}
			if n == 0 {
				finish <- true
				return
			}
			time.Sleep(time.Second)
		}
	}()
	select {
	case <-finish:
		log.Info("close comet finish")
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("close comet(server:%s push:%d room:%d broadcast:%d) timeout", c.serverID, len(c.pushChan), len(c.roomChan), len(c.broadcastChan))
	}
	c.cancel()
	return
}

// NewComet new a comet.
func NewComet(data *naming.Instance, conf *conf.Config, options CometOptions) (*Comet, error) {
	c := &Comet{
		serverID:      data.Hostname,
		pushChan:      make([]chan *client.PushMsgReq, options.RoutineSize),
		roomChan:      make([]chan *client.BroadcastRoomReq, options.RoutineSize),
		broadcastChan: make(chan *client.BroadcastReq, options.RoutineSize),
		options:       options,
	}
	var grpcAddr string
	for _, addrs := range data.Addrs {
		u, err := url.Parse(addrs)
		if err == nil && u.Scheme == "grpc" {
			grpcAddr = u.Host
		}
	}
	if grpcAddr == "" {
		return nil, fmt.Errorf("invalid grpc address:%v", data.Addrs)
	}
	var err error
	if c.broadcastClient, err = client.NewClient(grpcAddr, conf.RPC); err != nil {
		return nil, err
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	for i := uint64(0); i < options.RoutineSize; i++ {
		c.pushChan[i] = make(chan *client.PushMsgReq, options.RoutineChan)
		c.roomChan[i] = make(chan *client.BroadcastRoomReq, options.RoutineChan)
		go c.process(c.pushChan[i], c.roomChan[i], c.broadcastChan)
	}
	return c, nil
}

package server

import (
	"sync"

	"go-common/app/service/main/broadcast/libs/bufio"
	"go-common/app/service/main/broadcast/model"
)

// Channel used by message pusher send msg to write goroutine.
type Channel struct {
	Room     *Room
	CliProto Ring
	signal   chan *model.Proto
	Writer   bufio.Writer
	Reader   bufio.Reader
	Next     *Channel
	Prev     *Channel

	Mid      int64
	Key      string
	IP       string
	Platform string
	watchOps map[int32]struct{}
	mutex    sync.RWMutex
	V1       bool
}

// NewChannel new a channel.
func NewChannel(cli, svr int) *Channel {
	c := new(Channel)
	c.CliProto.Init(cli)
	c.signal = make(chan *model.Proto, svr)
	c.watchOps = make(map[int32]struct{})
	return c
}

// Watch watch a operation.
func (c *Channel) Watch(accepts ...int32) {
	c.mutex.Lock()
	for _, op := range accepts {
		if op >= model.MinBusinessOp && op <= model.MaxBusinessOp {
			c.watchOps[op] = struct{}{}
		}
	}
	c.mutex.Unlock()
}

// UnWatch unwatch an operation
func (c *Channel) UnWatch(accepts ...int32) {
	c.mutex.Lock()
	for _, op := range accepts {
		delete(c.watchOps, op)
	}
	c.mutex.Unlock()
}

// NeedPush verify if in watch.
func (c *Channel) NeedPush(op int32, platform string) bool {
	if c.Platform != platform && platform != "" {
		return false
	}
	if op >= 0 && op < model.MinBusinessOp {
		return true
	}
	c.mutex.RLock()
	if _, ok := c.watchOps[op]; ok {
		c.mutex.RUnlock()
		return true
	}
	c.mutex.RUnlock()
	return false
}

// Push server push message.
func (c *Channel) Push(p *model.Proto) (err error) {
	// NOTE: 兼容v1弹幕推送，等播放器接后可以去掉
	if c.V1 && p.Operation != 5 {
		return
	}
	select {
	case c.signal <- p:
	default:
	}
	return
}

// Ready check the channel ready or close?
func (c *Channel) Ready() *model.Proto {
	return <-c.signal
}

// Signal send signal to the channel, protocol ready.
func (c *Channel) Signal() {
	c.signal <- model.ProtoReady
}

// Close close the channel.
func (c *Channel) Close() {
	c.signal <- model.ProtoFinish
}

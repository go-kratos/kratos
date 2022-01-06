package v3

import (
	"github.com/go-kratos/kratos/v2/xds/resource"
	"google.golang.org/grpc"
)

func (c *Client) run() {
	go c.sendLoop()
	var stream adsStream

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		c.RecvResponse(stream)
	}
}

func (c *Client) sendLoop() {
	var stream adsStream
	var err error
	for {
		select {
		case <-c.ctx.Done():
			return
		case action := <-c.watchCh:
			err = c.WatchResource(stream, action.rType, action.resource, action.remove)
		case action := <-c.ackCh:
			if action.stream != stream {
				continue
			}
			err = c.SendAck(stream, action.rType, action.version, action.nonce, action.errMsg)
		}
		if err != nil {

		}
	}
}

func (c *Client) getStream() adsStream {
	c.lk.RLock()
	if c.stream == nil {
		c.lk.RUnlock()
		c.lk.Lock()
		if c.stream == nil {
			var err error
			c.stream, err = c.newStream(c.ctx, c.cc)
			if err != nil {

			}
		}
		c.lk.Unlock()

	}
}

type watchAction struct {
	rType    resource.ResourceType
	remove   bool // Whether this is to remove watch for the resource.
	resource string
}

type ackAction struct {
	rType   resource.ResourceType
	version string // NACK if version is an empty string.
	nonce   string
	errMsg  string // Empty unless it's a NACK.
	// ACK/NACK are tagged with the stream it's for. When the stream is down,
	// all the ACK/NACK for this stream will be dropped, and the version/nonce
	// won't be updated.
	stream grpc.ClientStream
}

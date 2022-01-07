package v3

import (
	"time"

	v3clusterpb "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	v3endpointpb "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	v3listenerpb "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	v3routepb "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	v3discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/xds/resource"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func (c *Client) run() {
	go c.sendLoop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		stream, err := c.getStream()
		if err != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		c.recv(stream)
	}
}

func (c *Client) recv(stream adsStream) error {

	for {
		r, err := c.RecvResponse(stream)
		if err != nil {
			c.resetStream(c.stream)
			return err
		}
		rType := resource.UnknownResource
		resp, ok := r.(*v3discoverypb.DiscoveryResponse)
		if !ok {
			// send nack
			c.ackCh <- &ackAction{
				rType:   rType,
				version: "",
				nonce:   resp.Nonce,
				errMsg:  err.Error(),
				stream:  stream,
			}
			continue
		}

		// Note that the xDS transport protocol is versioned independently of
		// the resource types, and it is supported to transfer older versions
		// of resource types using new versions of the transport protocol, or
		// vice-versa. Hence we need to handle v3 type_urls as well here.
		url := resp.GetTypeUrl()
		switch {
		case resource.IsListenerResource(url):
			rType = resource.ListenerResource
			for _, res := range resp.GetResources() {
				lis := &v3listenerpb.Listener{}
				err = proto.Unmarshal(res.GetValue(), lis)
				if err != nil {
					break
				}
			}
		case resource.IsRouteConfigResource(url):
			rType = resource.RouteConfigResource
			for _, res := range resp.GetResources() {
				rc := &v3routepb.RouteConfiguration{}
				err = proto.Unmarshal(res.GetValue(), rc)
				if err != nil {
					break
				}
			}
		case resource.IsClusterResource(url):
			rType = resource.ClusterResource
			for _, res := range resp.GetResources() {
				cluster := &v3clusterpb.Cluster{}
				err = proto.Unmarshal(res.GetValue(), cluster)
				if err != nil {
					break
				}
			}
		case resource.IsEndpointsResource(url):
			rType = resource.EndpointsResource
			for _, res := range resp.GetResources() {
				cla := &v3endpointpb.ClusterLoadAssignment{}
				err = proto.Unmarshal(res.GetValue(), cla)
				if err != nil {
					break
				}
			}
		default:
			// Unknown resource type
			continue
		}
		if err != nil {
			// send nack
			// send nack
			c.ackCh <- &ackAction{
				rType:   rType,
				version: "",
				nonce:   resp.Nonce,
				errMsg:  err.Error(),
				stream:  stream,
			}
		} else {
			c.ackCh <- &ackAction{
				rType:   rType,
				version: resp.GetVersionInfo(),
				nonce:   resp.GetNonce(),
				stream:  stream,
			}
		}

	}

}

func (c *Client) sendLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case action := <-c.watchCh:
			stream, err := c.getStream()
			if err != nil {
				c.logger.Errorf("xds client get stream failed!err:=%v", err)
				continue
			}
			err = c.WatchResource(stream, action.rType, action.resource, action.remove)
			if err != nil {
				c.logger.Errorf("xds client watch resource failed!err:=%v", err)
				if errors.IsServiceUnavailable(err) {
					c.resetStream(stream)
				}
			}
		case action := <-c.ackCh:
			stream, err := c.getStream()
			if err != nil {
				c.logger.Errorf("xds client get stream failed!err:=%v", err)
				continue
			}
			if action.stream != stream {
				continue
			}
			err = c.SendAck(stream, action.rType, action.version, action.nonce, action.errMsg)
			if err != nil {
				c.logger.Errorf("xds client send ack failed!err:=%v", err)
				if errors.IsServiceUnavailable(err) {
					c.resetStream(stream)
				}
			}
		}
	}
}

func (c *Client) resetStream(origin adsStream) {
	c.lk.Lock()
	if origin == c.stream {
		c.stream = nil
	}
	c.lk.Unlock()
}

func (c *Client) getStream() (adsStream, error) {
	c.lk.RLock()
	if c.stream == nil {
		c.lk.RUnlock()

		c.lk.Lock()
		defer c.lk.Unlock()
		if c.stream == nil {
			var err error
			for i := 0; i < 3; i++ {
				c.stream, err = c.newStream(c.ctx, c.cc)
				if err == nil {
					return c.stream, nil
				}
				if i < 2 {
					time.Sleep(time.Millisecond * 250 * time.Duration(i+1))
				}
			}
			return nil, err
		}
		return c.stream, nil
	}
	defer c.lk.RUnlock()
	return c.stream, nil

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

/*
 *
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package v3 provides xDS v3 transport protocol specific functionality.
package v3

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/xds/resource"

	v3corepb "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v3adsgrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	v3discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

const (
	reasonSendError = "send_error"
)

// BuildOptions contains options to be passed to client builders.
type BuildOptions struct {
	// NodeProto contains the Node proto to be used in xDS requests. The actual
	// type depends on the transport protocol version used.
	NodeProto proto.Message
	// // Backoff returns the amount of time to backoff before retrying broken
	// // streams.
	// Backoff func(int) time.Duration
	// Logger provides enhanced logging capabilities.
	Logger *log.Helper
}

var (
	resourceTypeToURL = map[resource.ResourceType]string{
		resource.ListenerResource:    resource.V3ListenerURL,
		resource.RouteConfigResource: resource.V3RouteConfigURL,
		resource.ClusterResource:     resource.V3ClusterURL,
		resource.EndpointsResource:   resource.V3EndpointsURL,
	}
)

// NewClient new xds client
func NewClient(cc *grpc.ClientConn, opts BuildOptions) (*Client, error) {
	nodeProto, ok := opts.NodeProto.(*v3corepb.Node)
	if !ok {
		return nil, fmt.Errorf("xds: unsupported Node proto type: %T, want %T", opts.NodeProto, v3corepb.Node{})
	}
	c := &Client{
		nodeProto: nodeProto,
		logger:    opts.Logger,
		watchMp:   make(map[resource.ResourceType]map[string]bool),
		versionMp: make(map[resource.ResourceType]string),
		nonceMp:   make(map[resource.ResourceType]string),
		cc:        cc,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c, nil
}

type adsStream v3adsgrpc.AggregatedDiscoveryService_StreamAggregatedResourcesClient

// client performs the actual xDS RPCs using the xDS v3 API. It creates a
// single ADS stream on which the different types of xDS requests and responses
// are multiplexed.
type Client struct {
	nodeProto *v3corepb.Node
	logger    *log.Helper
	cc        *grpc.ClientConn // Connection to the management server.

	watchMp   map[resource.ResourceType]map[string]bool
	versionMp map[resource.ResourceType]string
	nonceMp   map[resource.ResourceType]string
	ackCh     chan *ackAction
	watchCh   chan *watchAction
	ctx       context.Context
	cancel    context.CancelFunc
	lk        sync.RWMutex
	stream    adsStream
}

func (c *Client) newStream(ctx context.Context, cc *grpc.ClientConn) (adsStream, error) {
	return v3adsgrpc.NewAggregatedDiscoveryServiceClient(cc).StreamAggregatedResources(ctx, grpc.WaitForReady(true))
}

func (c *Client) SendAck(stream adsStream, rType resource.ResourceType, version, nonce, errMsg string) error {
	c.lk.Lock()
	c.nonceMp[rType] = nonce
	s, ok := c.watchMp[rType]
	if !ok || len(s) == 0 {
		c.lk.Unlock()
		// We don't send the request ack if there's no active watch (this can be
		// either the server sends responses before any request, or the watch is
		// canceled while the ackAction is in queue), because there's no resource
		// name. And if we send a request with empty resource name list, the
		// server may treat it as a wild card and send us everything.
		return errors.NotFound(resource.UnknownResource.String(), rType.String())
	}
	target := mapToSlice(s)
	if version == "" {
		// This is a nack, get the previous acked version.
		version = c.versionMp[rType]
		// version will still be an empty string if rType isn't
		// found in versionMap, this can happen if there wasn't any ack
		// before.
	} else {
		c.versionMp[rType] = version
	}
	c.lk.Unlock()

	return c.sendRequest(stream, target, rType, version, nonce, errMsg)
}

func (c *Client) WatchResource(s adsStream, rType resource.ResourceType, resource string, remove bool) error {
	c.lk.Lock()

	var current map[string]bool
	current, ok := c.watchMp[rType]
	if !ok {
		current = make(map[string]bool)
		c.watchMp[rType] = current
	}
	if remove {
		delete(current, resource)
		if len(current) == 0 {
			delete(c.watchMp, rType)
		}
	} else {
		current[resource] = true
	}
	target := mapToSlice(current)
	// We don't reset version or nonce when a new watch is started. The version
	// and nonce from previous response are carried by the request unless the
	// stream is recreated.
	ver := c.versionMp[rType]
	nonce := c.nonceMp[rType]
	c.lk.Unlock()
	return c.sendRequest(s, target, rType, ver, nonce, "")
}

// sendRequest sends out a DiscoveryRequest for the given resourceNames, of type
// rType, on the provided stream.
//
// version is the ack version to be sent with the request
// - If this is the new request (not an ack/nack), version will be empty.
// - If this is an ack, version will be the version from the response.
// - If this is a nack, version will be the previous acked version (from
//   versionMap). If there was no ack before, it will be empty.
func (c *Client) sendRequest(stream adsStream, resourceNames []string, rType resource.ResourceType, version, nonce, errMsg string) error {
	req := &v3discoverypb.DiscoveryRequest{
		Node:          c.nodeProto,
		TypeUrl:       resourceTypeToURL[rType],
		ResourceNames: resourceNames,
		VersionInfo:   version,
		ResponseNonce: nonce,
	}
	if errMsg != "" {
		req.ErrorDetail = &statuspb.Status{
			Code: int32(codes.InvalidArgument), Message: errMsg,
		}
	}
	if err := stream.Send(req); err != nil {
		return errors.ServiceUnavailable(reasonSendError, fmt.Sprintf("xds: stream.Send(%+v) failed: %v", req, err))
	}
	c.logger.Debugf("ADS request sent: %v", (req))
	return nil
}

// RecvResponse blocks on the receipt of one response message on the provided
// stream.
func (c *Client) RecvResponse(stream adsStream) (proto.Message, error) {
	resp, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("xds: stream.Recv() failed: %v", err)
	}
	c.logger.Infof("ADS response received, type: %v", resp.GetTypeUrl())
	c.logger.Debugf("ADS response received: %+v", (resp))
	return resp, nil
}

func mapToSlice(m map[string]bool) []string {
	ret := make([]string, 0, len(m))
	for i := range m {
		ret = append(ret, i)
	}
	return ret
}

func (c *Client) Close() error {
	c.cancel()
	return c.cc.Close()
}

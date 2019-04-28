package hrpc

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/tsuna/gohbase/pb"
)

// ClusterStatus to represent a cluster status request
type ClusterStatus struct {
	base
}

// NewClusterStatus creates a new ClusterStatusStruct with default fields
func NewClusterStatus() *ClusterStatus {
	return &ClusterStatus{
		base{
			ctx:      context.Background(),
			table:    []byte{},
			resultch: make(chan RPCResult, 1),
		},
	}
}

// Name returns the name of the rpc function
func (c *ClusterStatus) Name() string {
	return "GetClusterStatus"
}

// ToProto returns the Protobuf message to be sent
func (c *ClusterStatus) ToProto() proto.Message {
	return &pb.GetClusterStatusRequest{}
}

// NewResponse returns the empty protobuf response
func (c *ClusterStatus) NewResponse() proto.Message {
	return &pb.GetClusterStatusResponse{}
}

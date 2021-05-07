package apimeta

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/api/kratos/api"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// Server is api meta server
type Server struct {
	api.UnimplementedMetadataServer
	s *grpc.Server
}

// NewServer create server instance
func NewServer(grpcSrv *grpc.Server) *Server {
	return &Server{s: grpcSrv}
}

// ListServices return all services
func (s *Server) ListServices(ctx context.Context, in *anypb.Any) (*api.ListServicesReply, error) {
	var reply api.ListServicesReply
	serviceInfo := s.s.GetServiceInfo()
	for svc := range serviceInfo {
		reply.Services = append(reply.Services, svc)
	}
	return &reply, nil
}

// GetService return service meta by name
func (s *Server) GetService(ctx context.Context, in *api.GetServiceRequest) (*api.GetServiceReply, error) {
	reply := api.GetServiceReply{ProtoSet: &dpb.FileDescriptorSet{}}
	serviceInfo := s.s.GetServiceInfo()
	if info, ok := serviceInfo[in.Name]; ok {
		fdenc, ok := parseMetadata(info.Metadata)
		if !ok {
			return nil, fmt.Errorf("invalid service %s meta", in.Name)
		}
		fd, err := decodeFileDesc(fdenc)
		if err != nil {
			return nil, err
		}
		reply.ProtoSet.File, err = allDependency(fd)
		if err != nil {
			return nil, err
		}
	}
	return &reply, nil
}

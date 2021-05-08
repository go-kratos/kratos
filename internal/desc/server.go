package desc

import (
	"context"

	"github.com/go-kratos/kratos/v2/api/kratos/api"
	"google.golang.org/grpc"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// Server is api meta server
type Server struct {
	api.UnimplementedDescriptionServer
	s *Service
}

// NewServer create server instance
func NewServer(grpcSrv ...*grpc.Server) *Server {
	return &Server{s: NewService(grpcSrv...)}
}

// ListServices return all services
func (s *Server) ListServices(ctx context.Context, in *anypb.Any) (*api.ListServicesReply, error) {
	var reply api.ListServicesReply
	var err error
	reply.Services, err = s.s.ListServices(ctx)
	return &reply, err
}

// GetServiceDesc return service meta by name
func (s *Server) GetServiceDesc(ctx context.Context, in *api.GetServiceDescRequest) (*api.GetServiceDescReply, error) {
	var reply api.GetServiceDescReply
	var err error
	reply.ProtoSet, err = s.s.GetServiceDesc(ctx, in.Name)
	return &reply, err
}

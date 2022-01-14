package health

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/health"
	pb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	pb.UnimplementedHealthServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	info, ok := kratos.FromContext(ctx)
	if !ok {
		return nil, errors.InternalServer("kratos.FromContext(ctx) failed", "no info found in context")
	}
	status, ok := info.Health().GetStatus(req.Service)
	if !ok {
		status = health.Status_UNKNOWN
	}
	var sv pb.HealthCheckResponse_ServingStatus
	switch status {
	case health.Status_SERVING:
		sv = pb.HealthCheckResponse_SERVING
	case health.Status_NOT_SERVING:
		sv = pb.HealthCheckResponse_NOT_SERVING
	case health.Status_SERVICE_UNKNOWN:
		sv = pb.HealthCheckResponse_SERVICE_UNKNOWN
	default:
		sv = pb.HealthCheckResponse_NOT_SERVING
	}
	return &pb.HealthCheckResponse{
		Status: sv,
	}, nil
}

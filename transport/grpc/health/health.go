package health

//go:generate protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. health.proto

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/health"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthCheckServer struct {
	UnimplementedHealthServer
}

func (s *HealthCheckServer) Check(ctx context.Context, req *HealthCheckRequest) (resp *HealthCheckResponse, err error) {
	info, _ := kratos.FromContext(ctx)
	status := info.Health().GetStatus()
	var v HealthCheckResponse_ServingStatus
	switch status {
	case health.Status_UNKNOWN:
		v = HealthCheckResponse_UNKNOWN
	case health.Status_SERVING:
		v = HealthCheckResponse_SERVING
	case health.Status_NOT_SERVING:
		v = HealthCheckResponse_NOT_SERVING
	}
	resp = &HealthCheckResponse{
		Status: v,
	}
	return
}
func (s *HealthCheckServer) Watch(req *HealthCheckRequest, server Health_WatchServer) (err error) {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

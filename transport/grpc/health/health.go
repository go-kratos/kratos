package health

//go:generate protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. health.proto

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/health"
)

type HealthCheckServer struct {
	UnimplementedHealthServer
}

func NewHealthCheckServer() *HealthCheckServer {
	return &HealthCheckServer{}
}

func (s *HealthCheckServer) Check(ctx context.Context, req *HealthCheckRequest) (resp *HealthCheckResponse, err error) {
	info, _ := kratos.FromContext(ctx)
	status, ok := info.Health().GetStatus(req.Service)
	var sv HealthCheckResponse_ServingStatus
	if !ok {
		sv = HealthCheckResponse_SERVICE_UNKNOWN
	}
	switch status {
	case health.Status_SERVING:
		sv = HealthCheckResponse_SERVING
	case health.Status_NOT_SERVING:
		sv = HealthCheckResponse_NOT_SERVING
	case health.Status_SERVICE_UNKNOWN:
		sv = HealthCheckResponse_SERVICE_UNKNOWN
	default:
		sv = HealthCheckResponse_NOT_SERVING
	}
	resp = &HealthCheckResponse{
		Status: sv,
	}
	return
}
func (s *HealthCheckServer) Watch(req *HealthCheckRequest, ss Health_WatchServer) (err error) {
	ctx := ss.Context()
	info, ok := kratos.FromContext(ctx)
	if !ok {
		return
	}
	info.Health().Watch(req.Service, func() {
		status, ok := info.Health().GetStatus(req.Service)
		var sv HealthCheckResponse_ServingStatus
		if !ok {
			sv = HealthCheckResponse_SERVICE_UNKNOWN
		}
		switch status {
		case health.Status_SERVING:
			sv = HealthCheckResponse_SERVING
		case health.Status_NOT_SERVING:
			sv = HealthCheckResponse_NOT_SERVING
		case health.Status_SERVICE_UNKNOWN:
			sv = HealthCheckResponse_SERVICE_UNKNOWN
		default:
			sv = HealthCheckResponse_NOT_SERVING
		}
		ss.Send(&HealthCheckResponse{Status: sv})
	})
	return
}

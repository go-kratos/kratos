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
func (s *HealthCheckServer) Watch(req *HealthCheckRequest, server Health_WatchServer) (err error) {
	info, _ := kratos.FromContext(ctx)
	info.Health().Watch(req.Service, func(status health.Status) {
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
		server.Send(&HealthCheckResponse{Status: sv})
	})

	return
}

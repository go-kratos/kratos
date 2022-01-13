package health

//go:generate protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. health.proto

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/health"
	"github.com/google/uuid"
)

type GrpcHealthCheckServer struct {
	UnimplementedHealthServer
}

func NewHealthCheckServer() *GrpcHealthCheckServer {
	return &GrpcHealthCheckServer{}
}

func (s *GrpcHealthCheckServer) Check(ctx context.Context, req *HealthCheckRequest) (resp *HealthCheckResponse, err error) {
	info, _ := kratos.FromContext(ctx)
	status, _ := info.Health().GetStatus(req.Service)
	var sv HealthCheckResponse_ServingStatus
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

func (s *GrpcHealthCheckServer) Watch(req *HealthCheckRequest, ss Health_WatchServer) (err error) {
	ctx := ss.Context()
	info, ok := kratos.FromContext(ctx)
	if !ok {
		return errors.InternalServer("get info failed", "")
	}
	uid, err := uuid.NewUUID()
	if err != nil {
		return errors.InternalServer("new uuid failed", err.Error())
	}
	update := info.Health().Watch(req.Service, uid.String())
	status, ok := info.Health().GetStatus(req.Service)
	if !ok {
		update <- health.Status_SERVICE_UNKNOWN
	} else {
		update <- status
	}

	var lastStatus health.Status = -1
	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-update:
			if lastStatus == status {
				continue
			}
			lastStatus = status
			var sv HealthCheckResponse_ServingStatus
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
			resp := &HealthCheckResponse{
				Status: sv,
			}
			if err := ss.Send(resp); err != nil {
				return err
			}
		}
	}
}

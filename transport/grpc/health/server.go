package health

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/health"
	"github.com/google/uuid"
	pb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	pb.UnimplementedHealthServer
}

func NewServerr() *Server {
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

func (s *Server) Watch(req *pb.HealthCheckRequest, ss pb.Health_WatchServer) (err error) {
	ctx := ss.Context()
	info, ok := kratos.FromContext(ctx)
	if !ok {
		return errors.InternalServer("get info failed", "")
	}
	uid, err := uuid.NewUUID()
	if err != nil {
		return errors.InternalServer("new uuid failed", err.Error())
	}
	update := info.Health().Update(req.Service, uid.String())
	defer info.Health().DelUpdate(req.Service, uid.String())
	status, ok := info.Health().GetStatus(req.Service)
	if !ok {
		update <- health.Status_SERVICE_UNKNOWN
	} else {
		update <- status
	}
	ticker := info.Health().Ticker()
	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-update:
			err := send(ss, status)
			if err != nil {
				return err
			}
		case <-ticker.C:
			err := send(ss, status)
			if err != nil {
				return err
			}
		}
	}
}

func send(ss pb.Health_WatchServer, status health.Status) error {
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
	reply := &pb.HealthCheckResponse{
		Status: sv,
	}
	return ss.Send(reply)
}

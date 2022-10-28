package grpc

import "google.golang.org/grpc/health/grpc_health_v1"

type HealthCheck interface {
	grpc_health_v1.HealthServer
	Resume()
	Shutdown()
}

package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"time"
	v1 "traces/api/blog/v1"

	pb "traces/api/user/v1"
)

type UserService struct {
	pb.UnimplementedUserServer
	tracer trace.TracerProvider
}

func NewUserService(tracer trace.TracerProvider) *UserService {
	return &UserService{tracer: tracer}
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	var res pb.GetUserReply
	var prop propagation.TraceContext
	res.Name = strconv.FormatInt(req.Id, 10)
	conn, err := grpc.DialInsecure(ctx,
		grpc.WithEndpoint("127.0.0.1:9012"),
		grpc.WithMiddleware(middleware.Chain(
			tracing.Client(tracing.WithTracerProvider(s.tracer),	tracing.WithPropagators(prop),
			),
			recovery.Recovery())),
		grpc.WithTimeout(2*time.Second),
	)
	if err != nil {
		return nil, err
	}

	blogC := v1.NewBlogClient(conn)
	reply, err := blogC.GetBlog(ctx, &v1.GetBlogRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	res.Name = reply.Blog
	return &res, nil
}

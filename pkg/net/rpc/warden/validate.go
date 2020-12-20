package warden

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

// Validate return a client interceptor validate incoming request per RPC call.
func (s *Server) validate() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if err = validate.Struct(req); err != nil {
			err = status.Error(codes.InvalidArgument, err.Error())
			return
		}
		resp, err = handler(ctx, req)
		return
	}
}

// RegisterValidation adds a validation Func to a Validate's map of validators denoted by the key
// NOTE: if the key already exists, the previous validation function will be replaced.
// NOTE: this method is not thread-safe it is intended that these all be registered prior to any validation
func (s *Server) RegisterValidation(key string, fn validator.Func) error {
	return validate.RegisterValidation(key, fn)
}

// GetValidate return the default validate
func (s *Server) GetValidate() *validator.Validate {
	return validate
}

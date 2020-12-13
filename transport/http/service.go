package http

import (
	"context"
	"net/http"
)

// ServiceRegistrar wraps a single method that supports service registration.
type ServiceRegistrar interface {
	RegisterService(desc *ServiceDesc, impl interface{})
}

// ServiceDesc represents a HTTP service's specification.
type ServiceDesc struct {
	ServiceName string
	HandlerType interface{}
	Methods     []MethodDesc
	Metadata    interface{}
}

type methodHandler func(srv interface{}, ctx context.Context, c Codec) ([]byte, error)

// MethodDesc represents a HTTP service's method specification.
type MethodDesc struct {
	Path         string
	Method       string
	Body         string
	ResponseBody string
	Handler      methodHandler
}

// RegisterService registers a service and its implementation to the HTTP server.
func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) {
	for _, method := range sd.Methods {
		m := method
		s.router.HandleFunc(m.Path, func(res http.ResponseWriter, req *http.Request) {

			ctx := req.Context()
			codec, err := codecForReq(req)
			if err != nil {
				s.encodeError(ctx, err, codec, res)
				return
			}
			reply, err := m.Handler(ss, ctx, codec)
			if err != nil {
				s.encodeError(ctx, err, codec, res)
				return
			}

			s.encodeResponse(ctx, reply, codec, res)

		}).Methods(m.Method)
	}
}

func (s *Server) encodeError(ctx context.Context, err error, codec Codec, res http.ResponseWriter) {
	s.opts.ErrorHandler(ctx, err, codec, res)
}

func (s *Server) encodeResponse(ctx context.Context, out interface{}, codec Codec, res http.ResponseWriter) {
	body, err := codec.Marshal(out)
	if err != nil {
		s.encodeError(ctx, ErrCodecMarshal(err.Error()), codec, res)
		return
	}
	res.Write(body)
}

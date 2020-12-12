package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
)

const baseContentType = "application"

func contentSubtype(contentType string) string {
	if contentType == baseContentType {
		return ""
	}
	if !strings.HasPrefix(contentType, baseContentType) {
		return ""
	}
	// guaranteed since != baseContentType and has baseContentType prefix
	switch contentType[len(baseContentType)] {
	case '/', ';':
		// this will return true for "application/grpc+" or "application/grpc;"
		// which the previous validContentType function tested to be valid, so we
		// just say that no content-subtype is specified in this case
		return contentType[len(baseContentType)+1:]
	default:
		return ""
	}
}

func codec(req *http.Request) (encoding.Codec, error) {
	contentType := req.Header.Get("content-type")
	codec := encoding.GetCodec(contentSubtype(contentType))
	if codec == nil {
		return nil, ErrUnknownCodec(contentType)
	}
	return codec, nil
}

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

type methodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error)

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
			codec, err := codec(req)
			if err != nil {
				s.encodeError(ctx, err, codec, res)
				return
			}
			// TODO Middleware
			reply, err := m.Handler(ss, ctx, func(v interface{}) error {
				return s.decodeRequest(ctx, v, codec, req)
			})
			if err != nil {
				s.encodeError(ctx, err, codec, res)
				return
			}

			s.encodeResponse(ctx, reply, codec, res)

		}).Methods(m.Method)
	}
}

func (s *Server) encodeError(ctx context.Context, err error, codec encoding.Codec, res http.ResponseWriter) {
	s.opts.ErrorHandler(ctx, err, codec, res)
}

func (s *Server) encodeResponse(ctx context.Context, out interface{}, codec encoding.Codec, res http.ResponseWriter) {
	body, err := codec.Marshal(out)
	if err != nil {
		s.encodeError(ctx, ErrCodecMarshal(err.Error()), codec, res)
		return
	}
	res.Write(body)
}

func (s *Server) decodeRequest(ctx context.Context, in interface{}, codec encoding.Codec, req *http.Request) error {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return ErrDataLoss(err.Error())
	}
	defer req.Body.Close()
	if err = codec.Unmarshal(data, in); err != nil {
		return ErrCodecUnmarshal(err.Error())
	}
	return nil
}

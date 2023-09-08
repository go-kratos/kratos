package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/redact"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// Server is an server logging middleware.
func Server(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
				method    string
				url       string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				if httpTransporter, ok := info.(http.Transporter); ok {
					httpRequest := httpTransporter.Request()
					if httpRequest != nil {
						method = httpRequest.Method
						url = httpRequest.URL.String()
					}
				}
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			_ = log.WithContext(ctx, logger).Log(level,
				"kind", "server",
				"method", method,
				"component", kind,
				"operation", operation,
				"url", url,
				"args", extractArgs(req, method),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Seconds(),
			)
			return
		}
	}
}

// Client is a client logging middleware.
func Client(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
				method    string
				url       string
			)
			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				if httpTransporter, ok := info.(http.Transporter); ok {
					httpRequest := httpTransporter.Request()
					if httpRequest != nil {
						method = httpRequest.Method
						url = httpRequest.URL.String()
					}
				}
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			_ = log.WithContext(ctx, logger).Log(level,
				"kind", "client",
				"method", method,
				"component", kind,
				"operation", operation,
				"url", url,
				"args", extractArgs(req, method),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Seconds(),
			)
			return
		}
	}
}

var redactMethods = []string{"POST", "PUT", "PATCH", ""}

// extractArgs returns the string of the req
func extractArgs(req interface{}, method string) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}

	if protoMessage, ok := req.(proto.Message); ok {
		if slices.Contains(redactMethods, method) {
			Redact(protoMessage)
			return protoMessage.(fmt.Stringer).String()
		}
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}

func Redact(pb proto.Message) {
	redactFields(pb.ProtoReflect())
}

func redactFields(m protoreflect.Message) {
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		opts := fd.Options().(*descriptorpb.FieldOptions)
		if proto.GetExtension(opts, redact.E_Sensitive).(bool) {
			switch fd.Kind() {
			case protoreflect.StringKind:
				m.Set(fd, protoreflect.ValueOfString("***"))
			case protoreflect.BytesKind:
				m.Set(fd, protoreflect.ValueOfBytes([]byte("***")))
			default:
				m.Clear(fd) // Redact the field by clearing its value
			}
		} else if fd.Kind() == protoreflect.MessageKind && v.IsValid() {
			redactValue(v)
		}
		return true
	})
}

func redactValue(v protoreflect.Value) {
	switch v.Interface().(type) {
	case protoreflect.Message:
		redactFields(v.Message())
	case protoreflect.List:
		redactList(v)
	case protoreflect.Map:
		redactMap(v)
	}
}

func redactList(v protoreflect.Value) {
	list := v.List()
	for i := 0; i < list.Len(); i++ {
		redactValue(list.Get(i))
	}
}

func redactMap(v protoreflect.Value) {
	mapVal := v.Map()
	mapVal.Range(func(key protoreflect.MapKey, val protoreflect.Value) bool {
		redactValue(val)
		return true
	})
}

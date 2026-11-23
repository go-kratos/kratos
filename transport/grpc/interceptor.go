package grpc

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	grpcmd "google.golang.org/grpc/metadata"

	ic "github.com/go-kratos/kratos/v3/internal/context"
	"github.com/go-kratos/kratos/v3/internal/matcher"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"
)

// unaryServerInterceptor is a gRPC unary server interceptor
func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, cancel := ic.Merge(ctx, s.baseCtx)
		defer cancel()
		md, _ := grpcmd.FromIncomingContext(ctx)
		replyHeader := grpcmd.MD{}
		tr := &Transport{
			operation:   info.FullMethod,
			reqHeader:   headerCarrier(md),
			replyHeader: headerCarrier(replyHeader),
		}
		if s.endpoint != nil {
			tr.endpoint = s.endpoint.String()
		}
		ctx = transport.NewServerContext(ctx, tr)
		if s.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req any) (any, error) {
			return handler(ctx, req)
		}
		// The built-in health service is probed without user credentials.
		if !strings.HasPrefix(info.FullMethod, builtinHealthPrefix) {
			if next := s.middleware.Match(tr.Operation()); len(next) > 0 {
				h = middleware.Chain(next...)(h)
			}
		}
		reply, err := h(ctx, req)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return reply, err
	}
}

// builtinHealthPrefix is the standard gRPC health service. The Kratos client
// watches it with no user credentials, so user stream middleware must not see it.
const builtinHealthPrefix = "/grpc.health.v1.Health/"

// wrappedStream is rewrite grpc stream's context
type wrappedStream struct {
	grpc.ServerStream
	ctx        context.Context
	middleware matcher.Matcher
	// opened is set when the stream interceptor already ran the middleware once.
	// SendMsg and RecvMsg then pass through, so the chain is not run again per message.
	opened bool
}

func NewWrappedStream(ctx context.Context, stream grpc.ServerStream, m matcher.Matcher) grpc.ServerStream {
	return &wrappedStream{
		ServerStream: stream,
		ctx:          ctx,
		middleware:   m,
	}
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// streamServerInterceptor is a gRPC stream server interceptor
func (s *Server) streamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, cancel := ic.Merge(ss.Context(), s.baseCtx)
		defer cancel()
		md, _ := grpcmd.FromIncomingContext(ctx)
		replyHeader := grpcmd.MD{}
		tr := &Transport{
			operation:   info.FullMethod,
			reqHeader:   headerCarrier(md),
			replyHeader: headerCarrier(replyHeader),
		}
		if s.endpoint != nil {
			tr.endpoint = s.endpoint.String()
		}
		ctx = transport.NewServerContext(ctx, tr)
		ctx = context.WithValue(ctx, stream{
			ServerStream:     ss,
			streamMiddleware: s.streamMiddleware,
		}, ss)

		// Run the chain once, as unaryServerInterceptor does, and let the handler
		// read the derived context from stream.Context(). The built chain used to
		// be discarded, and SendMsg/RecvMsg ran it again per message without
		// keeping that context.
		h := func(callCtx context.Context, _ any) (any, error) {
			ws := &wrappedStream{
				ServerStream: ss,
				ctx:          callCtx,
				middleware:   s.streamMiddleware,
				opened:       true,
			}
			return nil, handler(srv, ws)
		}
		if !strings.HasPrefix(info.FullMethod, builtinHealthPrefix) {
			if next := s.streamMiddleware.Match(info.FullMethod); len(next) > 0 {
				h = middleware.Chain(next...)(h)
			}
		}
		_, err := h(ctx, nil)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return err
	}
}

type stream struct {
	grpc.ServerStream
	streamMiddleware matcher.Matcher
}

func GetStream(ctx context.Context) grpc.ServerStream {
	return ctx.Value(stream{}).(grpc.ServerStream)
}

func (w *wrappedStream) SendMsg(m any) error {
	if w.opened {
		return w.ServerStream.SendMsg(m)
	}
	h := func(ctx context.Context, req any) (any, error) {
		w.ctx = ctx
		return req, w.ServerStream.SendMsg(m)
	}

	info, ok := transport.FromServerContext(w.ctx)
	if !ok {
		return fmt.Errorf("transport value stored in ctx returns: %v", ok)
	}

	if next := w.middleware.Match(info.Operation()); len(next) > 0 {
		h = middleware.Chain(next...)(h)
	}

	_, err := h(w.ctx, m)
	return err
}

func (w *wrappedStream) RecvMsg(m any) error {
	if w.opened {
		return w.ServerStream.RecvMsg(m)
	}
	h := func(ctx context.Context, req any) (any, error) {
		w.ctx = ctx
		return req, w.ServerStream.RecvMsg(m)
	}

	info, ok := transport.FromServerContext(w.ctx)
	if !ok {
		return fmt.Errorf("transport value stored in ctx returns: %v", ok)
	}

	if next := w.middleware.Match(info.Operation()); len(next) > 0 {
		h = middleware.Chain(next...)(h)
	}

	_, err := h(w.ctx, m)
	return err
}

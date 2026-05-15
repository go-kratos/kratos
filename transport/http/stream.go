package http

import (
	"bufio"
	"bytes"
	"context"
	stderrors "errors"
	"fmt"
	"io"
	stdhttp "net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-kratos/kratos/v3/encoding"
	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/internal/httputil"
	"github.com/go-kratos/kratos/v3/selector"
	"github.com/go-kratos/kratos/v3/transport"
)

const (
	sseContentType = "text/event-stream"

	websocketControlPrefix = "\x1e"
	websocketControlEnd    = websocketControlPrefix + "end"
	websocketControlError  = websocketControlPrefix + "error:"
)

type streamMode int

const (
	streamModeSSE streamMode = iota + 1
	streamModeWebSocket
)

// ServerStream adapts HTTP streaming transports to grpc generated stream interfaces.
type ServerStream interface {
	grpc.ServerStream
	Send(any) error
	Recv(any) error
	SendAndClose(any) error
	Close(error) error
	SetContext(context.Context)
}

// ClientStream adapts HTTP streaming clients to grpc generated stream interfaces.
type ClientStream interface {
	grpc.ClientStream
	Send(any) error
	Recv(any) error
	CloseAndRecv(any) error
}

type serverStream struct {
	ctx      context.Context
	req      *stdhttp.Request
	res      stdhttp.ResponseWriter
	mode     streamMode
	conn     *websocket.Conn
	header   metadata.MD
	trailer  metadata.MD
	encoder  encoding.Codec
	decoder  encoding.Codec
	started  bool
	writeMu  sync.Mutex
	upgrader websocket.Upgrader
}

// NewServerSentEventServerStream returns a stream that writes server messages as SSE events.
func NewServerSentEventServerStream(ctx Context) ServerStream {
	s := &serverStream{
		ctx:  ctx,
		req:  ctx.Request(),
		res:  ctx.Response(),
		mode: streamModeSSE,
	}
	s.encoder = streamCodecFromHeaders(s.req.Header, "Accept", "Content-Type")
	s.decoder = streamCodecFromHeaders(s.req.Header, "Content-Type", "Accept")
	return s
}

// NewWebSocketServerStream upgrades the current request and returns a WebSocket stream.
func NewWebSocketServerStream(ctx Context) (ServerStream, error) {
	s := &serverStream{
		ctx:  ctx,
		req:  ctx.Request(),
		res:  ctx.Response(),
		mode: streamModeWebSocket,
	}
	s.encoder = streamCodecFromHeaders(s.req.Header, "Accept", "Content-Type")
	s.decoder = streamCodecFromHeaders(s.req.Header, "Content-Type", "Accept")
	conn, err := s.upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return nil, err
	}
	s.conn = conn
	return s, nil
}

func (s *serverStream) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *serverStream) SetHeader(md metadata.MD) error {
	s.header = metadata.Join(s.header, md)
	if s.mode == streamModeSSE && !s.started {
		copyMetadataToHeader(s.res.Header(), md)
	}
	return nil
}

func (s *serverStream) SendHeader(md metadata.MD) error {
	if err := s.SetHeader(md); err != nil {
		return err
	}
	if s.mode == streamModeSSE {
		s.startSSE()
	}
	return nil
}

func (s *serverStream) SetTrailer(md metadata.MD) {
	s.trailer = metadata.Join(s.trailer, md)
}

func (s *serverStream) Context() context.Context {
	if s.ctx == nil {
		return context.Background()
	}
	return s.ctx
}

func (s *serverStream) Send(m any) error {
	return s.SendMsg(m)
}

func (s *serverStream) Recv(m any) error {
	if err := s.RecvMsg(m); err != nil {
		return err
	}
	if s.req != nil {
		if err := DefaultRequestQuery(s.req, m); err != nil {
			return err
		}
		if err := DefaultRequestVars(s.req, m); err != nil {
			return err
		}
	}
	return nil
}

func (s *serverStream) SendAndClose(m any) error {
	return s.SendMsg(m)
}

func (s *serverStream) SendMsg(m any) error {
	switch s.mode {
	case streamModeSSE:
		return s.sendSSE("message", m)
	case streamModeWebSocket:
		return s.writeWebSocketMessage(m)
	default:
		return stderrors.New("unknown HTTP stream mode")
	}
}

func (s *serverStream) RecvMsg(m any) error {
	if s.mode != streamModeWebSocket {
		return io.EOF
	}
	return readWebSocketMessage(s.conn, m, s.decoder)
}

func (s *serverStream) Close(err error) error {
	switch s.mode {
	case streamModeSSE:
		if err == nil {
			return nil
		}
		if !s.started {
			return err
		}
		_ = s.sendSSE("error", kerrors.FromError(err))
		return nil
	case streamModeWebSocket:
		if s.conn == nil {
			return err
		}
		if err != nil {
			_ = s.writeWebSocketControl(websocketControlError + err.Error())
			_ = s.writeWebSocketClose(websocket.CloseInternalServerErr, err.Error())
			_ = s.conn.Close()
			return nil
		}
		_ = s.writeWebSocketClose(websocket.CloseNormalClosure, "")
		return s.conn.Close()
	default:
		return err
	}
}

func (s *serverStream) startSSE() {
	if s.started {
		return
	}
	h := s.res.Header()
	h.Set("Content-Type", sseContentType)
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no")
	copyMetadataToHeader(h, s.header)
	s.res.WriteHeader(stdhttp.StatusOK)
	s.started = true
}

func (s *serverStream) sendSSE(event string, v any) error {
	data, err := marshalStreamMessage(v, s.encoder)
	if err != nil {
		return err
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	s.startSSE()
	if _, err = fmt.Fprintf(s.res, "event: %s\n", event); err != nil {
		return err
	}
	for _, line := range bytes.Split(data, []byte("\n")) {
		if _, err = fmt.Fprintf(s.res, "data: %s\n", line); err != nil {
			return err
		}
	}
	if _, err = io.WriteString(s.res, "\n"); err != nil {
		return err
	}
	if flusher, ok := s.res.(stdhttp.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

func (s *serverStream) writeWebSocketMessage(m any) error {
	data, err := marshalStreamMessage(m, s.encoder)
	if err != nil {
		return err
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.conn.WriteMessage(websocket.TextMessage, data)
}

func (s *serverStream) writeWebSocketControl(message string) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (s *serverStream) writeWebSocketClose(code int, text string) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	msg := websocket.FormatCloseMessage(code, text)
	return s.conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
}

type sseClientStream struct {
	ctx     context.Context
	res     *stdhttp.Response
	scanner *bufio.Scanner
	decoder encoding.Codec
}

func newSSEClientStream(ctx context.Context, res *stdhttp.Response, decoder encoding.Codec) ClientStream {
	scanner := bufio.NewScanner(res.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return &sseClientStream{ctx: ctx, res: res, scanner: scanner, decoder: decoder}
}

func (s *sseClientStream) Header() (metadata.MD, error) {
	return metadataFromHeader(s.res.Header), nil
}

func (s *sseClientStream) Trailer() metadata.MD {
	return metadataFromHeader(s.res.Trailer)
}

func (s *sseClientStream) CloseSend() error {
	return s.res.Body.Close()
}

func (s *sseClientStream) Context() context.Context {
	if s.ctx == nil {
		return context.Background()
	}
	return s.ctx
}

func (s *sseClientStream) Send(any) error {
	return stderrors.New("SSE client stream does not support Send")
}

func (s *sseClientStream) Recv(m any) error {
	return s.RecvMsg(m)
}

func (s *sseClientStream) CloseAndRecv(any) error {
	return stderrors.New("SSE client stream does not support CloseAndRecv")
}

func (s *sseClientStream) SendMsg(any) error {
	return stderrors.New("SSE client stream does not support SendMsg")
}

func (s *sseClientStream) RecvMsg(m any) error {
	for {
		event, data, err := s.readEvent()
		if err != nil {
			return err
		}
		switch event {
		case "", "message":
			return unmarshalStreamMessage(data, m, s.decoder)
		case "error":
			se := new(kerrors.Error)
			if err := unmarshalStreamMessage(data, se, s.decoder); err == nil {
				return se
			}
			return stderrors.New(string(data))
		}
	}
}

func (s *sseClientStream) readEvent() (string, []byte, error) {
	var (
		event string
		data  bytes.Buffer
	)
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if line == "" {
			if event == "" && data.Len() == 0 {
				continue
			}
			return event, bytes.TrimSuffix(data.Bytes(), []byte("\n")), nil
		}
		switch {
		case strings.HasPrefix(line, "event:"):
			event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		case strings.HasPrefix(line, "data:"):
			value := strings.TrimPrefix(line, "data:")
			if strings.HasPrefix(value, " ") {
				value = strings.TrimPrefix(value, " ")
			}
			data.WriteString(value)
			data.WriteByte('\n')
		}
	}
	if err := s.scanner.Err(); err != nil {
		return "", nil, err
	}
	return "", nil, io.EOF
}

type websocketClientStream struct {
	ctx     context.Context
	conn    *websocket.Conn
	header  stdhttp.Header
	done    func(error)
	encoder encoding.Codec
	decoder encoding.Codec
	closed  bool
	writeMu sync.Mutex
}

func (s *websocketClientStream) Header() (metadata.MD, error) {
	return metadataFromHeader(s.header), nil
}

func (s *websocketClientStream) Trailer() metadata.MD {
	return nil
}

func (s *websocketClientStream) CloseSend() error {
	if s.closed {
		return nil
	}
	return s.writeControl(websocketControlEnd)
}

func (s *websocketClientStream) Context() context.Context {
	if s.ctx == nil {
		return context.Background()
	}
	return s.ctx
}

func (s *websocketClientStream) Send(m any) error {
	return s.SendMsg(m)
}

func (s *websocketClientStream) Recv(m any) error {
	return s.RecvMsg(m)
}

func (s *websocketClientStream) CloseAndRecv(m any) error {
	if err := s.CloseSend(); err != nil {
		return err
	}
	defer s.close(nil)
	return s.RecvMsg(m)
}

func (s *websocketClientStream) SendMsg(m any) error {
	data, err := marshalStreamMessage(m, s.encoder)
	if err != nil {
		return err
	}
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.conn.WriteMessage(websocket.TextMessage, data)
}

func (s *websocketClientStream) RecvMsg(m any) error {
	if err := readWebSocketMessage(s.conn, m, s.decoder); err != nil {
		_ = s.close(err)
		return err
	}
	return nil
}

func (s *websocketClientStream) writeControl(message string) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (s *websocketClientStream) close(err error) error {
	if s.closed {
		return nil
	}
	s.closed = true
	if s.done != nil {
		s.done(err)
	}
	_ = s.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
	return s.conn.Close()
}

// ServerSentEvent opens an HTTP server-streaming call and receives replies as SSE events.
func (client *Client) ServerSentEvent(ctx context.Context, method, path string, args any, opts ...CallOption) (ClientStream, error) {
	var (
		contentType string
		body        io.Reader
	)
	c := defaultCallInfo(path)
	for _, o := range opts {
		if err := o.before(&c); err != nil {
			return nil, err
		}
	}
	if args != nil {
		data, err := client.opts.encoder(ctx, c.contentType, args)
		if err != nil {
			return nil, err
		}
		contentType = c.contentType
		body = bytes.NewReader(data)
	} else if c.contentTypeSet {
		contentType = c.contentType
	}
	url := fmt.Sprintf("%s://%s%s", client.target.Scheme, client.target.Authority, path)
	req, err := stdhttp.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	prepareClientRequest(client, req, contentType, c)
	ctx = transport.NewClientContext(ctx, &Transport{
		endpoint:     client.opts.endpoint,
		reqHeader:    headerCarrier(req.Header),
		operation:    c.operation,
		request:      req,
		pathTemplate: c.pathTemplate,
	})
	res, err := client.do(req.WithContext(ctx))
	if res != nil {
		cs := csAttempt{res: res}
		for _, o := range opts {
			o.after(&c, &cs)
		}
	}
	if err != nil {
		return nil, err
	}
	return newSSEClientStream(ctx, res, streamCodecFromCallInfo(c, "Accept", "Content-Type")), nil
}

// WebSocket opens an HTTP bidirectional streaming call over WebSocket.
func (client *Client) WebSocket(ctx context.Context, path string, opts ...CallOption) (ClientStream, error) {
	c := defaultCallInfo(path)
	for _, o := range opts {
		if err := o.before(&c); err != nil {
			return nil, err
		}
	}
	scheme := "ws"
	if client.target.Scheme == schemeHTTPS {
		scheme = "wss"
	}
	url := fmt.Sprintf("%s://%s%s", scheme, client.target.Authority, path)
	header := stdhttp.Header{}
	if c.headerCarrier != nil {
		header = *c.headerCarrier
	}
	if c.accept != "" {
		header.Set("Accept", c.accept)
	}
	if c.contentTypeSet {
		header.Set("Content-Type", c.contentType)
	}
	if client.opts.userAgent != "" {
		header.Set("User-Agent", client.opts.userAgent)
	}

	var done func(context.Context, selector.DoneInfo)
	if client.r != nil {
		node, doneFunc, err := client.selector.Select(ctx, selector.WithNodeFilter(client.opts.nodeFilters...))
		if err != nil {
			return nil, kerrors.ServiceUnavailable("NODE_NOT_FOUND", err.Error())
		}
		done = doneFunc
		if client.insecure {
			scheme = "ws"
		} else {
			scheme = "wss"
		}
		url = fmt.Sprintf("%s://%s%s", scheme, node.Address(), path)
	}

	req, err := stdhttp.NewRequestWithContext(ctx, stdhttp.MethodGet, url, nil)
	if err != nil {
		if done != nil {
			done(ctx, selector.DoneInfo{Err: err})
		}
		return nil, err
	}
	req.Header = header
	ctx = transport.NewClientContext(ctx, &Transport{
		endpoint:     client.opts.endpoint,
		reqHeader:    headerCarrier(req.Header),
		operation:    c.operation,
		request:      req,
		pathTemplate: c.pathTemplate,
	})
	dialer := websocket.Dialer{
		Proxy:            stdhttp.ProxyFromEnvironment,
		HandshakeTimeout: client.opts.timeout,
		TLSClientConfig:  client.opts.tlsConf,
	}
	conn, res, err := dialer.DialContext(ctx, url, header)
	if res != nil {
		cs := csAttempt{res: res}
		for _, o := range opts {
			o.after(&c, &cs)
		}
	}
	if err != nil {
		if done != nil {
			done(ctx, selector.DoneInfo{Err: err})
		}
		return nil, err
	}
	return &websocketClientStream{
		ctx:     ctx,
		conn:    conn,
		header:  res.Header,
		encoder: streamCodecFromCallInfo(c, "Content-Type", "Accept"),
		decoder: streamCodecFromCallInfo(c, "Accept", "Content-Type"),
		done: func(err error) {
			if done != nil {
				done(ctx, selector.DoneInfo{Err: err})
			}
		},
	}, nil
}

func prepareClientRequest(client *Client, req *stdhttp.Request, contentType string, c callInfo) {
	if c.headerCarrier != nil {
		req.Header = *c.headerCarrier
	}
	if contentType != "" {
		req.Header.Set("Content-Type", c.contentType)
	}
	if c.accept != "" {
		req.Header.Set("Accept", c.accept)
	}
	if client.opts.userAgent != "" {
		req.Header.Set("User-Agent", client.opts.userAgent)
	}
}

func marshalStreamMessage(v any, codec encoding.Codec) ([]byte, error) {
	if body, ok := HTTPBody(v); ok {
		return body.GetData(), nil
	}
	if codec == nil {
		codec = defaultStreamCodec()
	}
	return codec.Marshal(v)
}

func unmarshalStreamMessage(data []byte, v any, codec encoding.Codec) error {
	if body, ok := HTTPBody(v); ok {
		body.Data = data
		return nil
	}
	if codec == nil {
		codec = defaultStreamCodec()
	}
	return codec.Unmarshal(data, v)
}

func readWebSocketMessage(conn *websocket.Conn, m any, codec encoding.Codec) error {
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return io.EOF
			}
			return err
		}
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}
		text := string(data)
		switch {
		case text == websocketControlEnd:
			return io.EOF
		case strings.HasPrefix(text, websocketControlError):
			return stderrors.New(strings.TrimPrefix(text, websocketControlError))
		default:
			return unmarshalStreamMessage(data, m, codec)
		}
	}
}

func streamCodecFromCallInfo(c callInfo, names ...string) encoding.Codec {
	header := stdhttp.Header{}
	if c.accept != "" {
		header.Set("Accept", c.accept)
	}
	if c.contentTypeSet {
		header.Set("Content-Type", c.contentType)
	}
	return streamCodecFromHeaders(header, names...)
}

func streamCodecFromHeaders(header stdhttp.Header, names ...string) encoding.Codec {
	for _, name := range names {
		for _, values := range header.Values(name) {
			for _, value := range strings.Split(values, ",") {
				contentType := strings.TrimSpace(value)
				if codec := encoding.GetCodec(httputil.ContentSubtype(contentType)); codec != nil {
					return codec
				}
			}
		}
	}
	return defaultStreamCodec()
}

func defaultStreamCodec() encoding.Codec {
	if codec := encoding.GetCodec("protojson"); codec != nil {
		return codec
	}
	return encoding.GetCodec("json")
}

func copyMetadataToHeader(h stdhttp.Header, md metadata.MD) {
	for k, values := range md {
		for _, v := range values {
			h.Add(k, v)
		}
	}
}

func metadataFromHeader(h stdhttp.Header) metadata.MD {
	md := metadata.MD{}
	for k, values := range h {
		for _, v := range values {
			md.Append(k, v)
		}
	}
	return md
}

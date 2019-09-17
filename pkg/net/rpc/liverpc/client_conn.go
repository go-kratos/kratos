package liverpc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/net/trace"

	"github.com/gogo/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// ClientConn connect represent a real client connection to a rpc server
type ClientConn struct {
	addr        string
	network     string
	rwc         io.ReadWriteCloser
	Timeout     time.Duration
	DialTimeout time.Duration
	callInfo    *callInfo

	// appID here is used to log or something
	// which will not affect any rpc behavior
	// appID is the discovery_id
	appID string
}

type fullReqMsg struct {
	Header *Header     `json:"header"`
	HTTP   *HTTP       `json:"http"`
	Body   interface{} `json:"body"`
}

// Dial dial a rpc server
func Dial(ctx context.Context,
	network, addr string,
	timeout time.Duration, connTimeout time.Duration, appID string) (*ClientConn, error) {
	c := &ClientConn{
		appID:       appID,
		addr:        addr,
		network:     network,
		Timeout:     timeout,
		DialTimeout: connTimeout,
	}
	conn, err := net.DialTimeout(c.network, c.addr, c.DialTimeout)
	if err != nil {
		return nil, err
	}
	c.rwc = conn
	return c, err
}

// Close close the caller connection.
func (c *ClientConn) Close() error {
	if c.rwc != nil {
		return c.rwc.Close()
	}
	return nil
}

func (c *ClientConn) writeRequest(ctx context.Context, req *protoReq) (err error) {
	var (
		headerBuf = make([]byte, _headerLen)
		header    = req.Header
		body      = req.Body
	)
	binary.BigEndian.PutUint32(headerBuf[0:4], header.magic)
	binary.BigEndian.PutUint32(headerBuf[4:8], header.timestamp)
	binary.BigEndian.PutUint32(headerBuf[8:12], header.checkSum)
	binary.BigEndian.PutUint32(headerBuf[12:16], header.version)
	binary.BigEndian.PutUint32(headerBuf[16:20], header.reserved)
	binary.BigEndian.PutUint32(headerBuf[20:24], header.seq)
	binary.BigEndian.PutUint32(headerBuf[24:28], uint32(len(body)))
	copy(headerBuf[28:60], header.cmd)
	if _, err = c.rwc.Write(headerBuf); err != nil {
		err = errors.Wrap(err, "write req header error")
		return
	}
	if log.V(2) {
		log.Info("liverpc body: %s", string(body))
	}
	if _, err = c.rwc.Write(body); err != nil {
		err = errors.Wrap(err, "write req body error")
		return
	}
	return
}

func (c *ClientConn) readResponse(ctx context.Context, resp *protoResp) (err error) {
	var (
		headerBuf = make([]byte, _headerLen)
		length    int
	)
	if _, err = c.rwc.Read(headerBuf); err != nil {
		err = errors.Wrap(err, "read resp header error")
		return
	}
	resp.Header.magic = binary.BigEndian.Uint32(headerBuf[0:4])
	resp.Header.timestamp = binary.BigEndian.Uint32(headerBuf[4:8])
	resp.Header.checkSum = binary.BigEndian.Uint32(headerBuf[8:12])
	resp.Header.version = binary.BigEndian.Uint32(headerBuf[12:16])
	resp.Header.reserved = binary.BigEndian.Uint32(headerBuf[16:20])
	resp.Header.seq = binary.BigEndian.Uint32(headerBuf[20:24])
	resp.Header.length = binary.BigEndian.Uint32(headerBuf[24:28])
	resp.Header.cmd = headerBuf[28:60]
	resp.Body = make([]byte, resp.Header.length)
	if length, err = io.ReadFull(c.rwc, resp.Body); err != nil {
		err = errors.Wrap(err, "read resp body error")
		return
	}
	if uint32(length) != resp.Header.length {
		err = errors.New("bad resp body data")
		return
	}
	return
}

func (c *ClientConn) composeReqPackHeader(reqPack *protoReq, version int, serviceMethod string) {
	reqPack.Header.magic = _magic
	reqPack.Header.checkSum = 0
	reqPack.Header.seq = 1
	reqPack.Header.timestamp = uint32(time.Now().Unix())
	reqPack.Header.reserved = 0
	reqPack.Header.version = uint32(version)
	// command: {message_type}controller.method
	reqPack.Header.cmd = make([]byte, 32)
	reqPack.Header.cmd[0] = _cmdReqType
	// serviceMethod: Room.room_init
	copy(reqPack.Header.cmd[1:], []byte(serviceMethod))
}

func (c *ClientConn) setupDeadline(ctx context.Context) error {
	var t time.Duration
	if c.callInfo.Timeout != 0 {
		t = c.callInfo.Timeout
	} else {
		t, _ = ctx.Value(KeyTimeout).(time.Duration)
	}
	if t == 0 {
		t = c.Timeout
	}

	conn := c.rwc.(net.Conn)
	if conn != nil {
		err := conn.SetDeadline(time.Now().Add(t))
		if err != nil {
			conn.Close()
			return err
		}
	}
	return nil
}

// CallRaw call the service method, waits for it to complete, and returns reply its error status.
// this is can be use without protobuf
// client: {service}
// serviceMethod: {version}|{controller.method}
// httpURL: /room/v1/Room/room_init
// httpURL: /{service}/{version}/{controller}/{method}
func (c *ClientConn) CallRaw(ctx context.Context, version int, serviceMethod string, in *Args) (out *Reply, err error) {
	var (
		reqPack  protoReq
		respPack protoResp
		code     = "0"
		now      = time.Now()
		uid      int64
	)

	path := getPath(serviceMethod, version, c.appID)
	tr := c.doTrace(ctx, path)
	defer func() {
		stats.Timing(path, int64(time.Since(now)/time.Millisecond))
		if code != "" {
			stats.Incr(path, code)
		}
		logging(ctx, path, c.addr, err, time.Since(now), uid)
		defer tr.Finish(&err)
	}()

	if err = c.setupDeadline(ctx); err != nil {
		return
	}
	// it is ok for request http field to be nil

	h := createHeader(ctx, tr)
	if hdr, _ := ctx.Value(KeyHeader).(*Header); hdr != nil {
		mergeHeader(h, hdr)
	}
	if c.callInfo.Header != nil {
		mergeHeader(h, c.callInfo.Header)
	}
	if in.Header != nil {
		mergeHeader(h, in.Header)
	}
	in.Header = h
	uid = in.Header.Uid

	if in.HTTP == nil {
		if c.callInfo.HTTP != nil {
			in.HTTP = c.callInfo.HTTP
		}
	}
	if in.Body == nil {
		in.Body = map[string]interface{}{}
	}

	c.composeReqPackHeader(&reqPack, version, serviceMethod)

	var reqBytes []byte
	if reqBytes, err = json.Marshal(in); err != nil {
		err = errors.Wrap(err, "CallRaw json marshal error")
		code = "marshalErr"
		return
	}
	reqPack.Body = reqBytes

	ch := make(chan error, 1)
	go func() {
		ch <- c.sendAndRecv(ctx, &reqPack, &respPack)
	}()
	select {
	case <-ctx.Done():
		err = errors.WithStack(ctx.Err())
		code = "canceled"
		return
	case err = <-ch:
		if err != nil {
			code = "ioErr"
			return
		}
	}

	out = &Reply{}
	if err = json.Unmarshal(respPack.Body, out); err != nil {
		err = errors.Wrap(err, "proto unmarshal error: "+string(respPack.Body))
		code = "unmarshalErr"
		return
	}
	return
}

func logging(ctx context.Context, path string, addr string, err error, ts time.Duration, uid int64) {
	var (
		errMsg string
	)
	logFunc := log.Infov
	if err != nil {
		if errors.Cause(err) == context.Canceled {
			logFunc = log.Warnv
		} else {
			logFunc = log.Errorv
		}
		errMsg = fmt.Sprintf("%+v", err)
	}
	logFunc(ctx,
		log.KVString("path", path),
		log.KVString("error", errMsg),
		log.KVString("addr", addr),
		log.KVFloat64("ts", ts.Seconds()),
		log.KVInt64("uid", uid),
		log.KVString("log", "LIVERPC"),
	)
}

func (c *ClientConn) sendAndRecv(ctx context.Context, reqPack *protoReq, respPack *protoResp) (err error) {
	if err = c.writeRequest(ctx, reqPack); err != nil {
		return
	}
	if err = c.readResponse(ctx, respPack); err != nil {
		return
	}
	return
}

// Call call the service method, waits for it to complete, and returns its error status.
// this is used with protobuf generated msg
// client: {service}
// serviceMethod: {version}|{controller.method}
// httpURL: /room/v1/Room/room_init
// httpURL: /{service}/{version}/{controller}/{method}
func (c *ClientConn) Call(ctx context.Context, version int, serviceMethod string, in, out proto.Message) (err error) {
	var (
		reqPack  protoReq
		respPack protoResp
		code     = "0"
		now      = time.Now()
		uid      int64
	)

	path := getPath(serviceMethod, version, c.appID)
	tr := c.doTrace(ctx, path)

	defer func() {
		stats.Timing(path, int64(time.Since(now)/time.Millisecond))
		if code != "" {
			stats.Incr(path, code)
		}
		logging(ctx, path, c.addr, err, time.Since(now), uid)
		tr.Finish(&err)
	}()

	if err = c.setupDeadline(ctx); err != nil {
		return
	}
	fullMsg := &fullReqMsg{}
	fullMsg.Header = createHeader(ctx, tr)
	if hdr, _ := ctx.Value(KeyHeader).(*Header); hdr != nil {
		mergeHeader(fullMsg.Header, hdr)
	}
	if c.callInfo.Header != nil {
		mergeHeader(fullMsg.Header, c.callInfo.Header)
	}
	uid = fullMsg.Header.Uid

	if c.callInfo.HTTP != nil {
		fullMsg.HTTP = c.callInfo.HTTP
	}
	fullMsg.Body = in

	// it is ok for request http field to be nil

	c.composeReqPackHeader(&reqPack, version, serviceMethod)

	var reqBody []byte
	if reqBody, err = json.Marshal(fullMsg); err != nil {
		err = errors.Wrap(err, "Call json marshal error")
		code = "marshalErr"
		return
	}
	reqPack.Body = reqBody

	ch := make(chan error, 1)
	go func() {
		ch <- c.sendAndRecv(ctx, &reqPack, &respPack)
	}()
	select {
	case <-ctx.Done():
		err = errors.WithStack(ctx.Err())
		code = "canceled"
		return
	case err = <-ch:
		if err != nil {
			code = "ioErr"
			return
		}
	}

	if err = jsoniter.Unmarshal(respPack.Body, out); err != nil {
		err = errors.Wrap(err, "proto unmarshal error: "+string(respPack.Body))
		code = "unmarshalErr"
		return
	}
	return
}

// gRPC like path
// eg. /live.room/v1/Room/get_info
func getPath(serviceMethod string, version int, appID string) string {
	return fmt.Sprintf("/%s/v%d/%s", appID, version, strings.Replace(serviceMethod, ".", "/", 1))
}

func (c *ClientConn) doTrace(ctx context.Context, path string) trace.Trace {
	tracer, ok := trace.FromContext(ctx)
	if !ok {
		tracer = trace.New(path)
	} else {
		tracer = tracer.Fork("", path)
	}
	tracer.SetTag(
		trace.String(trace.TagPeerAddress, c.addr),
		trace.String(trace.TagPeerService, c.appID),
	)
	return tracer
}

func createHeader(ctx context.Context, tracer trace.Trace) *Header {
	header := &Header{}
	header.UserIp = metadata.String(ctx, metadata.RemoteIP)
	header.Caller = strings.Replace(env.AppID, ".", "-", -1)
	if header.Caller == "" {
		header.Caller = "unknown"
	}
	err := trace.Inject(tracer, nil, header)
	if err != nil {
		log.Errorw(ctx, "log", fmt.Sprintf("liverpc Inject error: %+v", err))
	}
	mid, _ := metadata.Value(ctx, metadata.Mid).(int64)
	header.Uid = mid
	if color := metadata.String(ctx, metadata.Color); color != "" {
		header.SourceGroup = color
	} else {
		header.SourceGroup = env.Color
	}
	//header.Platform = ctx.Request.FormValue("platform")
	device, ok := metadata.Value(ctx, metadata.Device).(*blademaster.Device)
	if ok {
		header.Platform = device.RawPlatform
		header.Buvid = device.Buvid
	}
	return header
}

func mergeHeader(src, ext *Header) {
	if ext.Caller != "" {
		src.Caller = ext.Caller
	}
	if ext.Uid != 0 {
		src.Uid = ext.Uid
	}
	if ext.Platform != "" {
		src.Platform = ext.Platform
	}
	if ext.Src != "" {
		src.Src = ext.Src
	}
	if ext.TraceId != "" {
		src.TraceId = ext.TraceId
	}
	if ext.UserIp != "" {
		src.UserIp = ext.UserIp
	}
	if ext.SourceGroup != "" {
		src.SourceGroup = ext.SourceGroup
	}
	if ext.Buvid != "" {
		src.Buvid = ext.Buvid
	}
	if ext.Sessdata2 != "" {
		src.Sessdata2 = ext.Sessdata2
	}
	if ext.XOperateUserName != "" {
		src.XOperateUserName = ext.XOperateUserName
	}
	if ext.XOperateAccessTime != "" {
		src.XOperateAccessTime = ext.XOperateAccessTime
	}
}

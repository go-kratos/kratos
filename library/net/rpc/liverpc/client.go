package liverpc

import (
	"context"
	"fmt"
	"net/url"
	"sync/atomic"
	"time"

	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	"go-common/library/net/metadata"
	"go-common/library/net/trace"
	"go-common/library/stat"
	xtime "go-common/library/time"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// Key is ContextKey
type Key int

const (
	_ Key = iota
	// KeyHeader use this in context to pass rpc header field
	// Depreated 请使用HeaderOption来传递Header
	KeyHeader
	// KeyTimeout deprecated
	// Depreated 请使用HTTPOption来传递HTTP
	KeyTimeout
)

const (
	_scheme      = "liverpc"
	_dialRetries = 3
)

// Get Implement tracer carrier interface
func (m *Header) Get(key string) string {
	if key == trace.KeyTraceID {
		return m.TraceId
	}
	return ""
}

// Set Implement tracer carrier interface
func (m *Header) Set(key string, val string) {
	if key == trace.KeyTraceID {
		m.TraceId = val
	}
}

var (
	// ErrNoClient no RPC client.
	errNoClient     = errors.New("no rpc client")
	errGroupInvalid = errors.New("invalid group")

	stats = stat.RPCClient
)

// GroupAddrs a map struct storing addrs vary groups
type GroupAddrs map[string][]string

// ClientConfig client config.
type ClientConfig struct {
	AppID       string
	Group       string
	Timeout     xtime.Duration
	ConnTimeout xtime.Duration
	Addr        string // if addr is provided, it will use add, else, use discovery
}

// Client is a RPC client.
type Client struct {
	conf     *ClientConfig
	dis      naming.Resolver
	addrs    atomic.Value // GroupAddrs
	addrsIdx int64
}

// NewClient new a RPC client with discovery.
func NewClient(c *ClientConfig) *Client {
	if c.Timeout <= 0 {
		c.Timeout = xtime.Duration(time.Second)
	}
	if c.ConnTimeout <= 0 {
		c.ConnTimeout = xtime.Duration(time.Second)
	}
	cli := &Client{
		conf: c,
	}
	if c.Addr != "" {
		groupAddrs := make(GroupAddrs)
		groupAddrs[""] = []string{c.Addr}
		cli.addrs.Store(groupAddrs)
		return cli
	}

	cli.dis = discovery.Build(c.AppID)
	// discovery watch & fetch nodes
	event := cli.dis.Watch()
	select {
	case _, ok := <-event:
		if !ok {
			panic("刚启动就从discovery拉到了关闭的event")
		}
		cli.disFetch()
		fmt.Printf("开始创建：%s 的liverpc client，等待从discovery拉取节点：%s\n", c.AppID, time.Now().Format("2006-01-02 15:04:05"))
	case <-time.After(10 * time.Second):
		fmt.Printf("失败创建：%s 的liverpc client，竟然从discovery拉取节点超时了：%s\n", c.AppID, time.Now().Format("2006-01-02 15:04:05"))
	}
	go cli.disproc(event)
	return cli
}

func (c *Client) disproc(event <-chan struct{}) {
	for {
		_, ok := <-event
		if !ok {
			return
		}
		c.disFetch()
	}
}

func (c *Client) disFetch() {
	ins, ok := c.dis.Fetch(context.Background())
	if !ok {
		return
	}
	insZone, ok := ins[env.Zone]
	if !ok {
		return
	}
	addrs := make(GroupAddrs)
	for _, svr := range insZone {
		group, ok := svr.Metadata["color"]
		if !ok {
			group = ""
		}
		for _, addr := range svr.Addrs {
			u, err := url.Parse(addr)
			if err == nil && u.Scheme == _scheme {
				addrs[group] = append(addrs[group], u.Host)
			}
		}
	}
	if len(addrs) > 0 {
		c.addrs.Store(addrs)
	}
}

// pickConn pick conn by addrs
func (c *Client) pickConn(ctx context.Context, addrs []string, dialTimeout time.Duration) (*ClientConn, error) {
	var (
		lastErr error
	)
	if len(addrs) == 0 {
		lastErr = errors.New("addrs empty")
	} else {
		for i := 0; i < _dialRetries; i++ {
			idx := atomic.AddInt64(&c.addrsIdx, 1)
			addr := addrs[int(idx)%len(addrs)]
			if dialTimeout == 0 {
				dialTimeout = time.Duration(c.conf.ConnTimeout)
			}
			cc, err := Dial(ctx, "tcp", addr, time.Duration(c.conf.Timeout), dialTimeout)
			if err != nil {
				lastErr = errors.Wrapf(err, "Dial %s error", addr)
				continue
			}
			return cc, nil
		}
	}
	if lastErr != nil {
		return nil, errors.WithMessage(errNoClient, lastErr.Error())
	}
	return nil, errors.WithStack(errNoClient)
}

// fetchAddrs fetch addrs by different strategies
// source_group first, come from request header if exists, currently only CallRaw supports source_group
// then env group, come from os.env
// since no invalid group found, return error
func (c *Client) fetchAddrs(ctx context.Context, request interface{}) (addrs []string, err error) {
	var (
		args        *Args
		groupAddrs  GroupAddrs
		ok          bool
		sourceGroup string
		groups      []string
	)
	defer func() {
		if err != nil {
			err = errors.WithMessage(errGroupInvalid, err.Error())
		}
	}()
	// try parse request header and fetch source group
	if args, ok = request.(*Args); ok && args.Header != nil {
		sourceGroup = args.Header.SourceGroup
		if sourceGroup != "" {
			groups = append(groups, sourceGroup)
		}
	}
	metaColor := metadata.String(ctx, metadata.Color)
	if metaColor != "" && metaColor != sourceGroup {
		groups = append(groups, metaColor)
	}

	if env.Color != "" && env.Color != metaColor {
		groups = append(groups, env.Color)
	}

	groups = append(groups, "")

	if groupAddrs, ok = c.addrs.Load().(GroupAddrs); !ok {
		err = errors.New("addrs load error")
		return
	}
	if len(groupAddrs) == 0 {
		err = errors.New("group addrs empty")
		return
	}
	for _, group := range groups {
		if addrs, ok = groupAddrs[group]; ok {
			break
		}
	}
	if len(addrs) == 0 {
		err = errors.Errorf("addrs empty source(%s), metadata(%s), env(%s), default empty, allAddrs(%+v)",
			sourceGroup, metaColor, env.Color, groupAddrs)
		return
	}
	return
}

// Call call the service method, waits for it to complete, and returns its error status.
// client: {service}
// serviceMethod: {version}|{controller.method}
// httpURL: /room/v1/Room/room_init
// httpURL: /{service}/{version}/{controller}/{method}
func (c *Client) Call(ctx context.Context, version int, serviceMethod string, in proto.Message, out proto.Message, opts ...CallOption) (err error) {
	var (
		cc    *ClientConn
		addrs []string
	)
	isPickErr := true

	defer func() {
		if cc != nil {
			cc.Close()
		}
		if err != nil && isPickErr {
			log.Error("liverpc Call pick connection error, version %d, method: %s, error: %+v", version, serviceMethod, err)
		}
	}() // for now it is non-persistent connection
	var cInfo = &callInfo{}
	for _, o := range opts {
		o.before(cInfo)
	}
	addrs, err = c.fetchAddrs(ctx, in)
	if err != nil {
		return
	}
	cc, err = c.pickConn(ctx, addrs, cInfo.DialTimeout)
	if err != nil {
		return
	}
	isPickErr = false
	cc.callInfo = cInfo
	err = cc.Call(ctx, version, serviceMethod, in, out)
	if err != nil {
		return
	}
	for _, o := range opts {
		o.after(cc.callInfo)
	}
	return
}

// CallRaw call the service method, waits for it to complete, and returns reply its error status.
// this is can be use without protobuf
// client: {service}
// serviceMethod: {version}|{controller.method}
// httpURL: /room/v1/Room/room_init
// httpURL: /{service}/{version}/{controller}/{method}
func (c *Client) CallRaw(ctx context.Context, version int, serviceMethod string, in *Args, opts ...CallOption) (out *Reply, err error) {
	var (
		cc    *ClientConn
		addrs []string
	)
	isPickErr := true

	defer func() {
		if cc != nil {
			cc.Close()
		}
		if err != nil && isPickErr {
			log.Error("liverpc CallRaw pick connection error, version %d, method: %s, error: %+v", version, serviceMethod, err)
		}
	}() // for now it is non-persistent connection
	var cInfo = &callInfo{}
	for _, o := range opts {
		o.before(cInfo)
	}
	addrs, err = c.fetchAddrs(ctx, in)
	if err != nil {
		return
	}
	cc, err = c.pickConn(ctx, addrs, cInfo.DialTimeout)
	if err != nil {
		return
	}
	isPickErr = false
	cc.callInfo = cInfo
	out, err = cc.CallRaw(ctx, version, serviceMethod, in)
	if err != nil {
		return
	}
	for _, o := range opts {
		o.after(cc.callInfo)
	}
	return
}

//Close handle client exit
func (c *Client) Close() {
	if c.dis != nil {
		c.dis.Close()
	}
}

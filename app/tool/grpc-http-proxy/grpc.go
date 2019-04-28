package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/library/log"
	"go-common/library/naming/discovery"

	"github.com/gogo/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

// Reply .
type Reply struct {
	res []byte
}

// Reference https://jbrandhorst.com/post/grpc-json/
func init() {
	encoding.RegisterCodec(JSON{
		Marshaler: jsonpb.Marshaler{
			EmitDefaults: true,
			OrigName:     true,
		},
	})
}

func ipFromDiscovery(appid string) (ip string, err error) {
	d := discovery.New(nil)
	defer d.Close()
	b := d.Build(appid)
	defer b.Close()
	ch := b.Watch()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	select {
	case <-ch:
	case <-ctx.Done():
		err = fmt.Errorf("查找节点超时 请检查appid是否填写正确")
		return
	}
	ins, ok := b.Fetch(context.Background())
	if !ok {
		err = fmt.Errorf("discovery 拉取:%s 失败", appid)
		return
	}
	for _, vs := range ins {
		for _, v := range vs {
			for _, addr := range v.Addrs {
				if strings.Contains(addr, "grpc://") {
					ip = strings.Replace(addr, "grpc://", "", -1)
					return
				}
			}
		}
	}
	err = fmt.Errorf("discovery 找不到服务节点:%s", appid)
	return
}

func callGrpc(addr, method string, body []byte) (res []byte, err error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(JSON{}.Name())),
	}
	var conn *grpc.ClientConn
	if !strings.Contains(addr, ":") {
		addr, err = ipFromDiscovery(addr)
	}
	if err != nil {
		return
	}
	conn, err = grpc.Dial(addr, opts...)
	if err != nil {
		return
	}
	var reply Reply
	log.Info("callrpc method: %s body: %s", method, body)
	if err = conn.Invoke(context.Background(), method, body, &reply); err != nil {
		return
	}
	res = reply.res
	return
}

// JSON is impl of encoding.Codec
type JSON struct {
	jsonpb.Marshaler
	jsonpb.Unmarshaler
}

// Name is name of JSON
func (j JSON) Name() string {
	return "json"
}

// Marshal is json marshal
func (j JSON) Marshal(v interface{}) (out []byte, err error) {
	return v.([]byte), nil
}

// Unmarshal is json unmarshal
func (j JSON) Unmarshal(data []byte, v interface{}) (err error) {
	v.(*Reply).res = data
	return nil
}

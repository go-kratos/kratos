package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
)

// Reply for test
type Reply struct {
	res []byte
}

var data string
var file string
var method string
var addr string
var tlsCert string
var tlsServerName string

//Reference https://jbrandhorst.com/post/grpc-json/
func init() {
	encoding.RegisterCodec(JSON{
		Marshaler: jsonpb.Marshaler{
			EmitDefaults: true,
			OrigName:     true,
		},
	})
	flag.StringVar(&data, "data", `{"name":"longxia","age":19}`, `{"name":"longxia","age":19}`)
	flag.StringVar(&file, "file", ``, `./data.json`)
	flag.StringVar(&method, "method", "/testproto.Greeter/SayHello", `/testproto.Greeter/SayHello`)
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", `127.0.0.1:8080`)
	flag.StringVar(&tlsCert, "cert", "", `./cert.pem`)
	flag.StringVar(&tlsServerName, "server_name", "", `hello_server`)
}

// 该example因为使用的是json传输格式所以只能用于调试或测试，用于线上会导致性能下降
// 使用方法：
//  ./grpcDebug -data='{"name":"xia","age":19}' -addr=127.0.0.1:8080 -method=/testproto.Greeter/SayHello
//  ./grpcDebug -file=data.json -addr=127.0.0.1:8080 -method=/testproto.Greeter/SayHello
func main() {
	flag.Parse()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(JSON{}.Name())),
	}
	if tlsCert != "" {
		creds, err := credentials.NewClientTLSFromFile(tlsCert, tlsServerName)
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	if file != "" {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("ioutil.ReadFile %s failed!err:=%v", file, err)
			os.Exit(1)
		}
		if len(content) > 0 {
			data = string(content)
		}
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}
	var reply Reply
	err = grpc.Invoke(context.Background(), method, []byte(data), &reply, conn)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(reply.res))
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

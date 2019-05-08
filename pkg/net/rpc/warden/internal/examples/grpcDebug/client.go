package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
)

// Reply for test
type Reply struct {
	res []byte
}

type Discovery struct {
	HttpClient *http.Client
	Nodes      []string
}

var (
	data          string
	file          string
	method        string
	addr          string
	tlsCert       string
	tlsServerName string
	appID         string
	env           string
)

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
	flag.StringVar(&appID, "appid", "", `appid`)
	flag.StringVar(&env, "env", "", `env`)
}

// 该example因为使用的是json传输格式所以只能用于调试或测试，用于线上会导致性能下降
// 使用方法：
//  ./grpcDebug -data='{"name":"xia","age":19}' -addr=127.0.0.1:8080 -method=/testproto.Greeter/SayHello
//  ./grpcDebug -file=data.json -addr=127.0.0.1:8080 -method=/testproto.Greeter/SayHello
//  DEPLOY_ENV=uat ./grpcDebug -appid=main.community.reply-service -method=/reply.service.v1.Reply/ReplyInfoCache -data='{"rp_id"=1493769244}'
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
			fmt.Printf("ioutil.ReadFile(%s) error(%v)\n", file, err)
			os.Exit(1)
		}
		if len(content) > 0 {
			data = string(content)
		}
	}
	if appID != "" {
		addr = ipFromDiscovery(appID, env)
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

func ipFromDiscovery(appID, env string) string {
	d := &Discovery{
		Nodes:      []string{"discovery.bilibili.co", "api.bilibili.co"},
		HttpClient: http.DefaultClient,
	}
	deployEnv := os.Getenv("DEPLOY_ENV")
	if deployEnv != "" {
		env = deployEnv
	}
	return d.addr(appID, env, d.nodes())
}

func (d *Discovery) nodes() (addrs []string) {
	res := new(struct {
		Code int `json:"code"`
		Data []struct {
			Addr string `json:"addr"`
		} `json:"data"`
	})
	resp, err := d.HttpClient.Get(fmt.Sprintf("http://%s/discovery/nodes", d.Nodes[rand.Intn(len(d.Nodes))]))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		panic(err)
	}
	for _, data := range res.Data {
		addrs = append(addrs, data.Addr)
	}
	return
}

func (d *Discovery) addr(appID, env string, nodes []string) (ip string) {
	res := new(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    map[string]*struct {
			ZoneInstances map[string][]*struct {
				AppID string   `json:"appid"`
				Addrs []string `json:"addrs"`
			} `json:"zone_instances"`
		} `json:"data"`
	})
	host, _ := os.Hostname()
	resp, err := d.HttpClient.Get(fmt.Sprintf("http://%s/discovery/polls?appid=%s&env=%s&hostname=%s", nodes[rand.Intn(len(nodes))], appID, env, host))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		panic(err)
	}
	for _, data := range res.Data {
		for _, zoneInstance := range data.ZoneInstances {
			for _, instance := range zoneInstance {
				if instance.AppID == appID {
					for _, addr := range instance.Addrs {
						if strings.Contains(addr, "grpc://") {
							return strings.Replace(addr, "grpc://", "", -1)
						}
					}
				}
			}
		}
	}
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

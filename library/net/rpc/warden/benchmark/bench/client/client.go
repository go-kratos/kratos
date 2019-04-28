package main

import (
	"flag"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/net/netutil/breaker"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/rpc/warden/benchmark/bench/proto"
	xtime "go-common/library/time"

	goproto "github.com/gogo/protobuf/proto"
	"github.com/montanaflynn/stats"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	iws         = 65535 * 1000
	iwsc        = 65535 * 10000
	readBuffer  = 32 * 1024
	writeBuffer = 32 * 1024
)

var concurrency = flag.Int("c", 50, "concurrency")
var total = flag.Int("t", 500000, "total requests for all clients")
var host = flag.String("s", "127.0.0.1:8972", "server ip and port")
var isWarden = flag.Bool("w", true, "is warden or grpc client")
var strLen = flag.Int("l", 600, "the length of the str")

func wardenCli() proto.HelloClient {
	log.Println("start warden cli")
	client := warden.NewClient(&warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 10),
		Timeout: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(3 * time.Second),
			Sleep:   xtime.Duration(3 * time.Second),
			Bucket:  10,
			Ratio:   0.3,
			Request: 20,
		},
	},
		grpc.WithInitialWindowSize(iws),
		grpc.WithInitialConnWindowSize(iwsc),
		grpc.WithReadBufferSize(readBuffer),
		grpc.WithWriteBufferSize(writeBuffer))
	conn, err := client.Dial(context.Background(), *host)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	cli := proto.NewHelloClient(conn)
	return cli
}

func grpcCli() proto.HelloClient {
	log.Println("start grpc cli")
	conn, err := grpc.Dial(*host, grpc.WithInsecure(),
		grpc.WithInitialWindowSize(iws),
		grpc.WithInitialConnWindowSize(iwsc),
		grpc.WithReadBufferSize(readBuffer),
		grpc.WithWriteBufferSize(writeBuffer))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	cli := proto.NewHelloClient(conn)
	return cli
}

func main() {
	flag.Parse()
	c := *concurrency
	m := *total / c
	var wg sync.WaitGroup
	wg.Add(c)
	log.Printf("concurrency: %d\nrequests per client: %d\n\n", c, m)

	args := prepareArgs()
	b, _ := goproto.Marshal(args)
	log.Printf("message size: %d bytes\n\n", len(b))

	var trans uint64
	var transOK uint64
	d := make([][]int64, c)
	for i := 0; i < c; i++ {
		dt := make([]int64, 0, m)
		d = append(d, dt)
	}
	var cli proto.HelloClient
	if *isWarden {
		cli = wardenCli()
	} else {
		cli = grpcCli()
	}
	//warmup
	cli.Say(context.Background(), args)

	totalT := time.Now().UnixNano()
	for i := 0; i < c; i++ {
		go func(i int) {
			for j := 0; j < m; j++ {
				t := time.Now().UnixNano()
				reply, err := cli.Say(context.Background(), args)
				t = time.Now().UnixNano() - t
				d[i] = append(d[i], t)
				if err == nil && reply.Field1 == "OK" {
					atomic.AddUint64(&transOK, 1)
				}
				atomic.AddUint64(&trans, 1)
			}
			wg.Done()
		}(i)

	}
	wg.Wait()

	totalT = time.Now().UnixNano() - totalT
	totalT = totalT / 1e6
	log.Printf("took %d ms for %d requests\n", totalT, *total)
	totalD := make([]int64, 0, *total)
	for _, k := range d {
		totalD = append(totalD, k...)
	}
	totalD2 := make([]float64, 0, *total)
	for _, k := range totalD {
		totalD2 = append(totalD2, float64(k))
	}

	mean, _ := stats.Mean(totalD2)
	median, _ := stats.Median(totalD2)
	max, _ := stats.Max(totalD2)
	min, _ := stats.Min(totalD2)
	tp99, _ := stats.Percentile(totalD2, 99)
	tp999, _ := stats.Percentile(totalD2, 99.9)

	log.Printf("sent     requests    : %d\n", *total)
	log.Printf("received requests_OK : %d\n", atomic.LoadUint64(&transOK))
	log.Printf("throughput  (TPS)    : %d\n", int64(c*m)*1000/totalT)
	log.Printf("mean: %v ms, median: %v ms, max: %v ms, min: %v ms, p99: %v ms, p999:%v ms\n", mean/1e6, median/1e6, max/1e6, min/1e6, tp99/1e6, tp999/1e6)

}

func prepareArgs() *proto.BenchmarkMessage {
	b := true
	var i int32 = 120000
	var i64 int64 = 98765432101234
	var s = "许多往事在眼前一幕一幕，变的那麼模糊"
	repeat := *strLen / (8 * 54)
	if repeat == 0 {
		repeat = 1
	}
	var str string
	for i := 0; i < repeat; i++ {
		str += s
	}
	var args proto.BenchmarkMessage

	v := reflect.ValueOf(&args).Elem()
	num := v.NumField()
	for k := 0; k < num; k++ {
		field := v.Field(k)
		if field.Type().Kind() == reflect.Ptr {
			switch v.Field(k).Type().Elem().Kind() {
			case reflect.Int, reflect.Int32:
				field.Set(reflect.ValueOf(&i))
			case reflect.Int64:
				field.Set(reflect.ValueOf(&i64))
			case reflect.Bool:
				field.Set(reflect.ValueOf(&b))
			case reflect.String:
				field.Set(reflect.ValueOf(&str))
			}
		} else {
			switch field.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				field.SetInt(9876543)
			case reflect.Bool:
				field.SetBool(true)
			case reflect.String:
				field.SetString(str)
			}
		}
	}
	return &args
}

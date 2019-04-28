package main

// Start Commond eg: ./broadcast 10 10000 127.0.0.1:7831
// first routine count
// second parameter: running time
// third parameter: service server ip

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"
)

var (
	httpClient *bm.Client
)

const TestContent = "{\"test\":\"test push braocast\"}"

func init() {
	httpClient = bm.NewClient(&bm.ClientConfig{
		App: &bm.App{
			Key:    "6aa4286456d16b97",
			Secret: "test",
		},
		Dial:      xtime.Duration(time.Second),
		Timeout:   xtime.Duration(time.Second),
		KeepAlive: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:    xtime.Duration(time.Second),
			Sleep:     xtime.Duration(time.Second),
			Bucket:    10,
			Ratio:     0.8,
			Request:   100,
			SwitchOff: false,
		},
	})
}

func main() {
	rountineNum, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	t, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	addr := os.Args[3]

	speed := os.Args[4]

	gap := time.Second / time.Duration(rountineNum)
	delay := time.Duration(0)
	time.AfterFunc(time.Duration(t)*time.Second, stop)

	go run(addr, time.Duration(0)*time.Second, speed)
	for i := 0; i < rountineNum-1; i++ {
		go run(addr, delay, speed)
		delay += gap
		fmt.Println("delay:", delay)
	}
	time.Sleep(9999 * time.Hour)
}

func stop() {
	os.Exit(-1)
}

func run(addr string, delay time.Duration, speed string) {
	time.Sleep(delay)
	for {
		go post(addr, speed)
		time.Sleep(time.Second)
	}
}

func post(addr string, speed string) {
	params := url.Values{}
	params.Set("operation", "9")
	params.Set("speed", speed)
	params.Set("message", TestContent)
	if err := httpClient.Get(context.Background(), "http://"+addr+"/x/internal/broadcast/push/all", "", params, nil); err != nil {
		fmt.Printf("Error: bm.get() error(%v)\n", err.Error())
		return
	}
}

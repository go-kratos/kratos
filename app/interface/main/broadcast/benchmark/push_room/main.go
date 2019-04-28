package main

// Start Commond eg: ./push_room test://test_room 10 100 127.0.0.1:7831
// first parameter: room id
// second parameter: routine count
// third parameter: running time
// fourth parameter: service server ip

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

const TestContent = "{\"test\":\"test push room\"}"

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
	rountineNum, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	t, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}
	addr := os.Args[4]

	time.AfterFunc(time.Duration(t)*time.Second, stop)

	gap := time.Second / time.Duration(rountineNum)
	delay := time.Duration(0)

	go run(addr, time.Duration(0)*time.Second)
	for i := 0; i < rountineNum-1; i++ {
		go run(addr, delay)
		delay += gap
		fmt.Println("delay:", delay)
	}
	time.Sleep(9999 * time.Hour)
}

func run(addr string, delay time.Duration) {
	time.Sleep(delay)
	i := int64(0)
	for {
		go post(addr, i)
		time.Sleep(time.Second)
		i++
	}
}

func stop() {
	os.Exit(-1)
}

func post(addr string, i int64) {
	params := url.Values{}
	params.Set("room", os.Args[1])
	params.Set("operation", "9")
	params.Set("message", TestContent)
	if err := httpClient.Get(context.Background(), "http://"+addr+"/x/internal/broadcast/push/room", "", params, nil); err != nil {
		fmt.Printf("Error: bm.post() error(%s)\n", err.Error())
		return
	}
}

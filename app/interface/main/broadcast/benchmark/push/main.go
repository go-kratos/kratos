package main

// Start Commond eg:./push 1 500 127.0.0.1:7831 100
// first parameter：beginning port
// second parameter: end port
// third parameter: comet server ip
// fourth parameter: runing time

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"
)

var (
	lg         *log.Logger
	httpClient *bm.Client
	t          int
)

const TestContent = "{\"test\":1}"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	infoLogfi, err := os.OpenFile("./pushs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	lg = log.New(infoLogfi, "", log.LstdFlags|log.Lshortfile)

	begin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	end, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	t, err = strconv.Atoi(os.Args[4])
	if err != nil {
		panic(err)
	}

	num := runtime.NumCPU() * 2
	lg.Printf("start routine num:%d", num)

	l := end / num
	b, e := begin, begin+l
	time.AfterFunc(time.Duration(t)*time.Second, stop)
	for i := begin; i <= end; i++ {
		this := i
		go func() {
			for {
				startPush(this, num)
				time.Sleep(time.Second)
			}
		}()
	}

	for i := 0; i < num; i++ {
		go startPush(b, e)
		b += l
		e += l
	}

	time.Sleep(9999 * time.Hour)
}

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

func stop() {
	os.Exit(-1)
}

func startPush(b, e int) {
	lg.Printf("start Push key from %d to %d", b, e)
	for i := b; i < b+e; i++ {
		params := url.Values{}
		params.Set("operation", "9")
		params.Set("keys", strconv.Itoa(b))
		params.Set("message", TestContent)
		err := httpClient.Get(context.Background(), fmt.Sprintf("http://%s/x/internal/broadcast/push/keys", os.Args[3]), "", params, nil)
		if err != nil {
			lg.Printf("get error (%v)", err)
			return
		}
	}
}

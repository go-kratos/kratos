package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"time"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"
)

const (
	_apiPushKey  = "http://172.22.33.126:7831/x/internal/broadcast/push/keys"
	_apiPushMid  = "http://172.22.33.126:7831/x/internal/broadcast/push/mids"
	_apiPushRoom = "http://172.22.33.126:7831/x/internal/broadcast/push/room"
	_apiPushAll  = "http://172.22.33.126:7831/x/internal/broadcast/push/all"
)

var (
	cmd      string
	op       string
	key      string
	mid      string
	room     string
	platform string
	message  string

	httpClient = bm.NewClient(&bm.ClientConfig{
		App: &bm.App{
			Key:    "6a29f8ed87407c11",
			Secret: "d3c5a85f5b895a03735b5d20a273bc57",
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
)

func init() {
	flag.StringVar(&cmd, "cmd", "", "cmd=key/mid/room/all")
	flag.StringVar(&op, "op", "", "op=1000,1002,1003")
	flag.StringVar(&key, "key", "", "client key")
	flag.StringVar(&mid, "mid", "", "mid")
	flag.StringVar(&room, "room", "", "room")
	flag.StringVar(&platform, "platform", "", "platform")
	flag.StringVar(&message, "message", "", "message content")
}

func main() {
	flag.Parse()
	if op == "" {
		panic("please input the op=1000/1002/1003")
	}
	switch cmd {
	case "key":
		pushKey(op, key, message)
	case "mid":
		pushMid(op, mid, message)
	case "room":
		pushRoom(op, room, message)
	case "all":
		pushAll(op, platform, message)
	default:
		log.Printf("unknown cmd=%s\n", cmd)
		return
	}
}

func pushKey(op, key, content string) (err error) {
	params := url.Values{}
	params.Set("operation", op)
	params.Set("keys", key)
	params.Set("message", content)
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = httpClient.Post(context.Background(), _apiPushKey, "", params, &res); err != nil {
		log.Printf("http error(%v)", err)
		return
	}
	log.Printf("sent op[%s] key[%s] message:%s\n result:(%d,%s)\n", op, key, message, res.Code, res.Msg)
	return
}

func pushMid(op, mid, content string) (err error) {
	params := url.Values{}
	params.Set("operation", op)
	params.Set("mids", mid)
	params.Set("message", content)
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = httpClient.Post(context.Background(), _apiPushMid, "", params, &res); err != nil {
		log.Printf("http error(%v)\n", err)
		return
	}
	log.Printf("sent op[%s] mid[%s] message:%s\n, result:(%d,%s)\n", op, mid, message, res.Code, res.Msg)
	return
}

func pushRoom(op, room, content string) (err error) {
	params := url.Values{}
	params.Set("operation", op)
	params.Set("room", room)
	params.Set("message", content)
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = httpClient.Post(context.Background(), _apiPushRoom, "", params, nil); err != nil {
		log.Printf("http error(%v)\n", err)
		return
	}
	log.Printf("sent op[%s] room[%s] message:%s\n, result:(%d,%s)\n", op, room, message, res.Code, res.Msg)
	return
}

func pushAll(op, platform, content string) (err error) {
	params := url.Values{}
	params.Set("operation", op)
	params.Set("platform", platform)
	params.Set("message", content)
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = httpClient.Post(context.Background(), _apiPushAll, "", params, &res); err != nil {
		log.Printf("http error(%v)\n", err)
		return
	}
	log.Printf("sent op[%s] platform[%s] message:%s\n, result:(%d,%s)\n", op, platform, message, res.Code, res.Msg)
	return
}

package main

// Start Commond eg: ./client 1 5000 localhost:8080
// first parameter：beginning userId
// second parameter: amount of clients
// third parameter: comet server ip

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"go-common/app/service/main/broadcast/libs/bufio"
	"go-common/app/service/main/broadcast/model"
)

var (
	lg *log.Logger
)

const (
	_heartTime = 10 * time.Second
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	_, err := os.OpenFile("./client.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	lg = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	flag.Parse()
	begin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	num, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	for i := begin; i < begin+num; i++ {
		go client(fmt.Sprintf("%d", i))
	}
	var exit chan bool
	<-exit
}

func client(key string) {
	time.Sleep(time.Duration(rand.Intn(5000)) * time.Microsecond)
	for {
		conn, err := net.DialTimeout("tcp", os.Args[3], time.Second)
		if err != nil {
			lg.Printf("err.net.Dial(%s) error(%v)\n", os.Args[3], err)
			return
		}
		defer func() {
			conn.Close()
			lg.Printf("err.conn.Close() key:%s", key)
		}()
		wr := bufio.NewWriter(conn)
		rd := bufio.NewReader(conn)
		// timeout
		if err = conn.SetDeadline(time.Now().Add(time.Second * 5)); err != nil {
			lg.Printf("err.conn.SetDeadline() error(%v)\n", err)
			return
		}
		body := fmt.Sprintf(`{"device_id":"1231231","access_key":"test","platform":"android","mobi_app":"android","build":50000,"accepts":[10001]}`)
		fmt.Println(body)
		// auth first
		proto := &model.Proto{
			Operation: model.OpAuth,
			Body:      []byte(body),
		}
		if err = proto.WriteTCP(wr); err != nil {
			lg.Printf("err.auth write error(%v)\n", err)
			return
		}
		if err = wr.Flush(); err != nil {
			lg.Printf("err.auth flush error(%v)\n", err)
			return
		}
		go heartbeat(conn, wr, key)
		for {
			if err = proto.ReadTCP(rd); err != nil {
				lg.Printf("err.read a proto key:%v error(%v)\n", key, err)
				return
			}
			switch proto.Operation {
			case model.OpHeartbeatReply:
				if err = conn.SetDeadline(time.Now().Add(_heartTime * 5)); err != nil {
					lg.Printf("err.conn.SetDeadline() key:%v error(%v)\n", key, err)
					return
				}

				lg.Printf("info.read key:%v heartbeat proto.body(%s) seq: %d\n", key, string(proto.Body), proto.SeqId)
			default:
				lg.Printf("info.read a key:%v proto(%+v)  op:%d\n", key, proto, proto.Operation)
			}
		}
	}
}

func heartbeat(conn net.Conn, wr *bufio.Writer, key string) {
	var (
		a   int32
		err error
	)
	defer conn.Close()
	for {
		a++
		proto := &model.Proto{
			Operation: model.OpHeartbeat,
			SeqId:     a,
		}
		if err = proto.WriteTCP(wr); err != nil {
			lg.Printf("err.heartbeat write error(%v)\n", err)
			return
		}
		if err = wr.Flush(); err != nil {
			lg.Printf("err.heartbeat key:%s flush error(%v)\n", key, err)
			return
		}
		lg.Printf("info.write key:%s heartbeat proto(%+v)\n", key, proto)
		a++
		// op := ChangeRoomReq{
		// 	RoomID: "test://2233",
		// }
		// body, err := json.Marshal(op)
		// if err != nil {
		// 	panic(err)
		// }
		// proto = &model.Proto{
		// 	Operation: model.OpChangeRoom,
		// 	SeqId:     a,
		// 	Body:      body,
		// }
		// if err = proto.WriteTCP(wr); err != nil {
		// 	lg.Printf("err.heartbeat write error(%v)\n", err)
		// 	return
		// }
		// if err = wr.Flush(); err != nil {
		// 	lg.Printf("err.heartbeat key:%s flush error(%v)\n", key, err)
		// 	return
		// }
		// lg.Printf("info.write key:%s heartbeat proto(%+v)\n", key, proto)

		time.Sleep(_heartTime)
	}
}

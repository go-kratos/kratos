package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"go-common/app/service/main/broadcast/libs/bufio"
	"go-common/app/service/main/broadcast/model"
	"go-common/library/xstr"
)

var (
	op        string
	key       string
	room      string
	accessKey string
	platform  string
	mobiApp   string
	addr      string
	build     int
)

const (
	_heartTime = 30 * time.Second
)

func init() {
	flag.StringVar(&addr, "addr", "172.22.33.126:7821", "server ip:port")
	flag.StringVar(&key, "device_id", "test_"+fmt.Sprint(time.Now().Unix()), "client key")
	flag.StringVar(&room, "room_id", "", "client room")
	flag.StringVar(&accessKey, "access_key", "", "client access_key")
	flag.StringVar(&platform, "platform", "", "client platform")
	flag.StringVar(&mobiApp, "mobi_app", "", "client moni_app")
	flag.StringVar(&op, "op", "", "op=1000,1002,1003")
	flag.IntVar(&build, "build", 0, "client build")
}

func main() {
	flag.Parse()
	if op == "" {
		panic("please input the op=1000/1002/1003")
	}
	if platform == "" {
		panic("please input the platform=android/ios/web")
	}
	log.Printf("addr=%s op=%s key=%s\n", addr, op, key)

	_, err := xstr.SplitInts(op)
	if err != nil {
		panic(err)
	}
	client()
}

func client() {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		log.Printf("err.net.Dial(%s) error(%v)\n", addr, err)
		return
	}
	defer func() {
		conn.Close()
		log.Printf("err.conn.Close() key:%s\n", key)
	}()
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)
	// timeout
	if err = conn.SetDeadline(time.Now().Add(time.Second * 5)); err != nil {
		log.Printf("err.conn.SetDeadline() error(%v)\n", err)
		return
	}
	body := fmt.Sprintf(`{"device_id":"%s","access_key":"%s","platform":"%s","mobi_app":"%s","build":%d,"accepts":[%s],"room_id":"%s"}`, key, accessKey, platform, mobiApp, build, op, room)
	log.Printf("auth:%s\n", body)
	// auth first
	proto := &model.Proto{
		Operation: model.OpAuth,
		Body:      []byte(body),
	}
	if err = proto.WriteTCP(wr); err != nil {
		log.Printf("err.auth write error(%v)\n", err)
		return
	}
	if err = wr.Flush(); err != nil {
		log.Printf("err.auth flush error(%v)\n", err)
		return
	}
	go heartbeat(conn, wr, key)
	for {
		if err = proto.ReadTCP(rd); err != nil {
			log.Printf("err.read a proto key:%v error(%v)\n", key, err)
			return
		}
		switch proto.Operation {
		case model.OpHeartbeatReply:
			if err = conn.SetDeadline(time.Now().Add(_heartTime * 5)); err != nil {
				log.Printf("err.conn.SetDeadline() key:%v error(%v)\n", key, err)
				return
			}

			log.Printf("info.read key:%v heartbeat proto.body(%s) seq: %d\n", key, string(proto.Body), proto.SeqId)
		default:
			log.Printf("info.read a key:%v proto(%+v)  op:%d\n", key, proto, proto.Operation)
		}
	}
}

func heartbeat(conn net.Conn, wr *bufio.Writer, key string) {
	var (
		seq int32
		err error
	)
	defer conn.Close()
	for {
		seq++
		proto := &model.Proto{
			Operation: model.OpHeartbeat,
			SeqId:     seq,
		}
		if err = proto.WriteTCP(wr); err != nil {
			log.Printf("err.heartbeat write error(%v)\n", err)
			return
		}
		if err = wr.Flush(); err != nil {
			log.Printf("err.heartbeat key:%s flush error(%v)\n", key, err)
			return
		}
		log.Printf("info.write key:%s heartbeat proto(%+v)\n", key, proto)
		time.Sleep(_heartTime)
	}
}

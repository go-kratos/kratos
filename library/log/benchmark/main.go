package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync/atomic"
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"
)

var isClient bool

func init() {
	os.Setenv("LOG_AGENT", "tcp://127.0.0.1:8080")
	flag.BoolVar(&isClient, "c", false, "-c=true")
}

func main() {

	flag.Parse()
	if isClient {
		go func() {
			fmt.Println(http.ListenAndServe("localhost:6060", nil))
		}()
		log.Init(&log.Config{
			Stdout: false,
			Agent: &log.AgentConfig{
				Proto:   "tcp",
				Addr:    "127.0.0.1:8080",
				Chan:    2048,
				Buffer:  1,
				Timeout: xtime.Duration(time.Second),
			}})
		for i := 0; i < 3; i++ {
			go func() {
				arg := `area:"reply" 
		message:"\345\233\236\345\244\215 @\347\231\276\350\215\211\345\221\263\346\235\245\344\274\212\344\273\275 :\345\223\210\357\274\214\344\275\240\345\260\261\345\217\252\347\234\213\345\210\260\344\272\206pick me up\346\262\241\347\234\213\345\210\260\346\222\236\357\274\237\347\262\211\344\270\235\350\247\243\351\207\212\344\272\206\346\227\240\346\225\260\351\201\215\344\273\245\345\211\215\351\202\243\344\270\252\345\245\263\345\233\242\344\270\200\344\270\252\346\234\210\345\260\261\346\262\241\345\255\246\345\225\245\350\210\236\357\274\214\347\277\273\346\235\245\350\246\206\345\216\273\350\267\263\357\274\214\345\217\202\345\212\240101\347\232\204\346\227\266\345\200\2313\345\244\251\344\270\200\344\270\252\350\210\236\350\277\230\345\234\250\350\277\231\350\257\264\343\200\202\346\234\215\344\272\206" 
		id:1075824609 oid:31813039 mid:76839481 `
				for {
					log.Infov(context.Background(),
						log.KV("user", "test_user"),
						log.KV("ip", "127.0.0.1:8080"),
						log.KV("path", "127.0.0.1:8080/test_user"),
						log.KV("ret", 0),
						log.KV("args", arg),
						log.KV("stack", "nil"),
						log.KV("error", ""),
						log.KV("ts", time.Second.Seconds()),
					)
				}
			}()
		}
		time.Sleep(time.Hour)
	} else {
		startServer()
	}
}

func startServer() {
	go calcuQPS()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		go serve(conn)
	}
}

var total int64

func calcuQPS() {
	var lastTotal int64
	var lastUpdated time.Time
	for {
		time.Sleep(time.Second * 3)
		t := atomic.LoadInt64(&total)
		n := time.Now()
		change := t - lastTotal
		gap := n.Sub(lastUpdated)
		fmt.Println("bw:", change/(gap.Nanoseconds()/1e7))
		lastTotal = t
		lastUpdated = n
	}
}

func serve(conn net.Conn) {
	p := make([]byte, 4096)
	for {
		nRead, err := conn.Read(p)
		if err != nil {
			conn.Close()
			break
		}
		atomic.AddInt64(&total, int64(nRead))
	}
	fmt.Println(string(p))
}

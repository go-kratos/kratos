package aqm

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func init() {
	log.Init(nil)
}

func TestAQM(t *testing.T) {
	var group sync.WaitGroup
	rand.Seed(time.Now().Unix())
	eng := bm.Default()
	router := eng.Use(New(nil).Limit())
	router.GET("/aqm", testaqm)
	go eng.Run(":9999")
	var errcount int64
	for i := 0; i < 100; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for j := 0; j < 300; j++ {
				_, err := http.Get("http://127.0.0.1:9999/aqm")
				if err != nil {
					atomic.AddInt64(&errcount, 1)
				}
			}
		}()
	}
	group.Wait()
	fmt.Println("errcount", errcount)
}

func testaqm(ctx *bm.Context) {
	count := rand.Intn(100)
	time.Sleep(time.Millisecond * time.Duration(count))
}

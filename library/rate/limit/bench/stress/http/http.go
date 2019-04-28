package http

import (
	"hash/crc32"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/rate"
	"go-common/library/rate/limit"
	"go-common/library/rate/limit/bench/stress/conf"
	"go-common/library/rate/limit/bench/stress/service"
	"go-common/library/rate/vegas"
)

var (
	svc *service.Service
	req int64
	qps int64
)

// Init init
func Init(c *conf.Config) {
	rand.Seed(time.Now().Unix())
	initService(c)
	// init router
	engineInner := bm.DefaultServer(c.BM.Inner)
	outerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
	engineLocal := bm.DefaultServer(c.BM.Local)
	localRouter(engineLocal)
	if err := engineLocal.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
	if log.V(1) {
		go calcuQPS()
	}
}

// initService init services.
func initService(c *conf.Config) {
	//	idfSvc = identify.New(c.Identify)
	svc = service.New(c)
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	v := vegas.New()
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		defer ticker.Stop()
		for {
			<-ticker.C
			m := v.Stat()
			log.Info("vegas: limit(%d) inFlight(%d) minRtt(%v) rtt(%v)", m.Limit, m.InFlight, m.MinRTT, m.LastRTT)
		}
	}()
	l := limit.New(nil)
	//init api
	e.GET("/monitor/ping", ping)
	group := e.Group("/stress")
	group.GET("/normal", aqmTest)
	group.GET("/vegas", func(c *bm.Context) {
		start := time.Now()
		done, success := v.Acquire()
		if !success {
			done(time.Time{}, rate.Ignore)
			c.AbortWithStatus(509)
			return
		}
		defer done(start, rate.Success)
		c.Next()
	}, aqmTest)
	group.GET("/attack", func(c *bm.Context) {
		done, err := l.Allow(c)
		defer done(rate.Success)
		if err != nil {
			c.AbortWithStatus(509)
			return
		}
		c.Next()
	}, aqmTest)

}

func calcuQPS() {
	var creq, breq int64
	for {
		time.Sleep(time.Second * 5)
		creq = atomic.LoadInt64(&req)
		delta := creq - breq
		atomic.StoreInt64(&qps, delta/5)
		breq = creq
		log.Info("HTTP QPS:%d", atomic.LoadInt64(&qps))
	}

}
func aqmTest(c *bm.Context) {
	params := c.Request.Form
	sleep, err := strconv.ParseInt(params.Get("sleep"), 10, 64)
	if err == nil {
		time.Sleep(time.Millisecond * time.Duration(sleep))
	}
	atomic.AddInt64(&req, 1)
	for i := 0; i < 3000+rand.Intn(3000); i++ {
		crc32.Checksum([]byte(`testasdwfwfsddsfgwddcscsc
			http://git.bilibili.co/platform/go-common/merge_requests/new?merge_request%5Bsource_branch%5D=stress%2Fcodel`), crc32.IEEETable)
	}
}

// ping check server ok.
func ping(c *bm.Context) {
}

// innerRouter init local router api path.
func localRouter(e *bm.Engine) {
	//init api
	e.GET("/monitor/ping", ping)
	group := e.Group("/x/main/stress")
	{
		group.GET("", aqmTest)
	}
}

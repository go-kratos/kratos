package blademaster

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"

	"github.com/go-kratos/kratos/pkg/conf/dsn"

	"github.com/pkg/errors"
)

var (
	_perfOnce sync.Once
	_perfDSN  string
)

func init() {
	v := os.Getenv("HTTP_PERF")
	flag.StringVar(&_perfDSN, "http.perf", v, "listen http perf dsn, or use HTTP_PERF env variable.")
}

func startPerf(engine *Engine) {
	_perfOnce.Do(func() {
		if os.Getenv("HTTP_PERF") == "" {
			prefixRouter := engine.Group("/debug/pprof")
			{
				prefixRouter.GET("/", pprofHandler(pprof.Index))
				prefixRouter.GET("/cmdline", pprofHandler(pprof.Cmdline))
				prefixRouter.GET("/profile", pprofHandler(pprof.Profile))
				prefixRouter.POST("/symbol", pprofHandler(pprof.Symbol))
				prefixRouter.GET("/symbol", pprofHandler(pprof.Symbol))
				prefixRouter.GET("/trace", pprofHandler(pprof.Trace))
				prefixRouter.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
				prefixRouter.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
				prefixRouter.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
				prefixRouter.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
				prefixRouter.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
				prefixRouter.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
			}
			return
		}

		go func() {
			d, err := dsn.Parse(_perfDSN)
			if err != nil {
				panic(errors.Errorf("blademaster: http perf dsn must be tcp://$host:port, %s:error(%v)", _perfDSN, err))
			}
			if err := http.ListenAndServe(d.Host, nil); err != nil {
				panic(errors.Errorf("blademaster: listen %s: error(%v)", d.Host, err))
			}
		}()
	})
}

func pprofHandler(h http.HandlerFunc) HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

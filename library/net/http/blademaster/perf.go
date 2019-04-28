package blademaster

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"

	"go-common/library/conf/dsn"

	"github.com/pkg/errors"
)

var (
	_perfOnce sync.Once
	_perfDSN  string
)

func init() {
	v := os.Getenv("HTTP_PERF")
	if v == "" {
		v = "tcp://0.0.0.0:2333"
	}
	flag.StringVar(&_perfDSN, "http.perf", v, "listen http perf dsn, or use HTTP_PERF env variable.")
}

func startPerf() {
	_perfOnce.Do(func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

		go func() {
			d, err := dsn.Parse(_perfDSN)
			if err != nil {
				panic(errors.Errorf("blademaster: http perf dsn must be tcp://$host:port, %s:error(%v)", _perfDSN, err))
			}
			if err := http.ListenAndServe(d.Host, mux); err != nil {
				panic(errors.Errorf("blademaster: listen %s: error(%v)", d.Host, err))
			}
		}()
	})
}

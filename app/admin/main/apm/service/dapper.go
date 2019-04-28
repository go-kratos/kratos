package service

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	_apmDapperPrefix = "/x/admin/apm/dapper"
	_dapperPrxfix    = "/x/internal/dapper"
)

type dapperProxy struct {
	*httputil.ReverseProxy
}

func newDapperProxy(host string) (*dapperProxy, error) {
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	target, err := url.Parse(host)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid dapper host %s", host)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(r *http.Request) {
		r.Host = target.Host
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.URL.Path = strings.Replace(r.URL.Path, _apmDapperPrefix, _dapperPrxfix, 1)
	}
	return &dapperProxy{ReverseProxy: proxy}, nil
}

// DapperProxy dapper proxy.
func (s *Service) DapperProxy(rw http.ResponseWriter, r *http.Request) {
	s.dapperProxy.ServeHTTP(rw, r)
}

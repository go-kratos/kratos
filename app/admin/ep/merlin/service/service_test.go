package service

import (
	"context"
	"flag"
	"path/filepath"
	"strings"
	"time"

	"go-common/app/admin/ep/merlin/conf"

	"gopkg.in/h2non/gock.v1"
)

var (
	s *Service
	c context.Context
)

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	c = context.Background()
	time.Sleep(time.Second)
	s.dao.SetProxy()
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

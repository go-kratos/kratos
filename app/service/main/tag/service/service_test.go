package service

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/service/main/tag/conf"
	"go-common/library/cache/redis"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ps          = 1
	pn          = 1
	mid   int64 = 1
	tid   int64 = 3
	ip          = ""
	tids        = []int64{1, 2, 3}
	order       = 1
	svr   *Service
)

func CleanCache() {
	c := context.Background()
	pool := redis.NewPool(conf.Conf.Redis.Config)
	pool.Get(c).Do("FLUSHDB")
}

func TestMain(m *testing.M) {
	dir, _ := filepath.Abs("../cmd/tag.toml")
	flag.Set("conf", dir)
	conf.Init()
	fmt.Println(conf.Conf)
	log.Init(conf.Conf.Log)
	svr = New(conf.Conf)
	os.Exit(m.Run())
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(svr)
	}
}

func Test_Ranking(t *testing.T) {
	Convey("Test_Ranking    service ", t, WithService(func(s *Service) {
		var (
			c         = context.Background()
			rid int64 = 2090
		)
		s.RankingHot(c)
		s.RankingBangumi(c)
		s.RankingRegion(c, rid)
	}))
}

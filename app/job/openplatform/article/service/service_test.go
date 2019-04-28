package service

import (
	"flag"
	"path/filepath"
	"sync"
	"time"

	"go-common/app/job/openplatform/article/conf"
	"go-common/app/job/openplatform/article/dao"
	"go-common/library/queue/databus"
)

// func WithService(f func(s *Service)) func() {
// 	return func() {
// 		dir, _ := filepath.Abs("../goconvey.toml")
// 		flag.Set("conf", dir)
// 		conf.Init()
// 		s := New(conf.Conf)
// 		// s.dao = dao.New(conf.Conf)
// 		f(s)
// 	}
// }

func WithoutProcService(f func(s *Service)) func() {
	return func() {
		dir, _ := filepath.Abs("../goconvey.toml")
		flag.Set("conf", dir)
		conf.Init()
		s := &Service{
			c:      conf.Conf,
			dao:    dao.New(conf.Conf),
			waiter: new(sync.WaitGroup),
			// articleRPC:       artrpc.New(conf.Conf.ArticleRPC),
			articleSub:       databus.New(conf.Conf.ArticleSub),
			articleStatSub:   databus.New(conf.Conf.ArticleStatSub),
			updateDbInterval: int64(time.Duration(conf.Conf.Job.UpdateDbInterval) / time.Second),
		}
		f(s)
	}
}

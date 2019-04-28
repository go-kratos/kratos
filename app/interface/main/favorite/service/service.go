package service

import (
	"context"

	"go-common/app/interface/main/favorite/conf"
	musicDao "go-common/app/interface/main/favorite/dao/music"
	topicDao "go-common/app/interface/main/favorite/dao/topic"
	videoDao "go-common/app/interface/main/favorite/dao/video"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	favpb "go-common/app/service/main/favorite/api"
	"go-common/library/log"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

// Service define fav service
type Service struct {
	conf *conf.Config
	// dao
	videoDao *videoDao.Dao
	topicDao *topicDao.Dao
	musicDao *musicDao.Dao
	// cache chan
	cache *fanout.Fanout
	// prom
	prom *prom.Prom
	// rpc
	favClient favpb.FavoriteClient
	arcRPC    *arcrpc.Service2
}

// New return fav service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf: c,
		// dao
		videoDao: videoDao.New(c),
		topicDao: topicDao.New(c),
		musicDao: musicDao.New(c),
		// cache
		cache: fanout.New("cache"),
		// prom
		prom: prom.New().WithTimer("fav_add_video", []string{"method"}),
		// rpc
		arcRPC: arcrpc.New2(c.RPCClient2.Archive),
	}
	favClient, err := favpb.New(c.RPCClient2.FavClient)
	if err != nil {
		panic(err)
	}
	s.favClient = favClient
	return
}

// Ping check service health
func (s *Service) Ping(c context.Context) (err error) {
	return s.videoDao.Ping(c)
}

// Close close service
func (s *Service) Close() {
	s.videoDao.Close()
}

// PromError stat and log.
func (s *Service) PromError(name string, format string, args ...interface{}) {
	prom.BusinessErrCount.Incr(name)
	log.Error(format, args...)
}

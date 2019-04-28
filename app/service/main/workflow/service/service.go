package service

import (
	"context"
	"time"

	archive "go-common/app/service/main/archive/api"
	"go-common/app/service/main/workflow/conf"
	"go-common/app/service/main/workflow/dao"
	"go-common/app/service/main/workflow/dao/sobot"
	"go-common/app/service/main/workflow/model"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"

	"github.com/pkg/errors"
)

// Service is service.
type Service struct {
	c *conf.Config
	// dao
	dao   *dao.Dao
	sobot *sobot.Dao
	// tags cache
	tagsCache *model.TagsCache
	arcClient archive.ArchiveClient //archive rpc client
	cache     *fanout.Fanout
}

// New is videoup-admin service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		dao:   dao.New(c),
		sobot: sobot.New(c),
		tagsCache: &model.TagsCache{
			TagMap:   make(map[int8]map[int32]*model.Tag),
			TagSlice: make(map[int8][]*model.Tag),
			TagMap3:  make(map[int64]map[int64][]*model.Tag3),
		},
		cache: fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
	}
	var err error
	if s.arcClient, err = archive.NewClient(c.ArchiveRPC); err != nil {
		panic(errors.Wrap(err, "archive.NewClient failed"))
	}
	// load cache
	go s.cacheproc()
	return
}

// cacheproc goroutine
func (s *Service) cacheproc() {
	for {
		s.loadTags()
		s.loadTags3()
		time.Sleep(5 * time.Minute)
	}
}

// Close  consumer close.
func (s *Service) Close() {
	s.dao.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.Ping(c); err != nil {
		log.Error("Ping() err(%v)", err)
		return
	}
	return
}

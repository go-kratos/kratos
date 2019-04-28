package v1

import (
	"context"
	v1pb "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/library/net/metadata"

	expModel "go-common/app/service/live/xuser/model/exp"

	"go-common/library/log"
)

func (s *UserExpService) asyncCLearExpCache(ctx context.Context, req *v1pb.UserExpChunk) (err error) {
	err = s.dao.DelExpFromMemCache(ctx, req.Uid)
	if err != nil {
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_ASYNC_CACHE_FAIL3+"|clear UserExp cache(%d) error(%v)", req.Uid, err)
	}
	return
}

func (s *UserExpService) asyncSetExpCache(ctx context.Context, req map[int64]*expModel.LevelInfo) error {
	c := metadata.WithContext(ctx)
	f := func(c context.Context) {
		if err := s.dao.SetExpListCache(c, req); err != nil {
			log.Error(_errorServiceLogPrefix+"|"+_ERROR_ASYNC_CACHE_FAIL+"|asyncSetExpCache|error(%v),missedUIDs(%v)", err)
		}
	}
	if runErr := s.addExpCache.Save(func() {
		f(c)
	}); runErr != nil {
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_ASYNC_CACHE_FAIL2+"|asyncSetExpCache|error(%v),run cache is full(%v)", runErr)
		f(c)
	}
	return nil
}

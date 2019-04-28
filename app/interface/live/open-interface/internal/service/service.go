package service

import (
	"context"

	"go-common/app/interface/live/open-interface/internal/dao"
	titansSdk "go-common/app/service/live/resource/sdk"
	"go-common/library/conf/paladin"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao *dao.Dao
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  ac,
		dao: dao.New(),
	}

	dao.InitGrpc()
	InitTitan()
	return s
}

//InitTitan 初始化kv配置
func InitTitan() {
	conf := &titansSdk.Config{
		TreeId: 82686,
		Expire: 1,
	}
	titansSdk.Init(conf)
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}

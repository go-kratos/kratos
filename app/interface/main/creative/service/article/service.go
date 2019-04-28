package article

import (
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/activity"
	"go-common/app/interface/main/creative/dao/article"
	"go-common/app/interface/main/creative/dao/bfs"
	"go-common/app/interface/main/creative/service"
)

//Service struct.
type Service struct {
	c   *conf.Config
	art *article.Dao
	acc *account.Dao
	bfs *bfs.Dao
	act *activity.Dao
}

//New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:   c,
		art: rpcdaos.Art,
		acc: rpcdaos.Acc,
		bfs: bfs.New(c),
		act: activity.New(c),
	}
	return s
}

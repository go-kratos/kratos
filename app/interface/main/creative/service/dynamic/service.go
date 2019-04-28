package dynamic

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/lottery"
	lMdl "go-common/app/interface/main/creative/model/lottery"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

//Service struct
type Service struct {
	c       *conf.Config
	lottery *lottery.Dao
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:       c,
		lottery: lottery.New(c),
	}
	return s
}

// LotteryUserCheck  fn
func (s *Service) LotteryUserCheck(c context.Context, mid int64) (check bool, err error) {
	var (
		ret int
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	if ret, err = s.lottery.UserCheck(c, mid, ip); err != nil {
		log.Error("s.lottery.UserCheck err(%v) | mid(%d), ip(%s)", err, mid, ip)
		return
	}
	if ret == 1 {
		check = true
	}
	return
}

// LotteryNotice  fn
func (s *Service) LotteryNotice(c context.Context, aid, mid int64) (bind bool, err error) {
	var ret *lMdl.Notice
	if ret, err = s.lottery.Notice(c, aid, mid, metadata.String(c, metadata.RemoteIP)); err != nil {
		log.Error("s.lottery.Notice err(%v) |aid(%d)|mid(%d)", err, aid, mid)
		return
	}
	if ret != nil && ret.BizID > 0 {
		bind = true
	}
	return
}

package medal

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/medal"
	mdMdl "go-common/app/interface/main/creative/model/medal"
	"go-common/app/interface/main/creative/model/newcomer"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
)

//Service struct
type Service struct {
	c     *conf.Config
	medal *medal.Dao
	//task
	p *service.Public
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos, p *service.Public) *Service {
	s := &Service{
		c:     c,
		medal: medal.New(c),
		p:     p,
	}
	return s
}

// Medal get medal.
func (s *Service) Medal(c context.Context, mid int64) (t *mdMdl.Medal, err error) {
	t, err = s.medal.Medal(c, mid)
	return
}

// OpenMedal open medal.
func (s *Service) OpenMedal(c context.Context, mid int64, name string) (err error) {
	if err = s.medal.OpenMedal(c, mid, name); err != nil {
		log.Error("s.medal.OpenMedal mid(%d) error(%v)", mid, err)
		return
	}
	s.p.TaskPub(mid, newcomer.MsgForOpenFansMedal, newcomer.MsgFinishedCount)
	return
}

// RecentFans get recent fans.
func (s *Service) RecentFans(c context.Context, mid int64) (res []*mdMdl.RecentFans, err error) {
	res, err = s.medal.RecentFans(c, mid)
	return
}

// CheckMedal check medal name valid
func (s *Service) CheckMedal(c context.Context, mid int64, name string) (valid int, err error) {
	valid, err = s.medal.CheckMedal(c, mid, name)
	return
}

// Rank get fans rank list
func (s *Service) Rank(c context.Context, mid int64) (rank []*mdMdl.FansRank, err error) {
	rank, err = s.medal.Rank(c, mid)
	return
}

// Rename rename medal name
func (s *Service) Rename(c context.Context, mid int64, name, ak, ck string) (err error) {
	err = s.medal.Rename(c, mid, name, ak, ck)
	return
}

package elec

import (
	"context"
	"go-common/library/log"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/elec"
	elecMdl "go-common/app/interface/main/creative/model/elec"
	"go-common/app/interface/main/creative/service"
)

//Service struct.
type Service struct {
	c    *conf.Config
	elec *elec.Dao
	acc  *account.Dao
	arc  *archive.Dao
}

//New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:    c,
		elec: elec.New(c),
		acc:  rpcdaos.Acc,
		arc:  rpcdaos.Arc,
	}
	return s
}

// UserState get user elec state.
func (s *Service) UserState(c context.Context, mid int64, ip, ak, ck string) (data *elecMdl.UserState, err error) {
	data, err = s.elec.UserState(c, mid, ip)
	return
}

// ArchiveState get arc elec state.
func (s *Service) ArchiveState(c context.Context, aid, mid int64, ip string) (data *elecMdl.ArcState, err error) {
	data, err = s.elec.ArchiveState(c, aid, mid, ip)
	return
}

// CheckIsFriend check paymid state.
func (s *Service) CheckIsFriend(c context.Context, data []*elecMdl.Rank, mid int64, ip string) (res []*elecMdl.Rank, err error) {
	var mids []int64
	for _, v := range data {
		mids = append(mids, v.PayMID)
	}
	richRel, err := s.acc.RichRelation(c, mid, mids, ip)
	if err != nil {
		log.Error("s.acc.RichRelation error(%d, %v)", mid, err)
		return
	}
	if len(richRel) > 0 {
		for _, v := range data {
			if richRel[v.PayMID] == 3 || richRel[v.PayMID] == 4 {
				v.IsFriend = true
			} else {
				v.IsFriend = false
			}
		}
	}
	res = data
	return
}

package usersuit

import (
	"context"
	"sync"
	"time"

	"go-common/app/interface/main/account/model"
	accmdl "go-common/app/service/main/account/model"
	usmdl "go-common/app/service/main/usersuit/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_batch            = 20
	_fetchInfoTimeout = time.Second * 10
)

var (
	_emptyRichInvites = make([]*model.RichInvite, 0)
	_emptyInfoMap     = make(map[int64]*accmdl.Info)
)

// Buy buy invite code.
func (s *Service) Buy(c context.Context, mid int64, num int64) (res []*model.RichInvite, err error) {
	var invs []*usmdl.Invite
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &usmdl.ArgBuy{Mid: mid, Num: num, IP: ip}
	if invs, err = s.usRPC.Buy(c, arg); err != nil {
		log.Error("service.userserviceRPC.Buy(%v) error(%v)", arg, err)
		return
	}
	res = make([]*model.RichInvite, 0)
	for _, inv := range invs {
		res = append(res, model.NewRichInvite(inv, nil))
	}
	return
}

// Apply apply invite code.
func (s *Service) Apply(c context.Context, mid int64, code string, cookie string) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &usmdl.ArgApply{Mid: mid, Code: code, Cookie: cookie, IP: ip}
	if err = s.usRPC.Apply(c, arg); err != nil {
		log.Error("service.userserviceRPC.Apply(%v) error(%v)", arg, err)
	}
	return
}

// Stat get user's invite code stat.
func (s *Service) Stat(c context.Context, mid int64) (res *model.RichInviteStat, err error) {
	var st *usmdl.InviteStat
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &usmdl.ArgStat{Mid: mid, IP: ip}
	if st, err = s.usRPC.Stat(c, arg); err != nil {
		log.Error("service.userserviceRPC.Stat(%v) error(%v)", arg, err)
		return
	}
	res = &model.RichInviteStat{
		Mid:           st.Mid,
		CurrentLimit:  st.CurrentLimit,
		CurrentBought: st.CurrentBought,
		TotalBought:   st.TotalBought,
		TotalUsed:     st.TotalUsed,
		InviteCodes:   s.fillInviteeInfo(c, st.InviteCodes, ip),
	}
	return
}

func (s *Service) fillInviteeInfo(c context.Context, invs []*usmdl.Invite, ip string) []*model.RichInvite {
	if len(invs) == 0 {
		return _emptyRichInvites
	}
	imidm := make(map[int64]int)
	for _, inv := range invs {
		if inv.Status == usmdl.StatusUsed {
			imidm[inv.Imid] = 1
		}
	}
	infom := _emptyInfoMap
	if len(imidm) > 0 {
		imids := make([]int64, 0, len(imidm))
		for imid := range imidm {
			imids = append(imids, imid)
		}
		var err1 error
		if infom, err1 = s.fetchInfos(c, imids, ip, _fetchInfoTimeout); err1 != nil {
			log.Error("service.fetchInfos(%v, %s, %v) error(%v)", imids, ip, _fetchInfoTimeout, err1)
		}
	}
	rinvs := make([]*model.RichInvite, 0)
	for _, inv := range invs {
		rinvs = append(rinvs, model.NewRichInvite(inv, infom[inv.Imid]))
	}
	return rinvs
}

func (s *Service) fetchInfos(c context.Context, mids []int64, ip string, timeout time.Duration) (res map[int64]*accmdl.Info, err error) {
	if len(mids) == 0 {
		res = _emptyInfoMap
		return
	}
	batches := len(mids)/_batch + 1
	tc, cancel := context.WithTimeout(c, timeout)
	defer cancel()
	eg, errCtx := errgroup.WithContext(tc)
	bms := make([]map[int64]*accmdl.Info, batches)
	mu := sync.Mutex{}
	for i := 0; i < batches; i++ {
		idx := i
		end := (idx + 1) * _batch
		if idx == batches-1 {
			end = len(mids)
		}
		ids := mids[idx*_batch : end]
		eg.Go(func() error {
			m, err1 := s.accRPC.Infos3(errCtx, &accmdl.ArgMids{Mids: ids})
			mu.Lock()
			bms[idx] = m
			mu.Unlock()
			return err1
		})
	}
	err = eg.Wait()
	res = make(map[int64]*accmdl.Info)
	for _, bm := range bms {
		for mid, info := range bm {
			res[mid] = info
		}
	}
	return
}

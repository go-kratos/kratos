package relation

import (
	"context"
	"fmt"
	"hash/crc32"
	"sort"

	"go-common/app/interface/main/app-interface/model/live"
	model "go-common/app/interface/main/app-interface/model/relation"
	account "go-common/app/service/main/account/model"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var (
	_emptyFollowings = []*model.Following{}
	_emptyTagInfos   = []*model.Tag{}
)

// Followings get user followings list.
func (s *Service) Followings(c context.Context, vmid, mid int64, pn, ps int, version uint64, order string) (f []*model.Following, crc32v uint32, total int, err error) {
	res, err := s.relDao.Followings(c, vmid)
	if err != nil {
		return
	}
	total = len(res)
	if order == "asc" {
		sort.Sort(model.ByMTime(res))
	}
	start, end := (pn-1)*ps, pn*ps
	switch {
	case start >= len(res):
		res = res[:0]
	case end >= len(res):
		res = res[start:]
	default:
		res = res[start:end]
	}
	if len(res) == 0 {
		f = _emptyFollowings
		return
	}
	data := []byte(fmt.Sprintf("%s", res))
	crc32v = crc32.Checksum(data, crc32.IEEETable)
	if uint64(crc32v) == version {
		err = ecode.NotModified
		return
	}
	var fids []int64
	for _, r := range res {
		fids = append(fids, r.Mid)
	}
	var (
		flm   map[int64]*relation.Following
		cardm map[int64]*account.Card
		sm    map[int64]*live.Status
	)
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		cardm, err = s.accDao.Cards3(ctx, fids)
		return
	})
	// !self, compute !self user and up's followings' attr
	if mid != vmid && mid != 0 {
		g.Go(func() (err error) {
			if flm, err = s.relDao.Relations(ctx, mid, fids); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	g.Go(func() error {
		sm, _ = s.liveDao.StatusInfo(ctx, fids)
		return nil
	})
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	f = make([]*model.Following, 0, len(res))
	for _, r := range res {
		fl := &model.Following{Following: r}
		if _, ok := cardm[r.Mid]; ok {
			fl.Face = cardm[r.Mid].Face
			fl.Uname = cardm[r.Mid].Name
			fl.Sign = cardm[r.Mid].Sign
			if mid != vmid {
				if f, ok := flm[r.Mid]; ok {
					fl.Attribute = f.Attribute
				} else {
					fl.Attribute = relation.AttrNoRelation
				}
			}
			fl.OfficialVerify = cardm[r.Mid].Official
			fl.Vip.Type = int(cardm[r.Mid].Vip.Type)
			fl.Vip.VipStatus = int(cardm[r.Mid].Vip.Status)
			fl.Vip.DueDate = cardm[r.Mid].Vip.DueDate
			if status, ok := sm[r.Mid]; ok {
				fl.Live = status.LiveStatus
			}
		}
		f = append(f, fl)
	}
	return
}

func (s *Service) Tag(c context.Context, mid, tid int64, pn, ps int) (t []*model.Tag, err error) {
	mids, err := s.relDao.Tag(c, mid, tid)
	if err != nil {
		return
	}
	start, end := (pn-1)*ps, pn*ps
	switch {
	case start >= len(mids):
		mids = mids[:0]
	case end >= len(mids):
		mids = mids[start:]
	default:
		mids = mids[start:end]
	}
	if len(mids) == 0 {
		t = _emptyTagInfos
		return
	}
	var (
		cardm map[int64]*account.Card
		sm    map[int64]*live.Status
	)
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		cardm, err = s.accDao.Cards3(ctx, mids)
		return
	})
	g.Go(func() error {
		sm, _ = s.liveDao.StatusInfo(ctx, mids)
		return nil
	})
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}

	t = make([]*model.Tag, 0, len(mids))
	for _, mid = range mids {
		tag := &model.Tag{Mid: mid}
		if _, ok := cardm[mid]; ok {
			tag.Mid = mid
			tag.Face = cardm[mid].Face
			tag.Uname = cardm[mid].Name
			tag.Sign = cardm[mid].Sign
			tag.OfficialVerify = cardm[mid].Official
			tag.Vip.Type = int(cardm[mid].Vip.Type)
			tag.Vip.VipStatus = int(cardm[mid].Vip.Status)
			tag.Vip.DueDate = cardm[mid].Vip.DueDate
		}
		if status, ok := sm[mid]; ok {
			tag.Live = status.LiveStatus
		}
		t = append(t, tag)
	}
	return
}

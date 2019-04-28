package service

import (
	"context"
	"sort"

	"go-common/app/admin/main/creative/model/whitelist"
	accapi "go-common/app/service/main/account/api"
	"go-common/library/log"

	"golang.org/x/sync/errgroup"
)

// Cards fn
func (s *Service) Cards(c context.Context, wls []*whitelist.Whitelist) (wlsWithAcc []*whitelist.Whitelist, err error) {
	wlsWithAcc = []*whitelist.Whitelist{}
	var (
		g errgroup.Group
	)
	ch := make(chan *whitelist.Whitelist, len(wls))
	for _, wl := range wls {
		id := wl.ID
		mid := wl.MID
		adminMid := wl.AdminMID
		comment := wl.Comment
		state := wl.State
		tp := wl.Type
		ctime := wl.Ctime
		mtime := wl.Mtime
		g.Go(func() (err error) {
			pfl, err := s.dao.ProfileStat(c, mid)
			if err != nil {
				log.Error("s.dao.Card mid(%+v)|err(%+v)", mid, err)
				return
			}
			var name string
			if pfl.Profile != nil {
				name = pfl.Profile.Name
			}
			ch <- &whitelist.Whitelist{
				ID:           id,
				MID:          mid,
				AdminMID:     adminMid,
				Comment:      comment,
				State:        state,
				Type:         tp,
				Fans:         pfl.Follower,
				CurrentLevel: pfl.LevelInfo.Cur,
				Name:         name,
				Ctime:        ctime,
				Mtime:        mtime,
			}
			return
		})
	}
	g.Wait()
	close(ch)
	for c := range ch {
		wlsWithAcc = append(wlsWithAcc, c)
	}
	sort.Slice(wlsWithAcc, func(i, j int) bool { return wlsWithAcc[i].Ctime > wlsWithAcc[j].Ctime })
	return
}

// ProfileStat fn
func (s *Service) ProfileStat(c context.Context, mid int64) (pfl *accapi.ProfileStatReply, err error) {
	if pfl, err = s.dao.ProfileStat(c, mid); err != nil {
		log.Error("s.dao.Profile mid(%+v)|err(%+v)", mid, err)
	}
	return
}

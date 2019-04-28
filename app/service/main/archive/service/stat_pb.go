package service

import (
	"context"

	"go-common/app/service/main/archive/api"
)

// Stat3 get archive stat.
func (s *Service) Stat3(c context.Context, aid int64) (st *api.Stat, err error) {
	st, err = s.arc.Stat3(c, aid)
	return
}

// Click3 get archive click.
func (s *Service) Click3(c context.Context, aid int64) (clk *api.Click, err error) {
	clk, err = s.arc.Click3(c, aid)
	return
}

// Stats3 get archive stat.
func (s *Service) Stats3(c context.Context, aids []int64) (stm map[int64]*api.Stat, err error) {
	stm, err = s.arc.Stats3(c, aids)
	return
}

// SetStat set all stat
func (s *Service) SetStat(c context.Context, st *api.Stat) (err error) {
	s.arc.SetStat3(c, &api.Stat{
		Aid:     st.Aid,
		View:    int32(st.View),
		Danmaku: int32(st.Danmaku),
		Reply:   int32(st.Reply),
		Fav:     int32(st.Fav),
		Coin:    int32(st.Coin),
		Share:   int32(st.Share),
		NowRank: int32(st.NowRank),
		HisRank: int32(st.HisRank),
		Like:    int32(st.Like),
		DisLike: 0,
	})
	return
}

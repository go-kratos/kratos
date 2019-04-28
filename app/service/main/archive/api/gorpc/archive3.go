package gorpc

import (
	"context"

	"go-common/app/service/main/archive/api"
	model "go-common/app/service/main/archive/model/archive"
)

const (
	_archive3           = "RPC.Archive3"
	_archives3          = "RPC.Archives3"
	_view3              = "RPC.View3"
	_views3             = "RPC.Views3"
	_stat3              = "RPC.Stat3"
	_stats3             = "RPC.Stats3"
	_click3             = "RPC.Click3"
	_upArcs3            = "RPC.UpArcs3"
	_upsArcs3           = "RPC.UpsArcs3"
	_page3              = "RPC.Page3"
	_recommend3         = "RPC.Recommend3"
	_rankArcs3          = "RPC.RankArcs3"
	_ranksArcs3         = "RPC.RanksArcs3"
	_rankTopArcs3       = "RPC.RankTopArcs3"
	_rankAllArcs3       = "RPC.RankAllArcs3"
	_video3             = "RPC.Video3"
	_archivesWithPlayer = "RPC.ArchivesWithPlayer"
	_maxAid             = "RPC.MaxAID"
)

// MaxAID get max aid
func (s *Service2) MaxAID(c context.Context) (id int64, err error) {
	err = s.client.Call(c, _maxAid, _noArg, &id)
	return
}

// Archive3 Get receive aid, then init archive info.
func (s *Service2) Archive3(c context.Context, arg *model.ArgAid2) (res *api.Arc, err error) {
	res = new(api.Arc)
	err = s.client.Call(c, _archive3, arg, res)
	return
}

// Archives3 receive aids, then init archives info.
func (s *Service2) Archives3(c context.Context, arg *model.ArgAids2) (res map[int64]*api.Arc, err error) {
	err = s.client.Call(c, _archives3, arg, &res)
	return
}

// View3 get archive info and view pages.
func (s *Service2) View3(c context.Context, arg *model.ArgAid2) (res *model.View3, err error) {
	res = new(model.View3)
	err = s.client.Call(c, _view3, arg, res)
	return
}

// Views3 get archives info and view pages.
func (s *Service2) Views3(c context.Context, arg *model.ArgAids2) (res map[int64]*model.View3, err error) {
	err = s.client.Call(c, _views3, arg, &res)
	return
}

// Stat3 get archive stat
func (s *Service2) Stat3(c context.Context, arg *model.ArgAid2) (res *api.Stat, err error) {
	err = s.client.Call(c, _stat3, arg, &res)
	return
}

// ArchivesWithPlayer archives witch player
func (s *Service2) ArchivesWithPlayer(c context.Context, arg *model.ArgPlayer) (res map[int64]*model.ArchiveWithPlayer, err error) {
	err = s.client.Call(c, _archivesWithPlayer, arg, &res)
	return
}

// Stats3 get archive stat
func (s *Service2) Stats3(c context.Context, arg *model.ArgAids2) (res map[int64]*api.Stat, err error) {
	err = s.client.Call(c, _stats3, arg, &res)
	return
}

// Click3 get archive click
func (s *Service2) Click3(c context.Context, arg *model.ArgAid2) (res *api.Click, err error) {
	err = s.client.Call(c, _click3, arg, &res)
	return
}

// UpsArcs3 get archives of upper.
func (s *Service2) UpsArcs3(c context.Context, arg *model.ArgUpsArcs2) (res map[int64][]*api.Arc, err error) {
	err = s.client.Call(c, _upsArcs3, arg, &res)
	return
}

// UpArcs3 get archives of upper.
func (s *Service2) UpArcs3(c context.Context, arg *model.ArgUpArcs2) (res []*api.Arc, err error) {
	err = s.client.Call(c, _upArcs3, arg, &res)
	return
}

// Page3 get videos by aid
func (s *Service2) Page3(c context.Context, arg *model.ArgAid2) (res []*api.Page, err error) {
	err = s.client.Call(c, _page3, arg, &res)
	return
}

// Recommend3 from archive_recommend by aid
func (s *Service2) Recommend3(c context.Context, arg *model.ArgAid2) (res []*api.Arc, err error) {
	err = s.client.Call(c, _recommend3, arg, &res)
	return
}

// RankArcs3 get rank archives by type.
func (s *Service2) RankArcs3(c context.Context, arg *model.ArgRank2) (res *model.RankArchives3, err error) {
	res = new(model.RankArchives3)
	err = s.client.Call(c, _rankArcs3, arg, res)
	return
}

// RanksArcs3 get rank archives by types.
func (s *Service2) RanksArcs3(c context.Context, arg *model.ArgRanks2) (res map[int16]*model.RankArchives3, err error) {
	err = s.client.Call(c, _ranksArcs3, arg, &res)
	return
}

// RankTopArcs3 get top region archives by reid
func (s *Service2) RankTopArcs3(c context.Context, arg *model.ArgRankTop2) (res []*api.Arc, err error) {
	err = s.client.Call(c, _rankTopArcs3, arg, &res)
	return
}

// RankAllArcs3 get left 7 days all archives
func (s *Service2) RankAllArcs3(c context.Context, arg *model.ArgRankAll2) (res *model.RankArchives3, err error) {
	err = s.client.Call(c, _rankAllArcs3, arg, &res)
	return
}

// Video3 get video by aid & cid.
func (s *Service2) Video3(c context.Context, arg *model.ArgVideo2) (res *api.Page, err error) {
	res = new(api.Page)
	err = s.client.Call(c, _video3, arg, res)
	return
}

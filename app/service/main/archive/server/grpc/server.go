package grpc

import (
	"context"

	v1 "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/service"
	"go-common/library/net/rpc/warden"
)

type server struct {
	srv *service.Service
}

// New grpc server
func New(cfg *warden.ServerConfig, srv *service.Service) (wsvr *warden.Server, err error) {
	wsvr = warden.NewServer(cfg)
	v1.RegisterArchiveServer(wsvr.Server(), &server{srv: srv})
	wsvr, err = wsvr.Start()
	return
}

// Types get all types
func (s *server) Types(c context.Context, noArg *v1.NoArgRequest) (resp *v1.TypesReply, err error) {
	types := s.srv.AllTypes(c)
	resp = new(v1.TypesReply)
	resp.Types = make(map[int32]*v1.Tp)
	for _, tp := range types {
		resp.Types[int32(tp.ID)] = &v1.Tp{
			ID:   int32(tp.ID),
			Pid:  int32(tp.Pid),
			Name: tp.Name,
		}
	}
	return
}

// Arc get archive
func (s *server) Arc(c context.Context, req *v1.ArcRequest) (resp *v1.ArcReply, err error) {
	resp = new(v1.ArcReply)
	a, err := s.srv.Archive3(c, req.Aid)
	if err != nil {
		return
	}
	resp.Arc = a
	return
}

// Arcs get archives
func (s *server) Arcs(c context.Context, req *v1.ArcsRequest) (resp *v1.ArcsReply, err error) {
	resp = new(v1.ArcsReply)
	resp.Arcs = make(map[int64]*v1.Arc)
	as, err := s.srv.Archives3(c, req.Aids)
	if err != nil {
		return
	}
	if len(as) == 0 {
		return
	}
	for aid, a := range as {
		resp.Arcs[aid] = a
	}
	return
}

// View get archive and page
func (s *server) View(c context.Context, req *v1.ViewRequest) (resp *v1.ViewReply, err error) {
	resp = new(v1.ViewReply)
	v, err := s.srv.View3(c, req.Aid)
	if err != nil {
		return
	}
	resp = v
	return
}

// Views get archives and pages
func (s *server) Views(c context.Context, req *v1.ViewsRequest) (resp *v1.ViewsReply, err error) {
	resp = new(v1.ViewsReply)
	resp.Views = make(map[int64]*v1.ViewReply)
	vs, err := s.srv.Views3(c, req.Aids)
	if err != nil {
		return
	}
	if len(vs) == 0 {
		return
	}
	resp.Views = vs
	return
}

func (s *server) Stat(c context.Context, req *v1.StatRequest) (resp *v1.StatReply, err error) {
	resp = new(v1.StatReply)
	stat, err := s.srv.Stat3(c, req.Aid)
	if err != nil {
		return
	}
	resp.Stat = stat
	return
}

func (s *server) Stats(c context.Context, req *v1.StatsRequest) (resp *v1.StatsReply, err error) {
	resp = new(v1.StatsReply)
	resp.Stats = make(map[int64]*v1.Stat)
	stats, err := s.srv.Stats3(c, req.Aids)
	if err != nil {
		return
	}
	resp.Stats = stats
	return
}

func (s *server) Click(c context.Context, req *v1.ClickRequest) (resp *v1.ClickReply, err error) {
	resp = new(v1.ClickReply)
	resp.Click, err = s.srv.Click3(c, req.Aid)
	return
}

func (s *server) Page(c context.Context, req *v1.PageRequest) (resp *v1.PageReply, err error) {
	resp = new(v1.PageReply)
	resp.Pages, err = s.srv.Page3(c, req.Aid)
	return
}

func (s *server) Recommend(c context.Context, req *v1.RecommendRequest) (resp *v1.RecommendReply, err error) {
	resp = new(v1.RecommendReply)
	resp.Arcs, err = s.srv.UpperReommend(c, req.Aid)
	return
}

func (s *server) Video(c context.Context, req *v1.VideoRequest) (resp *v1.VideoReply, err error) {
	resp = new(v1.VideoReply)
	resp.Page, err = s.srv.Video3(c, req.Aid, req.Cid)
	return
}

func (s *server) MaxAid(c context.Context, req *v1.NoArgRequest) (resp *v1.MaxAidReply, err error) {
	resp = new(v1.MaxAidReply)
	resp.Aid, err = s.srv.MaxAID(c)
	return
}

func (s *server) ArcFieldCache(c context.Context, req *v1.ArcFieldCacheRequest) (resp *v1.NoReply, err error) {
	resp = new(v1.NoReply)
	err = s.srv.FieldCacheUpdate(c, req.Aid, int16(req.OldTypeID), int16(req.TypeID))
	return
}

func (s *server) ArcCache(c context.Context, req *v1.ArcCacheRequest) (resp *v1.NoReply, err error) {
	resp = new(v1.NoReply)
	err = s.srv.CacheUpdate(c, req.Aid, req.Tp, req.OldMid)
	return
}

func (s *server) DelVideoCache(c context.Context, req *v1.DelVideoCacheRequest) (resp *v1.NoReply, err error) {
	resp = new(v1.NoReply)
	err = s.srv.DelVideo(c, req.Aid, req.Cid)
	return
}

func (s *server) UpVideoCache(c context.Context, req *v1.UpVideoCacheRequest) (resp *v1.NoReply, err error) {
	resp = new(v1.NoReply)
	err = s.srv.UpVideo(c, req.Aid, req.Cid)
	return
}

func (s *server) Description(c context.Context, req *v1.DescriptionRequest) (resp *v1.DescriptionReply, err error) {
	resp = new(v1.DescriptionReply)
	resp.Desc, err = s.srv.Description(c, req.Aid)
	return
}

func (s *server) SetStat(c context.Context, req *v1.SetStatRequest) (resp *v1.NoReply, err error) {
	resp = new(v1.NoReply)
	err = s.srv.SetStat(c, req.Stat)
	return
}

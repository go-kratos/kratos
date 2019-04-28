package service

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/sync/errgroup"
)

// UpVideo update video by aid & cid
func (s *Service) UpVideo(c context.Context, aid, cid int64) (err error) {
	err = s.arc.UpVideo3(c, aid, cid)
	return
}

// DelVideo delete video by aid & cid
func (s *Service) DelVideo(c context.Context, aid, cid int64) (err error) {
	err = s.arc.DelVideoCache3(c, aid, cid)
	return
}

// Description update page by aid & cid
func (s *Service) Description(c context.Context, aid int64) (des string, err error) {
	des, err = s.arc.Description(c, aid)
	return
}

// Page3 get page by aid
func (s *Service) Page3(c context.Context, aid int64) (ps []*api.Page, err error) {
	ps, err = s.arc.Videos3(c, aid)
	return
}

// View3 get archive with video.
func (s *Service) View3(c context.Context, aid int64) (v *api.ViewReply, err error) {
	var (
		a  *api.Arc
		ps []*api.Page
	)
	if a, err = s.arc.Archive3(c, aid); err != nil {
		return
	}
	if ps, err = s.arc.Videos3(c, aid); err != nil {
		return
	}
	v = &api.ViewReply{Arc: a, Pages: ps}
	return
}

// Views3 get archives with videos.
func (s *Service) Views3(c context.Context, aids []int64) (vm map[int64]*api.ViewReply, err error) {
	if len(aids) == 0 {
		err = ecode.RequestErr
		return
	}
	var (
		am map[int64]*api.Arc
		vs map[int64][]*api.Page
	)
	var g errgroup.Group
	g.Go(func() (err error) {
		am, err = s.arc.Archives3(c, aids)
		return
	})
	g.Go(func() (err error) {
		vs, err = s.arc.VideosByAids3(c, aids)
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	vm = make(map[int64]*api.ViewReply, len(aids))
	for _, a := range am {
		vm[a.Aid] = &api.ViewReply{Arc: a, Pages: vs[a.Aid]}
	}
	return
}

// Video3 get video by aid & cid
func (s *Service) Video3(c context.Context, aid, cid int64) (video *api.Page, err error) {
	video, err = s.arc.Video3(c, aid, cid)
	return
}

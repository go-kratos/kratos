package ugc

import (
	"context"
	"fmt"

	ugcmdl "go-common/app/job/main/tv/model/ugc"
	arccli "go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

// VideoApi calls the rpc of video, pick videos of an archive
func (s *Service) videoPick(c context.Context, aid int64) (resp *arccli.ViewReply, err error) {
	if err = Retry(func() (err error) {
		if resp, err = s.arcClient.View(c, &arccli.ViewRequest{
			Aid: aid,
		}); err != nil {
			log.Error("ArcRPC For Aid: %d, Error: %v", aid, err)
		}
		return
	}, _arcRetry, _sleep); err != nil {
		log.Error("upArchives Error %+v", err)
		return
	}
	return
}

// arcAllowImport tells whether the archive is allowed to import into TV database
func (s *Service) arcAllowImport(arc *ugcmdl.ArcAllow) (allowed bool) {
	if !arc.CanPlay() {
		log.Warn("arcAllowImport Aid %d Not allowed Due to State %d", arc.Aid, arc.State)
		return
	}
	if arc.Ugcpay == arcmdl.AttrYes {
		log.Warn("arcAllowImport Aid %d Not allowed Due to Ugcpay %d", arc.Aid, arc.Ugcpay)
		return
	}
	if s.hitPGC(arc.Typeid) {
		log.Warn("arcAllowImport Aid %d Not allowed Due to HitPGC %d", arc.Aid, arc.Typeid)
		return
	}
	if !arc.IsOrigin() {
		log.Warn("arcAllowImport Aid %d Not allowed Due to Not Origin", arc.Aid, arc.Copyright)
	}
	allowed = true
	return
}

// Archive calls the api of Archive, pick the archive data
func (s *Service) arcPick(c context.Context, aid int64) (arc *arccli.Arc, err error) {
	var arcReply *arccli.ArcReply
	for i := 0; i < _arcRetry; i++ {
		if arcReply, err = s.arcClient.Arc(c, &arccli.ArcRequest{Aid: aid}); err == nil {
			break
		}
	}
	if err != nil {
		log.Error("upArchives Aid %d Error %v", aid, err)
		return
	}
	if arcReply == nil || arcReply.Arc == nil {
		err = ecode.NothingFound
		return
	}
	arc = arcReply.Arc
	return
}

// ArcCount counts the mid's archive, pick the number
func (s *Service) arcCount(mid int64) (count int, err error) {
	if err = Retry(func() (err error) {
		if count, err = s.arcRPC.UpCount2(ctx, &arcmdl.ArgUpCount2{
			Mid: mid,
		}); err != nil {
			log.Error("ArcCount For Mid: %d, Error: %v", mid, err)
		}
		return
	}, _arcRetry, _sleep); err != nil {
		log.Error("upArchives Error %+v", err)
	}
	return
}

// arcViews picks the views of the given page of aids
func (s *Service) arcViews(aids []int64) (res map[int64]*arccli.ViewReply, err error) {
	var resp *arccli.ViewsReply
	if err = Retry(func() (err error) {
		if resp, err = s.arcClient.Views(ctx, &arccli.ViewsRequest{
			Aids: aids,
		}); err != nil {
			log.Error("%+v", err)
		}
		return
	}, _arcRetry, _sleep); err != nil {
		log.Error("upArchives Error %+v", err)
		return
	} else if len(resp.Views) == 0 {
		err = fmt.Errorf("result empty")
		return
	}
	res = resp.Views
	return
}

// VideoApi calls the rpc of video, pick videos of an archive
func (s *Service) pagePick(c context.Context, cid int64, aid int64, ip string) (res *arccli.Page, err error) {
	if err = Retry(func() (err error) {
		if res, err = s.arcRPC.Video3(c, &arcmdl.ArgVideo2{
			Aid:    aid,
			Cid:    cid,
			RealIP: ip,
		}); err != nil {
			log.Error("ArcRPC For Aid: %d, Cid: %d, Error: %v", aid, cid, err)
		}
		return
	}, _arcRetry, _sleep); err != nil {
		log.Error("upArchives Error %+v", err)
		return
	} else if res == nil {
		err = fmt.Errorf("result empty")
	}
	return
}

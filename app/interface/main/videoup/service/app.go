package service

import (
	"bytes"
	"context"
	"net"
	"strings"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

// AppEdit edit archive by appclient.
func (s *Service) AppEdit(c context.Context, ap *archive.ArcParam, mid int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	ap.IPv6 = net.ParseIP(ip)
	if err = s.checkIdentify(c, mid, ip); err != nil {
		log.Error("s.CheckIdentify mid(%d) ap(%+v) error(%v)", mid, ap, err)
		return
	}
	var (
		a  = &archive.Archive{}
		vs = []*archive.Video{}
	)
	if a, vs, err = s.arc.View(c, ap.Aid, ip); err != nil {
		log.Error("s.arc.View err(%v) | aid(%d) ip(%s)", err, ap.Aid, ip)
		return
	}
	if a == nil {
		log.Error("s.arc.View(%d) not found", mid)
		err = ecode.ArchiveNotExist
		return
	}
	// pre check
	if err = s.preEdit(c, mid, a, vs, ap, ip, ap.UpFrom); err != nil {
		return
	}
	// edit
	if err = s.arc.Edit(c, ap, ip); err != nil {
		return
	}
	g := &errgroup.Group{}
	ctx := context.TODO()
	g.Go(func() error {
		s.dealElec(ctx, ap.OpenElec, ap.Aid, mid, ip)
		return nil
	})
	g.Wait()
	return
}

// AppUpCover main app upload cover.
func (s *Service) AppUpCover(c context.Context, fileType string, body []byte, mid int64) (url string, err error) {
	if len(body) == 0 {
		err = ecode.FileNotExists
		log.Error("AppEcode FileNotExists mid(%d) error(%v)", mid, err)
		return
	}
	if len(body) > s.c.Bfs.MaxFileSize {
		err = ecode.FileTooLarge
		log.Error("AppEcode FileTooLarge mid(%d) error(%v)", mid, err)
		return
	}
	url, err = s.bfs.Upload(c, fileType, bytes.NewReader(body))
	if err != nil {
		log.Error("AppEcode s.bfs.Upload error(%v)", err)
	}
	return
}

func (s *Service) freshAppMissionByFirstTag(ap *archive.ArcParam) (res *archive.ArcParam) {
	if ap.MissionID == 0 {
		firstTag := strings.Split(ap.Tag, ",")[0]
		if missionID, ok := s.missTagsCache[firstTag]; ok {
			ap.MissionID = missionID
		}
	}
	res = ap
	return
}

// AppAdd add archive by main app.
func (s *Service) AppAdd(c context.Context, mid int64, ap *archive.ArcParam, ar *archive.AppRequest) (aid int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	ap.IPv6 = net.ParseIP(ip)
	defer func() {
		if err != nil && err != ecode.VideoupCanotRepeat {
			s.acc.DelSubmitCache(c, ap.Mid, ap.Title)
		}
	}()
	ap = s.freshAppMissionByFirstTag(ap)
	if err = s.checkIdentify(c, mid, ip); err != nil {
		log.Error("s.CheckIdentify mid(%d) ap(%+v) error(%v)", mid, ap, err)
		return
	}
	// pre check
	if err = s.preAdd(c, mid, ap, ip, ap.UpFrom); err != nil {
		return
	}
	if ap.PoiObj != nil {
		log.Warn("poi_object is not nil, mid(%d),upfrom(%d),poi_object(%+v)", mid, ap.UpFrom, ap.PoiObj)
	}
	if aid, err = s.arc.Add(c, ap, ip); err != nil || aid == 0 {
		return
	}
	ap.Aid = aid
	g := &errgroup.Group{}
	ctx := context.TODO()
	g.Go(func() error {
		s.dealOrder(ctx, mid, aid, ap.OrderID, ip)
		return nil
	})
	g.Go(func() error {
		s.dealWaterMark(ctx, mid, ap.Watermark, ip)
		return nil
	})
	g.Go(func() error {
		s.freshFavs(ctx, mid, ap, ip)
		return nil
	})
	g.Go(func() error {
		s.dealElec(ctx, 1, aid, mid, ip)
		return nil
	})
	g.Go(func() error {
		s.uploadVideoEditInfo(ctx, ap, aid, mid, ip)
		return nil
	})
	g.Go(func() error {
		s.lotteryBind(ctx, ap.LotteryID, aid, mid, ip)
		return nil
	})
	g.Go(func() error {
		s.addFollowing(ctx, mid, ap.FollowMids, ap.UpFrom, ip)
		return nil
	})
	g.Go(func() error {
		s.VideoInfoc(ctx, ap, ar)
		return nil
	})
	g.Wait()
	return
}

// AppEditFull fn
func (s *Service) AppEditFull(c context.Context, ap *archive.ArcParam, mid, buildNum int64, ar *archive.AppRequest) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	ap.IPv6 = net.ParseIP(ip)
	platform := ar.Platform
	if err = s.checkIdentify(c, mid, ip); err != nil {
		log.Error("s.CheckIdentify mid(%d) ap(%+v) error(%v)", mid, ap, err)
		return
	}
	var (
		a  = &archive.Archive{}
		vs = []*archive.Video{}
	)
	if a, vs, err = s.arc.View(c, ap.Aid, ip); err != nil {
		log.Error("s.arc.View err(%v) | aid(%d) ip(%s)", err, ap.Aid, ip)
		return
	}
	if a == nil {
		log.Error("s.arc.View(%d) not found", mid)
		err = ecode.ArchiveNotExist
		return
	}
	if nvsCnt := s.checkVideosMaxLimitForEdit(vs, ap.Videos); nvsCnt > s.c.MaxAddVsCnt {
		log.Error("checkVideosMaxLimitForEdit, vsCnt(%d), limit(%d), nvsCnt(%d)", len(vs), s.c.MaxAddVsCnt, nvsCnt)
		err = ecode.VideoupVideosMaxLimit
		return
	}
	ap = s.protectFeatureForApp(ap, a, buildNum, platform)
	// pre check
	if err = s.preEdit(c, mid, a, vs, ap, ip, ap.UpFrom); err != nil {
		return
	}
	// edit
	if err = s.arc.Edit(c, ap, ip); err != nil {
		return
	}
	g := &errgroup.Group{}
	ctx := context.TODO()
	g.Go(func() error {
		s.dealElec(ctx, ap.OpenElec, ap.Aid, mid, ip)
		return nil
	})
	g.Go(func() error {
		s.uploadVideoEditInfo(ctx, ap, ap.Aid, mid, ip)
		return nil
	})
	g.Go(func() error {
		s.addFollowing(ctx, mid, ap.FollowMids, ap.UpFrom, ip)
		return nil
	})
	g.Go(func() error {
		s.VideoInfoc(ctx, ap, ar)
		return nil
	})
	g.Wait()
	return
}

// protectFeatureForApp fn
// feature list: porder,order,desc_format_id
func (s *Service) protectFeatureForApp(origin *archive.ArcParam, a *archive.Archive, buildNum int64, platform string) (res *archive.ArcParam) {
	res = origin
	res.Porder = a.Porder
	res.OrderID = a.OrderID
	res.DescFormatID = a.DescFormatID
	// android except 5.26 5260000
	if buildNum < 5260000 && platform == "android" {
		res.Dynamic = a.Dynamic
		res.MissionID = a.MissionID
		// ios include 5.25.1 6680
	} else if buildNum <= 6680 && platform == "ios" {
		res.Dynamic = a.Dynamic
		res.MissionID = a.MissionID
	}
	return
}

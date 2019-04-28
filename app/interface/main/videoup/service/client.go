package service

import (
	"bytes"
	"context"
	"net"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

// ClientAdd add archive by client.
func (s *Service) ClientAdd(c context.Context, mid int64, ap *archive.ArcParam) (aid int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	ap.IPv6 = net.ParseIP(ip)

	defer func() {
		if err != nil && err != ecode.VideoupCanotRepeat {
			s.acc.DelSubmitCache(c, ap.Mid, ap.Title)
		}
	}()
	if err = s.checkIdentify(c, mid, ip); err != nil {
		log.Error("s.CheckIdentify mid(%d) ap(%+v) error(%v)", mid, ap, err)
		return
	}
	if err = s.tagsCheck(c, mid, ap.Tag, ip); err != nil {
		log.Error("s.tagsCheck mid(%d) ap(%+v) error(%v)", mid, ap.Tag, err)
		return
	}
	// pre check
	if err = s.preAdd(c, mid, ap, ip, archive.UpFromWindows); err != nil {
		return
	}
	// add
	if aid, err = s.arc.Add(c, ap, ip); err != nil || aid == 0 {
		return
	}
	g := &errgroup.Group{}
	ctx := context.TODO()
	g.Go(func() error {
		s.dealOrder(ctx, mid, aid, ap.OrderID, ip)
		return nil
	})
	g.Go(func() error {
		s.freshFavs(ctx, mid, ap, ip)
		return nil
	})
	g.Go(func() error {
		s.dealElec(ctx, ap.OpenElec, aid, mid, ip)
		return nil
	})
	g.Wait()
	return
}

// ClientEdit edit archive by client.
func (s *Service) ClientEdit(c context.Context, ap *archive.ArcParam, mid int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	ap.IPv6 = net.ParseIP(ip)

	if err = s.checkIdentify(c, mid, ip); err != nil {
		log.Error("s.CheckIdentify mid(%d) ap(%+v) error(%v)", mid, ap, err)
		return
	}
	if err = s.tagsCheck(c, mid, ap.Tag, ip); err != nil {
		log.Error("s.tagsCheck mid(%d) ap(%+v) error(%v)", mid, ap.Tag, err)
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
	// pre check
	if err = s.preEdit(c, mid, a, vs, ap, ip, archive.UpFromWindows); err != nil {
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

// ClientUpCover client upload cover.
func (s *Service) ClientUpCover(c context.Context, fileType string, body []byte, mid int64) (url string, err error) {
	if len(body) == 0 {
		err = ecode.FileNotExists
		return
	}
	if len(body) > s.c.Bfs.MaxFileSize {
		err = ecode.FileTooLarge
		return
	}
	url, err = s.bfs.Upload(c, fileType, bytes.NewReader(body))
	if err != nil {
		log.Error("s.bfs.Upload error(%v)", err)
	}
	return
}

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

// WebAdd add archive by web.
func (s *Service) WebAdd(c context.Context, mid int64, ap *archive.ArcParam, validated bool) (aid int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	ap.IPv6 = net.ParseIP(ip)
	g := &errgroup.Group{}
	ctx := context.TODO()
	defer func() {
		// VideoupCanotRepeat high level but basic
		if err != nil && err != ecode.VideoupCanotRepeat && err != ecode.VideoupAddLimitHalfMin {
			g.Go(func() error {
				s.acc.DelSubmitCache(ctx, mid, ap.Title)
				return nil
			})
		}
		if err != nil && err != ecode.VideoupAddLimitHalfMin {
			g.Go(func() error {
				s.acc.DelHalfMin(ctx, mid)
				return nil
			})
		}
	}()
	if !validated && !s.allowHalfMin(c, mid) {
		log.Warn("VideoupAddLimitHalfMin mid(%d) ap(%+v) validated(%+v)", mid, ap, validated)
		err = ecode.VideoupAddLimitHalfMin
		return
	}
	if err = s.checkIdentify(c, mid, ip); err != nil {
		log.Error("s.CheckIdentify mid(%d) ap(%+v) error(%v)", mid, ap, err)
		return
	}
	if err = s.checkAddPay(c, ap, ip); err != nil {
		log.Error("s.checkAddPay mid(%d) ap(%+v) error(%v)", mid, ap, err)
		return
	}
	// 检查联合投稿
	if err = s.checkAddStaff(c, ap, mid, ip); err != nil {
		return
	}
	// pre check
	if err = s.preAdd(c, mid, ap, ip, archive.UpFromWeb); err != nil {
		return
	}
	// add
	if aid, err = s.arc.Add(c, ap, ip); err != nil || aid == 0 {
		return
	}
	g.Go(func() error {
		s.dealOrder(ctx, mid, aid, ap.OrderID, ip)
		return nil
	})
	g.Go(func() error {
		s.freshFavs(ctx, mid, ap, ip)
		return nil
	})
	g.Go(func() error {
		s.acc.AddHalfMin(ctx, mid)
		return nil
	})
	g.Go(func() error {
		s.dealElec(ctx, ap.OpenElec, aid, mid, ip)
		return nil
	})
	g.Go(func() error {
		s.dealSubtitle(ctx, ap.Subtitle, aid, mid, ip)
		return nil
	})
	// same to edit go func, 当且仅当付费设置有且开启
	g.Go(func() error {
		if ap.Pay != nil && ap.Pay.Open == 1 && !ap.Pay.RefuseUpdate {
			if err = s.dealAddPay(ctx, ap.Pay, aid, mid, ip); err != nil {
				//异步重试队列
				s.asyncCh <- func() error {
					return s.dealAddPay(ctx, ap.Pay, aid, mid, ip)
				}
			}
			err = nil
		}
		return nil
	})
	g.Wait()
	return
}

// WebEdit edit archive by web.
func (s *Service) WebEdit(c context.Context, ap *archive.ArcParam, mid int64) (err error) {
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
	if nvsCnt := s.checkVideosMaxLimitForEdit(vs, ap.Videos); nvsCnt > s.c.MaxAddVsCnt {
		log.Error("checkVideosMaxLimitForEdit, vsCnt(%d), limit(%d), nvsCnt(%d)", len(vs), s.c.MaxAddVsCnt, nvsCnt)
		err = ecode.VideoupVideosMaxLimit
		return
	}
	if err = s.checkEditStaff(c, ap, mid, a, ip); err != nil {
		return
	}
	// pre check
	if err = s.preEdit(c, mid, a, vs, ap, ip, archive.UpFromWeb); err != nil {
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
		s.dealSubtitle(ctx, ap.Subtitle, ap.Aid, mid, ip)
		return nil
	})
	g.Go(func() error {
		// web端60天之后无脑调价 && web端开放过，在60天内，不允许修改价格，但是可以修改稿件其他信息
		if ap.Pay != nil && ap.Pay.Open == 1 && !ap.Pay.RefuseUpdate {
			if err = s.dealAdjustPay(ctx, ap.Pay, ap.Aid, mid, ip); err != nil {
				//异步重试队列
				s.asyncCh <- func() error {
					return s.dealAdjustPay(ctx, ap.Pay, ap.Aid, mid, ip)
				}
			}
			err = nil
		}
		return nil
	})
	g.Wait()
	return
}

// WebUpCover client upload cover.
func (s *Service) WebUpCover(c context.Context, fileType string, body []byte, mid int64) (url string, err error) {
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

// WebCmAdd add archive by web client and from business order.
func (s *Service) WebCmAdd(c context.Context, mid int64, ap *archive.ArcParam) (aid int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	ap.IPv6 = net.ParseIP(ip)

	defer func() {
		if err != nil && err != ecode.VideoupCanotRepeat {
			s.acc.DelSubmitCache(c, ap.Mid, ap.Title)
		}
	}()
	if ap.TypeID != archive.AdvertisingTypeID {
		err = ecode.VideoupTypeidErr
		log.Error("ap.TypeID is not AdvertisingTypeID mid(%d),type(%d),err(%v) ", mid, ap.TypeID, err)
		return
	}
	if err = s.checkIdentify(c, mid, ip); err != nil {
		log.Error("s.CheckIdentify mid(%d) ap(%+v) error(%v)", mid, ap, err)
		return
	}
	// pre check
	if err = s.preAdd(c, mid, ap, ip, archive.UpFromCM); err != nil {
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
	g.Wait()
	return
}

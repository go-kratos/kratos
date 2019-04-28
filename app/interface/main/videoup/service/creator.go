package service

import (
	"bytes"
	"context"
	"go-common/library/net/metadata"
	"hash/crc32"
	"net"
	"strconv"
	"strings"

	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// CreatorEdit edit archive by creator.
func (s *Service) CreatorEdit(c context.Context, mid int64, cp *archive.CreatorParam) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)

	if err = s.checkIdentify(c, mid, ip); err != nil {
		log.Error("s.CheckIdentify mid(%d) ap(%+v) error(%v)", mid, cp, err)
		return
	}
	var (
		a  = &archive.Archive{}
		vs = make([]*archive.Video, 0)
	)
	if a, vs, err = s.arc.View(c, cp.Aid, ip); err != nil {
		log.Error("s.arc.View err(%v) | aid(%d) ip(%s)", err, cp.Aid, ip)
		return
	}
	if a == nil || vs == nil {
		log.Error("s.arc.View(%d) not found", mid)
		err = ecode.ArchiveNotExist
		return
	}
	ap := &archive.ArcParam{
		Aid:      cp.Aid,
		Tag:      cp.Tag,
		Title:    cp.Title,
		Desc:     cp.Desc,
		OpenElec: cp.OpenElec,
		// ------ diff values ----- //
		Mid:          a.Mid,
		Author:       a.Author,
		TypeID:       a.TypeID,
		Cover:        coverURL(a.Cover),
		Copyright:    a.Copyright,
		NoReprint:    a.NoReprint,
		OrderID:      a.OrderID,
		Source:       a.Source,
		Attribute:    a.Attribute,
		UpFrom:       archive.UpFromCreator,
		DTime:        a.DTime,
		DescFormatID: a.DescFormatID,
		Dynamic:      a.Dynamic,
		IPv6:         net.ParseIP(ip),
		MissionID:    int(a.MissionID),
	}
	for _, vp := range vs {
		ap.Videos = append(ap.Videos, &archive.VideoParam{
			Title:    vp.Title,
			Desc:     vp.Desc,
			Filename: vp.Filename,
		})
	}
	if only := onlyChangeTagArc(cp, a); only {
		ap.Tag = s.removeDupTag(ap.Tag)
		if !s.allowTag(ap.Tag) {
			log.Error("s.allowTag mid(%d) ap.Tag(%s) tag name or number too large or Empty", mid, ap.Tag)
			err = ecode.VideoupTagErr
			return
		}
		if err = s.tagsCheck(c, mid, ap.Tag, ip); err != nil {
			log.Error("s.tagsCheck mid(%d) ap(%+v) error(%v)", mid, ap.Tag, err)
			return
		}
		if a.Tag != ap.Tag {
			s.arc.TagUp(c, ap.Aid, ap.Tag, ip)
		}

		s.dealTag(c, mid, ap.Aid, a.Tag, ap.Tag, ip, ap.TypeID)
	} else {
		if err = s.preEdit(c, mid, a, vs, ap, ip, archive.UpFromCreator); err != nil {
			log.Error("s.preCreatorEdit mid(%d) ap(%+v) error(%v)", mid, ap, err)
			return
		}
		if err = s.arc.Edit(c, ap, ip); err != nil {
			return
		}
	}
	s.dealElec(c, ap.OpenElec, ap.Aid, mid, ip)
	return
}

// 判断是否只在这种特殊的情况下，开放浏览/待审的稿件只修改了Tag信息
func onlyChangeTagArc(cp *archive.CreatorParam, a *archive.Archive) (only bool) {
	st := a.State == archive.StateForbidSubmit ||
		a.State == archive.StateForbidUserDelay ||
		a.State == archive.StateOpen ||
		a.State == archive.StateOrange ||
		a.State == archive.StateForbidWait
	ch := a.Title == cp.Title &&
		a.Desc == cp.Desc &&
		a.Tag != cp.Tag

	if st && ch {
		only = true
	}
	return
}

// CreatorAdd add archive by creator.
func (s *Service) CreatorAdd(c context.Context, mid int64, ap *archive.ArcParam) (aid int64, err error) {
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
	// pre check
	if err = s.preAdd(c, mid, ap, ip, archive.UpFromCreator); err != nil {
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
		s.dealElec(ctx, ap.OpenElec, ap.Aid, mid, ip)
		return nil
	})
	g.Wait()
	return
}

// CreatorUpCover creator upload cover.
func (s *Service) CreatorUpCover(c context.Context, fileType string, body []byte, mid int64) (url string, err error) {
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

// coverURL convert cover url to full url.
func coverURL(uri string) (cover string) {
	if uri == "" {
		//cover = "http://static.hdslb.com/images/transparent.gif"
		return
	}
	cover = uri
	if strings.Index(uri, "http://") == 0 {
		return
	}
	if len(uri) >= 10 && uri[:10] == "/templets/" {
		return
	}
	if strings.HasPrefix(uri, "group1") {
		cover = "http://i0.hdslb.com/" + uri
		return
	}
	if pos := strings.Index(uri, "/uploads/"); pos != -1 && (pos == 0 || pos == 3) {
		cover = uri[pos+8:]
	}
	cover = strings.Replace(cover, "{IMG}", "", -1)
	cover = "http://i" + strconv.FormatInt(int64(crc32.ChecksumIEEE([]byte(cover)))%3, 10) + ".hdslb.com" + cover
	return
}

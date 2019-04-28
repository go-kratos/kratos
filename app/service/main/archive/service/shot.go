package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// Videoshot get video shot.
func (s *Service) Videoshot(c context.Context, aid, cid int64) (shot *archive.Videoshot, err error) {
	// check archive&video state
	var attr int32
	if aid == 0 {
		var vsm map[int64]map[int64]*api.Page
		if vsm, err = s.arc.VideosByCids(c, []int64{cid}); err != nil {
			err = errors.WithStack(err)
			return
		}
		if len(vsm) == 0 {
			err = ecode.VideoshotNotExist
			return
		}
		for aid := range vsm {
			var a *api.Arc
			if a, err = s.arc.Archive3(c, aid); err != nil {
				err = errors.WithStack(err)
				return
			}
			if !a.IsNormal() {
				err = ecode.VideoshotNotExist
				return
			}
		}
	} else {
		var a *api.Arc
		if a, err = s.arc.Archive3(c, aid); err != nil {
			err = errors.WithStack(err)
			return
		}
		if !a.IsNormal() {
			err = ecode.VideoshotNotExist
			return
		}
		attr = a.Attribute
	}
	// video shot
	v, err := s.shot.Videoshot(c, cid)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if v == nil || v.Count == 0 {
		err = ecode.VideoshotNotExist
		return
	}
	shot = &archive.Videoshot{
		XLen:  10,
		YLen:  10,
		XSize: 160,
		YSize: 90,
		Image: make([]string, 0, v.Count),
		Attr:  attr,
	}
	var sign = func(fn string, ver int) (uri string) {
		h := hmac.New(sha1.New, []byte(s.c.Videoshot.Key))
		h.Write([]byte(fn))
		uri = fmt.Sprintf("%s%s?vsign=%s&ver=%d", s.c.Videoshot.URI, fn, fmt.Sprintf("%x", h.Sum(nil)), ver)
		return
	}
	shot.PvData = sign(fmt.Sprintf("%d.bin", cid), v.Version())
	for i := 0; i < v.Count; i++ {
		if i == 0 {
			shot.Image = append(shot.Image, sign(fmt.Sprintf("%d.jpg", cid), v.Version()))
			continue
		}
		shot.Image = append(shot.Image, sign(fmt.Sprintf("%d-%d.jpg", cid, i), v.Version()))
	}
	return
}

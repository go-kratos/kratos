package service

import (
	"context"
	"time"

	"go-common/app/admin/main/growup/model"
)

// Banners Get Banners
func (s *Service) Banners(c context.Context, from, limit int64) (total int64, bs []*model.Banner, err error) {
	total, err = s.dao.TotalBannerCount(c)
	if err != nil {
		return
	}
	bs, err = s.dao.Banners(c, from, limit)
	return
}

// AddBanner Add Banner
func (s *Service) AddBanner(c context.Context, image, link string, startAt, endAt int64) (dup int64, err error) {
	if startAt > endAt {
		return
	}
	dup, err = s.dao.DupBanner(c, startAt, endAt, time.Now().Unix())
	if err != nil {
		return
	}
	if dup > 0 {
		return
	}
	_, err = s.dao.InsertBanner(c, image, link, startAt, endAt)
	return
}

// EditBanner Update Banner
func (s *Service) EditBanner(c context.Context, id, startAt, endAt int64, image, link string) (dup int64, err error) {
	if startAt > endAt {
		return
	}
	dup, err = s.dao.DupEditBanner(c, startAt, endAt, time.Now().Unix(), id)
	if err != nil {
		return
	}
	if dup > 0 {
		return
	}
	_, err = s.dao.UpdateBanner(c, image, link, startAt, endAt, id)
	return
}

// Off set end time to now
func (s *Service) Off(c context.Context, endAt, id int64) (err error) {
	_, err = s.dao.UpdateBannerEndAt(c, endAt, id)
	return
}

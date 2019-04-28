package service

import (
	"context"

	"go-common/app/admin/main/upload/model"
	"go-common/library/log"
)

// ListBucket .
func (s *Service) ListBucket(c context.Context, lbp *model.ListBucketParam) (bucketPage *model.BucketListPage, err error) {
	var (
		buckets []*model.Bucket
		count   int
	)
	if err = s.orm.Table("bucket").Order("id desc").Limit(lbp.PS).Offset((lbp.PN - 1) * lbp.PS).Find(&buckets).Error; err != nil {
		log.Error("read bucket error(%v)", err)
		return
	}
	if err = s.orm.Table("bucket").Count(&count).Error; err != nil {
		log.Error("bucket count error(%v)", err)
		err = nil
	}
	bucketPage = &model.BucketListPage{
		Items: buckets,
		Page: &model.Page{
			PS:    lbp.PS,
			PN:    lbp.PN,
			Total: count,
		},
	}
	return
}

// ListPublicBucket .
func (s *Service) ListPublicBucket(c context.Context, lbp *model.ListBucketParam) (bucketPage *model.BucketListPage, err error) {
	var (
		buckets []*model.Bucket
		count   int
	)
	if err = s.orm.Table("bucket").Order("id desc").Limit(lbp.PS).Offset((lbp.PN-1)*lbp.PS).Where("property != ? AND property != ?", model.PrivateRead, model.PrivateReadWrite).Find(&buckets).Error; err != nil {
		log.Error("read bucket error(%v)", err)
		return
	}
	if err = s.orm.Table("bucket").Count(&count).Error; err != nil {
		log.Error("bucket count error(%v)", err)
		err = nil
	}
	bucketPage = &model.BucketListPage{
		Items: buckets,
		Page: &model.Page{
			PS:    lbp.PS,
			PN:    lbp.PN,
			Total: count,
		},
	}
	return
}

// AddBucket .
func (s *Service) AddBucket(c context.Context, abp *model.AddBucketParam) (err error) {
	b := &model.Bucket{BucketName: abp.Name}
	if err = s.orm.Table("bucket").Where("bucket_name=?", abp.Name).
		Assign(map[string]interface{}{
			"key_id":        abp.KeyID,
			"key_secret":    abp.KeySecret,
			"purge_cdn":     abp.PurgeCDN,
			"property":      abp.Property,
			"cache_control": abp.CacheControl,
			"domain":        abp.Domain,
		}).FirstOrCreate(b).Error; err != nil {
		log.Error("Failed to add bucket (%+v): %v", b, err)
		return
	}
	return s.dao.CreateTable(context.Background(), abp.Name)
}

// DetailBucket .
func (s *Service) DetailBucket(c context.Context, bucketName string) (bucket *model.Bucket, err error) {
	var (
		limits []*model.DirLimit
	)
	bucket = &model.Bucket{}
	if err = s.orm.Table("bucket").Where("bucket_name = ?", bucketName).Find(bucket).Error; err != nil {
		return
	}
	if err = s.orm.Table("dir_limit").Where("bucket_name = ?", bucket.BucketName).Find(&limits).Error; err != nil {
		return
	}
	if len(limits) == 0 {
		return
	}
	bucket.DirLimit = limits
	return
}

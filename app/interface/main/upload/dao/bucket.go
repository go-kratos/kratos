package dao

import (
	"encoding/json"
	"strings"

	"go-common/app/interface/main/upload/model"
	"go-common/library/log"
)

// Buckets all bucket info from database.
func (d *Dao) Buckets() (bucketMap map[string]*model.Bucket, err error) {
	var (
		buckets  []*model.Bucket
		limitMap map[string]map[string]*model.DirConfig
	)
	if err = d.orm.Table("bucket").Find(&buckets).Error; err != nil {
		log.Error("orm.Table(bucket) error(%v)", err)
		return
	}
	if limitMap, err = d.dirLimits(); err != nil {
		return
	}
	bucketMap = make(map[string]*model.Bucket)
	for _, b := range buckets {
		v, ok := limitMap[b.Name]
		if ok {
			b.DirLimit = v
		}
		bucketMap[b.Name] = b
	}
	return
}

// dirLimits directory limit from database.
func (d *Dao) dirLimits() (limitMap map[string]map[string]*model.DirConfig, err error) {
	limits := make([]*model.DirLimit, 0)
	if err = d.orm.Table("dir_limit").Find(&limits).Error; err != nil {
		return
	}
	limitMap = make(map[string]map[string]*model.DirConfig)
	for _, l := range limits {
		var (
			pic  model.DirPicConfig
			rate model.DirRateConfig
		)
		if err = json.Unmarshal([]byte(l.DirPicConfig), &pic); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", l.DirPicConfig, err)
			err = nil
			continue
		}
		if pic.AllowType != "" {
			pic.AllowTypeSlice = strings.Split(pic.AllowType, ",")
		}
		if err = json.Unmarshal([]byte(l.DirRateConfig), &rate); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", l.DirRateConfig, err)
			err = nil
			continue
		}
		if _, ok := limitMap[l.BucketName]; !ok {
			limitMap[l.BucketName] = make(map[string]*model.DirConfig)
		}
		// NOTE empty dir is also in limit map
		l.Dir = strings.Trim(l.Dir, "/")
		limitMap[l.BucketName][l.Dir] = &model.DirConfig{Pic: pic, Rate: rate}
	}
	return
}

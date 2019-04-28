package dao

import (
	"math/rand"
	"time"

	"go-common/app/interface/bbq/app-bbq/model"
)

// GetRandVideoList .
func (d *Dao) GetRandVideoList(mid int64, limit int) []*model.SvInfo {
	var result []*model.SvInfo
	r := rand.New(rand.NewSource(time.Now().Unix()))

	mask := len(d.redundanceVideos) - limit
	cursor := r.Int() % mask

	for _, v := range d.redundanceVideos[cursor : cursor+limit] {
		result = append(result, &model.SvInfo{
			SVID: v.Svid,
			AVID: v.Avid,
			CID:  v.Cid,
			MID:  mid,
		})
	}

	return result
}

// GetRandSvList .
func (d *Dao) GetRandSvList(limit int) []int64 {
	result := make([]int64, limit)
	r := rand.New(rand.NewSource(time.Now().Unix()))

	mask := len(d.redundanceVideos) - limit
	cursor := r.Int() % mask

	for _, v := range d.redundanceVideos[cursor : cursor+limit] {
		result = append(result, v.Svid)
	}

	return result
}

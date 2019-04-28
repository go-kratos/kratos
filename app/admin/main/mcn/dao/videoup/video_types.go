package videoup

import (
	"context"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/log"
)

// refreshUpTypeAsync refresh in a goroutine
func (d *Dao) refreshUpTypeAsync() {
	for {
		time.Sleep(1 * time.Hour)
		d.refreshUpType()
	}
}

// videoTypes .
func (d *Dao) refreshUpType() {
	var (
		err error
		res struct {
			Code int                  `json:"code"`
			Data map[int]archive.Type `json:"data"`
		}
	)
	if err = d.client.Get(context.Background(), d.videTypeURL, "", nil, &res); err != nil {
		log.Error("refresh videoup types fail, err=%v", err)
		return
	}
	if res.Code != 0 {
		log.Error("videoTypes d.client.Get(%d)", res.Code)
	}
	d.videoUpTypeCache = res.Data
}

// GetTidName get tid name
func (d *Dao) GetTidName(tids []int64) (tpNames map[int64]string) {
	tpNames = make(map[int64]string, len(tids))
	for _, tid := range tids {
		info, ok := d.videoUpTypeCache[int(tid)]
		if !ok {
			continue
		}
		tpNames[tid] = info.Name
	}
	return
}

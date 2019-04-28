package weeklyhonor

import (
	"context"
	"go-common/library/log"
	"time"
)

// HonorInfoc log honor msg send status
func (d *Dao) HonorInfoc(c context.Context, mid int64, success int32) (err error) {
	ctime := time.Now().Format("20060102 15:04:05")
	i := map[string]interface{}{
		"mid":     mid,
		"ctime":   ctime,
		"exc":     "",
		"success": success,
	}
	log.Warn("infocproc create infoc(%v)", i)
	err = d.infoc.Info(mid, ctime, "", success)
	return
}

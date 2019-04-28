package pgc

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_style = "style_label"
)

// GetLabelCache .
func (d *Dao) GetLabelCache(ctx context.Context) (res map[int64][]*model.ParamStyle, err error) {
	var (
		conn = d.mc.Get(ctx)
		key  = _style
		rp   *memcache.Item
	)
	res = make(map[int64][]*model.ParamStyle)
	defer conn.Close()
	if rp, err = conn.Get(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
		} else {
			log.Error("mc.Get(%s) error(%v)", key, err)
		}
		return
	}
	if err = conn.Scan(rp, &res); err != nil {
		log.Error("conn.Scan error(%v)", err)
	}
	return
}

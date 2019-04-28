package dao

import (
	"context"
	"strconv"

	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_mcnSign = "mcn_s_"
)

// user mcn sign key.
func mcnSignKey(mcnMid int64) string {
	return _mcnSign + strconv.FormatInt(mcnMid, 10)
}

// DelMcnSignCache del mcn sign cache info.
func (d *Dao) DelMcnSignCache(c context.Context, mcnMid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(mcnSignKey(mcnMid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%d) error(%v)", mcnMid, err)
	}
	return
}

package up

import (
	"context"
	"fmt"
	"strconv"

	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_mcnSign  = "mcn_s_"
	_mcnUpper = "mcn_upperm_%d_%d"
)

// user mcn sign key.
func mcnSignKey(mcnMid int64) string {
	return _mcnSign + strconv.FormatInt(mcnMid, 10)
}

// user upper key.
func mcnUpperKey(signID, upMid int64) string {
	return fmt.Sprintf(_mcnUpper, signID, upMid)
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

// DelMcnUpperCache del mcn upper cache.
func (d *Dao) DelMcnUpperCache(c context.Context, signID, upMid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(mcnUpperKey(signID, upMid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%d, %d) error(%v)", signID, upMid, err)
	}
	return
}

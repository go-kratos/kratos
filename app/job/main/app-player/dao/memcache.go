package dao

import (
	"context"
	"strconv"

	"go-common/app/interface/main/app-player/model/archive"
	"go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

const (
	_prefixArc = "p_"
)

func keyArc(aid int64) string {
	return _prefixArc + strconv.FormatInt(aid, 10)
}

// AddArchiveCache add archive cache.
func (d *Dao) AddArchiveCache(c context.Context, aid int64, arc *archive.Info) (err error) {
	conn := d.mc.Get(c)
	key := keyArc(aid)
	item := &memcache.Item{Key: key, Object: arc, Flags: memcache.FlagProtobuf, Expiration: 0}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%v)", item)
	}
	conn.Close()
	return
}

package archive

import (
	"context"
	"strconv"

	"go-common/app/interface/main/app-player/model/archive"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_prefixArc = "p_"
)

func keyArc(aid int64) string {
	return _prefixArc + strconv.FormatInt(aid, 10)
}

func (d *Dao) archiveCache(c context.Context, aid int64) (arcMc *archive.Info, err error) {
	conn := d.arcMc.Get(c)
	key := keyArc(aid)
	defer conn.Close()
	r, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}
	arcMc = &archive.Info{}
	if err = conn.Scan(r, arcMc); err != nil {
		err = errors.Wrapf(err, "conn.Scan(%s)", r.Value)
	}
	return
}

// addArchiveCache add archive cache.
func (d *Dao) addArchiveCache(c context.Context, aid int64, arc *archive.Info) (err error) {
	conn := d.arcMc.Get(c)
	key := keyArc(aid)
	item := &memcache.Item{Key: key, Object: arc, Flags: memcache.FlagProtobuf}
	if err = conn.Set(item); err != nil {
		log.Error("conn.Set(%s, %v) error(%v)", key, arc, err)
	}
	conn.Close()
	return
}

func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.arcMc.Get(c)
	err = conn.Set(&memcache.Item{Key: "ping", Value: []byte{1}, Flags: memcache.FlagRAW, Expiration: 0})
	conn.Close()
	return
}

package card

import (
	"context"
	"fmt"
	"time"

	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_channelCardKey = "ccmd5_buvid_%v_%d"
)

func channelCardKey(buvid string, channelID int64) string {
	return fmt.Sprintf(_channelCardKey, buvid, channelID)
}

// AddChannelCardCache add user buvid and cardID cache
func (d *Dao) AddChannelCardCache(c context.Context, buvid, md5 string, channelID int64, now time.Time) (err error) {
	var (
		key            = channelCardKey(buvid, channelID)
		conn           = d.mc.Get(c)
		currenttimeSec = int32((now.Hour() * 60 * 60) + (now.Minute() * 60) + now.Second())
		overtime       int32
	)
	if overtime = d.expire - currenttimeSec; overtime < 1 {
		overtime = d.expire
	}
	if err = conn.Set(&memcache.Item{Key: key, Object: md5, Flags: memcache.FlagJSON, Expiration: overtime}); err != nil {
		log.Error("AddChannelCardCache d.mc.Set(%s,%v) error(%v)", key, channelID, err)
	}
	conn.Close()
	return
}

// ChannelCardCache user buvid channel card
func (d *Dao) ChannelCardCache(c context.Context, buvid string, channelID int64) (md5 string, err error) {
	var (
		key  = channelCardKey(buvid, channelID)
		conn = d.mc.Get(c)
		r    *memcache.Item
	)
	defer conn.Close()
	if r, err = conn.Get(key); err != nil {
		return
	}
	if err = conn.Scan(r, &md5); err != nil {
		log.Error("r.Scan(%s) error(%v)", r.Value, err)
	}
	return
}

package archive

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

func vcoverKey(filename string) string {
	return fmt.Sprintf("%s_%s", "vcover_", filename)
}

// AddRdsCovers fn.
func (d *Dao) AddRdsCovers(c context.Context, covers []*archive.Cover) (ok bool, err error) {
	var conn = d.coverRds.Get(c)
	defer conn.Close()
	var key string
	log.Info("AddRdsCovers info:(%+v)", covers)
	for _, c := range covers {
		if key == "" {
			key = vcoverKey(c.Filename)
		}
		if err = conn.Send("SADD", key, c.BFSPath); err != nil {
			log.Error("conn.Do(SETEX, %s, %s, %d, %d) error(%v)", c.Filename, c.BFSPath, d.coverExpire, time.Now().Unix(), err)
		}
	}
	if err = conn.Send("EXPIRE", key, d.coverExpire); err != nil {
		log.Error("conn.Send(EXPIRE, %s, %d) error(%v)", key, d.coverExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(covers); i++ {
		if ok, err = redis.Bool(conn.Receive()); err != nil {
			log.Error("conn.Receive error(%v)", err)
		}
	}
	return
}

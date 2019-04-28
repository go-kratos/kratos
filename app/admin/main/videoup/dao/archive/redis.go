package archive

import (
	"context"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_lockingVideo = "videoup_admin_locking_video"
)

// IsLockingVideo 是否正在自动锁定视频
func (d *Dao) IsLockingVideo(c context.Context) (locking bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if locking, err = redis.Bool(conn.Do("EXISTS", _lockingVideo)); err != nil {
		log.Error("conn.Do(EXISTS,%s) error(%v)", _lockingVideo, err)
	}
	return
}

// LockingVideo 设置是否正在自动锁定视频
func (d *Dao) LockingVideo(c context.Context, v int8) (err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()
	if v == 1 {
		if _, err = conn.Do("SET", _lockingVideo, v); err != nil {
			log.Error("conn.Do(SET, %s,%d) error(%v)", _lockingVideo, v, err)
		}
	} else {
		if _, err = conn.Do("DEL", _lockingVideo); err != nil {
			log.Error("conn.Do(SET, %s) error(%v)", _lockingVideo, err)
		}
	}

	return
}

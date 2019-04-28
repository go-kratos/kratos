package space

import (
	"context"
	"go-common/app/job/main/app/model"
	"strconv"

	"github.com/pkg/errors"
)

const (
	_prefixContribute = "cb_"
)

func keyContribute(vmid int64) string {
	return _prefixContribute + strconv.FormatInt(vmid, 10)
}

func (d *Dao) DelContributeIDCache(c context.Context, vmid, id int64, gt string) (err error) {
	conn := d.redis.Get(c)
	key := keyContribute(vmid)
	member := model.FormatKey(id, gt)
	if _, err = conn.Do("ZREM", key, member); err != nil {
		err = errors.Wrapf(err, "conn.Do(ZREM,%s,%d)", key, id)
	}
	conn.Close()
	return
}

func (d *Dao) DelContributeCache(c context.Context, vmid int64) (err error) {
	conn := d.redis.Get(c)
	key := keyContribute(vmid)
	if _, err = conn.Do("DEL", key); err != nil {
		err = errors.Wrapf(err, "conn.Do(DEL,%s)", key)
	}
	conn.Close()
	return
}

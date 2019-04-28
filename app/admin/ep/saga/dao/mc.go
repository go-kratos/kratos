package dao

import (
	"context"

	"go-common/app/admin/ep/saga/model"
	"go-common/library/cache/memcache"

	"github.com/pkg/errors"
)

func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 0}); err != nil {
		err = errors.Wrap(err, "conn.Store(set,ping,1)")
	}
	return
}

// SetData set data info to memcache
func (d *Dao) SetData(c context.Context, key string, dataMap map[string]*model.TeamDataResp) (err error) {
	var (
		conn = d.mc.Get(c)
		item *memcache.Item
		//dataMap = make(map[string]*model.TeamDataResp)
	)
	defer conn.Close()

	item = &memcache.Item{Key: key, Object: dataMap, Expiration: 0, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%v)", key, dataMap)
		return
	}
	return
}

// GetData get data info from memcache
func (d *Dao) GetData(c context.Context, key string, dataMap *map[string]*model.TeamDataResp) (err error) {
	var (
		conn  = d.mc.Get(c)
		reply *memcache.Item
	)
	defer conn.Close()

	reply, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s)", key)
		return
	}

	if err = conn.Scan(reply, dataMap); err != nil {
		err = errors.Wrapf(err, "reply.Scan(%s)", string(reply.Value))
		return
	}

	return
}

// DeleteData delete data info in memcache
func (d *Dao) DeleteData(c context.Context, key string) (err error) {
	var (
		conn = d.mc.Get(c)
	)
	defer conn.Close()

	err = conn.Delete(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Delete(%s)", key)
	}

	return
}

// SetPipeline set pipeline info info to memcache
func (d *Dao) SetPipeline(c context.Context, key string, pipeline *model.PipelineDataResp) (err error) {
	var (
		conn = d.mc.Get(c)
		item *memcache.Item
	)
	defer conn.Close()
	item = &memcache.Item{Key: key, Object: pipeline, Expiration: 0, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%v)", key, pipeline)
		return
	}
	return
}

// GetPipeline get pipeline info from memcache
func (d *Dao) GetPipeline(c context.Context, key string) (pipeline *model.PipelineDataResp, err error) {
	var (
		conn  = d.mc.Get(c)
		reply *memcache.Item
	)
	defer conn.Close()

	reply, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			return
		}
		err = errors.Wrapf(err, "conn.Get(%s,%v)", key, pipeline)
		return
	}

	pipeline = new(model.PipelineDataResp)
	if err = conn.Scan(reply, pipeline); err != nil {
		err = errors.Wrapf(err, "reply.Scan(%s)", string(reply.Value))
	}

	return
}

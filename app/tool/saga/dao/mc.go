package dao

import (
	"context"
	"fmt"

	"go-common/app/tool/saga/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	requiredViableUsersKey = "saga_wechat_require_visible_users_key"
)

func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mcMR.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 0}); err != nil {
		err = errors.Wrap(err, "conn.Store(set,ping,1)")
	}
	return
}

func mrRecordKey(mrID int) string {
	return fmt.Sprintf("saga_mr_%d", mrID)
}

// MRRecordCache get MRRecord from mc
func (d *Dao) MRRecordCache(c context.Context, mrID int) (record *model.MRRecord, err error) {
	var (
		key   = mrRecordKey(mrID)
		conn  = d.mcMR.Get(c)
		reply *memcache.Item
	)
	defer conn.Close()
	reply, err = conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(get,%s)", key)
		return
	}
	record = &model.MRRecord{}
	if err = conn.Scan(reply, record); err != nil {
		err = errors.Wrapf(err, "reply.Scan(%s)", string(reply.Value))
	}
	return
}

// SetMRRecordCache set MRRecord to mc
func (d *Dao) SetMRRecordCache(c context.Context, record *model.MRRecord) (err error) {
	var (
		key  = mrRecordKey(record.MRID)
		conn = d.mcMR.Get(c)
	)
	defer conn.Close()

	if err = conn.Set(&memcache.Item{Key: key, Object: record, Expiration: 0, Flags: memcache.FlagJSON}); err != nil {
		err = errors.Wrapf(err, "conn.Add(%s,%v)", key, record)
		return
	}
	return
}

func weixinTokenKey(key string) string {
	return fmt.Sprintf("saga_weixin_token_%s", key)
}

// AccessToken get access token from mc
func (d *Dao) AccessToken(c context.Context, key string) (token string, err error) {
	var (
		wkey  = weixinTokenKey(key)
		conn  = d.mcMR.Get(c)
		reply *memcache.Item
	)
	defer conn.Close()

	reply, err = conn.Get(wkey)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "conn.Get(get,%s)", wkey)
		return
	}

	if err = conn.Scan(reply, &token); err != nil {
		err = errors.Wrapf(err, "reply.Scan(%s)", string(reply.Value))
	}
	return
}

// SetAccessToken set the access token to mc
func (d *Dao) SetAccessToken(c context.Context, key string, token string, expire int32) (err error) {
	var (
		wkey = weixinTokenKey(key)
		conn = d.mcMR.Get(c)
		item *memcache.Item
	)
	defer conn.Close()

	item = &memcache.Item{Key: wkey, Object: token, Expiration: expire, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Add(%s,%v)", wkey, token)
		return
	}
	return
}

// RequireVisibleUsers get wechat require visible users from memcache
func (d *Dao) RequireVisibleUsers(c context.Context, userMap *map[string]model.RequireVisibleUser) (err error) {
	var (
		conn  = d.mcMR.Get(c)
		reply *memcache.Item
	)
	defer conn.Close()

	reply, err = conn.Get(requiredViableUsersKey)
	if err != nil {
		if err == memcache.ErrNotFound {
			log.Info("no such key (%s) in cache, err (%s)", requiredViableUsersKey, err.Error())
			err = nil
		}
		return
	}

	if err = conn.Scan(reply, userMap); err != nil {
		err = errors.Wrapf(err, "reply.Scan(%s)", string(reply.Value))
	}

	return
}

// SetRequireVisibleUsers set wechat require visible users to memcache
func (d *Dao) SetRequireVisibleUsers(c context.Context, contactInfo *model.ContactInfo) (err error) {
	var (
		conn    = d.mcMR.Get(c)
		item    *memcache.Item
		userMap = make(map[string]model.RequireVisibleUser)
	)
	defer conn.Close()

	if err = d.RequireVisibleUsers(c, &userMap); err != nil {
		log.Error("get require visible user error(%v)", err)
		return
	}

	user := model.RequireVisibleUser{
		UserName: contactInfo.UserName,
		NickName: contactInfo.NickName,
	}
	userMap[contactInfo.UserID] = user

	item = &memcache.Item{Key: requiredViableUsersKey, Object: userMap, Expiration: 0, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%v)", requiredViableUsersKey, userMap)
		return
	}
	return
}

// DeleteRequireVisibleUsers delete the wechat require visible key in memcache
func (d *Dao) DeleteRequireVisibleUsers(c context.Context) (err error) {
	var (
		conn = d.mcMR.Get(c)
	)
	defer conn.Close()

	err = conn.Delete(requiredViableUsersKey)
	if err != nil {
		err = errors.Wrapf(err, "conn.Delete(%s)", requiredViableUsersKey)
	}

	return
}

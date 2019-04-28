package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/admin/ep/saga/model"
	"go-common/library/cache/redis"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	requiredViableUsersKeyRedis = "saga_wechat_require_visible_users_key"
)

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	defer conn.Close()
	return
}

func weixinTokenKeyRedis(key string) string {
	return fmt.Sprintf("saga_weixin_token_%s", key)
}

// AccessTokenRedis get access token from redis
func (d *Dao) AccessTokenRedis(c context.Context, key string) (token string, err error) {
	var (
		wkey  = weixinTokenKeyRedis(key)
		conn  = d.redis.Get(c)
		value []byte
	)
	defer conn.Close()
	if value, err = redis.Bytes(conn.Do("GET", wkey)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	if err = json.Unmarshal(value, &token); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// SetAccessTokenRedis set the access token to redis
func (d *Dao) SetAccessTokenRedis(c context.Context, key, token string, expire int32) (err error) {
	var (
		wkey = weixinTokenKeyRedis(key)
		conn = d.redis.Get(c)
		item []byte
	)
	defer conn.Close()
	if item, err = json.Marshal(token); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = conn.Send("SET", wkey, item, "EX", expire); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		return
	}
	return
}

// RequireVisibleUsersRedis get wechat require visible users from redis
func (d *Dao) RequireVisibleUsersRedis(c context.Context, userMap *map[string]model.RequireVisibleUser) (err error) {
	var (
		conn  = d.redis.Get(c)
		reply []byte
	)
	defer conn.Close()

	if reply, err = redis.Bytes(conn.Do("GET", requiredViableUsersKeyRedis)); err != nil {
		if err == redis.ErrNil {
			log.Info("no such key (%s) in cache, err (%s)", requiredViableUsersKeyRedis, err.Error())
			err = nil
		}
		return
	}
	if err = json.Unmarshal(reply, &userMap); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// SetRequireVisibleUsersRedis set wechat require visible users to redis
func (d *Dao) SetRequireVisibleUsersRedis(c context.Context, contactInfo *model.ContactInfo) (err error) {
	var (
		conn    = d.redis.Get(c)
		item    []byte
		userMap = make(map[string]model.RequireVisibleUser)
	)
	defer conn.Close()

	if err = d.RequireVisibleUsersRedis(c, &userMap); err != nil {
		log.Error("get require visible user error(%v)", err)
		return
	}

	user := model.RequireVisibleUser{
		UserName: contactInfo.UserName,
		NickName: contactInfo.NickName,
	}
	userMap[contactInfo.UserID] = user

	if item, err = json.Marshal(userMap); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = conn.Send("SET", requiredViableUsersKeyRedis, item); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		err = errors.Wrapf(err, "conn.Set(%s,%v)", requiredViableUsersKeyRedis, userMap)
		return
	}
	return
}

// DeleteRequireVisibleUsersRedis delete the wechat require visible key in redis
func (d *Dao) DeleteRequireVisibleUsersRedis(c context.Context) (err error) {
	var (
		conn = d.redis.Get(c)
	)
	defer conn.Close()

	if err = conn.Send("DEL", requiredViableUsersKeyRedis); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if _, err = conn.Receive(); err != nil {
		err = errors.Wrapf(err, "conn.Delete(%s)", requiredViableUsersKeyRedis)
		return
	}
	return
}

// SetItemRedis ...
func (d *Dao) SetItemRedis(c context.Context, key string, value interface{}, ttl int) (err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()

	if bs, err = json.Marshal(value); err != nil {
		return errors.WithStack(err)
	}
	if ttl == 0 {
		if _, err = conn.Do("SET", key, bs); err != nil {
			return errors.Wrapf(err, "conn.Do(SET %s, %s) error(%v)", key, value, err)
		}
	} else {
		if _, err = conn.Do("SET", key, bs, "EX", ttl); err != nil {
			return errors.Wrapf(err, "conn.Do(SET %s, %s) error(%v)", key, value, err)
		}
	}
	return
}

// ItemRedis ...
func (d *Dao) ItemRedis(c context.Context, key string, value interface{}) (err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()

	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			return
		}
		return errors.Wrapf(err, "conn.Get(%s)", key)
	}
	if err = json.Unmarshal(bs, &value); err != nil {
		return errors.WithStack(err)
	}
	return
}

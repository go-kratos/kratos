package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/interface/main/up-rating/model"
	"go-common/library/cache/redis"

	"github.com/pkg/errors"
)

// GetUpRatingCache ...
func (d *Dao) GetUpRatingCache(c context.Context, mid int64) (rating *model.Rating, err error) {
	var (
		key  = upRatingKey(mid)
		conn = d.redis.Get(c)
		data []byte
	)
	defer conn.Close()
	if data, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		}
		return
	}
	rating = new(model.Rating)
	if err = json.Unmarshal(data, &rating); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return
}

// ExpireUpRatingCache ...
func (d *Dao) ExpireUpRatingCache(c context.Context, mid int64) (err error) {
	var (
		key  = upRatingKey(mid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("EXPIRE", key, 0); err != nil {
		return
	}
	return
}

// SetUpRatingCache ...
func (d *Dao) SetUpRatingCache(c context.Context, mid int64, rating *model.Rating) (err error) {
	var (
		key  = upRatingKey(mid)
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	if bs, err = json.Marshal(rating); err != nil {
		return errors.WithStack(err)
	}
	if err = conn.Send("SET", key, bs); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", key, d.upRatingExpire); err != nil {
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

func upRatingKey(mid int64) string {
	return fmt.Sprintf("up_rating_detail_mid_%d", mid)
}

package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/container/pool"
	xtime "go-common/library/time"
)

var p *redis.Pool
var config *redis.Config
var d *Dao
var ctx = context.TODO()

func init() {
	config = getConfig()
	p = redis.NewPool(config)
	d = &Dao{redis: p}
}

func getConfig() (c *redis.Config) {
	c = &redis.Config{
		Name:         "test",
		Proto:        "tcp",
		Addr:         "127.0.0.1:6379",
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	c.Config = &pool.Config{
		Active:      20,
		Idle:        2,
		IdleTimeout: xtime.Duration(90 * time.Second),
	}
	return
}

func TestDao_Get(t *testing.T) {
	v, e := d.Get(ctx, "golang")
	fmt.Println(v, e)
}

func TestDao_Set(t *testing.T) {
	fmt.Println(d.Set(ctx, "golang", 22))
}

func TestDao_SetEx(t *testing.T) {
	got, err := d.SetEx(ctx, "golang", "b", 10)
	fmt.Println(got, err)

}

func TestDao_Del(t *testing.T) {
	got, err := d.Del(ctx, "golang")
	fmt.Println(got, err)

}

func TestDao_HMSet(t *testing.T) {
	m := map[string]interface{}{"id": 1, "uid": 2, "content": "sss"}

	got, err := d.HMSet(ctx, "myhash", m)
	fmt.Println(got, err)

	v, _ := d.HGetAll(ctx, "myhash")

	fmt.Printf("%+v\n", v)
}

func TestDao_Expire(t *testing.T) {
	got, err := d.Expire(ctx, "ttl", 10)
	fmt.Println(got, err)
}

func TestDao_HGetAll(t *testing.T) {
	got, err := d.HGetAll(ctx, "myhash")
	if err != nil && err != ErrEmptyMap {
		fmt.Println(err.Error())
	}

	fmt.Printf("%+v ,%+v", got, err)

}

func TestDao_SetWithNxEx(t *testing.T) {
	fmt.Println(d.SetWithNxEx(ctx, "python", "1", 100))
}

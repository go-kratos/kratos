package dao

import (
	"context"
	"strconv"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_expire        = 30 * 60 // 30 minutes
	_prefixCaptcha = "recovery:ca_"
)

// SetLinkMailCount set linkMail expire time.
func (d *Dao) SetLinkMailCount(c context.Context, linkMail string) (state int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	//用时间戳去减去时间 当天0点过期，第二天又可以发送10封邮件
	i, _ := redis.Int(conn.Do("incr", linkMail))
	conn.Do("expire", linkMail, getSubtime())
	if i >= 11 { //第10封邮件之后
		state = 10 //当天邮件发送到达最大次数
		return
	}
	return
}

func getSubtime() (subtime int64) {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	t2 := t.AddDate(0, 0, 1).Unix()
	curr := time.Now().Unix() //当前时间
	subtime = t2 - curr
	return
}

// keyCaptcha
func keyCaptcha(mid int64, linkMail string) string {
	return _prefixCaptcha + strconv.FormatInt(mid, 10) + "_" + linkMail
}

// SetCaptcha set linkMail expire time.
func (d *Dao) SetCaptcha(c context.Context, code string, mid int64, linkMail string) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := keyCaptcha(mid, linkMail)
	//验证码30分钟内有效
	if _, err = conn.Do("SETEX", key, _expire, code); err != nil {
		log.Error("conn.Do(SETEX, %d, %v, %s) error(%v)", mid, _expire, code, err)
	}
	return
}

// GetEMailCode get captcha from redis
func (d *Dao) GetEMailCode(c context.Context, mid int64, linkMail string) (code string, err error) {
	key := keyCaptcha(mid, linkMail)
	conn := d.redis.Get(c)
	defer conn.Close()
	code, err = redis.String(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("conn.Do(GET, %s, ), err (%v)", key, err)
		return
	}
	return
}

// DelEMailCode del captcha from redis 提交：校验验证之后就删除验证码(保证只能提交一次)
func (d *Dao) DelEMailCode(c context.Context, mid int64, linkMail string) (err error) {
	key := keyCaptcha(mid, linkMail)
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("conn.Do(DEL, %s, ), err (%v)", key, err)
		return
	}
	return
}

// PingRedis check connection success.
func (d *Dao) PingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("GET", "PING")
	return
}

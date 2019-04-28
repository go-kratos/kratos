package email

import (
	"context"
	"encoding/json"
	"go-common/app/job/main/aegis/conf"
	"go-common/app/job/main/aegis/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"gopkg.in/gomail.v2"
	"time"
)

type Dao struct {
	c     *conf.Config
	redis *redis.Pool
	email *gomail.Dialer
}

const (
	// MoniEmailKey 监控邮件队列key
	MoniEmailKey = "monitor_stats_email"
)

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:     c,
		email: gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.Username, c.Mail.Password),
		redis: redis.NewPool(c.Redis),
	}
	return d
}

// MonitorEmailAsync 异步发送监控邮件
func (d *Dao) MonitorEmailAsync(c context.Context, members []string, title, content string) (err error) {
	var (
		conn = d.redis.Get(c)
		bs   []byte
	)
	defer conn.Close()
	temp := &model.MoniTemp{
		From:    d.c.Mail.Username,
		Members: members,
		Subject: title,
		Body:    content,
	}
	if bs, err = json.Marshal(temp); err != nil {
		log.Error("d.MonitorEmailAsync() json.Marshal(%+v) error(%v) key(%s)", temp, err, MoniEmailKey)
		return
	}
	if _, err = conn.Do("RPUSH", MoniEmailKey, bs); err != nil {
		log.Error("d.MonitorEmailAsync() conn.Do(RPUSH, %s, %s) error(%v)", MoniEmailKey, bs, err)
	}
	return
}

// MonitorEmailProc 发送监控邮件
func (d *Dao) MonitorEmailProc() (err error) {
	var (
		bs      []byte
		temp    *model.MoniTemp
		headers map[string][]string
	)
	headers = make(map[string][]string)
	bs, err = d.PopRedis(context.TODO(), MoniEmailKey)
	if err != nil || bs == nil {
		log.Warn("d.MonitorEmailProc() warn:%v content:%s", err, bs)
		time.Sleep(5 * time.Second)
		return
	}
	if err = json.Unmarshal(bs, &temp); err != nil {
		log.Error("d.MonitorEmailProc() json.unmarshal error(%v) content(%s)", err, bs)
		return
	}
	msg := gomail.NewMessage()
	headers["From"] = []string{d.c.Mail.Username}
	headers["To"] = temp.Members
	headers["Subject"] = []string{temp.Subject}
	msg.SetHeaders(headers)
	msg.SetBody("text/html", temp.Body)
	if err = d.email.DialAndSend(msg); err != nil {
		log.Error("d.email.DialAndSend(%+v) error:%v", msg, err)
		return
	}
	return
}

// PopRedis lpop fail item from redis
func (d *Dao) PopRedis(c context.Context, key string) (bs []byte, err error) {
	var conn = d.redis.Get(c)
	defer conn.Close()

	if bs, err = redis.Bytes(conn.Do("LPOP", key)); err != nil && err != redis.ErrNil {
		log.Error("d.PopRedis() redis.Bytes(conn.Do(LPOP, %s)) error(%v)", key, err)
	}
	return
}

package email

import (
	"crypto/tls"

	"go-common/app/job/main/videoup-report/conf"
	"go-common/app/job/main/videoup-report/model/email"
	"go-common/library/cache/redis"
	gomail "gopkg.in/gomail.v2"
)

// Dao is redis dao.
type Dao struct {
	c           *conf.Config
	redis       *redis.Pool
	email       *gomail.Dialer
	FansAddr    map[int16][]string
	emailAddr   map[string][]string
	PrivateAddr map[string][]string
	//fast behavior detector
	detector *email.FastDetector
	//快速通道token
	fastChan chan int
	//邮件发送api的频率token，发送邮件5s后插入
	controlChan chan int64
}

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	emailAddr := make(map[string][]string)
	for _, v := range c.Mail.Addr {
		emailAddr[v.Type] = v.Addr
	}

	privateMail := make(map[string][]string)
	for _, v := range c.Mail.PrivateAddr {
		privateMail[v.Type] = v.Addr
	}

	d = &Dao{
		c:           c,
		redis:       redis.NewPool(c.Redis.Mail),
		email:       gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.Username, c.Mail.Password),
		emailAddr:   emailAddr,
		PrivateAddr: privateMail,
		detector:    email.NewFastDetector(c.Mail.SpeedThreshold, c.Mail.OverspeedThreshold),
		fastChan:    make(chan int, 10240),
		controlChan: make(chan int64, 1),
	}

	d.email.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	d.fastChan <- 1
	d.controlChan <- 1
	return d
}

//Close close
func (d *Dao) Close() (err error) {
	if d.redis != nil {
		err = d.redis.Close()
	}
	return
}

//FastChan get fast channel
func (d *Dao) FastChan() <-chan int {
	return d.fastChan
}

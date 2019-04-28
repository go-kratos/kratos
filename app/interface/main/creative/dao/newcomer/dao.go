package newcomer

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

// Dao  define
type Dao struct {
	c  *conf.Config
	db *sql.DB
	// http
	client       *httpx.Client
	mallURI      string
	bPayURI      string
	pendantURI   string
	bigMemberURI string
	msgNotifyURI string
	// redis
	redis *redis.Pool
}

const (
	_mall      = "/mall-marketing/coupon_code/create"   //会员购
	_bpay      = "/api/coupon/add"                      //B币券
	_pendant   = "/x/internal/pendant/multiGrantByMid"  //挂件：批量发放挂件(多个MID对应一个挂件)
	_bigmember = "/x/internal/coupon/allowance/receive" //大会员代金券
	_notify    = "/api/notify/send.user.notify.do"      //发送用户通知消息接口
)

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		db:           sql.NewMySQL(c.DB.Creative),
		client:       httpx.NewClient(c.HTTPClient.Slow),
		mallURI:      c.Host.Mall + _mall,
		bPayURI:      c.Host.BPay + _bpay,
		pendantURI:   c.Host.Pendant + _pendant,
		bigMemberURI: c.Host.BigMember + _bigmember,
		msgNotifyURI: c.Host.Notify + _notify,
		redis:        redis.NewPool(c.Redis.Cover.Config),
	}
	return
}

// Ping db
func (d *Dao) Ping(c context.Context) (err error) {
	if d.db != nil {
		d.db.Ping(c)
	}

	if d.redis != nil {
		conn := d.redis.Get(c)
		if _, err = conn.Do("SET", "ping", "pong"); err != nil {
			log.Error("conn.Do(SET) error(%v)", err)
		}
		conn.Close()
	}
	return
}

// Close db
func (d *Dao) Close() (err error) {
	if d.db != nil {
		d.db.Close()
	}
	if d.redis != nil {
		d.redis.Close()
	}
	return
}

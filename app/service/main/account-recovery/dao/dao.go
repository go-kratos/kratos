package dao

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"

	"go-common/app/service/main/account-recovery/conf"
	account "go-common/app/service/main/account/api"
	location "go-common/app/service/main/location/rpc/client"
	member "go-common/app/service/main/member/api"
	"go-common/library/cache/redis"
	"go-common/library/database/elastic"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	"gopkg.in/gomail.v2"
)

// Dao dao
type Dao struct {
	c     *conf.Config
	redis *redis.Pool
	db    *xsql.DB
	// httpClient
	httpClient *bm.Client

	// email
	email *gomail.Dialer
	es    *elastic.Elastic

	// rpc
	locRPC *location.Service

	// grpc
	memberClient  member.MemberClient
	accountClient account.AccountClient

	hashSalt []byte
	AESBlock cipher.Block
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		redis: redis.NewPool(c.Redis),
		db:    xsql.NewMySQL(c.MySQL),
		// httpClient
		httpClient: bm.NewClient(c.HTTPClientConfig),

		email:    gomail.NewDialer(c.MailConf.Host, c.MailConf.Port, c.MailConf.Username, c.MailConf.Password),
		es:       elastic.NewElastic(c.Elastic),
		locRPC:   location.New(c.LocationRPC),
		hashSalt: []byte(c.AESEncode.Salt),
	}
	dao.email.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	dao.AESBlock, _ = aes.NewCipher([]byte(c.AESEncode.AesKey))

	var err error
	if dao.memberClient, err = member.NewClient(c.MemberGRPC); err != nil {
		panic(err)
	}
	if dao.accountClient, err = account.NewClient(c.AccountGRPC); err != nil {
		panic(err)
	}

	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		return
	}
	if err = d.PingRedis(c); err != nil {
		return
	}
	// TODO: if you need use mc,redis, please add
	return
}

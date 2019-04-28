package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/passport-game/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Dao dao
type Dao struct {
	c             *conf.Config
	getMemberStmt []*sql.Stmt
	cloudDB       *sql.DB
	otherRegion   *sql.DB
	mc            *memcache.Pool
	mcExpire      int32
	client        *bm.Client
	myInfoURI     string
	loginURI      string
	getKeyURI     string
	regV3URI      string
	regV2URI      string
	regURI        string
	byTelURI      string
	captchaURI    string
	sendSmsURI    string
}

const (
	_myInfoURI  = "/api/myinfo"
	_loginURI   = "/api/login"
	_getKeyURI  = "/api/login/get_key"
	_regV3URI   = "/api/reg/regV3"
	_regV2URI   = "/api/reg/regV2"
	_regURI     = "/api/reg/reg"
	_byTelURI   = "/api/reg/byTelGame"
	_captchaURI = "/bilicaptcha/token"
	_sendSmsURI = "/api/sms/sendCaptcha"
)

// New dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:           c,
		cloudDB:     sql.NewMySQL(c.DB.Cloud),
		otherRegion: sql.NewMySQL(c.DB.OtherRegion),
		mc:          memcache.NewPool(c.Memcache.Config),
		mcExpire:    int32(time.Duration(c.Memcache.Expire) / time.Second),
		client:      bm.NewClient(c.HTTPClient),
		myInfoURI:   c.AccountURI + _myInfoURI,
		loginURI:    c.PassportURI + _loginURI,
		getKeyURI:   c.PassportURI + _getKeyURI,
		regV3URI:    c.PassportURI + _regV3URI,
		regV2URI:    c.PassportURI + _regV2URI,
		regURI:      c.PassportURI + _regURI,
		byTelURI:    c.PassportURI + _byTelURI,
		captchaURI:  c.PassportURI + _captchaURI,
		sendSmsURI:  c.PassportURI + _sendSmsURI,
	}
	d.getMemberStmt = make([]*sql.Stmt, _memberShard)
	for i := 0; i < _memberShard; i++ {
		d.getMemberStmt[i] = d.cloudDB.Prepared(fmt.Sprintf(_getMemberInfoSQL, i))
	}
	return
}

// Ping verify server is ok.
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.cloudDB.Ping(c); err != nil {
		log.Info("dao.cloudDB.Ping() error(%v)", err)
	}
	if err = d.otherRegion.Ping(c); err != nil {
		log.Info("dao.otherRegion.Ping() error(%v)", err)
	}
	return d.pingMC(c)
}

// Close close connections of mc, cloudDB.
func (d *Dao) Close() (err error) {
	if d.cloudDB != nil {
		d.cloudDB.Close()
	}
	if d.otherRegion != nil {
		d.otherRegion.Close()
	}
	if d.mc != nil {
		d.mc.Close()
	}
	return
}

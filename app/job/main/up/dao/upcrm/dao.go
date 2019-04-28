package upcrm

import (
	"context"
	"go-common/app/job/main/up/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

//Dao upcrm dao
type Dao struct {
	// config
	conf *conf.Config
	// db
	crmdb      *gorm.DB
	httpClient *bm.Client
}

//New new dao
func New(c *conf.Config) *Dao {
	var d = &Dao{
		conf: c,
	}
	crmdb, err := gorm.Open("mysql", c.Upcrm.DSN)
	if crmdb == nil {
		log.Error("connect to db fail, err=%v", err)
		return nil
	}
	crmdb.SingularTable(true)
	d.crmdb = crmdb
	d.crmdb.LogMode(c.IsTest)
	return d
}

//SetHTTPClient set http client
func (d *Dao) SetHTTPClient(client *bm.Client) {
	d.httpClient = client
}

//GetDb get current gorm db
func (d *Dao) GetDb() *gorm.DB {
	return d.crmdb
}

// Close close resource.
func (d *Dao) Close() {
	d.crmdb.Close()
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}

package dao

import (
	"context"

	"go-common/app/service/main/workflow/conf"
	"go-common/app/service/main/workflow/model"
	"go-common/library/database/elastic"
	"go-common/library/database/orm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

const (
	_mngTagURL     = "/x/admin/manager/internal/tag/list"
	_mngControlURL = "/x/admin/manager/internal/control/list"
)

// Dao tag dao
type Dao struct {
	c *conf.Config
	// db *sql.DB
	DB *gorm.DB

	callback    *bm.Client
	callbackMap map[int8]string

	ReadClient    *bm.Client
	MngTagURL     string
	MngControlURL string
	es            *elastic.Elastic
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:             c,
		DB:            orm.NewMySQL(c.ORM.Write),
		callback:      bm.NewClient(c.HTTPClient.Write),
		callbackMap:   make(map[int8]string),
		ReadClient:    bm.NewClient(c.HTTPClient.Read),
		MngTagURL:     c.Host.ManagerURI + _mngTagURL,
		MngControlURL: c.Host.ManagerURI + _mngControlURL,
		es:            elastic.NewElastic(c.Elastic),
	}
	d.initORM()
	d.initCallback()
	return
}

func (d *Dao) initORM() {
	d.DB.LogMode(true)
}

func (d *Dao) initCallback() {
	callbacks := []model.Callback{}
	if err := d.DB.Where("state =?", model.Enabled).Find(&callbacks).Error; err != nil {
		log.Error("d.CallbackSetting() error(%v)", err)
		panic(err)
	}
	for _, callback := range callbacks {
		d.callbackMap[callback.Business] = callback.URL
	}
}

// Close close dao.
func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	if d.DB != nil {
		err = d.DB.DB().PingContext(c)
	}
	return
}

package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/service/openplatform/ticket-item/conf"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/cache/redis"
	"go-common/library/database/orm"
	"go-common/library/log"

	"go-common/library/database/elastic"

	"go-common/library/sync/pipeline/fanout"

	"github.com/jinzhu/gorm"
)

// Expire time
const (
	CacheTimeout    = 120
	_expireHalfhour = 1800 // 半小时过期
)

// Dao dao
type Dao struct {
	c     *conf.Config
	redis *redis.Pool
	cache *fanout.Fanout
	// DB
	db     *gorm.DB
	expire int32
	es     *elastic.Elastic
}

func keyItem(id int64) string {
	return "open_item_" + strconv.FormatInt(id, 10)
}

func keyItemDetail(id int64) string {
	return "open_item_detail_" + strconv.FormatInt(id, 10)
}

func keyItemTicket(id int64) string {
	return "open_itemticket_" + strconv.FormatInt(id, 10)
}

func keyTicket(id int64) string {
	return "open_ticket_" + strconv.FormatInt(id, 10)
}

func keyVenue(id int64) string {
	return "open_venue_" + strconv.FormatInt(id, 10)
}

func keyPlace(id int64) string {
	return "open_place_" + strconv.FormatInt(id, 10)
}

func keyItemScreen(id int64) string {
	return "open_itemscreen_" + strconv.FormatInt(id, 10)
}

func keyScreen(id int64) string {
	return "open_screen_" + strconv.FormatInt(id, 10)
}

func keyBannerList(order int32, districtID string, position int32, subPosition int32) string {
	return fmt.Sprintf("BANNERLISTV3:%d:%s:%d:%d", order, districtID, position, subPosition)
}

func keyBannerInfo(bannerID int64) string {
	return fmt.Sprintf("%d:BANNERINFOV2", bannerID)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -nullcache=&model.Item{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	Items(c context.Context, pid []int64) (info map[int64]*model.Item, err error)
	// cache: -nullcache=&model.ItemDetail{ProjectID:-1} -check_null_code=$!=nil&&$.ProjectID==-1
	ItemDetails(c context.Context, pid []int64) (details map[int64]*model.ItemDetail, err error)
	// cache: -nullcache=[]*model.TicketInfo{{TicketPrice:model.TicketPrice{ProjectID:-1}}} -check_null_code=len($)==1&&$[0].ProjectID==-1
	TkListByItem(c context.Context, pid []int64) (info map[int64][]*model.TicketInfo, err error)
	// cache: -nullcache=&model.Venue{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	Venues(c context.Context, id []int64) (venues map[int64]*model.Venue, err error)
	// cache: -nullcache=&model.Place{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	Place(c context.Context, id int64) (place *model.Place, err error)
	// cache: -nullcache=[]*model.Screen{{ProjectID:-1}} -check_null_code=len($)==1&&$[0].ProjectID==-1
	ScListByItem(c context.Context, pid []int64) (info map[int64][]*model.Screen, err error)
	// cache: -nullcache=&model.Screen{ProjectID:-1} -check_null_code=$!=nil&&$.ProjectID==-1
	ScList(c context.Context, sids []int64) (info map[int64]*model.Screen, err error)
	// cache: -nullcache=&model.TicketInfo{TicketPrice:model.TicketPrice{ProjectID:-1}} -check_null_code=$!=nil&&$.ProjectID==-1
	TkList(c context.Context, tids []int64) (info map[int64]*model.TicketInfo, err error)
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c: c,
		// orm
		db:     orm.NewMySQL(c.ORM),
		redis:  redis.NewPool(c.Redis.Master),
		expire: int32(time.Duration(c.Redis.Expire) / time.Second),
		cache:  fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		es: elastic.NewElastic(&elastic.Config{
			Host:       c.URL.ElasticHost,
			HTTPClient: c.HTTPClient.Read,
		}),
	}
	return
}

// Ping ping 方法
func (d *Dao) Ping(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("PING")
	if err != nil {
		return
	}
	return d.db.DB().PingContext(c)
}

// Close 关闭redis 和 db 连接
func (d *Dao) Close() (err error) {
	d.redis.Close()
	d.db.Close()
	return
}

// BeginTran 开启事务
func (d *Dao) BeginTran(c context.Context) (tx *gorm.DB, err error) {
	tx = d.db.Begin()
	if tx.Error != nil {
		err = tx.Error
		tx = nil
		log.Error("开启事务失败:%s", err)
	}
	return
}

// CommitTran 提交事务
func (d *Dao) CommitTran(c context.Context, tx *gorm.DB) (err error) {
	if err = tx.Commit().Error; err != nil {
		tx = nil
		log.Error("提交事务失败:%s", err)
	}
	return
}

package dao

import (
	"context"
	"time"

	"go-common/app/common/openplatform/encoding"

	acc "go-common/app/service/main/account/api"
	itemv1 "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/conf"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/library/cache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/naming/discovery"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/rpc/warden/resolver"
	"go-common/library/queue/databus"
)

//Dao 数据操作层结构体
type Dao struct {
	c           *conf.Config
	db          *sql.DB
	redis       *redis.Pool
	expire      int32
	cache       *cache.Cache
	httpClientR *bm.Client
	itemClient  itemv1.ItemClient
	accClient   acc.AccountClient
	databus     *databus.Databus
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -nullcache=&model.Promotion{PromoID:-1} -check_null_code=$!=nil&&$.PromoID==-1
	Promo(c context.Context, promoID int64) (*model.Promotion, error)
	// cache: -nullcache=&model.PromotionGroup{GroupID:-1} -check_null_code=$!=nil&&$.GroupID==-1
	PromoGroup(c context.Context, groupID int64) (*model.PromotionGroup, error)
	// cache: -nullcache=&model.PromotionOrder{OrderID:-1} -check_null_code=$!=nil&&$.OrderID==-1
	PromoOrder(c context.Context, orderID int64) (*model.PromotionOrder, error)
	// cache: -nullcache=[]*model.PromotionOrder{{GroupID:-1}} -check_null_code=len($)==1&&$[0].GroupID==-1
	PromoOrders(c context.Context, groupID int64) ([]*model.PromotionOrder, error)

	// cache: -nullcache=[]*model.OrderMain{{OrderID:-1}} -check_null_code=len($)==1&&$[0].OrderID==-1
	Orders(context.Context, *model.OrderMainQuerier) ([]*model.OrderMain, error)
	// cache: -nullcache=-1 -check_null_code=$==-1
	OrderCount(context.Context, *model.OrderMainQuerier) (int64, error)
	// cache: -nullcache=&model.OrderDetail{OrderID:0} -check_null_code=$!=nil&&$.OrderID==0
	OrderDetails(context.Context, []int64) (map[int64]*model.OrderDetail, error)
	// cache: -nullcache=[]*model.OrderSKU{{OrderID:-1}} -check_null_code=len($)==1&&$[0].OrderID==-1
	OrderSKUs(context.Context, []int64) (map[int64][]*model.OrderSKU, error)
	// cache: -nullcache=&model.OrderPayCharge{ChargeID:""} -check_null_code=$!=nil&&$.ChargeID==""
	OrderPayCharges(context.Context, []int64) (map[int64]*model.OrderPayCharge, error)

	// cache:
	SkuByItemID(c context.Context, itemID int64) (map[string]*model.SKUStock, error)
	// cache:
	GetSKUs(c context.Context, skuIds []int64, withNewStock bool) (map[int64]*model.SKUStock, error)
	// cache:
	Stocks(c context.Context, keys []int64, isLocked bool) (res map[int64]int64, err error)

	// cache: -nullcache=[]*model.Ticket{{ID:-1}} -check_null_code=len($)==1&&$[0].ID==-1
	TicketsByOrderID(c context.Context, orderID int64) (res []*model.Ticket, err error)
	// cache: -nullcache=[]*model.Ticket{{ID:-1}} -check_null_code=len($)==1&&$[0].ID==-1
	TicketsByScreen(c context.Context, screenID int64, UID int64) (res []*model.Ticket, err error)
	// cache: -nullcache=&model.Ticket{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	TicketsByID(c context.Context, id []int64) (res map[int64]*model.Ticket, err error)
	// cache: -nullcache=&model.TicketSend{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	TicketSend(c context.Context, SendTID []int64, TIDType string) (res map[int64]*model.TicketSend, err error)
}

// FIXME this just a example
func newItemClient(cfg *warden.ClientConfig) itemv1.ItemClient {
	cc, err := warden.NewClient(cfg).Dial(context.Background(), "discovery://default/ticket.service.item")
	if err != nil {
		panic(err)
	}
	return itemv1.NewItemClient(cc)
}

func newAccClient(cfg *warden.ClientConfig) acc.AccountClient {
	cc, err := warden.NewClient(cfg).Dial(context.Background(), "discovery://default/account.service")
	if err != nil {
		panic(err)
	}
	return acc.NewAccountClient(cc)
}

//New 根据配置文件 生成一个 Dao struct
func New(c *conf.Config) (d *Dao) {
	resolver.Register(discovery.New(nil))

	d = &Dao{
		c:           c,
		db:          sql.NewMySQL(c.DB.Master),
		redis:       redis.NewPool(c.Redis.Master),
		cache:       cache.New(1, 1024),
		expire:      int32(time.Duration(c.Redis.Expire) / time.Second),
		httpClientR: bm.NewClient(c.HTTPClient.Read),
		itemClient:  newItemClient(c.GRPCClient["item"]),
		accClient:   newAccClient(c.GRPCClient["account"]),
		databus:     databus.New(c.Databus["update"]),
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
	err = d.db.Ping(c)
	return
}

//Close 关闭redis 和 db 连接
func (d *Dao) Close() (err error) {
	d.redis.Close()
	d.db.Close()
	return
}

//BeginTx 开启事务
func (d *Dao) BeginTx(c context.Context) (conn *sql.Tx, err error) {
	return d.db.Begin(c)
}

//LogX 记录日志，args:请求参数 res:正常返回 err:错误返回 ld:更多日志项
func LogX(ctx context.Context, args interface{}, res interface{}, err error, ld ...log.D) {
	l := len(ld)
	u := 0
	ld1 := make([]log.D, l)
	for i := 0; i < l; i++ {
		if ld[i].Key != "" {
			ld1[u] = ld[i]
			u++
		}
	}
	ld1 = ld1[:u]
	if args != nil {
		ld1 = append(ld1, log.KV("args", encoding.JSON(args)))
	}
	if err != nil {
		ld1 = append(ld1, log.KV("log", err.Error()))
		log.Errorv(ctx, ld1...)
	} else if log.V(3) {
		if res == nil {
			log.Infov(ctx, ld1...)
			return
		}
		ld1 = append(ld1, log.KV("log", encoding.JSON(res)))
		log.Infov(ctx, ld1...)
	}
}

//DatabusPub 向databus发布消息
func (d *Dao) DatabusPub(ctx context.Context, action string, data interface{}) error {
	type input struct {
		Action string      `json:"action"`
		Data   interface{} `json:"data"`
	}
	err := d.databus.Send(ctx, action, &input{action, data})
	if err != nil {
		log.Error("pub databus failed action:%s, data:%+v", action, data)
	}
	return err
}

package newbiedao

import (
	"context"
	"time"

	"go-common/app/interface/main/growup/conf"
	"go-common/app/interface/main/growup/model"
	accApi "go-common/app/service/main/account/api"

	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	// Categories cache categories
	Categories map[int64]*model.Category
	// RecommendUpList cache recommend up list
	RecommendUpList map[int64]map[int64]*model.RecommendUp
)

// Dao def dao struct
type Dao struct {
	c  *conf.Config
	db *sql.DB
	// search
	httpRead *bm.Client
	// grpc
	accGRPC accApi.AccountClient
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:  c,
		db: sql.NewMySQL(c.DB.Growup),
		// search
		httpRead: bm.NewClient(c.HTTPClient.Read),
	}
	var err error
	if d.accGRPC, err = accApi.NewClient(c.AccCliConf); err != nil {
		panic(err)
	}

	d.loadCache()
	go func() {
		t := time.Tick(10 * time.Minute)
		for {
			d.loadCache()
			<-t
		}
	}()
	return d
}

// Ping ping db
func (d *Dao) Ping(c context.Context) (err error) {
	if err = d.db.Ping(c); err != nil {
		log.Error("d.db.Ping error(%v)", err)
		return
	}
	return
}

// Close close db conn
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

// loodCache load cache
func (d *Dao) loadCache() {
	_ = d.GetCategories(context.Background())
	log.Info("refresh categories cache: %+v", Categories)

	_ = d.GetRecommendUpList(context.Background())
	log.Info("refresh recommend up list cache: %+v", RecommendUpList)
}

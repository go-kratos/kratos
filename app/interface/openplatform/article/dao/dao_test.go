package dao

import (
	"context"
	"flag"
	"path/filepath"
	"strings"

	"go-common/app/interface/openplatform/article/conf"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/cache/redis"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

var (
	dataMID    = int64(1)
	noDataMID  = int64(10000)
	_noData    = int64(1000000)
	d          *Dao
	categories = []*artmdl.Category{
		&artmdl.Category{Name: "游戏", ID: 1},
		&artmdl.Category{Name: "动漫", ID: 2},
	}
	art = artmdl.Article{
		Meta: &artmdl.Meta{
			ID:              100,
			Category:        categories[0],
			Title:           "隐藏于时区记忆中的,是希望还是绝望!",
			Summary:         "说起日本校服,第一个浮现在我们脑海中的必然是那象征着青春阳光 蓝白色相称的水手服啦. 拉色短裙配上洁白的直袜",
			BannerURL:       "http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg",
			TemplateID:      1,
			State:           0,
			Author:          &artmdl.Author{Mid: 123, Name: "爱蜜莉雅", Face: "http://i1.hdslb.com/bfs/face/5c6109964e78a84021299cdf71739e21cd7bc208.jpg"},
			Reprint:         0,
			ImageURLs:       []string{"http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg"},
			OriginImageURLs: []string{"http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg", "http://i2.hdslb.com/bfs/archive/b5727f244d5c7a34c1c0e78f49765d09ff30c129.jpg"},
			PublishTime:     1495784507,
			Tags:            []*artmdl.Tag{},
			Stats:           &artmdl.Stats{Favorite: 100, Like: 10, View: 500, Dislike: 1, Share: 99},
		},
		Content: "content",
	}
)

func CleanCache() {
	c := context.TODO()
	pool := redis.NewPool(conf.Conf.Redis)
	pool.Get(c).Do("FLUSHDB")
}

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	d.httpClient.SetTransport(gock.DefaultTransport)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(d)
	}
}

func WithMysql(f func(d *Dao)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(d)
	}
}

func WithCleanCache(f func()) func() {
	return func() {
		Reset(func() { CleanCache() })
		f()
	}
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func ctx() context.Context {
	return context.Background()
}

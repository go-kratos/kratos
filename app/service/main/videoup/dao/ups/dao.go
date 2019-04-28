package ups

import (
	"context"
	"net/url"

	"go-common/app/service/main/videoup/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

const (
	_upSpecialURL = "/x/internal/uper/special"
)

// Dao is redis dao.
type Dao struct {
	c            *conf.Config
	httpW        *bm.Client
	upSpecialURI string
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		httpW:        bm.NewClient(c.HTTPClient.Write),
		upSpecialURI: c.Host.APICO + _upSpecialURL,
	}
	return d
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return
}

func (d *Dao) upSpecial(c context.Context, groupid string) (ups map[int64]int64, err error) {
	params := url.Values{}
	params.Set("group_id", groupid)
	ups = make(map[int64]int64)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			Mid int64 `json:"mid"`
		} `json:"data"`
	}
	if err = d.httpW.Get(c, d.upSpecialURI, "", params, &res); err != nil {
		log.Error("upSpecial error(%v) | upSpecialURI(%s)", err, d.upSpecialURI)
		return
	}
	log.Info("upSpecial url(%+v)|res(%+v)", d.upSpecialURI, res)
	if res.Code != 0 {
		log.Error("upSpecial api url(%s) res(%v);, code(%d)", d.upSpecialURI, res, res.Code)
		return
	}
	for _, item := range res.Data {
		ups[item.Mid] = item.Mid
	}
	return
}

// GrayCheckUps GrayCheckUps, type = 8
func (d *Dao) GrayCheckUps(c context.Context, gid int) (checkUps map[int64]int64, err error) {
	return d.upSpecial(c, strconv.Itoa(gid))
}

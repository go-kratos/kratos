package account

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/interface/main/answer/conf"
	"go-common/app/interface/main/answer/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// Dao is elec dao.
type Dao struct {
	c        *conf.Config
	mc       *memcache.Pool
	client   *bm.Client
	pendant  string
	beFormal string
	extraIds string
}

const (
	_wasFormal = -659
)

// New pendant dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:        c,
		mc:       memcache.NewPool(c.Memcache.Config),
		client:   bm.NewClient(c.HTTPClient.Normal),
		pendant:  c.Host.API + _multiGivePendant,
		beFormal: c.Host.Account + _beFormal,
		extraIds: c.Host.ExtraIds,
	}
	return
}

const (
	_multiGivePendant = "/x/internal/pendant/multiGrantByMid"
	_beFormal         = "/api/internal/member/beFormal"
)

// GivePendant send user pendant
func (d *Dao) GivePendant(c context.Context, mid int64, pid int64, days int, ip string) (err error) {
	log.Info(" GivePendant (%d,%d,%d) ", mid, pid, days)
	params := url.Values{}
	params.Set("mids", strconv.FormatInt(mid, 10))
	params.Set("pid", strconv.FormatInt(pid, 10))
	params.Set("expire", strconv.FormatInt(int64(days), 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.pendant, ip, params, &res); err != nil {
		log.Error("pendant url(%s) error(%v)", d.pendant+"?"+params.Encode(), err)
		log.Error("GivePendant(%d,%d),err:%+v", mid, pid, err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("pendant GivePendant failed(%v)", res.Code)
		log.Error(" d.client.Get(%s) error(%v)", d.pendant+"?"+params.Encode(), err)
		return
	}
	return
}

// BeFormal become a full member
func (d *Dao) BeFormal(c context.Context, mid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.beFormal, ip, params, &res); err != nil {
		err = errors.Wrapf(err, "beFormal url(%s)", d.beFormal+"?"+params.Encode())
		log.Error("BeFormal(%d),err:%+v", mid, err)
		return
	}
	if res.Code != 0 && res.Code != _wasFormal {
		err = errors.WithStack(fmt.Errorf("beFormal(%d) failed(%v)", mid, res.Code))
		log.Error("BeFormal(%d),res:%+v", mid, res)
		return
	}
	log.Info("beFormal suc(%d) ", mid)
	return
}

// ExtraIds BigData Extra Question ids.
func (d *Dao) ExtraIds(c context.Context, mid int64, ip string) (done []int64, pend []int64, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data *model.ExtraBigData
	}
	if err = d.client.Get(c, d.extraIds, ip, params, &res); err != nil {
		log.Error("ExtraIds url(%s) error(%v)", d.extraIds+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data == nil {
		err = fmt.Errorf("ExtraIds failed(%v)", res.Code)
		log.Error(" d.client.Get(%s) res(%v) error(%v)", d.extraIds+"?"+params.Encode(), res, err)
		return
	}
	done = res.Data.Done
	pend = res.Data.Pend
	return
}

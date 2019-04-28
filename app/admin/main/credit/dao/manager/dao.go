package manager

import (
	"context"
	"net/url"
	"strconv"
	"sync"

	"go-common/app/admin/main/credit/conf"
	"go-common/app/admin/main/credit/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/errgroup"
)

const (
	_managersURI     = "/x/admin/manager/users"
	_managerTotalURI = "/x/admin/manager/users/total"
)

// Dao struct info of Dao.
type Dao struct {
	// http
	client *bm.Client
	// conf
	c               *conf.Config
	managersURL     string
	managerTotalURL string
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// conf
		c: c,
		// http client
		client: bm.NewClient(c.HTTPClient),
	}
	d.managersURL = c.Host.Manager + _managersURI
	d.managerTotalURL = c.Host.Manager + _managerTotalURI
	return
}

// Managers get manager users info.
func (d *Dao) Managers(c context.Context) (manMap map[int64]string, err error) {
	var (
		count, page int64
		g           errgroup.Group
		l           sync.RWMutex
	)
	if count, err = d.ManagerTotal(c); err != nil {
		log.Error("d.ManagerTotal error(%v)", err)
		return
	}
	if count <= 0 {
		return
	}
	manMap = make(map[int64]string, count)
	ps := int64(500)
	pageNum := count / ps
	if count%ps != 0 {
		pageNum++
	}
	for page = 1; page <= pageNum; page++ {
		tmpPage := page
		g.Go(func() (err error) {
			mi, err := d.Manager(c, tmpPage, ps)
			if err != nil {
				log.Error("d.Manager(%d,%d) error(%v) ", tmpPage, ps, err)
				err = nil
				return
			}
			for _, v := range mi {
				l.Lock()
				manMap[v.OID] = v.Uname
				l.Unlock()
			}
			return
		})
	}
	g.Wait()
	return
}

// Manager  get manager users.
func (d *Dao) Manager(c context.Context, pn, ps int64) (mi []*model.MangerInfo, err error) {
	params := url.Values{}
	params.Set("pn", strconv.FormatInt(pn, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Items []*model.MangerInfo `json:"items"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.managersURL, "", params, &res); err != nil {
		log.Error("Manager(%s) error(%v)", d.managersURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Warn("Manager(%s) code(%d) data(%+v)", d.managersURL+"?"+params.Encode(), res.Code, res.Data)
		return
	}
	if res.Data != nil {
		mi = res.Data.Items
	}
	return
}

// ManagerTotal get manager user total.
func (d *Dao) ManagerTotal(c context.Context) (count int64, err error) {
	params := url.Values{}
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Total int64 `json:"total"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.managerTotalURL, "", params, &res); err != nil {
		log.Error("ManagerTotal(%s) error(%v)", d.managerTotalURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Warn("ManagerTotal(%s) code(%d) data(%+v)", d.managerTotalURL+"?"+params.Encode(), res.Code, res.Data)
		return
	}
	if res.Data != nil {
		count = res.Data.Total
	}
	return
}

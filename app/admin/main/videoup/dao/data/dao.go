package data

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/videoup/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
	"gopkg.in/h2non/gock.v1"
	"strings"
	"strconv"
)

// data.bilibili.co/recsys/related?key=XXAVID
const (
	_relatedURL  = "/recsys/related"
	_moniOidsURL = "/x/internal/aegis/monitor/result/oids"
)

// Dao is search dao
type Dao struct {
	c           *bm.ClientConfig
	httpClient  *bm.Client
	relatedURI  string
	moniOidsURI string
}

var (
	d *Dao
)

// New new search dao
func New(c *conf.Config) *Dao {
	return &Dao{
		c:           c.HTTPClient.Read,
		httpClient:  bm.NewClient(c.HTTPClient.Read),
		relatedURI:  c.Host.Data + _relatedURL,
		moniOidsURI: c.Host.API + _moniOidsURL,
	}
}

// ArchiveRelated get related archive from ai
func (d *Dao) ArchiveRelated(c context.Context, aidarr []int64) (aids string, err error) {
	params := url.Values{}
	params.Set("key", xstr.JoinInts(aidarr))

	res := new(struct {
		Code int `json:"code"`
		Data []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"data"`
	})
	if err = d.httpClient.Get(c, d.relatedURI, "", params, res); err != nil || res == nil {
		log.Error(" d.httpClient.Get error(%v)", err)
		return
	}
	log.Info("ArchiveRelated aids(%v) res(%+v)", aids, res)
	if res.Code != 0 {
		err = fmt.Errorf("data.bilibili.co错误(%d)", res.Code)
		log.Error(" d.httpClient.Get res(%+v)", res)
		return
	}
	if len(res.Data) > 0 {
		for _, item := range res.Data {
			if len(item.Value) > 0 {
				if len(aids) == 0 {
					aids = item.Value
				} else {
					aids += "," + item.Value
				}
			}
		}
	}
	return
}

// MonitorOids 获取监控的id
func (d *Dao) MonitorOids(c context.Context, id int64) (oidMap map[int64]int, err error) {
	oidMap = make(map[int64]int)
	params := url.Values{}
	params.Set("id", strconv.Itoa(int(id)))

	res := new(struct {
		Code int `json:"code"`
		Data []struct {
			OID  int64 `json:"oid"`
			Time int   `json:"time"`
		} `json:"data"`
	})
	if err = d.httpClient.Get(c, d.moniOidsURI, "", params, res); err != nil || res == nil {
		log.Error("d.MonitorOids() d.httpClient.Get(%s,%v) error(%v)", d.moniOidsURI, params, err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("monitor return code(%d)", res.Code)
		log.Error("d.MonitorOids() d.httpClient.Get(%s,%v) res(%v)", d.moniOidsURI, params, res)
		return
	}
	for _, v := range res.Data {
		oidMap[v.OID] = v.Time
	}
	return
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.httpClient.SetTransport(gock.DefaultTransport)
	return r
}

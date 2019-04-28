package alarm

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go-common/app/service/main/resource/conf"
	"go-common/app/service/main/resource/model"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

// Dao is redis dao.
type Dao struct {
	c          *conf.Config
	netClient  *http.Client
	httpClient *httpx.Client
}

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:          c,
		httpClient: httpx.NewClient(c.HTTPClient),
		netClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
	return d
}

func (d *Dao) CheckURL(originURL string, wis []*model.ResWarnInfo) {
	var (
		url  string
		req  *http.Request
		resp *http.Response
		err  error
	)
	if strings.HasPrefix(originURL, "https://") {
		log.Info("CheckURL url(%s) is https ,replace to http", originURL)
		url = strings.Replace(originURL, "https://", "http://", -1)
	} else if !strings.HasPrefix(originURL, "http://") {
		log.Info("CheckURL url(%s) don't have https and http", originURL)
		url = "http://" + originURL
	} else {
		url = originURL
	}
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Error("CheckURL NewRequest(%v) error(%v)", url, err)
		return
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	resp, err = d.netClient.Do(req)
	if err != nil {
		log.Error("CheckURL url(%s) originURL(%s) error(%v)", url, originURL, err)
	} else if resp.StatusCode != http.StatusOK {
		log.Error("CheckURL url(%s) originURL(%s) code(%v) not OK ", url, originURL, resp.StatusCode)
		var sends = make(map[string][]*model.ResWarnInfo)
		for _, wi := range wis {
			sends[wi.UserName] = append(sends[wi.UserName], wi)
		}
		for userName, send := range sends {
			d.sendWeChartURL(context.TODO(), resp.StatusCode, userName, send)
		}
	}
}

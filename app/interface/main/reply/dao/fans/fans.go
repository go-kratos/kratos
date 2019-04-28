package fans

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

// Dao Dao
type Dao struct {
	fansReceivedListURL        string
	fansReceivedListHTTPClient *httpx.Client
}

// New New
func New(c *conf.Config) *Dao {
	d := &Dao{
		fansReceivedListURL:        c.Reply.FansReceivedListURL,
		fansReceivedListHTTPClient: httpx.NewClient(c.HTTPClient),
	}
	return d
}

// Fetch Fetch
func (dao *Dao) Fetch(c context.Context, uids []int64, mid int64, now time.Time) (map[int64]*reply.FansDetail, error) {
	fansMap := make(map[int64]*reply.FansDetail)
	if len(uids) == 0 {
		return fansMap, nil
	}
	params := url.Values{}
	params.Set("target_id", strconv.FormatInt(mid, 10))
	params.Set("source", strconv.FormatInt(2, 10))
	for index := range uids {
		params.Add("uid[]", strconv.FormatInt(uids[index], 10))
	}
	var res struct {
		Code    int                 `json:"code"`
		Message string              `json:"msg"`
		Data    []*reply.FansDetail `json:"data"`
	}
	if err := dao.fansReceivedListHTTPClient.Get(c, dao.fansReceivedListURL, "", params, &res); err != nil {
		log.Error("fansFetch url(%v),err (%v)", dao.fansReceivedListURL+"?"+params.Encode(), err)
		return fansMap, err
	}
	if res.Code != 0 {
		return fansMap, nil
	}
	for _, d := range res.Data {
		fansMap[d.UID] = d
	}
	return fansMap, nil
}

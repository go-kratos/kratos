package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"go-common/app/job/main/reply/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_appIDRecord = "reply_record"
	_updateURL   = "/x/internal/search/reply/update"
)

// Dao Dao
type Dao struct {
	searchHTTPClient *bm.Client
	searchUpdateURI  string
}

// New New
func New(c *conf.Config) *Dao {
	return &Dao{
		searchHTTPClient: bm.NewClient(c.HTTPClient),
		searchUpdateURI:  c.Host.API + _updateURL,
	}
}

// DelReply DelReply
func (dao *Dao) DelReply(c context.Context, rpid, oid, mid int64, state int8) (err error) {
	return dao.update(c, rpid, oid, mid, state)
}

func (dao *Dao) update(c context.Context, rpid, oid, mid int64, state int8) (err error) {
	type updateRecord struct {
		ID    int64 `json:"id"`
		Oid   int64 `json:"oid"`
		Mid   int64 `json:"mid"`
		State int8  `json:"state"`
	}
	var (
		res struct {
			Code int    `json:"code"`
			Msg  string `json:"message"`
		}
	)
	records := make([]*updateRecord, 0)
	record := &updateRecord{}
	record.ID = rpid
	record.Oid = oid
	record.Mid = mid
	record.State = state
	records = append(records, record)
	recordsStr, _ := json.Marshal(records)
	params := url.Values{}
	params.Set("appid", _appIDRecord)
	params.Set("data", string(recordsStr))
	if err = dao.searchHTTPClient.Post(c, dao.searchUpdateURI, "", params, &res); err != nil {
		log.Error("bm.Post(%s) failed error(%v)", dao.searchUpdateURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = fmt.Errorf("update reply es records failed")
	}
	log.Info("updateSearch: %s post:%s ret:%v", dao.searchUpdateURI, params.Encode(), res)
	return
}

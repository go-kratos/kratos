package reply

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/reply/conf"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

// CreditUserDao CreditUserDao
type CreditUserDao struct {
	requestURL string
	httpClient *httpx.Client
}

// NewCreditDao NewCreditDao
func NewCreditDao(c *conf.Config) *CreditUserDao {
	d := &CreditUserDao{
		httpClient: httpx.NewClient(c.HTTPClient),
		requestURL: c.Reply.CreditUserURL,
	}
	return d
}

// Result Result
type Result struct {
	Code int `json:"code"`
	Data map[int64]struct {
		Mid     int64 `json:"mid"`
		Expired int64 `json:"expired"`
		Status  int   `json:"status"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// IsCreditUser IsCreditUser
func (dao *CreditUserDao) IsCreditUser(c context.Context, mid int64) (bool, error) {
	var res Result
	params := url.Values{}
	params.Set("mids", fmt.Sprintf("%d", mid))
	if err := dao.httpClient.Get(c, dao.requestURL, "", params, &res); err != nil {
		log.Error("call 风纪委身份验证 url(%s) error(%v)", dao.requestURL+"?"+params.Encode(), err)
		return false, err
	}
	if res.Code != 0 {
		err := fmt.Errorf("call 风纪委身份验证 url(%s) error(%v)", dao.requestURL+"?"+params.Encode(), res.Code)
		log.Error("%v", err)
		return false, err
	}
	log.Info("call 风纪委身份验证 successful. url(%s)", dao.requestURL+"?"+params.Encode())
	return res.Data[mid].Status == 1, nil
}

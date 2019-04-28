package realname

import (
	"context"
	"net/http"
	"net/url"

	"go-common/app/interface/main/account/conf"
	"go-common/library/log"

	"github.com/pkg/errors"
)

type respAlipay struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

func (r *respAlipay) Error() error {
	if r.Code == "10000" {
		return nil
	}
	return errors.Errorf("alipay response failed , code : %s, msg : %s, sub_code : %s, sub_msg : %s", r.Code, r.Msg, r.SubCode, r.SubMsg)
}

// AlipayInit .
func (d *Dao) AlipayInit(c context.Context, param url.Values) (bizno string, err error) {
	var (
		req *http.Request
	)
	url := conf.Conf.Realname.Alipay.Gateway + "?" + param.Encode()
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		err = errors.Wrapf(err, "http.NewRequest(GET,%s)", url)
		return
	}
	var resp struct {
		Resp struct {
			respAlipay
			Bizno string `json:"biz_no"`
		} `json:"zhima_customer_certification_initialize_response"`
		Sign string `json:"sign"`
	}
	if err = d.client.Do(c, req, &resp); err != nil {
		return
	}
	log.Info("Realname alipay init \n\tparam : %+v \n\tresp : %+v", param, resp)
	if err = resp.Resp.Error(); err != nil {
		return
	}
	bizno = resp.Resp.Bizno
	return
}

// AlipayQuery .
func (d *Dao) AlipayQuery(c context.Context, param url.Values) (pass bool, reason string, err error) {
	var (
		req *http.Request
	)
	url := conf.Conf.Realname.Alipay.Gateway + "?" + param.Encode()
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		err = errors.Wrapf(err, "http.NewRequest(GET,%s)", url)
		return
	}
	var resp struct {
		Resp struct {
			respAlipay
			Passed          string `json:"passed"`
			FailedReason    string `json:"failed_reason"`
			IdentityInfo    string `json:"identity_info"`
			AttributeInfo   string `json:"attribute_info"`
			ChannelStatuses string `json:"channel_statuses"`
		} `json:"zhima_customer_certification_query_response"`
		Sign string `json:"sign"`
	}
	if err = d.client.Do(c, req, &resp); err != nil {
		return
	}
	log.Info("Realname alipay query \n\tparam : %+v \n\tresp : %+v", param, resp)
	if err = resp.Resp.Error(); err != nil {
		return
	}
	if resp.Resp.Passed == "true" {
		pass = true
	} else {
		pass = false
	}
	reason = resp.Resp.FailedReason
	return
}

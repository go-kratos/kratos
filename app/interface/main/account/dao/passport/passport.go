package passport

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/dao/account"
	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	fastRegUri    = "/intranet/acc/user/source"
	useridUri     = "/intranet/acc/userid"
	testUsername  = "/api/reg/testUserName"
	queryByMidUri = "/intranet/acc/queryByMid"
	updateName    = "/intranet/acc/updateUname"
)

// Dao dao
type Dao struct {
	c      *conf.Config
	client *bm.Client
}

// New new
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: bm.NewClient(c.HTTPClient.Normal),
	}
	return
}

func (d *Dao) FastReg(c context.Context, mid int64, ip string) (isRegFast bool, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			FastReg bool `json:"fastReg"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err = d.client.Get(c, d.c.Host.Passport+fastRegUri, ip, params, &res); err != nil {
		log.Error("fastReg url(%s) error(%v)", fastRegUri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("fastReg url(%s) res(%v)", fastRegUri+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	isRegFast = res.Data.FastReg
	return
}

//UserID user id.
func (d *Dao) UserID(c context.Context, mid int64, ip string) (userID string, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			UserID string `json:"userid"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err = d.client.Get(c, d.c.Host.Passport+useridUri, ip, params, &res); err != nil {
		log.Error("UserID url(%s) error(%v)", useridUri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("userID url(%s) res(%v)", useridUri+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	userID = res.Data.UserID
	return
}

// TestUserName is.
func (d *Dao) TestUserName(c context.Context, name string, mid int64, ip string) error {
	params := url.Values{}
	params.Set("user_name", name)
	params.Set("mid", strconv.FormatInt(mid, 10))

	var res struct {
		Code int `json:"code"`
	}
	if err := d.client.Get(c, d.c.Host.Passport+testUsername, ip, params, &res); err != nil {
		log.Error("Failed to test username: %+v: %+v", params, err)
		return err
	}
	if res.Code != 0 {
		log.Error("Failed to test username with code: %+v: %d", params, res.Code)
		return account.ParseJavaCode(res.Code)
	}
	return nil
}

// QueryByMid is.
func (d *Dao) QueryByMid(c context.Context, mid int64, ip string) (*model.PassportProfile, error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))

	var res struct {
		Code int                    `json:"code"`
		Data *model.PassportProfile `json:"data"`
	}
	if err := d.client.Get(c, d.c.Host.Passport+queryByMidUri, ip, params, &res); err != nil {
		log.Error("Failed to query by mid: %+v: %+v", params, err)
		return nil, err
	}
	if res.Code != 0 {
		log.Error("Failed to query by mid with code: %+v: %d", params, res.Code)
		return nil, account.ParseJavaCode(res.Code)
	}
	return res.Data, nil
}

// UpdateName is.
func (d *Dao) UpdateName(c context.Context, mid int64, name, ip string) error {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("uname", name)
	var res struct {
		Code int `json:"code"`
	}
	if err := d.client.Post(c, d.c.Host.Passport+updateName, ip, params, &res); err != nil {
		log.Error("Failed to update uname params: %+v: %v", params, err)
		return err
	}
	if res.Code != 0 {
		log.Error("Failed to update uname, params: %+v,code: %d", params, res.Code)
		return account.ParseJavaCode(res.Code)
	}
	return nil
}

package account

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"

	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_smsSendURI = "/x/internal/sms/send"
	_offTcode   = "acc_20"
	_vcExpire   = 30 * 60 // 30min
)

func vcodeKey(mid int64, mobile string) string {
	return fmt.Sprintf("mv_%d_%s", mid, mobile)
}

// SendMobileVerify is.
func (d *Dao) SendMobileVerify(ctx context.Context, vcode, country int64, mobile, ip string) error {
	params := url.Values{}
	// params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("mobile", mobile)
	params.Set("country", strconv.FormatInt(country, 10))
	params.Set("tcode", _offTcode)

	tparam, err := OffiVerifyParam(vcode)
	if err != nil {
		return err
	}
	params.Set("tparam", tparam)

	var res struct {
		Code int `json:"code"`
	}
	if err := d.client.Post(ctx, d.smsSendURI, ip, params, &res); err != nil {
		log.Error("Failed to send sms send request: %+v: %+v", params, err)
		return err
	}
	if res.Code != 0 {
		log.Error("Failed to send sms send requets: %+v: %d", params, res.Code)
		return ecode.Int(res.Code)
	}
	return nil
}

// OffiVerifyParam is.
func OffiVerifyParam(vcode int64) (string, error) {
	p := map[string]string{
		"identify_code": strconv.FormatInt(vcode, 10),
	}
	ps, err := json.Marshal(p)
	return string(ps), err
}

// GenVerifyCode is.
func (d *Dao) GenVerifyCode(ctx context.Context, mid int64, mobile string) (int64, error) {
	vcode := int64(rand.Intn(9999-1000) + 1000)
	key := vcodeKey(mid, mobile)
	conn := d.mc.Get(ctx)
	vcs := strconv.FormatInt(vcode, 10)
	defer conn.Close()
	if err := conn.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(vcs),
		Expiration: _vcExpire,
	}); err != nil {
		log.Error("conn.Set(%s, %d) error(%v)", key, vcode, err)
		return 0, err
	}
	return vcode, nil
}

// GetVerifyCode is.
func (d *Dao) GetVerifyCode(ctx context.Context, mid int64, mobile string) (int64, error) {
	key := vcodeKey(mid, mobile)
	conn := d.mc.Get(ctx)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(item.Value), 10, 64)
}

// DelVerifyCode is.
func (d *Dao) DelVerifyCode(ctx context.Context, mid int64, mobile string) error {
	key := vcodeKey(mid, mobile)
	conn := d.mc.Get(ctx)
	defer conn.Close()
	return conn.Delete(key)
}

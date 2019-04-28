package account

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/account/model"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// CommercialSign is
// 商业产品部专用签名
func CommercialSign(params url.Values) (query string, err error) {
	if len(params) == 0 {
		return
	}
	if params.Get("appkey") == "" {
		err = errors.New("utils http get must have parameter appkey")
		return
	}
	if params.Get("appsecret") == "" {
		err = errors.New("utils http get must have parameter appsecret")
		return
	}
	if params.Get("sign") != "" {
		err = errors.New("utils http get must have not parameter sign")
		return
	}
	// sign
	secret := params.Get("appsecret")
	params.Del("appsecret")
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(tmp + secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	query = params.Encode()
	return
}

func (d *Dao) isBusinessAccount(ctx context.Context, mids []int64) (map[int64]bool, error) {
	cm := d.c.HTTPClient.Normal.Host[d.c.Host.CM]
	params := url.Values{}
	params.Set("mids", xstr.JoinInts(mids))
	params.Set("appkey", cm.Key)
	params.Set("appsecret", cm.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	query, err := CommercialSign(params)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(d.c.Host.CM)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	u.Path = _cmIsBusinessAccount
	u.RawQuery = query
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("Accept", "application/json")

	resp := &struct {
		Code    int             `json:"int"`
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    map[string]bool `json:"data"`
	}{}
	if err := d.client.Do(ctx, req, resp); err != nil {
		return nil, errors.WithStack(err)
	}
	if !resp.Success {
		return nil, errors.Errorf("Failed to request cm api: %+v", resp)
	}

	res := make(map[int64]bool, len(resp.Data))
	for smid, isBus := range resp.Data {
		mid, err := strconv.ParseInt(smid, 10, 64)
		if err != nil {
			log.Warn("Failed to parse mid: %s: %+v", smid, err)
			continue
		}
		res[mid] = isBus
	}
	return res, nil
}

// IsBusinessAccount is
func (d *Dao) IsBusinessAccount(ctx context.Context, mid int64) bool {
	res, err := d.isBusinessAccount(ctx, []int64{mid})
	if err != nil {
		log.Error("Failed to check is business account with mid: %d: %+v", mid, err)
		return false
	}
	return res[mid]
}

// BusinessAccountInfo is
func (d *Dao) BusinessAccountInfo(ctx context.Context, mid int64) (*model.CMAccountInfo, error) {
	cm := d.c.HTTPClient.Normal.Host[d.c.Host.CM]
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("appkey", cm.Key)
	params.Set("appsecret", cm.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	query, err := CommercialSign(params)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(d.c.Host.CM)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	u.Path = _cmBusinessAccountInfo
	u.RawQuery = query
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("Accept", "application/json")

	res := &struct {
		Code    int                 `json:"int"`
		Success bool                `json:"success"`
		Message string              `json:"message"`
		Data    model.CMAccountInfo `json:"data"`
	}{}
	if err := d.client.Do(ctx, req, res); err != nil {
		return nil, errors.WithStack(err)
	}
	if !res.Success {
		return nil, errors.Errorf("Failed to request cm api: %+v", res)
	}

	return &res.Data, nil
}

package vip

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	xhttp "net/http"
	"net/url"
	"strings"

	"go-common/library/conf/env"
	bm "go-common/library/net/http/blademaster"

	pkgerr "github.com/pkg/errors"
)

//Clientl Client is http client, for url sign.
type Clientl struct {
	client *bm.Client
	conf   *bm.ClientConfig
}

// NewClientl new a http client.
func NewClientl(c *bm.ClientConfig, client *bm.Client) *Clientl {
	// check appkey
	if c.Key == "" || c.Secret == "" {
		panic("http client must config appkey and appsecret")
	}
	cl := new(Clientl)
	cl.client = client
	cl.conf = c
	return cl
}

const (
	_httpHeaderRemoteIP = "x-backend-bili-real-ip"
	_noKickUserAgent    = "haoguanwei@bilibili.com "
)

func (cl *Clientl) get(c context.Context, host, path, realIP string, params url.Values, res interface{}) (err error) {
	req, err := cl.newGetRequest(host, path, realIP, params)
	if err != nil {
		return
	}
	return cl.client.Do(c, req, res)
}

// newGetRequest new get http request with  host, path, ip, values and headers, without sign.
func (cl *Clientl) newGetRequest(host, path, realIP string, params url.Values) (req *xhttp.Request, err error) {
	if params == nil {
		params = url.Values{}
	}
	params.Add("client_id", cl.conf.App.Key)
	enc := sign(params, path, cl.conf.App.Secret)
	ru := host + path
	if enc != "" {
		ru = ru + "?" + enc
	}
	req, err = xhttp.NewRequest(xhttp.MethodGet, ru, nil)
	if err != nil {
		err = pkgerr.Wrapf(err, "uri:%s", ru)
		return
	}
	const (
		_userAgent = "User-Agent"
	)
	if realIP != "" {
		req.Header.Set(_httpHeaderRemoteIP, realIP)
	}
	req.Header.Set(_userAgent, _noKickUserAgent+" "+env.AppID)
	return
}

func sign(params url.Values, path, secret string) string {
	if params == nil {
		params = url.Values{}
	}
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(path)
	b.WriteString("?")
	b.WriteString(tmp)
	b.WriteString(secret)
	mh := md5.Sum(b.Bytes())
	// query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(hex.EncodeToString(mh[:]))
	return qb.String()
}

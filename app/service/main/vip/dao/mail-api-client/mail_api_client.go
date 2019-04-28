package client

import (
	"context"
	"encoding/json"
	xhttp "net/http"
	"strings"

	bm "go-common/library/net/http/blademaster"

	pkgerr "github.com/pkg/errors"
)

const (
	_contentTypeJSON = "application/json"
)

//Client client is http client.
type Client struct {
	client *bm.Client
}

// NewClient new a http client.
func NewClient(client *bm.Client) *Client {
	cl := new(Client)
	cl.client = client
	return cl
}

// Get a json req http get.
func (client *Client) Get(c context.Context, uri string, params interface{}, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodGet, uri, params)
	if err != nil {
		return
	}
	return client.client.Do(c, req, res)
}

// Post a json req http post.
func (client *Client) Post(c context.Context, uri string, params interface{}, res interface{}) (err error) {
	req, err := client.NewRequest(xhttp.MethodPost, uri, params)
	if err != nil {
		return
	}
	return client.client.Do(c, req, res)
}

// NewRequest new http request with method, uri, ip, values and headers.
func (client *Client) NewRequest(method, uri string, params interface{}) (req *xhttp.Request, err error) {
	marshal, err := json.Marshal(params)
	if err != nil {
		err = pkgerr.Wrapf(err, "marshal:%v", params)
		return
	}
	req, err = xhttp.NewRequest(method, uri, strings.NewReader(string(marshal)))
	if err != nil {
		err = pkgerr.Wrapf(err, "method:%s,uri:%s", method, string(marshal))
		return
	}
	req.Header.Set("Content-Type", _contentTypeJSON)
	return
}

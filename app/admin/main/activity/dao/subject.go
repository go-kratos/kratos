package dao

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

// AddTags add tags from http request.
func (d *Dao) AddTags(c context.Context, tags string, ip string) (err error) {
	var res struct {
		Code int `json:"code"`
	}
	params := url.Values{}
	params.Set("tag_name", tags)
	if err = d.client.Post(c, d.actURLAddTags, ip, params, &res); err != nil {
		err = errors.Wrapf(err, "d.client.Post(%s)", d.actURLAddTags)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("res code(%v)", res)
	}
	return
}

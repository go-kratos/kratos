package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_wasFormal = -659
)

// BeFormal become a full member
func (d *Dao) BeFormal(c context.Context, mid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.xclient.Post(c, d.beFormal, ip, params, &res); err != nil {
		err = errors.Wrapf(err, "beFormal url(%s)", d.beFormal+"?"+params.Encode())
		return
	}
	if res.Code != 0 && res.Code != _wasFormal {
		err = errors.WithStack(fmt.Errorf("beFormal(%d) failed(%v)", mid, res.Code))
		return
	}
	log.Info("beFormal suc(%d) ", mid)
	return
}

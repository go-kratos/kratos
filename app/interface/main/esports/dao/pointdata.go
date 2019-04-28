package dao

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go-common/library/log"

	"github.com/pkg/errors"
)

// Leida get leidata.
func (d *Dao) Leida(c context.Context, url string) (rs []byte, err error) {
	if rs, err = d.ThirdGet(c, url); err != nil {
		log.Error("d.ThirdGet  url(%s) error(%+v)", url, err)
	}
	return
}

// ThirdGet get.
func (d *Dao) ThirdGet(c context.Context, url string) (res []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		err = errors.Wrapf(err, "ThirdGet http.NewRequest(%s)", url)
		return
	}
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.Leidata.Timeout))
	defer cancel()
	req = req.WithContext(ctx)
	if resp, err = d.ldClient.Do(req); err != nil {
		err = errors.Wrapf(err, "ThirdGet d.ldClient.Do(%s)", url)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("ThirdGet url(%s) resp.StatusCode(%v)", url, resp.StatusCode)
		return
	}
	res, err = ioutil.ReadAll(resp.Body)
	return
}

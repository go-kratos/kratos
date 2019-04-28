package dao

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// PvData get binary data from pvdata url
func (d *Dao) PvData(c context.Context, pvURL string) (res []byte, err error) {
	var (
		req    *http.Request
		resp   *http.Response
		cancel func()
	)
	if req, err = http.NewRequest("GET", pvURL, nil); err != nil {
		err = errors.Wrapf(err, "PvData http.NewRequest(%s)", pvURL)
		return
	}
	c, cancel = context.WithTimeout(c, time.Duration(d.c.Rule.VsTimeout))
	defer cancel()
	req = req.WithContext(c)
	if resp, err = d.vsClient.Do(req); err != nil {
		err = errors.Wrapf(err, "httpClient.Do(%s)", pvURL)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("PvData url(%s) resp.StatusCode(%v)", pvURL, resp.StatusCode)
		return
	}
	res, err = ioutil.ReadAll(resp.Body)
	return
}

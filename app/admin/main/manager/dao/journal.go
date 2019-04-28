package dao

import (
	"bufio"
	"net/http"
	"strings"

	"go-common/app/admin/main/manager/model"
	bm "go-common/library/net/http/blademaster"
)

const (
	_searchAuditURL  = "http://manager.bilibili.co/x/admin/search/log/audit"
	_searchActionURL = "http://manager.bilibili.co/x/admin/search/log/user_action"
)

// SearchLogAudit .
func (d *Dao) SearchLogAudit(c *bm.Context) (res *model.LogRes, err error) {
	res, err = d.combination(c, _searchAuditURL)
	return
}

// SearchLogAction .
func (d *Dao) SearchLogAction(c *bm.Context) (res *model.LogRes, err error) {
	res, err = d.combination(c, _searchActionURL)
	return
}

// combination .
func (d *Dao) combination(c *bm.Context, preURL string) (res *model.LogRes, err error) {
	params := c.Request.URL.RawQuery
	url := preURL + "?" + params
	cookie := c.Request.Header.Get("Cookie")
	req, err := http.NewRequest(http.MethodGet, url, bufio.NewReader(strings.NewReader(params)))
	if err != nil {
		return
	}
	req.Header.Set("Cookie", cookie)
	err = d.httpClient.Do(c, req, &res)
	return
}

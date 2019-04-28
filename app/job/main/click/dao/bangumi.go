package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/log"
)

const (
	_aid2epid  = "http://bangumi.bilibili.co/ext/internal/archive/aid2epid"
	_epidExist = "http://bangumi.bilibili.co/ext/internal/archive/aid/play"
	_isLegal   = int(1)
)

var seasonType = []int{1, 2, 3, 4, 5}

// LoadAllBangumi load all bangumi epid -> aid to map
func (d *Dao) LoadAllBangumi(c context.Context) (etam map[int64]int64, err error) {
	etam = make(map[int64]int64)
	for _, t := range seasonType {
		pageNum := 1
		for {
			var resp struct {
				Code    int             `json:"code"`
				Message string          `json:"message"`
				Result  map[int64]int64 `json:"result"`
			}
			p := url.Values{}
			p.Set("build", "0")
			p.Set("platform", "golang")
			p.Set("season_type", strconv.Itoa(t))
			p.Set("page_size", "1000")
			p.Set("page_no", strconv.Itoa(pageNum))
			// one time error,all return,wait for next update
			if err = d.client.Get(c, _aid2epid, "", p, &resp); err != nil {
				log.Error("d.client.Get(%s) error(%v)", _aid2epid+"?"+p.Encode(), err)
				return
			}
			// record the page number when result is empty
			if len(resp.Result) == 0 {
				log.Info("bangumi seasonType(%d) pageNo(%d) is end", t, pageNum)
				break
			}
			for epid, aid := range resp.Result {
				etam[epid] = aid
			}
			pageNum++
		}
	}
	return
}

// IsLegal check legal by aid epID seasonType
func (d *Dao) IsLegal(c context.Context, aid, epID int64, seasonType int) (islegal bool, err error) {
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Result  struct {
			Status int `json:"status"`
		} `json:"result"`
	}
	p := url.Values{}
	p.Set("build", "0")
	p.Set("platform", "golang")
	p.Set("season_type", strconv.Itoa(seasonType))
	p.Set("epid", strconv.FormatInt(epID, 10))
	p.Set("aid", strconv.FormatInt(aid, 10))
	if err = d.client.Get(c, _epidExist, "", p, &resp); err != nil {
		log.Error("d.client.Get(%s) error(%v)", _epidExist+"?"+p.Encode(), err)
		return
	}
	if resp.Result.Status != _isLegal {
		log.Error("aid(%d) epid(%d) seasonType(%d) is unlegal", aid, epID, seasonType)
		return
	}
	islegal = true
	return
}

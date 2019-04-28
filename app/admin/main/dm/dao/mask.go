package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-common/library/log"
)

const _mask = "/dl/api/masks/v1"

// GenerateMask ask AI to generate dm mask
func (d *Dao) GenerateMask(c context.Context, cid, mid int64, plat int8, force int8, aid int64, duration int64, typeID int32) (err error) {
	var (
		res struct {
			Code    int64  `json:"errcode"`
			Message string `json:"errmsg"`
		}
		params = make(map[string]interface{})
	)
	params["cid"] = cid
	params["mask_platform"] = plat
	params["force"] = force
	params["mid"] = mid
	params["aid"] = aid
	params["duration"] = duration
	params["region_2"] = typeID
	data, err := json.Marshal(params)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", params, err)
		return
	}
	reader := bytes.NewReader(data)
	req, err := http.NewRequest("POST", d.maskURI, reader)
	if err != nil {
		log.Error("http.NewRequest error(%v)", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	for i := 0; i < 3; i++ {
		if err = d.httpCli.Do(c, req, &res); err != nil {
			continue
		}
		if res.Code != 200 {
			err = fmt.Errorf("uri:%s,code:%d", d.maskURI, res.Code)
			continue
		}
		break
	}
	if err != nil {
		log.Error("d.GenerateMask(cid:%d, plat:%d) error(%v)", cid, plat, err)
	}
	return
}

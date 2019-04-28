package ugc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	appDao "go-common/app/job/main/tv/dao/app"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_reported     = 1
	_finishReport = "UPDATE ugc_video SET mark = %d WHERE cid IN (%s)"
)

// RepCidBatch reports cid info to VideoCloud api
func (d *Dao) RepCidBatch(c context.Context, cidReq []*ugcmdl.CidReq) (err error) {
	var (
		reportURL = d.conf.UgcSync.Cfg.ReportCidURL
		resp      = ugcmdl.CidResp{}
		jsonBody  []byte
	)
	if jsonBody, err = json.Marshal(cidReq); err != nil {
		log.Error("json.Marchal(%v) error(%v)", cidReq, err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, reportURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error("http.NewRequest err - %v", err)
		return
	}
	if err = d.client.Do(c, req, &resp); err != nil {
		log.Error("ReportCid jsonBody (%s) error(%v)", string(jsonBody), err)
		return
	}
	if resp.Code != 0 {
		return fmt.Errorf("ReportCid Error Code:%d, Msg:%s", resp.Code, resp.Message)
	}
	appDao.PromInfo("ReportCid:Succ")
	return
}

// FinishReport change's the ugc_video's state from 0 to 1, means it has already been reported
func (d *Dao) FinishReport(c context.Context, cids []int64) (err error) {
	if _, err = d.DB.Exec(c, fmt.Sprintf(_finishReport, _reported, xstr.JoinInts(cids))); err != nil {
		log.Error("FinishReport Error: %v", cids, err)
	}
	return
}

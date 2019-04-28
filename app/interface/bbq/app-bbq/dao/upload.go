package dao

import (
	"context"
	"go-common/library/ecode"
	"go-common/library/log"
	xhttp "net/http"
	"net/url"
	"strconv"
)

//MergeUploadReq ..
func (d *Dao) MergeUploadReq(c context.Context, rurl string, uploadID string, profile string, svid int64, auth string) (err error) {
	var (
		req *xhttp.Request
	)

	rurl = d.c.Upload.HTTPSchema + rurl
	param := make(url.Values)
	param.Set("uploadId", uploadID)
	param.Set("profile", profile)
	param.Set("biz_id", strconv.FormatInt(svid, 10))
	rurl = rurl + "?" + param.Encode()
	req, err = d.httpslowClient.NewRequest(xhttp.MethodPost, rurl, "", param)

	if err != nil {
		log.Errorw(c, "event", "MergeUploadReq d.httpClient.NewRequest err", "err", err)
		err = ecode.UploadFailed
		return
	}
	req.Header.Add("X-Upos-Auth", auth)
	if err = d.httpClient.Do(c, req, nil); err != nil {
		log.Errorw(c, "event", "MergeUploadReq d.httpClient.Do err", "err", err)
		err = ecode.UploadFailed
		return
	}
	return
}

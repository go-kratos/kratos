package lic

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_normalCode = "0000"
)

func (d *Dao) callLic(c context.Context, url string, xmlBody string) (result *model.Document, err error) {
	var resp []byte
	result = &model.Document{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(xmlBody)))
	if err != nil {
		log.Error("http.NewRequest err - %v", err)
		return
	}
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	if resp, err = d.client.Raw(c, req); err != nil {
		log.Error("ClientGet error[%v]", err)
		return
	}
	if err = xml.Unmarshal(resp, result); err != nil {
		log.Error("XML Unmarshal %s, Error %v", string(resp), err)
		return
	}
	if result == nil || result.Response == nil {
		err = errors.Wrap(err, "Response Empty Error")
		return
	}
	if result.Response.ResponseCode != _normalCode {
		err = fmt.Errorf("Response Code Error %s", result.Response.ResponseCode)
		return
	}
	if result.Response.ErrorList != nil && result.Response.ErrorList.Error != nil {
		log.Warn("Response Error %v", result.Response.ErrorList.Error)
	}
	return
}

// CallRetry retries the xml call
func (d *Dao) CallRetry(c context.Context, url string, xmlBody string) (result *model.Document, err error) {
	log.Info("callLic URL: %s, Body %s", url, xmlBody)
	for i := 0; i < 3; i++ {
		result, err = d.callLic(c, url, xmlBody)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Second) // 5 seconds gap for each retrial
	}
	return
}

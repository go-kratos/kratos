package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-common/app/service/main/tv/internal/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	ystUrl = "https://%s/%s"

	ystCreateOrderPath = "getBiliOrder"
	ystRenewOrderPath  = "papPayApply"
	ystOrderStatePath  = "getBiliOrderState"
)

// YstClient represents http client for sending requests to yst.
type YstClient struct {
	cli *bm.Client
}

// NewYstClient news a yst client.
func NewYstClient(cli *bm.Client) *YstClient {
	return &YstClient{cli: cli}
}

// Post sends one post request.
func (yc *YstClient) Post(c context.Context, uri string, data interface{}, res interface{}) error {
	req, err := yc.NewRequest(uri, data)
	if err != nil {
		return err
	}
	return yc.cli.Do(c, req, res)
}

// NewRequest news a post request.
func (yc *YstClient) NewRequest(uri string, data interface{}) (request *http.Request, err error) {
	var dataBytes []byte
	dataBytes, err = json.Marshal(data)
	if err != nil {
		return
	}
	request, err = http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func resultCode2err(code string) error {
	switch code {
	case model.YstResultFail:
		return ecode.TVIPYstRequestErr
	case model.YstResultSysErr:
		return ecode.TVIPYstSystemErr
	default:
		return ecode.TVIPYstUnknownErr
	}
}

// CreateYstOrder creates one yst order.
func (d *Dao) CreateYstOrder(c context.Context, req *model.YstCreateOrderReq) (res *model.YstCreateOrderReply, err error) {
	req.Sign, err = d.signer.Sign(req)
	if err != nil {
		return
	}
	res = &model.YstCreateOrderReply{}
	uri := fmt.Sprintf(ystUrl, d.c.YST.Domain, ystCreateOrderPath)
	if err = d.ystCli.Post(c, uri, req, res); err != nil {
		log.Error("d.ystCli.Post(%s, %+v) err(%+v)", uri, req, err)
		return
	}
	if res.ResultCode != model.YstResultSuccess {
		log.Error("d.ystCli.Post(%s, %+v) res(%+v)", uri, req, res)
		return nil, resultCode2err(res.ResultCode)
	}
	log.Info("dao.CreateYstOrder(%+v) res(%+v)", req, res)
	return
}

// RenewYstOrder creates one renew order.
func (d *Dao) RenewYstOrder(c context.Context, req *model.YstRenewOrderReq) (res *model.YstRenewOrderReply, err error) {
	req.Sign, err = d.signer.Sign(req)
	if err != nil {
		return
	}
	res = &model.YstRenewOrderReply{}
	uri := fmt.Sprintf(ystUrl, d.c.YST.Domain, ystRenewOrderPath)
	if err = d.ystCli.Post(c, uri, req, res); err != nil {
		log.Error("d.ystCli.Post(%s, %+v) err(%+v)", uri, req, err)
		return
	}
	if res.ResultCode != model.YstResultSuccess {
		log.Error("d.ystCli.Post(%s, %+v) res(%+v)", uri, req, res)
		return nil, resultCode2err(res.ResultCode)
	}
	log.Info("dao.RenewYstOrder(%+v) res(%+v)", req, res)
	return
}

// YstOrderState queries order details from yst.
func (d *Dao) YstOrderState(c context.Context, req *model.YstOrderStateReq) (res *model.YstOrderStateReply, err error) {
	req.Sign, err = d.signer.Sign(req)
	if err != nil {
		return
	}
	res = &model.YstOrderStateReply{}
	uri := fmt.Sprintf(ystUrl, d.c.YST.Domain, ystOrderStatePath)
	if err = d.ystCli.Post(c, uri, req, res); err != nil {
		log.Error("d.ystCli.Post(%s, %+v) err(%+v)", uri, req, err)
		return
	}
	if res.Result != model.YstResultSuccess {
		log.Error("d.ystCli.Post(%s, %+v) res(%+v)", uri, req, res)
		return nil, resultCode2err(res.Result)
	}
	log.Info("dao.YstOrderState(%+v) res(%+v)", req, res)
	return
}

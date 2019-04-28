package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"net/http"
	"net/url"
	"time"
)

const (
	ModifyUrl            = "http://api.bilibili.co/x/internal/v1/coin/user/modify"
	InfoUrl              = "http://api.bilibili.co/x/internal/v1/coin/user/count"
	ExchangeSilverReason = "兑换直播银瓜子 %d"
	ExchangeMetalReason  = "银瓜子兑换硬币"
)

var (
	respCodeError = errors.New("query response code err")
	paramError    = errors.New("param error")
)

func (d *Dao) GetMetal(c context.Context, uid int64) (metal float64, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", uid))
	now := time.Now().Unix()
	params.Set("ts", fmt.Sprintf("%d", now))
	var req *http.Request
	if req, err = d.httpClient.NewRequest("GET", InfoUrl, "", params); err != nil {
		log.Error("wallet metal newRequest err:%s", err.Error())
		return
	}
	queryUrl := req.URL.String()
	var res struct {
		Code int             `json:"code"`
		Data model.MetalData `json:"data"`
	}
	err = d.httpClient.Do(c, req, &res)
	if err != nil {
		log.Error("wallet query metal err  url: %s,err:%s", queryUrl, err.Error())
		return
	}
	resStr, _ := json.Marshal(res)
	log.Info("wallet query metal success url:%s,res:%s", queryUrl, resStr)

	if res.Code != 0 {
		err = respCodeError
	} else {
		metal = res.Data.Count
	}

	return
}

func (d *Dao) ModifyMetal(c context.Context, uid int64, coins int64, seeds int64, reason interface{}) (success bool, code int, err error) {
	if coins == 0 {
		err = paramError
		return
	}

	if coins < 0 {
		// 检查是否足够
		metal, _ := d.GetMetal(c, uid)
		if metal < float64(0-coins) {
			err = ecode.CoinNotEnough
			return
		}
	}

	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", uid))
	params.Set("count", fmt.Sprintf("%d", coins))
	now := time.Now().Unix()
	params.Set("ts", fmt.Sprintf("%d", now))

	var realReason string
	switch reason.(type) {
	case string:
		realReason = reason.(string)
	default:
		if coins < 0 {
			realReason = fmt.Sprintf(ExchangeSilverReason, seeds)
		} else {
			realReason = ExchangeMetalReason
		}
	}
	log.Info("user %d consume or recharge metal %d by reason %s", uid, coins, realReason)
	params.Set("reason", realReason)

	var req *http.Request
	req, err = d.httpClient.NewRequest("POST", ModifyUrl, "", params)
	if err != nil {
		log.Error("wallet metal newRequest err:%s", err.Error())
		return
	}

	queryUrl := req.URL.String()
	var res struct {
		Code int `json:"code"`
	}
	err = d.httpClient.Do(c, req, &res)
	if err != nil {
		log.Error("Metal#wallet query metal err  url: %s,err:%s uid:%d, count:%d", queryUrl, err.Error(), uid, coins)
		// 认为成功
		err = nil
		success = true
		code = 0
		return
	}
	resStr, _ := json.Marshal(res)
	log.Info("wallet query metal success url:%s, uid:%d, count:%d res:%s", queryUrl, uid, coins, resStr)

	code = res.Code
	if res.Code == 0 {
		success = true
	} else {
		if res.Code == ecode.LackOfCoins.Code() {
			err = ecode.CoinNotEnough
		} else {
			log.Error("Metal#wallet query metal code failed : uid:%d, count:%d,code:%d", uid, coins, res.Code)
			err = respCodeError
		}
	}

	return

}

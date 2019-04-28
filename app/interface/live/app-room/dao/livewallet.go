package dao

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/pkg/errors"
	"go-common/app/interface/live/app-room/model"
	"go-common/library/ecode"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//LiveWallet 获取用户直播钱包数据
func (d *Dao) LiveWallet(c context.Context, mid int64, platform string) (wallet *model.LiveWallet, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("appkey", d.c.HTTPClient.LiveRpc.Key)
	tmp := params.Encode() + d.c.HTTPClient.LiveRpc.Secret
	mh := md5.Sum([]byte(tmp))
	sign := hex.EncodeToString(mh[:])
	params.Set("sign", sign)
	req, _ := http.NewRequest("GET", d.liveWalletURL+"?"+params.Encode(), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("platform", platform)
	var wr struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Gold            string `json:"gold"`
			Silver          string `json:"silver"`
			GoldRechargeCnt string `json:"gold_recharge_cnt"`
			GoldPayCnt      string `json:"gold_pay_cnt"`
			SilverPayCnt    string `json:"silver_pay_cnt"`
		} `json:"data"`
	}
	if err = d.liveWalletClient.Do(c, req, &wr); err != nil {
		return
	}
	if wr.Code != 0 {
		err = errors.Wrap(ecode.Int(wr.Code), d.liveWalletURL+"?"+params.Encode())
		return
	}
	wallet = &model.LiveWallet{}
	wallet.Gold, _ = strconv.ParseInt(wr.Data.Gold, 10, 64)
	wallet.Silver, _ = strconv.ParseInt(wr.Data.Silver, 10, 64)
	wallet.GoldRechargeCnt, _ = strconv.ParseInt(wr.Data.GoldRechargeCnt, 10, 64)
	wallet.GoldPayCnt, _ = strconv.ParseInt(wr.Data.GoldPayCnt, 10, 64)
	wallet.SilverPayCnt, _ = strconv.ParseInt(wr.Data.SilverPayCnt, 10, 64)

	return
}

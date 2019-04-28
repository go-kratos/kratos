package newcomer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fmt"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
	"math/rand"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Mall receive mall coupon code.
func (d *Dao) Mall(c context.Context, mid int64, couponID, uname string) (err error) {
	type params struct {
		MID      int64  `json:"mid"`
		CouponID string `json:"couponId"`
		Uname    string `json:"uname"`
	}
	p := params{}
	p.MID = mid
	p.CouponID = couponID //优惠券id
	p.Uname = uname       //抱团购类型 uname必传
	paramJSON, err := json.Marshal(p)
	if err != nil {
		log.Error("Mall json.Marshal param(%+v) error(%v)", p, err)
		return
	}
	paramStr := string(paramJSON)

	var (
		req *http.Request
	)
	if req, err = http.NewRequest("POST", d.mallURI, strings.NewReader(paramStr)); err != nil {
		log.Error("Mall http.NewRequest url(%s) error(%v)", d.mallURI+"?"+paramStr, err)
		err = ecode.CreativeNewcomerMallAPIErr
		return
	}
	log.Info("Mall url(%s)", d.mallURI+"?"+paramStr)

	req.Header.Set("Content-Type", "application/json")
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("Mall d.client.Do url(%s) res(%v) err(%v)", d.mallURI+"?"+paramStr, res, err)
		err = ecode.CreativeNewcomerMallAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("Mall url(%s) res(%v)", d.mallURI+"?"+paramStr, res)
		err = ecode.Int(res.Code)
	}
	return
}

// BCoin receive b coin.
func (d *Dao) BCoin(c context.Context, mid int64, activityID string, money int64) (err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("activity_id", activityID)             //活动券id  uat-217 pre-266
	params.Set("money", strconv.FormatInt(money, 10)) //decimal 类型 领取的数量，最大保留两位小数
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))

	var (
		query, _ = tool.Sign(params)
		mallURL  = d.bPayURI
	)
	log.Info("BCoin url(%s)", mallURL+"?"+query)

	if req, err = http.NewRequest("POST", mallURL, strings.NewReader(params.Encode())); err != nil {
		log.Error("BCoin url(%s) error(%v)", mallURL, err)
		err = ecode.CreativeNewcomerBCoinAPIErr
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("BCoin d.client.Do url(%s) res(%+v) err(%v)", mallURL, res, err)
		err = ecode.CreativeNewcomerBCoinAPIErr
		return
	}

	if res.Code != 0 {
		log.Error("BCoin url(%s) res(%+v)", d.bPayURI+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
	}
	return
}

// Pendant receive pendant.
func (d *Dao) Pendant(c context.Context, mid int64, priceID string, expires int64) (err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)
	pid, err := strconv.ParseInt(priceID, 10, 64)
	if err != nil {
		return
	}
	params.Set("mids", strconv.FormatInt(mid, 10))
	params.Set("pid", strconv.FormatInt(pid, 10))        //挂件ID
	params.Set("expire", strconv.FormatInt(expires, 10)) //有效期(单位:天)
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var (
		query, _ = tool.Sign(params)
		pendURL  = d.pendantURI
	)
	log.Info("Pendant url(%s)", pendURL+"?"+query)
	if req, err = http.NewRequest("POST", pendURL, strings.NewReader(params.Encode())); err != nil {
		log.Error("Pendant url(%s) error(%v)", pendURL, err)
		err = ecode.CreativeNewcomerPendantAPIErr
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("Pendant d.client.Do url(%s) res(%+v) err(%v)", pendURL, res, err)
		err = ecode.CreativeNewcomerPendantAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("Pendant url(%s) res(%v)", pendURL, res)
		err = ecode.Int(res.Code)
	}
	return
}

// BigMemberCoupon receive Coupon
func (d *Dao) BigMemberCoupon(c context.Context, mid int64, batchToken string) (err error) {
	var (
		req    *http.Request
		params = url.Values{}
	)

	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("batch_token", batchToken)   //资源批次需提前申请，金额固定，不允许更改
	params.Set("order_no", genOrderNO(mid)) //每次领取的订单号不允许重复，同一订单同一业务方AppKey只允许领取一次
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var (
		query, _     = tool.Sign(params)
		bigMemberURL = d.bigMemberURI
	)
	log.Info("BigMemberCoupon url(%s)", bigMemberURL+"?"+query)

	if req, err = http.NewRequest("POST", bigMemberURL, strings.NewReader(params.Encode())); err != nil {
		log.Error("BigMemberCoupon url(%s) error(%v)", bigMemberURL, err)
		err = ecode.CreativeNewcomerBigMemberErr
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		Data string `json:"data"`
	}

	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("BigMemberCoupon d.client.Do url(%s) res(%+v) err(%v)", bigMemberURL, res, err)
		err = ecode.CreativeNewcomerBigMemberErr
		return
	}
	if res.Code != 0 {
		log.Error("BigMemberCoupon url(%s) res(%v)", bigMemberURL, res)
		err = ecode.Int(res.Code)
	}
	return
}

// generate orderNo
func genOrderNO(mid int64) string {
	s := fmt.Sprintf("%v%s", mid, time.Now().Format("20060102150405"))
	if len(s) >= 32 {
		return s[0:32]
	}
	size := 32 - len(s)
	bys := make([]byte, size)
	for i := 0; i < size; i++ {
		bys[i] = uint8(48 + rand.Intn(10))
	}
	return s + string(bys)
}

// SendNotify send msg notify user
func (d *Dao) SendNotify(c context.Context, mids []int64, mc, title, context string) (err error) {
	var (
		params = url.Values{}
		res    struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data *struct {
				TotalCount   int     `json:"total_count"`
				ErrorCount   int     `json:"error_count"`
				ErrorMidList []int64 `json:"error_mid_list"`
			} `json:"data"`
		}
	)
	params.Set("mc", mc)                        //消息码，用于识别消息类别
	params.Set("data_type", "4")                //消息类型：1、回复我的 2、@我 3、收到的爱 4、业务通知 5、系统公告
	params.Set("title", title)                  //消息标题
	params.Set("context", context)              //消息实体内容
	params.Set("mid_list", xstr.JoinInts(mids)) //用于接收该消息的用户mid列表，不超过1000个(半角逗号分割)

	log.Info("SendNotify params(%+v)|msgURI(%s)", params.Encode(), d.msgNotifyURI)
	if err = d.client.Post(c, d.msgNotifyURI, "", params, &res); err != nil {
		log.Error("d.httpClient.Post(%s,%v,%d)", d.msgNotifyURI, params, err)
		return
	}
	if res.Code != 0 {
		err = errors.New("code != 0")
		log.Error("d.httpClient.Post(%s,%v,%v,%d)", d.msgNotifyURI, params, err, res.Code)
	}
	if res.Data != nil {
		log.Info("SendNotify log total_count(%d) error_count(%d) error_mid_list(%v)", res.Data.TotalCount, res.Data.ErrorCount, res.Data.ErrorMidList)
	}
	return
}

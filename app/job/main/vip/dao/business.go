package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	xtime "time"

	"go-common/app/job/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	_cleanCache       = "/notify/cleanCache"
	_payOrder         = "/payplatform/pay/pay"
	_message          = "/api/notify/send.user.notify.do"
	_addBcoin         = "/api/coupon/regular/add"
	_sendVipbuyTicket = "/mall-marketing/coupon_code/create"
	_sendMedal        = "/api/nameplate/get/v2"
	_pushData         = "/x/internal/push-strategy/task/add"

	_retry   = 3
	_minRead = 1024 * 64

	_alreadySend = -804

	_alreadyGet = -663

	_alreadySendVipbuy = 83110005
	_ok                = 1

	//push
	appID = 1
)

//PushData http push data
func (d *Dao) PushData(c context.Context, mids []int64, pushData *model.VipPushData, curtime string) (rel *model.VipPushResq, err error) {
	var (
		pushTime   xtime.Time
		expireTime xtime.Time
		params     = url.Values{}
	)
	rel = new(model.VipPushResq)

	if pushTime, err = xtime.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v %v", curtime, pushData.PushStartTime), xtime.Local); err != nil {
		err = errors.WithStack(err)
		return
	}

	if expireTime, err = xtime.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v %v", curtime, pushData.PushEndTime), xtime.Local); err != nil {
		err = errors.WithStack(err)
		return
	}

	page := len(mids) / d.c.Property.SplitPush
	if len(mids)%d.c.Property.SplitPush != 0 {
		page++
	}
	for i := 0; i < page; i++ {
		startID := i * d.c.Property.SplitPush
		endID := (i + 1) * d.c.Property.SplitPush
		if endID > len(mids) {
			endID = len(mids)
		}
		tempMids := mids[startID:endID]
		params.Set("app_id", fmt.Sprintf("%v", appID))
		params.Set("business_id", fmt.Sprintf("%v", d.c.Property.BusinessID))
		params.Set("alert_title", pushData.Title)
		params.Set("alert_body", pushData.Content)
		params.Set("mids", xstr.JoinInts(tempMids))
		params.Set("link_type", fmt.Sprintf("%v", pushData.LinkType))
		params.Set("link_value", pushData.LinkURL)
		params.Set("builds", pushData.Platform)
		params.Set("group", pushData.GroupName)
		params.Set("uuid", uuid.New().String())
		params.Set("push_time", fmt.Sprintf("%v", pushTime.Unix()))
		params.Set("expire_time", fmt.Sprintf("%v", expireTime.Unix()))

		header := make(map[string]string)
		header["Authorization"] = fmt.Sprintf("token=%v", d.c.Property.PushToken)
		header["Content-Type"] = "application/x-www-form-urlencoded"

		for i := 0; i < _retry; i++ {
			if err = d.doNomalSend(c, d.c.URLConf.APICoURL, _pushData, "127.0.0.1", header, params, rel); err != nil {
				log.Error("send error(%v) url(%v)", err, d.c.URLConf.APICoURL+_pushData)
				return
			}

			if rel.Code == int64(ecode.OK.Code()) {
				log.Info("send url:%v params:%+v return:%+v error(%+v)", d.c.URLConf.APICoURL+_pushData, params, rel, err)
				break
			}

		}
	}

	return
}

//SendMedal send medal
func (d *Dao) SendMedal(c context.Context, mid, medalID int64) (status int64) {
	var (
		err error
	)
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%v", mid))
	params.Set("nid", fmt.Sprintf("%v", medalID))
	rel := new(struct {
		Code int64  `json:"code"`
		Data string `json:"data"`
	})
	defer func() {
		if err == nil {
			log.Info("send url:%+v params:%+v return:%+v", d.c.URLConf.AccountURL+_sendMedal, params, rel)
		}
	}()
	for i := 0; i < _retry; i++ {
		if err = d.client.Get(c, d.c.URLConf.AccountURL+_sendMedal, "127.0.0.1", params, rel); err != nil {
			log.Error("send error(%v) url(%v)", err, d.c.URLConf.AccountURL+_sendMedal)
			continue
		}
		if rel.Code == int64(ecode.OK.Code()) || rel.Code == _alreadyGet {
			status = 1
			if rel.Code == _alreadyGet {
				status = 2
			}
			return
		}
	}
	return
}

//SendVipBuyTicket send vipbuy ticket
func (d *Dao) SendVipBuyTicket(c context.Context, mid int64, couponID string) (status int64) {
	var (
		err error
	)
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	params := make(map[string]string)
	params["mid"] = fmt.Sprintf("%v", mid)
	params["couponId"] = fmt.Sprintf("%v", couponID)

	for i := 0; i < _retry; i++ {
		repl := new(struct {
			Code    int64  `json:"code"`
			Message string `json:"message"`
		})

		if err = d.doSend(c, d.c.URLConf.MallURL, _sendVipbuyTicket, "127.0.0.1", header, params, repl); err != nil {
			log.Error("send vip buy ticket(%+v) error(%+v)", params, err)
			continue
		}
		if repl.Code == int64(ecode.OK.Code()) || repl.Code == _alreadySendVipbuy {
			status = 1
			if repl.Code == _alreadySendVipbuy {
				status = 2
			}
			return
		}
	}
	return
}

//SendCleanCache clean cache
func (d *Dao) SendCleanCache(c context.Context, hv *model.HandlerVip) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(int64(hv.Mid), 10))
	if err = d.client.Get(c, d.c.VipURI+_cleanCache, "127.0.0.1", params, nil); err != nil {
		log.Error("SendCleanCache error(%v) url(%v)", err, d.c.VipURI+_cleanCache)
		return
	}
	return
}

//SendBcoin send bcoin http
func (d *Dao) SendBcoin(c context.Context, mids []int64, money int32, dueTime time.Time, ip string) (err error) {
	if len(mids) <= 0 {
		return
	}
	var midStrs []string
	for _, v := range mids {
		midStrs = append(midStrs, fmt.Sprintf("%v", v))
	}
	params := url.Values{}
	params.Add("activity_id", fmt.Sprintf("%v", d.c.Property.ActivityID))
	params.Add("mids", strings.Join(midStrs, ","))
	params.Add("money", fmt.Sprintf("%v", money*100))
	params.Add("due_time", dueTime.Time().Format("2006-01-02"))

	res := new(struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		TS      int64  `json:"ts"`
		Data    struct {
			CouponMoney int64 `json:"coupon_money"`
			Status      int8  `json:"status"`
		} `json:"data"`
	})
	if err = d.client.Post(c, d.c.URLConf.PayCoURL+_addBcoin, ip, params, res); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res.Code == int64(_alreadySend) {
		return
	}
	if int(res.Code) != ecode.OK.Code() {
		err = fmt.Errorf("发放B币失败 message: %s mid: %s resp(%+v)", res.Message, strings.Join(midStrs, ","), res)
		return
	}
	if res.Data.Status != _ok {
		err = fmt.Errorf("发放B币失败 message: %s mid: %s resp(%+v)", res.Message, strings.Join(midStrs, ","), res)
		return
	}
	log.Info("发放B币成功 mids %v resp(%+v)", mids, res)
	return
}

//SendAppCleanCache notice app clean cache
func (d *Dao) SendAppCleanCache(c context.Context, hv *model.HandlerVip, app *model.VipAppInfo) (err error) {
	params := url.Values{}
	params.Set("modifiedAttr", "updateVip")
	params.Set("mid", fmt.Sprintf("%v", hv.Mid))
	params.Set("status", fmt.Sprintf("%v", hv.Type))
	params.Set("buyMonths", fmt.Sprintf("%v", hv.Months))
	params.Set("days", fmt.Sprintf("%v", hv.Days))
	if err = d.client.Get(c, app.PurgeURL, "127.0.0.1", params, nil); err != nil {
		log.Error("SendAppCleanCache error(%v) url(%v) params(%v)", err, app.PurgeURL, params)
		return
	}
	return
}

//SendMultipMsg send multip msg
func (d *Dao) SendMultipMsg(c context.Context, mids, content, title, mc string, dataType int) (err error) {
	params := url.Values{}
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("context", content)
	params.Set("data_type", strconv.FormatInt(int64(dataType), 10))
	params.Set("mid_list", mids)
	defer func() {
		log.Info("SendMultipMsg(%v) params(%+v) error(%+v)", d.c.URLConf.MsgURL+_message, params, err)
	}()
	if err = d.client.Post(c, d.c.URLConf.MsgURL+_message, "127.0.0.1", params, nil); err != nil {
		log.Error("SendMultipMsg params(%+v) error(%v)", err, params)
		return
	}
	log.Info("cur send mid(%+v)", mids)
	return
}

// PayOrder pay order.
func (d *Dao) PayOrder(c context.Context, paramsMap map[string]interface{}) (err error) {
	params := make(map[string]string)
	for k, v := range paramsMap {
		params[k] = fmt.Sprintf("%v", v)
	}

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	success := false
	for i := 0; i < 3; i++ {
		repl := new(struct {
			ErrNo int64       `json:"errno"`
			Msg   string      `json:"msg"`
			Data  interface{} `json:"data"`
		})
		if err = d.doPaySend(c, d.c.URLConf.PayURL, _payOrder, "127.0.0.1", header, params, repl); err != nil {
			continue
		}
		if repl.ErrNo == 0 {
			success = true
			break
		}
	}
	if !success {
		err = fmt.Errorf("下单失败")
	}
	return
}

func (d *Dao) sortParamsKey(v map[string]string) string {
	if v == nil {
		return ""
	}

	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		buf.WriteString(vs)
	}
	return buf.String()
}

//PaySign pay sign
func (d *Dao) PaySign(params map[string]string, token string) (sign string) {

	tmp := d.sortParamsKey(params)

	var b bytes.Buffer
	b.WriteString(tmp)
	b.WriteString(fmt.Sprintf("&token=%s", token))
	log.Info("sign params:%v ", b.String())
	mh := md5.Sum(b.Bytes())
	// query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(hex.EncodeToString(mh[:]))
	sign = hex.EncodeToString(mh[:])
	log.Info("sign params(%v) and sign(%v)", b.String(), sign)
	return
}

func (d *Dao) doNomalSend(c context.Context, basePath, path, ip string, header map[string]string, params url.Values, data interface{}) (err error) {
	var (
		req     *http.Request
		client  = new(http.Client)
		resp    *http.Response
		bs      []byte
		marshal string
	)
	url := basePath + path

	if req, err = d.client.NewRequest(http.MethodPost, url, ip, params); err != nil {
		err = errors.WithStack(err)
		return
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}
	if resp, err = client.Do(req); err != nil {
		log.Error("call url:%v params:%v", basePath+path, string(marshal))
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()

	defer func() {
		log.Info("call url:%v params:(%v) result:(%+v) header:(%+v)", url, string(marshal), data, header)
	}()
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, req.URL.String())
		return
	}
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = errors.Wrapf(err, "host:%s, url:%s", req.URL.Host, req.URL.String())
		return
	}
	if err = json.Unmarshal(bs, data); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (d *Dao) doSend(c context.Context, basePath, path, IP string, header map[string]string, params map[string]string, data interface{}) (err error) {
	var (
		req    *http.Request
		client = new(http.Client)
		resp   *http.Response
		bs     []byte
	)
	url := basePath + path
	marshal, _ := json.Marshal(params)
	if req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(string(marshal))); err != nil {
		err = errors.WithStack(err)
		return
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	if resp, err = client.Do(req); err != nil {
		log.Error("call url:%v params:%v", basePath+path, string(marshal))
		err = errors.WithStack(err)
		return
	}
	defer resp.Body.Close()

	defer func() {
		log.Info("call url:%v params:(%v) result:(%+v) header:(%+v)", url, string(marshal), data, header)
	}()
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, req.URL.String())
		return
	}
	if bs, err = readAll(resp.Body, _minRead); err != nil {
		err = errors.Wrapf(err, "host:%s, url:%s", req.URL.Host, req.URL.String())
		return
	}
	if err = json.Unmarshal(bs, data); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (d *Dao) doPaySend(c context.Context, basePath, path, IP string, header map[string]string, params map[string]string, data interface{}) (err error) {
	sign := d.PaySign(params, d.c.PayConf.Token)
	params["sign"] = sign
	return d.doSend(c, basePath, path, IP, header, params, data)
}

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}

// SalaryCoupon salary coupon.
func (d *Dao) SalaryCoupon(c context.Context, mid int64, couponType int8, count int64, token string) (err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	params.Set("type", fmt.Sprintf("%d", couponType))
	params.Set("count", fmt.Sprintf("%d", count))
	params.Set("batch_no", token)
	var resp struct {
		Code int64 `json:"code"`
	}
	if err = d.client.Post(c, d.c.Property.SalaryCouponURL, "", params, &resp); err != nil {
		log.Error("message url(%s) error(%v)", d.c.Property.SalaryCouponURL+"?"+params.Encode(), err)
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("POST SalaryCoupon url resp(%v)", resp)
		return
	}
	return
}

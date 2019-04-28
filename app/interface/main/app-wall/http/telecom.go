package http

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go-common/app/interface/main/app-wall/model/telecom"
	"go-common/library/ecode"
	log "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// ordersSync
func telecomOrdersSync(c *bm.Context) {
	data, err := requestJSONTelecom(c.Request)
	if err != nil {
		telecomMessage(c, err.Error())
		return
	}
	switch v := data.(type) {
	case *telecom.TelecomOrderJson:
		if err = telecomSvc.InOrdersSync(c, metadata.String(c, metadata.RemoteIP), v); err != nil {
			log.Error("telecomSvc.InOrdersSync error (%v)", err)
			telecomMessage(c, err.Error())
			return
		}
		telecomMessage(c, "1")
	case *telecom.TelecomRechargeJson:
		if v == nil {
			telecomMessage(c, ecode.NothingFound.Error())
			return
		}
		if err = telecomSvc.InRechargeSync(c, metadata.String(c, metadata.RemoteIP), v.Detail); err != nil {
			log.Error("telecomSvc.InOrdersSync error (%v)", err)
			telecomMessage(c, err.Error())
			return
		}
		telecomMessage(c, "1")
	}
}

// telecomMsgSync
func telecomMsgSync(c *bm.Context) {
	data, err := requestJSONTelecomMsg(c.Request)
	if err != nil {
		telecomMessage(c, err.Error())
		return
	}
	if err = telecomSvc.TelecomMessageSync(c, metadata.String(c, metadata.RemoteIP), data); err != nil {
		log.Error("telecomSvc.TelecomMessageSync error (%v)", err)
		telecomMessage(c, err.Error())
		return
	}
	telecomMessage(c, "1")
}

func telecomPay(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phoneDES := params.Get("phone")
	phone, err := phoneDesToInt(phoneDES)
	if err != nil {
		log.Error("phoneDesToInt error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	isRepeatOrderStr := params.Get("is_repeat_order")
	isRepeatOrder, err := strconv.Atoi(isRepeatOrderStr)
	if err != nil {
		log.Error("isRepeatOrder strconv.Atoi error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	payChannelStr := params.Get("pay_channel")
	payChannel, err := strconv.Atoi(payChannelStr)
	if err != nil {
		log.Error("payChannel strconv.Atoi error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	payActionStr := params.Get("pay_action")
	payAction, err := strconv.Atoi(payActionStr)
	if err != nil {
		log.Error("payAction strconv.Atoi error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	orderIDStr := params.Get("orderid")
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 64)
	data, msg, err := telecomSvc.TelecomPay(c, phone, isRepeatOrder, payChannel, payAction, orderID, metadata.String(c, metadata.RemoteIP))
	if err != nil {
		log.Error("telecomSvc.TelecomPay error(%v)", err)
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	returnDataJSON(c, res, nil)
}

func cancelRepeatOrder(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phoneDES := params.Get("phone")
	phone, err := phoneDesToInt(phoneDES)
	if err != nil {
		log.Error("phoneDesToInt error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	msg, err := telecomSvc.CancelRepeatOrder(c, phone)
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["message"] = msg
	returnDataJSON(c, res, nil)
}

func orderList(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phoneDES := params.Get("phone")
	phone, err := phoneDesToInt(phoneDES)
	if err != nil {
		log.Error("phoneDesToInt error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	orderIDStr := params.Get("orderid")
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 64)
	data, msg, err := telecomSvc.OrderList(c, orderID, phone)
	if err != nil {
		log.Error("telecomSvc.OrderList error(%v)", err)
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	returnDataJSON(c, res, nil)
}

func phoneFlow(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phoneDES := params.Get("phone")
	phone, err := phoneDesToInt(phoneDES)
	if err != nil {
		log.Error("phoneDesToInt error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	orderIDStr := params.Get("orderid")
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 64)
	data, msg, err := telecomSvc.PhoneFlow(c, orderID, phone)
	if err != nil {
		log.Error("telecomSvc.PhoneFlow error(%v)", err)
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	returnDataJSON(c, res, nil)
}

func orderConsent(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phoneDES := params.Get("phone")
	phone, err := phoneDesToInt(phoneDES)
	if err != nil {
		log.Error("phoneDesToInt error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	captcha := params.Get("captcha")
	orderIDStr := params.Get("orderid")
	orderID, _ := strconv.ParseInt(orderIDStr, 10, 64)
	data, msg, err := telecomSvc.OrderConsent(c, phone, orderID, captcha)
	if err != nil {
		log.Error("telecomSvc.OrderConsent error(%v)", err)
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	returnDataJSON(c, res, nil)
}

func phoneSendSMS(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phoneDES := params.Get("phone")
	phone, err := phoneDesToInt(phoneDES)
	if err != nil {
		log.Error("phoneDesToInt error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	err = telecomSvc.PhoneSendSMS(c, phone)
	if err != nil {
		log.Error("telecomSvc.PhoneSendSMS error(%v)", err)
		res["code"] = err
		res["message"] = ""
		returnDataJSON(c, res, err)
		return
	}
	res["message"] = ""
	returnDataJSON(c, res, nil)
}

func phoneVerification(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phoneDES := params.Get("phone")
	phone, err := phoneDesToInt(phoneDES)
	if err != nil {
		log.Error("phoneDesToInt error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	captcha := params.Get("captcha")
	data, err, msg := telecomSvc.PhoneCode(c, phone, captcha, time.Now())
	if err != nil {
		log.Error("telecomSvc.PhoneCode error(%v)", err)
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	res["message"] = msg
	returnDataJSON(c, res, nil)
}

func orderState(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	orderIDStr := params.Get("orderid")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		log.Error("orderID strconv.ParseInt error(%v)", err)
		res["message"] = ""
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	data, msg, err := telecomSvc.OrderState(c, orderID)
	if err != nil {
		log.Error("telecomSvc.OrderState error(%v)", err)
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	res["message"] = msg
	returnDataJSON(c, res, nil)
}

// requestJSONTelecom
func requestJSONTelecom(request *http.Request) (res interface{}, err error) {
	var (
		telecomOrder    *telecom.TelecomOrderJson
		telecomRecharge *telecom.TelecomRechargeJson
	)
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Error("unicom_ioutil.ReadAll error (%v)", err)
		return
	}
	defer request.Body.Close()
	if len(body) == 0 {
		err = ecode.RequestErr
		return
	}
	log.Info("telecom orders json body(%s)", body)
	if err = json.Unmarshal(body, &telecomOrder); err == nil && telecomOrder != nil {
		if telecomOrder.ResultType != 2 {
			res = telecomOrder
			return
		}
	}
	if err = json.Unmarshal(body, &telecomRecharge); err == nil {
		res = telecomRecharge
		return
	}
	log.Error("telecom json.Unmarshal error (%v)", err)
	return
}

// requestJSONTelecomMsg
func requestJSONTelecomMsg(request *http.Request) (res *telecom.TelecomMessageJSON, err error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Error("unicom_ioutil.ReadAll error (%v)", err)
		return
	}
	defer request.Body.Close()
	if len(body) == 0 {
		err = ecode.RequestErr
		return
	}
	log.Info("telecom json msg body(%s)", body)
	if err = json.Unmarshal(body, &res); err != nil {
		log.Error("telecom Message json.Unmarshal error (%v)", err)
		return
	}
	return
}

// telecomMessage
func telecomMessage(c *bm.Context, code string) {
	// response header
	c.Writer.Header().Set("Content-Type", "text; charset=UTF-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Write([]byte(code))
}

// phoneDesToInt des to int
func phoneDesToInt(phoneDes string) (phoneInt int, err error) {
	var (
		_aesKey = []byte("6b7e8b8a")
	)
	bs, err := base64.StdEncoding.DecodeString(phoneDes)
	if err != nil {
		log.Error("base64.StdEncoding.DecodeString(%s) error(%v)", phoneDes, err)
		err = ecode.RequestErr
		return
	}
	if bs, err = telecomSvc.DesDecrypt(bs, _aesKey); err != nil {
		log.Error("phone s.DesDecrypt error(%v)", err)
		return
	}
	var phoneStr string
	if len(bs) > 11 {
		phoneStr = string(bs[:11])
	} else {
		phoneStr = string(bs)
	}
	phoneInt, err = strconv.Atoi(phoneStr)
	if err != nil {
		log.Error("phoneDesToInt phoneStr:%v error(%v)", phoneStr, err)
		err = ecode.RequestErr
		return
	}
	return
}

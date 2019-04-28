package http

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-wall/model"
	"go-common/app/interface/main/app-wall/model/unicom"
	"go-common/library/ecode"
	log "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// ordersSync
func ordersSync(c *bm.Context) {
	res := map[string]interface{}{}
	unicom, err := requestJSONToMap(c.Request)
	if err != nil {
		res["result"] = "1"
		res["errorcode"] = "100"
		log.Error("unicom_ioutil.ReadAll error (%v)", err)
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	usermob, err := url.QueryUnescape(unicom.Usermob)
	if err != nil {
		log.Error("unicom_url.QueryUnescape (%v) error (%v)", unicom.Usermob, err)
		res["result"] = "1"
		res["errorcode"] = "100"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	var (
		_aesKey = []byte("9ed226d9")
	)
	bs, err := base64.StdEncoding.DecodeString(usermob)
	if err != nil {
		log.Error("base64.StdEncoding.DecodeString(%s) error(%v)", usermob, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err = unicomSvc.DesDecrypt(bs, _aesKey)
	if err != nil {
		log.Error("unicomSvc.DesDecrypt error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var usermobStr string
	if len(bs) > 32 {
		usermobStr = string(bs[:32])
	} else {
		usermobStr = string(bs)
	}
	log.Info("unicomSvc.OrdersSync_usermob (%v) unicom (%v)", usermobStr, unicom)
	if err := unicomSvc.InOrdersSync(c, usermobStr, metadata.String(c, metadata.RemoteIP), unicom, time.Now()); err != nil {
		log.Error("unicomSvc.OrdersSync usermob (%v) unicom (%v) error (%v)", usermobStr, unicom, err)
		res["result"] = "1"
		res["errorcode"] = "100"
		returnDataJSON(c, res, err)
		return
	}
	res["result"] = "0"
	res["message"] = ""
	res["errorcode"] = ""
	returnDataJSON(c, res, nil)
}

// advanceSync
func advanceSync(c *bm.Context) {
	res := map[string]interface{}{}
	unicom, err := requestJSONToMap(c.Request)
	if err != nil {
		res["result"] = "1"
		res["errorcode"] = "100"
		log.Error("unicom_ioutil.ReadAll error (%v)", err)
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	usermob, err := url.QueryUnescape(unicom.Usermob)
	if err != nil {
		log.Error("unicom_url.QueryUnescape (%v) error (%v)", unicom.Usermob, err)
		res["result"] = "1"
		res["errorcode"] = "100"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	var (
		_aesKey = []byte("9ed226d9")
	)
	bs, err := base64.StdEncoding.DecodeString(usermob)
	if err != nil {
		log.Error("base64.StdEncoding.DecodeString(%s) error(%v)", usermob, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err = unicomSvc.DesDecrypt(bs, _aesKey)
	if err != nil {
		log.Error("unicomSvc.DesDecrypt error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var usermobStr string
	if len(bs) > 32 {
		usermobStr = string(bs[:32])
	} else {
		usermobStr = string(bs)
	}
	log.Info("unicomSvc.AdvanceSync_usermob (%v) unicom (%v)", usermobStr, unicom)
	if err := unicomSvc.InAdvanceSync(c, usermobStr, metadata.String(c, metadata.RemoteIP), unicom, time.Now()); err != nil {
		log.Error("unicomSvc.InAdvanceSync usermob (%v) unicom (%v) error (%v)", usermobStr, unicom, err)
		res["result"] = "1"
		res["errorcode"] = "100"
		returnDataJSON(c, res, err)
		return
	}
	res["result"] = "0"
	res["message"] = ""
	res["errorcode"] = ""
	returnDataJSON(c, res, nil)
}

// flowSync
func flowSync(c *bm.Context) {
	res := map[string]interface{}{}
	unicom, err := requestJSONToMap(c.Request)
	if err != nil {
		res["result"] = "1"
		res["errorcode"] = "100"
		log.Error("unicom_ioutil.ReadAll error (%v)", err)
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	var flowbyte float64
	if flowbyte, err = strconv.ParseFloat(unicom.FlowbyteStr, 64); err != nil {
		log.Error("unicom_flowbyte strconv.ParseFloat(%s) error(%v)", unicom.FlowbyteStr, err)
		res["result"] = "1"
		res["errorcode"] = "100"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	usermob, err := url.QueryUnescape(unicom.Usermob)
	if err != nil {
		log.Error("unicom_url.QueryUnescape (%v) error (%v)", unicom.Usermob, err)
		res["result"] = "1"
		res["errorcode"] = "100"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	var (
		_aesKey = []byte("9ed226d9")
	)
	bs, err := base64.StdEncoding.DecodeString(usermob)
	if err != nil {
		log.Error("unicom_base64.StdEncoding.DecodeString(%s) error(%v)", usermob, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err = unicomSvc.DesDecrypt(bs, _aesKey)
	if err != nil {
		log.Error("unicomSvc.DesDecrypt error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var usermobStr string
	if len(bs) > 32 {
		usermobStr = string(bs[:32])
	} else {
		usermobStr = string(bs)
	}
	if err := unicomSvc.FlowSync(c, int(flowbyte*1024), usermobStr, unicom.Time, metadata.String(c, metadata.RemoteIP), time.Now()); err != nil {
		log.Error("unicomSvc.FlowSync error (%v)", err)
		res["result"] = "1"
		res["errorcode"] = "100"
		returnDataJSON(c, res, err)
		return
	}
	res["result"] = "0"
	res["message"] = ""
	res["errorcode"] = ""
	returnDataJSON(c, res, nil)
}

// inIPSync
func inIPSync(c *bm.Context) {
	res := map[string]interface{}{}
	unicom, err := requestIPJSONToMap(c.Request)
	if err != nil {
		res["result"] = "1"
		res["errorcode"] = "100"
		log.Error("unicom_ioutil.ReadAll error (%v)", err)
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	log.Info("unicomSvc.InIpSync_unicom (%v)", unicom)
	if err := unicomSvc.InIPSync(c, metadata.String(c, metadata.RemoteIP), unicom, time.Now()); err != nil {
		log.Error("unicomSvc.InIpSync unicom (%v) error (%v)", unicom, err)
		res["result"] = "1"
		res["errorcode"] = "100"
		returnDataJSON(c, res, err)
		return
	}
	res["result"] = "0"
	res["message"] = ""
	res["errorcode"] = ""
	returnDataJSON(c, res, nil)
}

// userFlow
func userFlow(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	usermob := params.Get("usermob")
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	ipStr := metadata.String(c, metadata.RemoteIP)
	data, msg, err := unicomSvc.UserFlow(c, usermob, mobiApp, ipStr, build, time.Now())
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	returnDataJSON(c, res, nil)
}

// userFlowState
func userFlowState(c *bm.Context) {
	params := c.Request.Form
	usermob := params.Get("usermob")
	data, err := unicomSvc.UserFlowState(c, usermob, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"data": data,
	}
	returnDataJSON(c, res, nil)
}

// userState
func userState(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	usermob := params.Get("usermob")
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	ipStr := metadata.String(c, metadata.RemoteIP)
	data, msg, err := unicomSvc.UserState(c, usermob, mobiApp, ipStr, build, time.Now())
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	returnDataJSON(c, res, nil)
}

// unicomState
func unicomState(c *bm.Context) {
	params := c.Request.Form
	usermob := params.Get("usermob")
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	ipStr := metadata.String(c, metadata.RemoteIP)
	data, err := unicomSvc.UnicomState(c, usermob, mobiApp, ipStr, build, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"data": data,
	}
	returnDataJSON(c, res, nil)
}

// unicomStateM
func unicomStateM(c *bm.Context) {
	params := c.Request.Form
	usermob := params.Get("usermob")
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	ipStr := metadata.String(c, metadata.RemoteIP)
	var (
		_aesKey = []byte("9ed226d9")
	)
	bs, err := base64.StdEncoding.DecodeString(usermob)
	if err != nil {
		log.Error("base64.StdEncoding.DecodeString(%s) error(%v)", usermob, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, err = unicomSvc.DesDecrypt(bs, _aesKey)
	if err != nil {
		log.Error("unicomSvc.DesDecrypt error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var usermobStr string
	if len(bs) > 32 {
		usermobStr = string(bs[:32])
	} else {
		usermobStr = string(bs)
	}
	data, err := unicomSvc.UnicomState(c, usermobStr, mobiApp, ipStr, build, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	res := map[string]interface{}{
		"data": data,
	}
	returnDataJSON(c, res, nil)
}

// RequestJsonToMap
func requestJSONToMap(request *http.Request) (unicom *unicom.UnicomJson, err error) {
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
	log.Info("unicom orders json body(%s)", body)
	if err = json.Unmarshal(body, &unicom); err != nil {
		log.Error("json.Unmarshal UnicomJson(%v) error (%v)", unicom, err)
		return
	}
	if err = unicom.UnicomJSONChange(); err != nil {
		log.Error("unicom.UnicomJSONChange unicom (%v) error (%v)", unicom, err)
	}
	return
}

// RequestIpJsonToMap
func requestIPJSONToMap(request *http.Request) (unicom *unicom.UnicomIpJson, err error) {
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
	log.Info("unicom ip json body(%s)", body)
	json.Unmarshal(body, &unicom)
	return
}

func pack(c *bm.Context) {
	var mid int64
	res := map[string]interface{}{}
	params := c.Request.Form
	usermob := params.Get("usermob")
	midInter, _ := c.Get("mid")
	mid = midInter.(int64)
	msg, err := unicomSvc.Pack(c, usermob, mid, time.Now())
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["message"] = msg
	returnDataJSON(c, res, nil)
}

func isUnciomIP(c *bm.Context) {
	params := c.Request.Form
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	ipStr := metadata.String(c, metadata.RemoteIP)
	ip := model.InetAtoN(ipStr)
	if err := unicomSvc.IsUnciomIP(ip, ipStr, mobiApp, build, time.Now()); err != nil {
		log.Error("unicom_user_ip:%v", ipStr)
		c.JSON(nil, err)
	}
	res := map[string]interface{}{
		"code": ecode.OK,
	}
	returnDataJSON(c, res, nil)
}

func userUnciomIP(c *bm.Context) {
	params := c.Request.Form
	usermob := params.Get("usermob")
	mobiApp := params.Get("mobi_app")
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	ipStr := metadata.String(c, metadata.RemoteIP)
	ip := model.InetAtoN(ipStr)
	res := map[string]interface{}{}
	res["data"] = unicomSvc.UserUnciomIP(ip, ipStr, usermob, mobiApp, build, time.Now())
	returnDataJSON(c, res, nil)
}

func orderPay(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	usermob := params.Get("usermob")
	channel := params.Get("channel")
	ordertypeStr := params.Get("ordertype")
	ordertype, _ := strconv.Atoi(ordertypeStr)
	data, msg, err := unicomSvc.Order(c, usermob, channel, ordertype, time.Now())
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	returnDataJSON(c, res, nil)
}

func orderCancel(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	usermob := params.Get("usermob")
	data, msg, err := unicomSvc.CancelOrder(c, usermob)
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	returnDataJSON(c, res, nil)
}

func smsCode(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phone := params.Get("phone")
	msg, err := unicomSvc.UnicomSMSCode(c, phone, time.Now())
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["message"] = msg
	returnDataJSON(c, res, nil)
}

func addUnicomBind(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phone := params.Get("phone")
	codeStr := params.Get("code")
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		log.Error("code(%v) strconv.Atoi error(%v)", codeStr, err)
		res["message"] = "参数错误"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	var mid int64
	midInter, ok := c.Get("mid")
	if ok {
		mid = midInter.(int64)
	} else {
		res["message"] = "账号未登录"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	msg, err := unicomSvc.AddUnicomBind(c, phone, code, mid, time.Now())
	if err != nil {
		res["code"] = err
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["message"] = msg
	returnDataJSON(c, res, nil)
}

func releaseUnicomBind(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	phoneStr := params.Get("phone")
	phone, err := strconv.Atoi(phoneStr)
	if err != nil {
		log.Error("phone(%v) strconv.Atoi error(%v)", phoneStr, err)
		res["message"] = "参数错误"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	var mid int64
	midInter, ok := c.Get("mid")
	if ok {
		mid = midInter.(int64)
	} else {
		res["message"] = "账号未登录"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	msg, err := unicomSvc.ReleaseUnicomBind(c, mid, phone)
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["message"] = msg
	returnDataJSON(c, res, nil)
}

func userBind(c *bm.Context) {
	res := map[string]interface{}{}
	var mid int64
	midInter, ok := c.Get("mid")
	if ok {
		mid = midInter.(int64)
	} else {
		res["message"] = "账号未登录"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	data, msg, err := unicomSvc.UserBind(c, mid)
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["data"] = data
	res["message"] = ""
	returnDataJSON(c, res, nil)
}

func packList(c *bm.Context) {
	res := map[string]interface{}{
		"data": unicomSvc.UnicomPackList(),
	}
	returnDataJSON(c, res, nil)
}

func packReceive(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	idStr := params.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("id(%v) strconv.ParseInt error(%v)", idStr, err)
		c.JSON(nil, err)
		return
	}
	var mid int64
	midInter, ok := c.Get("mid")
	if ok {
		mid = midInter.(int64)
	} else {
		res["message"] = "账号未登录"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	msg, err := unicomSvc.UnicomPackReceive(c, mid, id, time.Now())
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["message"] = msg
	returnDataJSON(c, res, nil)
}

func flowPack(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	var mid int64
	midInter, ok := c.Get("mid")
	if ok {
		mid = midInter.(int64)
	} else {
		res["message"] = "账号未登录"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	flowID := params.Get("flow_id")
	msg, err := unicomSvc.UnicomFlowPack(c, mid, flowID, time.Now())
	if err != nil {
		res["message"] = msg
		returnDataJSON(c, res, err)
		return
	}
	res["message"] = msg
	returnDataJSON(c, res, nil)
}

func userPacksLog(c *bm.Context) {
	params := c.Request.Form
	starttimeStr := params.Get("starttime")
	pnStr := params.Get("pn")
	pn, err := strconv.Atoi(pnStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pn < 1 {
		pn = 1
	}
	timeLayout := "2006-01"
	loc, _ := time.LoadLocation("Local")
	startTime, err := time.ParseInLocation(timeLayout, starttimeStr, loc)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(unicomSvc.UserPacksLog(c, startTime, time.Now(), pn, metadata.String(c, metadata.RemoteIP)))
}

func userBindLog(c *bm.Context) {
	res := map[string]interface{}{}
	var mid int64
	midInter, ok := c.Get("mid")
	if ok {
		mid = midInter.(int64)
	} else {
		res["message"] = "账号未登录"
		returnDataJSON(c, res, ecode.RequestErr)
		return
	}
	data, err := unicomSvc.UserBindLog(c, mid, time.Now())
	if err != nil {
		c.JSON(nil, err)
		return
	}
	res["data"] = data
	returnDataJSON(c, res, nil)
}

func welfareBindState(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil || mid < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	data := unicomSvc.WelfareBindState(c, mid)
	res := map[string]interface{}{
		"state": data,
	}
	c.JSON(res, nil)
}

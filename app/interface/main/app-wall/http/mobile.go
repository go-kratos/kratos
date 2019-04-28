package http

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"

	"go-common/app/interface/main/app-wall/model/mobile"
	"go-common/library/ecode"
	log "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	mobileOrders = 0
	mobileFlow   = 1
)

var (
	emptyByte = []byte{}
)

// ordersSync
func ordersMobileSync(c *bm.Context) {
	mobile, mtype, err := requestXMLToMap(c.Request)
	if err != nil {
		log.Error("requestXMLToMap error (%v)", err)
		mobileMessage(c, err, mtype)
		return
	}
	b, _ := xml.Marshal(&mobile)
	ip := metadata.String(c, metadata.RemoteIP)
	switch mtype {
	case mobileOrders:
		log.Info("mobile user orders xml response(%s)", b)
		if mobile.Threshold == "" {
			mobile.Threshold = "100"
		}
		if err := mobileSvc.InOrdersSync(c, ip, mobile, time.Now()); err != nil {
			log.Error("mobileSvc.InOrdersSync error (%v)", err)
			mobileMessage(c, err, mtype)
			return
		}
	case mobileFlow:
		log.Info("mobile user flow xml response(%s)", b)
		if err := mobileSvc.FlowSync(c, mobile, ip); err != nil {
			log.Error("mobileSvc.FlowSync error (%v)", err)
			mobileMessage(c, err, mtype)
			return
		}
	}
	mobileMessage(c, ecode.OK, mtype)
}

// mobileActivation
func mobileActivation(c *bm.Context) {
	res := map[string]interface{}{}
	params := c.Request.Form
	usermob := params.Get("usermob")
	msg, err := mobileSvc.Activation(c, usermob, time.Now())
	if err != nil {
		if msg != "" {
			returnDataJSON(c, res, err)
			res["message"] = msg
		}

		returnDataJSON(c, res, err)
		return
	}
	res["message"] = ""
	returnDataJSON(c, res, nil)
}

// mobileState
func mobileState(c *bm.Context) {
	params := c.Request.Form
	usermob := params.Get("usermob")
	data := mobileSvc.MobileState(c, usermob, time.Now())
	res := map[string]interface{}{
		"data": data,
	}
	returnDataJSON(c, res, nil)
}

// userFlowState
func userMobileState(c *bm.Context) {
	params := c.Request.Form
	usermob := params.Get("usermob")
	data := mobileSvc.UserMobileState(c, usermob, time.Now())
	res := map[string]interface{}{
		"data": data,
	}
	returnDataJSON(c, res, nil)
}

// RequestXmlToMap
func requestXMLToMap(request *http.Request) (data *mobile.MobileXML, mobileType int, err error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Error("mobile_ioutil.ReadAll error (%v)", err)
		return
	}
	defer request.Body.Close()
	if len(body) == 0 {
		err = ecode.RequestErr
		return
	}
	var res interface{} = &mobile.OrderXML{}
	mobileType = mobileOrders
	log.Info("mobile orders xml body(%s)", body)
	if err = xml.Unmarshal(body, &res); err != nil {
		res = &mobile.FlowXML{}
		mobileType = mobileFlow
		if err = xml.Unmarshal(body, &res); err != nil {
			log.Error("xml.Unmarshal OrderXML(%v) error (%v)", res, err)
			err = ecode.RequestErr
			return
		}
	}
	if res == nil {
		err = ecode.NothingFound
		log.Error("xml.Unmarshal OrderXML is null")
		return
	}
	switch v := res.(type) {
	case *mobile.OrderXML:
		data = v.MobileXML
	case *mobile.FlowXML:
		data = v.MobileXML
	}
	return
}

func mobileMessage(c *bm.Context, err error, mobileType int) {
	// response header
	c.Writer.Header().Set("Content-Type", "text/xml; charset=UTF-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	var msg interface{}
	switch mobileType {
	case mobileOrders:
		msg = &mobile.OrderMsgXML{
			Msg: &mobile.Msg{
				Xmlns:   "http://www.monternet.com/dsmp/schemas/",
				MsgType: "SyncFlowPkgOrderResp",
				Version: "1.0.0",
				HRet:    err.Error(),
			},
		}
	case mobileFlow:
		msg = &mobile.FlowMsgXML{
			Msg: &mobile.Msg{
				Xmlns:   "http://www.monternet.com/dsmp/schemas/",
				MsgType: "SyncFlowPkgLeftQuotaResp",
				Version: "1.0.0",
				HRet:    err.Error(),
			},
		}
	}
	output, err := xml.MarshalIndent(msg, "  ", "    ")
	if err != nil {
		log.Error("xml.MarshalIndent (%v)", err)
		err = ecode.RequestErr
		output = emptyByte
	}
	c.Writer.Write([]byte(xml.Header))
	c.Writer.Write(output)
}

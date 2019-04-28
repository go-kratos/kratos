package mobile

import (
	"encoding/xml"
	"strconv"
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"
)

type OrderXML struct {
	XMLName xml.Name `xml:"SyncFlowPkgOrderReq"`
	*MobileXML
}

type FlowXML struct {
	XMLName xml.Name `xml:"SyncFlowPkgLeftQuotaReq"`
	*MobileXML
}

type MobileXML struct {
	Orderid        string `xml:"OrderID"`
	Userpseudocode string `xml:"UserPseudoCode"`
	Channelseqid   string `xml:"ChannelSeqId"`
	Price          string `xml:"Price"`
	Actiontime     string `xml:"ActionTime"`
	Actionid       string `xml:"ActionID"`
	Effectivetime  string `xml:"EffectiveTime"`
	Expiretime     string `xml:"ExpireTime"`
	Channelid      string `xml:"ChannelId"`
	Productid      string `xml:"ProductId"`
	Ordertype      string `xml:"OrderType"`
	Threshold      string `xml:"Threshold"`
	Resulttime     string `xml:"ResultTime"`
}

type Mobile struct {
	Orderid        string     `json:"-"`
	Userpseudocode string     `json:"-"`
	Channelseqid   string     `json:"-"`
	Price          int        `json:"-"`
	Actionid       int        `json:"actionid"`
	Effectivetime  xtime.Time `json:"starttime,omitempty"`
	Expiretime     xtime.Time `json:"endtime,omitempty"`
	Channelid      string     `json:"-"`
	Productid      string     `json:"productid,omitempty"`
	Ordertype      int        `json:"-"`
	Threshold      int        `json:"flow"`
	Resulttime     xtime.Time `json:"-"`
	MobileType     int        `json:"orderstatus,omitempty"`
	ProductType    int        `json:"product_type,omitempty"`
}

type MobileIP struct {
	IPStartUint uint32 `json:"-"`
	IPEndUint   uint32 `json:"-"`
}

type MobileUserIP struct {
	IPStr    string `json:"ip"`
	IsValide bool   `json:"is_valide"`
}

// MobileChange
func (u *Mobile) MobileChange() {
	if u.Effectivetime.Time().IsZero() {
		u.Effectivetime = 0
	}
	if u.Expiretime.Time().IsZero() {
		u.Expiretime = 0
	}
	switch u.Productid {
	case "100000000028":
		u.ProductType = 1
	case "100000000030":
		u.ProductType = 2
	}
}

type Msg struct {
	Xmlns   string `xml:"xmlns,attr"`
	MsgType string `xml:"MsgType"`
	Version string `xml:"Version"`
	HRet    string `xml:"hRet"`
}

type OrderMsgXML struct {
	XMLName xml.Name `xml:"SyncFlowPkgOrderResp"`
	*Msg
}

type FlowMsgXML struct {
	XMLName xml.Name `xml:"SyncFlowPkgLeftQuotaResp"`
	*Msg
}

// MobileXMLMobile
func (u *Mobile) MobileXMLMobile(uxml *MobileXML) {
	u.Actionid, _ = strconv.Atoi(uxml.Actionid)
	u.Effectivetime = timeStrToInt(uxml.Effectivetime)
	u.Expiretime = timeStrToInt(uxml.Expiretime)
	u.Threshold, _ = strconv.Atoi(uxml.Threshold)
	u.Productid = uxml.Productid
	u.MobileChange()
}

// timeStrToInt
func timeStrToInt(timeStr string) (timeInt xtime.Time) {
	var err error
	timeLayout := "20060102"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	if err = timeInt.Scan(theTime); err != nil {
		log.Error("timeInt.Scan error(%v)", err)
	}
	return
}

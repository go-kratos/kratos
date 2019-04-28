package unicom

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-wall/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

type Unicom struct {
	Id          int        `json:"-"`
	Spid        int        `json:"spid"`
	CardType    int        `json:"cardtype"`
	TypeInt     int        `json:"type"`
	Unicomtype  int        `json:"unicomtype,omitempty"`
	Ordertypes  int        `json:"-"`
	Channelcode int        `json:"-"`
	Usermob     string     `json:"-"`
	Cpid        string     `json:"-"`
	Ordertime   xtime.Time `json:"ordertime"`
	Canceltime  xtime.Time `json:"canceltime,omitempty"`
	Endtime     xtime.Time `json:"endtime,omitempty"`
	Province    string     `json:"-"`
	Area        string     `json:"-"`
	Videoid     string     `json:"-"`
	Time        xtime.Time `json:"-"`
	Flowbyte    int        `json:"flowbyte"`
}

type UnicomJson struct {
	Usermob     string `json:"usermob"`
	Userphone   string `json:"userphone"`
	Cpid        string `json:"cpid"`
	Spid        string `json:"spid"`
	TypeInt     string `json:"type"`
	Ordertime   string `json:"ordertime"`
	Canceltime  string `json:"canceltime"`
	Endtime     string `json:"endtime"`
	Channelcode string `json:"channelcode"`
	Province    string `json:"province"`
	Area        string `json:"area"`
	Ordertypes  string `json:"ordertype"`
	Videoid     string `json:"videoid"`
	Time        string `json:"time"`
	FlowbyteStr string `json:"flowbyte"`
}

type UnicomIpJson struct {
	Ipbegin   string `json:"ipbegin"`
	Ipend     string `json:"ipend"`
	Provinces string `json:"province"`
	Isopen    string `json:"isopen"`
	Opertime  string `json:"opertime"`
	Sign      string `json:"sign"`
}

type UnicomIP struct {
	Ipbegin     int    `json:"-"`
	Ipend       int    `json:"-"`
	IPStartUint uint32 `json:"-"`
	IPEndUint   uint32 `json:"-"`
}

type UnicomUserIP struct {
	IPStr    string `json:"ip"`
	IsValide bool   `json:"is_valide"`
}

type BroadbandOrder struct {
	Usermob string `json:"userid,omitempty"`
	Endtime string `json:"endtime,omitempty"`
	Channel string `json:"channel,omitempty"`
}

type UserBind struct {
	Usermob  string    `json:"usermob,omitempty"`
	Phone    int       `json:"phone"`
	Mid      int64     `json:"mid"`
	Name     string    `json:"name,omitempty"`
	State    int       `json:"state,omitempty"`
	Integral int       `json:"integral"`
	Flow     int       `json:"flow"`
	Monthly  time.Time `json:"monthly,omitempty"`
}

type UserPack struct {
	ID       int64  `json:"id"`
	Type     int    `json:"type"`
	Desc     string `json:"desc"`
	Amount   int    `json:"amount"`
	Capped   int8   `json:"capped"`
	Integral int    `json:"integral"`
	Param    string `json:"param"`
	State    int    `json:"state,omitempty"`
}

type UserPackLimit struct {
	IsLimit int `json:"is_limit"`
	Count   int `json:"count"`
}

type UnicomUserFlow struct {
	Phone      int    `json:"phone"`
	Mid        int64  `json:"mid"`
	Integral   int    `json:"integral"`
	Flow       int    `json:"flow"`
	Outorderid string `json:"outorderid"`
	Orderid    string `json:"orderid"`
	Desc       string `json:"desc"`
}

type UserPackLog struct {
	Phone     int    `json:"phone,omitempty"`
	Usermob   string `json:"usermob,omitempty"`
	Mid       int64  `json:"mid,omitempty"`
	RequestNo string `json:"request_no,omitempty"`
	Type      int    `json:"pack_type"`
	Desc      string `json:"-"`
	UserDesc  string `json:"pack_desc,omitempty"`
	Integral  int    `json:"integral,omitempty"`
}

type UserLog struct {
	Phone    int    `json:"phone,omitempty"`
	Integral int    `json:"integral,omitempty"`
	Desc     string `json:"pack_desc,omitempty"`
	Ctime    string `json:"ctime,omitempty"`
}

type UserBindInfo struct {
	MID    int64  `json:"mid"`
	Phone  int    `json:"phone"`
	Action string `json:"action"`
}

// UnicomChange
func (u *Unicom) UnicomChange() {
	if u.Canceltime.Time().IsZero() {
		u.Canceltime = 0
	}
	if u.Endtime.Time().IsZero() {
		u.Endtime = 0
	}
	switch u.Spid {
	case 10019:
		u.CardType = 1
	case 10020:
		u.CardType = 2
	case 10021:
		u.CardType = 3
	case 979:
		u.CardType = 4
	}
}

func (u *UnicomJson) UnicomJSONChange() (err error) {
	if u.Ordertypes != "" {
		if _, err = strconv.Atoi(u.Ordertypes); err != nil {
			log.Error("UnicomJsonChange error(%v)", u)
		}
	}
	return
}

func (u *UnicomIP) UnicomIPChange() {
	u.IPStartUint = u.unicomIPTOUint(u.Ipbegin)
	u.IPEndUint = u.unicomIPTOUint(u.Ipend)
}

func (u *UnicomIP) unicomIPTOUint(ip int) (ipUnit uint32) {
	var (
		ip1, ip2, ip3, ip4 int
		ipStr              string
	)
	var _initIP = "%d.%d.%d.%d"
	ip1 = ip / 1000000000
	ip2 = (ip / 1000000) % 1000
	ip3 = (ip / 1000) % 1000
	ip4 = ip % 1000
	ipStr = fmt.Sprintf(_initIP, ip1, ip2, ip3, ip4)
	ipUnit = model.InetAtoN(ipStr)
	return
}

// UnicomJSONTOUincom
func (u *Unicom) UnicomJSONTOUincom(usermob string, ujson *UnicomJson) {
	u.Spid, _ = strconv.Atoi(ujson.Spid)
	u.Ordertime = timeStrToInt(ujson.Ordertime)
	u.Canceltime = timeStrToInt(ujson.Canceltime)
	u.Endtime = timeStrToInt(ujson.Endtime)
	u.TypeInt, _ = strconv.Atoi(ujson.TypeInt)
	u.Ordertypes, _ = strconv.Atoi(ujson.Ordertypes)
	u.Channelcode, _ = strconv.Atoi(ujson.Channelcode)
	u.Usermob = usermob
	u.Cpid = ujson.Cpid
	u.Province = ujson.Province
	u.UnicomChange()
}

// timeStrToInt
func timeStrToInt(timeStr string) (timeInt xtime.Time) {
	var err error
	timeLayout := "20060102150405"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	if err = timeInt.Scan(theTime); err != nil {
		log.Error("timeInt.Scan error(%v)", err)
	}
	return
}

func (u *UnicomIP) UnicomIPStrToint(ipstart, ipend string) {
	u.Ipbegin = ipToInt(ipstart)
	u.Ipend = ipToInt(ipend)
}

// ipToint
func ipToInt(ipString string) (ipInt int) {
	tmp := strings.Split(ipString, ".")
	if len(tmp) < 4 {
		return
	}
	var ipStr string
	for _, tip := range tmp {
		var (
			ipLen = len(tip)
			last  int
			ip1   string
		)
		if ipLen < 3 {
			last = 3 - ipLen
			switch last {
			case 1:
				ip1 = "0" + tip
			case 2:
				ip1 = "00" + tip
			case 3:
				ip1 = "000"
			}
		} else {
			ip1 = tip
		}
		ipStr = ipStr + ip1
	}
	ipInt, _ = strconv.Atoi(ipStr)
	return
}

func (t *UserLog) UserLogJSONChange(jsonData string) (err error) {
	if err = json.Unmarshal([]byte(jsonData), &t); err != nil {
		return
	}
	return
}

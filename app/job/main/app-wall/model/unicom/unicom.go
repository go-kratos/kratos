package unicom

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/app-wall/model"
	xtime "go-common/library/time"
)

type UserBind struct {
	Usermob  string    `json:"usermob,omitempty"`
	Phone    int       `json:"phone"`
	Mid      int64     `json:"mid"`
	State    int       `json:"state,omitempty"`
	Integral int       `json:"integral"`
	Flow     int       `json:"flow"`
	Monthly  time.Time `json:"monthly"`
}

// type

type ClickMsg struct {
	Plat       int8
	AID        int64
	MID        int64
	Lv         int8
	BvID       string
	CTime      int64
	STime      int64
	IP         string
	KafkaBs    []byte
	EpID       int64
	SeasonType int
	UserAgent  string
}

type Unicom struct {
	Usermob   string     `json:"-"`
	Spid      int        `json:"spid"`
	TypeInt   int        `json:"type"`
	Ordertime xtime.Time `json:"ordertime"`
	Endtime   xtime.Time `json:"endtime,omitempty"`
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

type UnicomIP struct {
	Ipbegin     int    `json:"-"`
	Ipend       int    `json:"-"`
	IPStartUint uint32 `json:"-"`
	IPEndUint   uint32 `json:"-"`
}

type UserPackLog struct {
	Phone     int    `json:"-"`
	Usermob   string `json:"-"`
	Mid       int64  `json:"-"`
	RequestNo string `json:"-"`
	Type      int    `json:"-"`
	Desc      string `json:"-"`
	Integral  int    `json:"-"`
}

type UserIntegralLog struct {
	Phone      int    `json:"-"`
	Mid        int64  `json:"-"`
	UnicomDesc string `json:"-"`
	Type       int    `json:"-"`
	Integral   int    `json:"-"`
	Flow       int    `json:"-"`
	Desc       string `json:"-"`
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

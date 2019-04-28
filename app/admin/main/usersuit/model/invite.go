package model

import (
	"net"
	"strconv"

	accmdl "go-common/app/service/main/account/model"
	xtime "go-common/library/time"
)

const (
	// StatusOK status ok
	StatusOK = 0
	// StatusUsed status used
	StatusUsed = 1
	// StatusExpires status expires
	StatusExpires = 2
)

// Invite invite.
type Invite struct {
	Status  int64      `json:"status"`
	Mid     int64      `json:"mid"`
	Code    string     `json:"invite_code"`
	IP      uint32     `json:"-"` // legacy IP field
	IPng    []byte     `json:"-"`
	Ctime   xtime.Time `json:"buy_time"`
	Expires int64      `json:"expires"`
	Imid    int64      `json:"invited_mid,omitempty"`
	UsedAt  int64      `json:"used_at,omitempty"`
}

// BuyIPString is
func (inv *Invite) BuyIPString() string {
	if inv.IP != 0 {
		return inetNtoA(inv.IP)
	}
	return net.IP(inv.IPng).String()
}

func inetNtoA(sum uint32) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte((sum >> 24) & 0xFF)
	ip[1] = byte((sum >> 16) & 0xFF)
	ip[2] = byte((sum >> 8) & 0xFF)
	ip[3] = byte(sum & 0xFF)
	return ip.String()
}

// FillStatus fill status.
func (inv *Invite) FillStatus(now int64) {
	if inv.Used() {
		inv.Status = StatusUsed
		return
	}
	if inv.Expired(now) {
		inv.Status = StatusExpires
		return
	}
	inv.Status = StatusOK
}

// Used check if used.
func (inv *Invite) Used() bool {
	return inv.UsedAt > 0 && inv.Imid > 0
}

// Expired check if expired.
func (inv *Invite) Expired(now int64) bool {
	return now > inv.Expires
}

// RichInvite rich invite with invitee info.
type RichInvite struct {
	Status  int64      `json:"status"`
	Mid     int64      `json:"mid"`
	Code    string     `json:"invite_code"`
	BuyIP   string     `json:"buy_ip"`
	Ctime   xtime.Time `json:"buy_time"`
	Expires int64      `json:"expires"`
	Invitee *Invitee   `json:"invitee,omitempty"`
	UsedAt  int64      `json:"used_at,omitempty"`
}

// NewRichInvite new a rich invite.
func NewRichInvite(inv *Invite, info *accmdl.Info) *RichInvite {
	if inv == nil {
		return nil
	}
	var invt *Invitee
	if inv.Used() {
		if info != nil {
			invt = &Invitee{
				Mid:   inv.Imid,
				Uname: info.Name,
				Face:  info.Face,
			}
		} else {
			invt = &Invitee{
				Mid:   inv.Imid,
				Uname: "用户" + strconv.FormatInt(inv.Imid, 10),
				Face:  "http://static.hdslb.com/images/member/noface.gif",
			}
		}
	}
	return &RichInvite{
		Status:  inv.Status,
		Mid:     inv.Mid,
		Code:    inv.Code,
		Ctime:   inv.Ctime,
		Expires: inv.Expires,
		Invitee: invt,
		UsedAt:  inv.UsedAt,
		BuyIP:   inv.BuyIPString(),
	}
}

// Invitee invited.
type Invitee struct {
	Mid   int64  `json:"mid"`
	Uname string `json:"uname"`
	Face  string `json:"face"`
}

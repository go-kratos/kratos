package model

import (
	"strconv"

	accmdl "go-common/app/service/main/account/model"
	usmdl "go-common/app/service/main/usersuit/model"
	xtime "go-common/library/time"
)

// RichInviteStat rich invite stat.
type RichInviteStat struct {
	Mid           int64         `json:"mid"`
	CurrentLimit  int64         `json:"current_limit"`
	CurrentBought int64         `json:"current_bought"`
	TotalBought   int64         `json:"total_bought"`
	TotalUsed     int64         `json:"total_used"`
	InviteCodes   []*RichInvite `json:"invite_codes"`
}

// RichInvite rich invite.
type RichInvite struct {
	Status  int64      `json:"status"`
	Mid     int64      `json:"mid"`
	Code    string     `json:"invite_code"`
	Ctime   xtime.Time `json:"buy_time"`
	Expires int64      `json:"expires"`
	Invitee *Invitee   `json:"invitee,omitempty"`
	UsedAt  int64      `json:"used_at,omitempty"`
}

// NewRichInvite new a rich invite.
func NewRichInvite(inv *usmdl.Invite, info *accmdl.Info) *RichInvite {
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
	}
}

// Invitee invitee.
type Invitee struct {
	Mid   int64  `json:"mid"`
	Uname string `json:"uname"`
	Face  string `json:"face"`
}

package model

import (
	xtime "go-common/library/time"
)

const (
	// StatusOk ok
	StatusOk = 0
	// StatusUsed used
	StatusUsed = 1
	// StatusExpires expire
	StatusExpires = 2
)

// Invite invaite
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
	Mtime   xtime.Time `json:"-"`
}

// FillStatus fill status
func (inv *Invite) FillStatus(now int64) {
	if inv.Used() {
		inv.Status = StatusUsed
		return
	}
	if inv.Expired(now) {
		inv.Status = StatusExpires
		return
	}
	inv.Status = StatusOk
}

// Used use
func (inv *Invite) Used() bool {
	return inv.UsedAt > 0 && inv.Imid > 0
}

// Expired expire
func (inv *Invite) Expired(now int64) bool {
	return now > inv.Expires
}

// SortInvitesByCtimeDesc sort
type SortInvitesByCtimeDesc []*Invite

// Len len
func (invs SortInvitesByCtimeDesc) Len() int {
	return len(invs)
}

// Less less
func (invs SortInvitesByCtimeDesc) Less(i, j int) bool {
	return int64(invs[i].Ctime) > int64(invs[j].Ctime)
}

// Swap swap
func (invs SortInvitesByCtimeDesc) Swap(i, j int) {
	tmp := invs[i]
	invs[i] = invs[j]
	invs[j] = tmp
}

// InviteStat stat
type InviteStat struct {
	Mid           int64     `json:"mid"`
	CurrentLimit  int64     `json:"current_limit"`
	CurrentBought int64     `json:"current_bought"`
	TotalBought   int64     `json:"total_bought"`
	TotalUsed     int64     `json:"total_used"`
	InviteCodes   []*Invite `json:"invite_codes"`
}

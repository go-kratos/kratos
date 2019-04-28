package model

import (
	"fmt"
	"strconv"
	"strings"

	"go-common/library/ecode"
)

// audit related consts
const (
	_ugcPrefix = "ugc"
	_pgcPrefix = "xds"
	// the field valid options
	_hidden = 0
	_online = 1
	// pgc audit status
	_seasonReject = 0
	_seasonPass   = 1
	// ugc audit status
	_ugcPass   = 1
	_ugcReject = 2
	// audit type
	_reject = "1"
	// IDlist's type field
	_season = "1"
	// content type
	PgcSn    = "7"
	PgcEp    = "8"
	UgcArc   = "9"
	UgcVideo = "10"
)

// IDList def.
type IDList struct {
	Type     string `json:"type"`
	VID      string `json:"vid"`
	Action   string `json:"action"`
	AuditMsg string `json:"audit_msg"`
}

// IsReject def.
func (v *IDList) IsReject() bool {
	return v.Action == _reject
}

// IsShell tells whether it's about archive/season
func (v *IDList) IsShell() bool {
	return v.Type == _season
}

// AuditOp def.
type AuditOp struct {
	KID         int64 // aid/cid/sid/epid
	Result      int   // pgc sn: `check`, pgc ep: state, ugc: result
	Valid       int
	AuditMsg    string
	ContentType string // type
}

// ToMsg def
func (v *AuditOp) ToMsg() string {
	return fmt.Sprintf("audit_Type(%s)_KID(%d)", v.ContentType, v.KID)
}

// parse prefix and get the real ID
func parsePrefix(value string, prefix string) (res bool, vid int64) {
	if strings.Contains(value, prefix) {
		res = true
		ids := strings.Split(value, prefix)
		vid, _ = strconv.ParseInt(ids[1], 10, 64)
	}
	return
}

// FromIDList def.
func (v *AuditOp) FromIDList(req *IDList) (err error) {
	var (
		isUGC, isPGC bool
	)
	// auditMsg treatment
	v.AuditMsg = req.AuditMsg
	// KID & content type treatment
	if isPGC, v.KID = parsePrefix(req.VID, _pgcPrefix); !isPGC { // not pgc, try ugc
		if isUGC, v.KID = parsePrefix(req.VID, _ugcPrefix); !isUGC { // not ugc, unknown type
			return ecode.RequestErr // unknown type
		}
	}
	// Valid treatment
	if req.IsReject() { // decide the valid value
		v.Valid = _hidden
	} else {
		v.Valid = _online
	}
	// Result & Content Type treatment
	if isPGC { // pgc
		if req.IsShell() { // season
			v.ContentType = PgcSn
			if req.IsReject() {
				v.Result = _seasonReject
			} else {
				v.Result = _seasonPass
			}
		} else { // ep
			v.ContentType = PgcEp
			if req.IsReject() {
				v.Result = _epRejected
			} else {
				v.Result = _epPass
			}
		}
	} else { // ugc
		if req.IsShell() {
			v.ContentType = UgcArc
		} else {
			v.ContentType = UgcVideo
		}
		if req.IsReject() {
			v.Result = _ugcReject
		} else {
			v.Result = _ugcPass
		}
	}
	return
}

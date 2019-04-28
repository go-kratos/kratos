package search

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/admin/main/workflow/model"
	"go-common/library/xstr"
)

// AuditLogGroupSearchCond is the common condition model to send challenge search request
type AuditLogGroupSearchCond struct {
	Group      []string
	UNames     []string
	Uids       []int64
	Businesses []int64
	Type       []int64
	Oids       []int64
	Actions    []string
	CTimeFrom  string
	CTimeTo    string

	Int0     []int64
	Int0From int64
	Int0To   int64
	Int1     []int64
	Int1From int64
	Int1To   int64
	Int2     []int64
	Int2From int64
	Int2To   int64

	Str0 string
	Str1 string
	Str2 string

	PN    int64
	PS    int64
	Order string
	Sort  string
}

// AuditLogSearchResult is the model to parse search challenge common result
type AuditLogSearchResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int32  `json:"ttl"`

	Data struct {
		Order  string     `json:"order"`
		Sort   string     `json:"sort"`
		Page   model.Page `json:"page"`
		Result []struct {
			Action    string `json:"action"`
			Business  int64  `json:"business"`
			CTime     string `json:"ctime"`
			ExtraData string `json:"extra_data"`
			Str0      string `json:"str_0"`
			Str1      string `json:"str_1"`
			Str2      string `json:"str_2"`
			Oid       int64  `json:"oid"`
			Type      int64  `json:"type"`
			UID       int64  `json:"uid"`
			UName     string `json:"uname"`
		} `json:"result"`
	} `json:"data"`
}

// ArchiveAuditLogExtra archive audit log extra message
type ArchiveAuditLogExtra struct {
	Content struct {
		ID      int64  `json:"id"` //oid
		UID     int64  `json:"uid"`
		UName   string `json:"uname"`
		Note    string `json:"note"`
		Content string `json:"content"`
	} `json:"content"`
	Diff string `json:"diff"`
}

// Query make query for AuditLogGroupSearchCond
func (alsc *AuditLogGroupSearchCond) Query() (uv url.Values) {
	uv = url.Values{}

	// AppID
	uv.Set("appid", _auditLogSrhComID)

	// Common
	uv.Set("pn", "1")
	uv.Set("ps", "20")
	uv.Set("order", "ctime")
	uv.Set("sort", "desc")
	if alsc.PN != 0 {
		uv.Set("pn", fmt.Sprintf("%d", alsc.PN))
	}
	if alsc.PS != 0 {
		uv.Set("ps", fmt.Sprintf("%d", alsc.PS))
	}
	if alsc.Order != "" {
		uv.Set("order", alsc.Order)
	}
	if alsc.Sort != "" {
		uv.Set("sort", alsc.Sort)
	}
	if alsc.Int0From != 0 {
		uv.Set("int_0_from", strconv.FormatInt(alsc.Int0From, 10))
	}

	if alsc.Int0To != 0 {
		uv.Set("int_0_to", strconv.FormatInt(alsc.Int0To, 10))
	}

	if alsc.Int1From != 0 {
		uv.Set("int_1_from", strconv.FormatInt(alsc.Int1From, 10))
	}

	if alsc.Int1To != 0 {
		uv.Set("int_1_to", strconv.FormatInt(alsc.Int1To, 10))
	}

	if alsc.Int2From != 0 {
		uv.Set("int_2_from", strconv.FormatInt(alsc.Int2From, 10))
	}

	if alsc.Int2To != 0 {
		uv.Set("int_2_to", strconv.FormatInt(alsc.Int2To, 10))
	}

	// Group related
	uv.Set("group", strings.Join(alsc.Group, ","))
	uv.Set("uname", strings.Join(alsc.UNames, ","))
	uv.Set("uid", xstr.JoinInts(alsc.Uids))
	uv.Set("business", xstr.JoinInts(alsc.Businesses))
	uv.Set("type", xstr.JoinInts(alsc.Type))
	uv.Set("oid", xstr.JoinInts(alsc.Oids))
	uv.Set("action", strings.Join(alsc.Actions, ","))
	uv.Set("ctime_from", alsc.CTimeFrom)
	uv.Set("ctime_to", alsc.CTimeTo)

	uv.Set("int_0", xstr.JoinInts(alsc.Int0))
	uv.Set("int_1", xstr.JoinInts(alsc.Int1))
	uv.Set("int_2", xstr.JoinInts(alsc.Int2))

	uv.Set("str_0", alsc.Str0)
	uv.Set("str_1", alsc.Str1)
	uv.Set("str_2", alsc.Str2)

	return uv
}

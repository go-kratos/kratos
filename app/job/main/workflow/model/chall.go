package model

import (
	"fmt"
	"net/url"
	"time"
)

const (
	_challSrhComID = "workflow_chall_common"
	// QueueState .
	QueueState = 18
)

// Chall .
type Chall struct {
	ID            int64     `json:"id"`
	Business      int64     `json:"business"`
	DispatchState int       `json:"dispatch_state"`
	DispatchTime  time.Time `json:"dispatch_time"`
}

// ChallSearchRes .
type ChallSearchRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int32  `json:"ttl"`
	Data    struct {
		Order string `json:"order"`
		Sort  string `json:"sort"`
		Page  struct {
			Num   int64 `json:"num"`
			Size  int64 `json:"size"`
			Total int64 `json:"total"`
		} `json:"page"`
		Result []struct {
			ID int64 `json:"id"`
		} `json:"result"`
	} `json:"data"`
}

// ChallSearchParams .
type ChallSearchParams struct {
	Business            string
	States              string
	BusinessStates      string
	AssigneeAdminIDs    string
	AssigneeAdminIDsNot string
	MtimeTo             string
	PN                  int64
	PS                  int64
	Order               string
	Sort                string
}

// Serialize .
func (cp *ChallSearchParams) Serialize() (val url.Values) {
	val = url.Values{}
	val.Set("appid", _challSrhComID)
	val.Set("business", cp.Business)
	if cp.States != "" {
		val.Set("states", cp.States)
	}
	if cp.BusinessStates != "" {
		val.Set("business_states", cp.BusinessStates)
	}
	if cp.AssigneeAdminIDs != "" {
		val.Set("assignee_adminids", cp.AssigneeAdminIDs)
	}
	if cp.AssigneeAdminIDsNot != "" {
		val.Set("assignee_adminids_not", cp.AssigneeAdminIDsNot)
	}
	if cp.PN == 0 {
		val.Set("pn", "1")
	} else {
		val.Set("pn", fmt.Sprintf("%d", cp.PN))
	}
	if cp.PS == 0 {
		val.Set("ps", "200")
	} else {
		val.Set("ps", fmt.Sprintf("%d", cp.PS))
	}
	if cp.Order == "" {
		val.Set("order", "ctime")
	} else {
		val.Set("order", cp.Order)
	}
	if cp.Sort == "" {
		val.Set("sort", "desc")
	} else {
		val.Set("sort", cp.Sort)
	}
	if cp.MtimeTo != "" {
		val.Set("mtime_to", cp.MtimeTo)
	}
	return
}

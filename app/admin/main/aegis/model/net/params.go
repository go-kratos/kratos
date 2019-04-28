package net

import (
	"go-common/app/admin/main/aegis/model/common"
	xtime "go-common/library/time"
)

const (
	//StateAvailable .
	StateAvailable = "available"
	//StateDisable .
	StateDisable = "disable"
)

// ListNetParam .
type ListNetParam struct {
	BusinessID int64   `json:"business_id" form:"business_id" validte:"omitempty,gt=0"`
	ID         []int64 `json:"id" form:"id,split" validate:"omitempty,dive,gt=0"`
	State      string  `json:"state" form:"state"`
	Sort       string  `json:"sort" form:"sort" default:"desc"`
	Pn         int64   `json:"pn" form:"pn" default:"1" validate:"gt=0"`
	Ps         int64   `json:"ps" form:"ps" default:"30" validate:"gt=0"`
}

// ListNetItem .
type ListNetItem struct {
	ID           int64             `json:"id"`
	BusinessID   int64             `json:"business_id"`
	ChName       string            `json:"ch_name"`
	Description  string            `json:"description"`
	StartFlowID  int64             `json:"start_flow_id"`
	PID          int64             `json:"pid" gorm:"column:pid"`
	UID          int64             `json:"uid" gorm:"column:uid"`
	DisableTime  xtime.Time        `json:"disable_time"`
	Ctime        common.FormatTime `json:"ctime"`
	Mtime        common.FormatTime `json:"mtime"`
	BusinessName string            `json:"business_name"`
	Username     string            `json:"username"`
}

// ListNetRes .
type ListNetRes struct {
	Result []*ListNetItem `json:"result"`
	Pager  Pager          `json:"pager"`
}

// NetEditParam .
type NetEditParam struct {
	ID          int64  `form:"id" validate:"required,gt=0"`
	ChName      string `form:"ch_name" validate:"required,max=32"`
	Description string `form:"description" validate:"max=60"`
}

// ListTokenParam .
type ListTokenParam struct {
	NetID  int64   `json:"net_id" form:"net_id" validate:"required,gt=0"`
	Name   string  `json:"name" form:"name" validate:"omitempty,max=10"`
	ID     []int64 `json:"id" form:"id,split" validate:"omitempty,dive,gt=0"`
	Sort   string  `json:"sort" form:"sort" default:"desc"`
	Pn     int64   `json:"pn" form:"pn" default:"1" validate:"gt=0"`
	Ps     int64   `json:"ps" form:"ps" default:"30" validate:"gt=0"`
	Assign bool    `json:"assign" form:"assign"`
}

// GroupTokenParam .
type GroupTokenParam struct {
	BusinessID int64  `json:"business_id" form:"business_id" validate:"omitempty,gt=0"`
	NetID      int64  `json:"net_id" form:"net_id" validate:"omitempty,gt=0"`
	GroupName  string `json:"group_name" form:"group_name" validate:"omitempty,max=10"`
}

// ListTokenItem .
type ListTokenItem struct {
	ID       int64             `json:"id"`
	NetID    int64             `json:"net_id"`
	ChName   string            `json:"ch_name"`
	Name     string            `json:"name"`
	Compare  int8              `json:"compare"`
	Value    string            `json:"value"`
	Type     int8              `json:"type"`
	UID      int64             `json:"uid" gorm:"column:uid"`
	Ctime    common.FormatTime `json:"ctime"`
	Mtime    common.FormatTime `json:"mtime"`
	Username string            `json:"username"`
}

// ListTokenRes .
type ListTokenRes struct {
	Result []*ListTokenItem `json:"result"`
	Pager  Pager            `json:"pager"`
}

// ListNetElementParam .
type ListNetElementParam struct {
	NetID int64   `json:"net_id" form:"net_id" validate:"required,gt=0"`
	Name  string  `json:"name" form:"name" validate:"omitempty,max=10"`
	ID    []int64 `json:"id" form:"id,split" validate:"omitempty,dive,gt=0"`
	State string  `json:"state" form:"state"`
	Sort  string  `json:"sort" form:"sort" default:"desc"`
	Pn    int64   `json:"pn" form:"pn" default:"1" validate:"gt=0"`
	Ps    int64   `json:"ps" form:"ps" default:"30" validate:"gt=0"`
}

// ListFlowItem .
type ListFlowItem struct {
	ID          int64             `json:"id"`
	NetID       int64             `json:"net_id"`
	Name        string            `json:"name"`
	ChName      string            `json:"ch_name"`
	Description string            `json:"description"`
	UID         int64             `json:"uid" gorm:"column:uid"`
	DisableTime xtime.Time        `json:"disable_time"`
	Ctime       common.FormatTime `json:"ctime"`
	Mtime       common.FormatTime `json:"mtime"`
	Username    string            `json:"username"`
	Tokens      []string          `json:"tokens"`
	IsStart     bool              `json:"is_start"`
}

// ListFlowRes .
type ListFlowRes struct {
	Result []*ListFlowItem `json:"result"`
	Pager  Pager           `json:"pager"`
}

// ListTransitionItem .
type ListTransitionItem struct {
	ID          int64             `json:"id"`
	NetID       int64             `json:"net_id"`
	Trigger     int8              `json:"trigger"`
	Limit       int64             `json:"limit"`
	Name        string            `json:"name"`
	ChName      string            `json:"ch_name"`
	Description string            `json:"description"`
	UID         int64             `json:"uid" gorm:"column:uid"`
	DisableTime xtime.Time        `json:"disable_time"`
	Ctime       common.FormatTime `json:"ctime"`
	Mtime       common.FormatTime `json:"mtime"`
	Username    string            `json:"username"`
	Tokens      []string          `json:"tokens"`
}

// ListTransitionRes .
type ListTransitionRes struct {
	Result []*ListTransitionItem `json:"result"`
	Pager  Pager                 `json:"pager"`
}

// ListDirectionParam .
type ListDirectionParam struct {
	NetID        int64   `json:"net_id" form:"net_id" validate:"required,gt=0"`
	FlowID       int64   `json:"flow_id" form:"flow_id" validate:"omitempty,gt=0"`
	TransitionID int64   `json:"transition_id" form:"transition_id" validate:"omitempty,gt=0"`
	Direction    int8    `json:"direction" form:"direction"`
	ID           []int64 `json:"id" form:"id,split" validate:"omitempty,dive,gt=0"`
	State        string  `json:"state" form:"state"`
	Sort         string  `json:"sort" form:"sort" default:"desc"`
	Pn           int64   `json:"pn" form:"pn" default:"1" validate:"gt=0"`
	Ps           int64   `json:"ps" form:"ps" default:"30" validate:"gt=0"`
}

// ListDirectionItem .
type ListDirectionItem struct {
	ID             int64             `json:"id"`
	NetID          int64             `json:"net_id"`
	FlowID         int64             `json:"flow_id"`
	TransitionID   int64             `json:"transition_id"`
	Direction      int8              `json:"direction"`
	Order          int8              `json:"order"`
	Guard          string            `json:"guard"`
	Output         string            `json:"output"`
	UID            int64             `json:"uid" gorm:"column:uid"`
	DisableTime    xtime.Time        `json:"disable_time"`
	Ctime          common.FormatTime `json:"ctime"`
	Mtime          common.FormatTime `json:"mtime"`
	FlowName       string            `json:"flow_name"`
	TransitionName string            `json:"transition_name"`
	UserName       string            `json:"user_name"`
	NetName        string            `json:"net_name"`
}

// ListDirectionRes .
type ListDirectionRes struct {
	Result []*ListDirectionItem `json:"result"`
	Pager  Pager                `json:"pager"`
}

// SwitchParam .
type SwitchParam struct {
	ID      int64 `json:"id" form:"id" validate:"required,gt=0"`
	Disable bool  `json:"disable" form:"disable"`
}

// FlowEditParam .
type FlowEditParam struct {
	NetID       int64   `form:"net_id" validate:"omitempty,gt=0"`
	ID          int64   `form:"id" validate:"omitempty,gt=0"`
	Name        string  `form:"name" validate:"required,max=32"`
	ChName      string  `form:"ch_name" validate:"required,max=16"`
	Description string  `form:"description" validate:"max=60"`
	IsStart     bool    `form:"is_start" default:"false"`
	TokenIDList []int64 `form:"token_id_list,split" validate:"omitempty,dive,gt=0"`
}

// TokenBindParam .
type TokenBindParam struct {
	ID      int64  `json:"id,omitemtpy"`
	ChName  string `json:"ch_name,omitempty"`
	TokenID string `json:"token_id"`
	Type    int8   `json:"type"`
}

// TransitionEditParam .
type TransitionEditParam struct {
	NetID       int64             `json:"net_id" validate:"omitempty,gt=0"`
	ID          int64             `json:"id" validate:"omitempty,gt=0"`
	Trigger     int8              `json:"trigger" validate:"omitempty,min=1,max=1"`
	Limit       int64             `json:"limit"`
	Name        string            `json:"name" validate:"required,max=32"`
	ChName      string            `json:"ch_name" validate:"required,max=16"`
	Description string            `json:"description" validate:"max=60"`
	TokenList   []*TokenBindParam `json:"token_list" validate:"required,gt=0" `
}

// DirEditParam .
type DirEditParam struct {
	ID           int64 `form:"id" validate:"omitempty,gt=0"`
	NetID        int64 `form:"net_id" validate:"omitempty,gt=0"`
	FlowID       int64 `form:"flow_id" validate:"required,gt=0"`
	TransitionID int64 `form:"transition_id" validate:"required,gt=0"`
	Direction    int8  `form:"direction" validate:"required,min=1,max=2"`
	Order        int8  `form:"order" validate:"min=0,max=0"`
	Output       string
	Guard        string
}

// Pager .
type Pager struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// ShowFlowResult .
type ShowFlowResult struct {
	Flow    *Flow        `json:"flow"`
	Tokens  []*TokenBind `json:"tokens"`
	IsStart bool         `json:"is_start"`
}

// ShowTransitionResult .
type ShowTransitionResult struct {
	Transition *Transition  `json:"transition"`
	Tokens     []*TokenBind `json:"tokens"`
}

// ShowDirectionResult .
type ShowDirectionResult struct {
	Direction      *Direction `json:"direction"`
	FlowName       string     `json:"flow_name"`
	TransitionName string     `json:"transition_name"`
}

// TokenPackage .
type TokenPackage struct {
	Values      map[string]interface{} `json:"values"`
	TokenIDList []int64                `json:"token_id_list"`
	ChName      string                 `json:"ch_name"`
	HitAudit    bool                   `json:"hit_audit"`
}

// TriggerResult .
type TriggerResult struct {
	RID          int64         `json:"rid"`
	NetID        int64         `json:"net_id"`
	SubmitToken  *TokenPackage `json:"submit_token"`
	ResultToken  *TokenPackage `json:"result_token"`
	NewFlowID    int64         `json:"new_flow_id"`
	OldFlowID    string        `json:"old_flow_id"` //可能有多个旧节点
	TransitionID []int64       `json:"transition_id"`
	From         string        `json:"from"`
}

// TranOperation .
type TranOperation struct {
	BindIDList string `json:"bind_id_list"`
	ChName     string `json:"ch_name"`
}

// TranOperationArr .
type TranOperationArr []*TranOperation

func (a TranOperationArr) Len() int {
	return len(a)
}

func (a TranOperationArr) Less(i, j int) bool {
	return a[i].BindIDList < a[j].BindIDList
}

func (a TranOperationArr) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// TransitionInfo .
type TransitionInfo struct {
	RID        int64            `json:"rid"`
	FlowID     int64            `json:"flow_id"`
	Operations []*TranOperation `json:"operations"`
}

// JumpFlowResult .
type JumpFlowResult struct {
	RID         int64         `json:"rid"`
	NetID       int64         `json:"net_id"`
	SubmitToken *TokenPackage `json:"submit_token"`
	ResultToken *TokenPackage `json:"result_token"`
	NewFlowID   int64         `json:"new_flow_id"`
	OldFlowID   string        `json:"old_flow_id"`
}

// SimpleInfo .
type SimpleInfo struct {
	ID     int64  `json:"id"`
	ChName string `json:"ch_name"`
}

// JumpInfo .
type JumpInfo struct {
	Flows      []*SimpleInfo    `json:"flows"`
	Operations []*TranOperation `json:"operations"`
}

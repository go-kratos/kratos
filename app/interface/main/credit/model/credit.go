package model

import (
	"fmt"

	xtime "go-common/library/time"
)

// CreditInfo credit info.
type CreditInfo struct {
	Mid        int64      `json:"mid"`
	Status     int64      `json:"status"`
	PunishType int64      `json:"blockedType"`
	PunishEnd  xtime.Time `json:"punishEnd"`
	CTime      xtime.Time `json:"-"`
	MTime      xtime.Time `json:"mtime"`
}

// BlockedInfo blocked case info.
type BlockedInfo struct {
	ID                  int64      `json:"id"`
	Uname               string     `json:"uname"`
	Face                string     `json:"face"`
	UID                 int64      `json:"uid"`
	OriginTitle         string     `json:"originTitle"`
	BlockedRemark       string     `json:"blockedRemark"`
	OriginURL           string     `json:"originUrl"`
	OriginContent       string     `json:"originContent,omitempty"`
	OriginContentModify string     `json:"originContentModify"`
	OriginType          int64      `json:"originType"`
	OriginTypeName      string     `json:"originTypeName"`
	PunishTitle         string     `json:"punishTitle"`
	PunishTime          xtime.Time `json:"punishTime"`
	PunishType          int64      `json:"punishType"`
	PunishTypeName      string     `json:"punishTypeName"`
	MoralNum            int64      `json:"moralNum"`
	BlockedDays         int64      `json:"blockedDays"`
	PublishStatus       int64      `json:"publishStatus"`
	BlockedType         int64      `json:"blockedType"`
	BlockedForever      int64      `json:"blockedForever"`
	ReasonType          int64      `json:"reasonType"`
	ReasonTypeName      string     `json:"reasonTypeName"`
	OperatorName        string     `json:"operatorName"`
	CaseID              int64      `json:"caseId"`
	PublishTime         xtime.Time `json:"-"`
	CTime               xtime.Time `json:"ctime"`
	MTime               xtime.Time `json:"-"`
	CommentSum          int64      `json:"commentSum"`
	OID                 int64      `json:"-"`
}

// BlockedPublish  blocked publish info.
type BlockedPublish struct {
	ID            int64      `json:"id"`
	Title         string     `json:"title"`
	SubTitle      string     `json:"subTitle"`
	PublishStatus int8       `json:"publishStatus"`
	StickStatus   int8       `json:"stickStatus"`
	Content       string     `json:"content"`
	CTime         xtime.Time `json:"ctime"`
	MTime         xtime.Time `json:"mtime"`
	URL           string     `json:"url"`
}

// BlockedCase  blocked case info.
type BlockedCase struct {
	ID            int64      `json:"id"`
	MID           int64      `json:"mid"`
	Status        int8       `json:"status"`
	StatusTitle   string     `json:"statusTitle,omitempty"`
	OriginType    int8       `json:"originType"`
	ReasonType    int8       `json:"reasonType"`
	OriginContent string     `json:"originContent"`
	PunishResult  int64      `json:"punishResult"`
	PunishTitle   string     `json:"punishTitle,omitempty"`
	JudgeType     int8       `json:"judgeType"`
	OriginURL     string     `json:"originUrl"`
	BlockedDays   int32      `json:"blockedDays"`
	PutTotal      int64      `json:"putTotal"`
	VoteRule      int64      `json:"voteRule"`
	VoteBreak     int64      `json:"voteBreak"`
	VoteDelete    int64      `json:"voteDelete"`
	StartTime     xtime.Time `json:"startTime"`
	EndTime       xtime.Time `json:"endTime"`
	Operator      string     `json:"-"`
	CTime         xtime.Time `json:"ctime"`
	MTime         xtime.Time `json:"mtime"`
	OriginTitle   string     `json:"originTitle"`
	RelationID    string     `json:"relationId"`
	Face          string     `json:"face"`
	Uname         string     `json:"uname"`
	Vote          int8       `json:"vote"`
	VoteTime      xtime.Time `json:"voteTime,omitempty"`
	ExpiredMillis int64      `json:"expiredMillis,omitempty"`
	CaseType      int8       `json:"case_type"`
}

// BlockedCaseVote blocked_case_vote.
type BlockedCaseVote struct {
	ID      int64      `json:"id"`
	CID     int64      `json:"cid"`
	MID     int64      `json:"mid"`
	Vote    int8       `json:"vote"`
	CTime   xtime.Time `json:"ctime"`
	MTime   xtime.Time `json:"mtime"`
	Expired xtime.Time `json:"expired"`
}

// BlockedJury blocked jury info.
type BlockedJury struct {
	ID            int64      `json:"id"`
	MID           int64      `json:"mid"`
	Status        int8       `json:"status"`
	Expired       xtime.Time `json:"expired"`
	InvalidReason int8       `json:"invalidReason"`
	VoteTotal     int64      `json:"voteTotal"`
	CaseTotal     int64      `json:"caseTotal"`
	VoteRadio     int64      `json:"voteRadio"`
	CTime         xtime.Time `json:"ctime"`
	MTime         xtime.Time `json:"mtime"`
	VoteRight     int64      `json:"voteRight"`
	Black         int8       `json:"black"`
}

// JuryRequirement jury requirement info.
type JuryRequirement struct {
	Blocked bool `json:"blocked"`
	Cert    bool `json:"cert"`
	Level   bool `json:"level"`
	Rule    bool `json:"rule"`
}

//UserInfo jury interface
type UserInfo struct {
	CaseTotal  int64  `json:"caseTotal"`
	Face       string `json:"face"`
	RestDays   int64  `json:"restDays"`
	RightRadio int64  `json:"rightRadio"`
	Status     int8   `json:"status"`
	Uname      string `json:"uname"`
}

// BlockedConfig config struct.
type BlockedConfig struct {
	CaseGiveHours  string `json:"caseGiveHours"`
	CaseCheckHours string `json:"caseCheckHours"`
	JuryVoteRadio  string `json:"juryVoteRadio"`
	CaseJudgeRadio string `json:"caseJudgeRadio"`
	CaseVoteMin    string `json:"caseVoteMin"`
}

// VoteInfo vote info.
type VoteInfo struct {
	ID      int64      `json:"id"`
	MID     int64      `json:"mid"`
	Vote    int8       `json:"vote"`
	CID     int64      `json:"cid"`
	Expired xtime.Time `json:"expired"`
	Mtime   xtime.Time `json:"mtime"`
}

// CaseInfo struct
type CaseInfo struct {
	EndTime       xtime.Time `json:"endTime"`
	Face          string     `json:"face"`
	ID            int64      `json:"id"`
	OriginContent string     `json:"originContent"`
	OriginTitle   string     `json:"originTitle"`
	OriginURL     string     `json:"originUrl"`
	Status        int8       `json:"status"`
	UID           int64      `json:"uid"`
	Uname         string     `json:"uname"`
	VoteBreak     int64      `json:"voteBreak"`
	VoteRule      int64      `json:"voteRule"`
}

// Build set blocked title info.
func (bi *BlockedInfo) Build() {
	bi.OriginTypeName = _originType[int8(bi.OriginType)]
	bi.ReasonTypeName = _reasonType[int8(bi.ReasonType)]
	bi.PunishTypeName = _punishType[int8(bi.PunishType)]
	bi.PunishTitle = fmt.Sprintf("在%s中%s", bi.OriginTypeName, bi.ReasonTypeName)
}

// Build set blocked title info.
func (bc *BlockedCase) Build() {
	bc.PunishTitle = fmt.Sprintf("在%s中%s", _originType[int8(bc.OriginType)], _reasonType[int8(bc.ReasonType)])
}

// CheckVote check vote.
func CheckVote(vote int8) (ok bool) {
	if vote == VoteBanned || vote == VoteRule || vote == VoteAbstain || vote == VoteDel {
		ok = true
	}
	return
}

// IsCaseTypePublic is case type public.
func IsCaseTypePublic(caseType int8) (ok bool) {
	if caseType == JudeCaseTypePublic {
		ok = true
	}
	return
}

// GantMedalID .
func (bj *BlockedJury) GantMedalID() int64 {
	switch bj.CaseTotal {
	case GuardMedalPointC:
		return GuardMedalC
	case GuardMedalPointB:
		return GuardMedalB
	case GuardMedalPointA:
		return GuardMedalA
	}
	return GuardMedalNone
}

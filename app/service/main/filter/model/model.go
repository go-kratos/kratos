package model

import (
	"regexp"

	"go-common/app/service/main/filter/model/actriearea"
	"go-common/app/service/main/filter/model/actriekey"
	xregexp "go-common/app/service/main/filter/service/regexp"
)

// 规则状态
const (
	RuleStateNormal  = 0
	RuleStateDeleted = 1
	RuleStateExpired = 2
)

// 过滤模式
const (
	RegMode = int8(0)
	StrMode = int8(1)
)

// Filter struct.
type Filter struct {
	Regs    []*Regexp
	Matcher *actriearea.Matcher
}

// Regexp struct .
type Regexp struct {
	TypeIDs []int64
	Level   int8
	Reg     *xregexp.Regexp
	Fid     int64
	Area    string
}

// KeyFilter .
type KeyFilter struct {
	Regs    []*KeyRegxp
	Matcher *actriekey.Matcher
}

// KeyRegxp struct .
type KeyRegxp struct {
	FkID  int64
	Area  string
	Level int8
	Reg   *regexp.Regexp
}

// FilterTestResult struct .
type FilterTestResult struct {
	Area string `json:"area"`
	*actriearea.Rule
}

package net

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//..
const (
	TableToken     = "net_token"
	TableTokenBind = "net_token_bind"
	//BindTypeFlow 绑定flow
	BindTypeFlow = int8(1)
	//BindTypeTransition 绑定单个审核操作
	BindTypeTransition = int8(2)
	//BindTypeTranBatch 绑定批量审核操作
	BindTypeTranBatch = int8(3)
	//BindTypeTranHelp 绑定单个非审核操作，提交后不流转
	BindTypeTranHelp = int8(4)
	//BindTypeTranHelpBatch 绑定批量非审核操作，提交后不流转
	BindTypeTranHelpBatch = int8(5)

	BindLogTemp = "%s(%d)=%s" //中文名=表达式

	BindBatchPrefix = "批量"
	BindBatchOptAll = int8(0)
	BindBatchOnly   = int8(1)
	BindBatchNot    = int8(2)
)

//TokenCompareAssign 运算符.
const (
	TokenCompareAssign       = int8(0)
	TokenCompareEqual        = int8(1)
	TokenCompareInequal      = int8(2)
	TokenCompareGreater      = int8(3)
	TokenCompareGreaterEqual = int8(4)
	TokenCompareLower        = int8(5)
	TokenCompareLowerEqual   = int8(6)
	TokenCompareBetween      = int8(7)
	TokenCompareNotBetween   = int8(8)
	TokenCompareIn           = int8(9)
	TokenCompareNotIn        = int8(10)
	TokenCompareExist        = int8(11)
	TokenCompareNotExist     = int8(12)

	//TokenTypeInt8值类型 .
	TokenTypeInt8   = int8(0)
	TokenTypeInt16  = int8(1)
	TokenTypeInt32  = int8(2)
	TokenTypeInt64  = int8(3)
	TokenTypeString = int8(4)
	TokenTypeBool   = int8(5)
)

//BindTranType 变迁绑定令牌类型
var BindTranType = []int8{
	BindTypeTransition,
	BindTypeTranBatch,
	BindTypeTranHelp,
	BindTypeTranHelpBatch,
}

//BindTypes .
var BindTypes = []int8{
	BindTypeFlow,
	BindTypeTransition,
	BindTypeTranBatch,
	BindTypeTranHelp,
	BindTypeTranHelpBatch,
}

// TokenCompareDesc .
var TokenCompareDesc = map[int8]string{
	TokenCompareAssign: "=",
	//TokenCompareEqual:        "==",
	//TokenCompareInequal:      "!=",
	//TokenCompareGreater:      ">",
	//TokenCompareGreaterEqual: ">=",
	//TokenCompareLower:        "<",
	//TokenCompareLowerEqual:   "<=",
	//TokenCompareBetween:      "between",     //值有2个数字
	//TokenCompareNotBetween:   "not between", //值有2个数字
	//TokenCompareIn:           "in",          //a in []
	//TokenCompareNotIn:        "not in",      //a not in []
}

// TokenValueTypeDesc ..
var TokenValueTypeDesc = map[int8]string{
	TokenTypeInt8:   "int8",
	TokenTypeInt16:  "int16",
	TokenTypeInt32:  "int32",
	TokenTypeInt64:  "int64",
	TokenTypeString: "string",
	TokenTypeBool:   "bool",
}

// Token .
type Token struct {
	ID      int64     `gorm:"primary_key" json:"id" form:"id" validate:"omitempty,gt=0"`
	NetID   int64     `gorm:"column:net_id" json:"net_id" form:"net_id" validate:"omitempty,gt=0"`
	ChName  string    `gorm:"column:ch_name" json:"ch_name" form:"ch_name" validate:"required,max=16"`
	Name    string    `gorm:"column:name" json:"name" form:"name" validate:"required,max=10"`
	Compare int8      `gorm:"column:compare" json:"compare" form:"compare" validate:"omitempty,min=0,max=10" default:"-1"`
	Value   string    `gorm:"column:value" json:"value" form:"value"`
	Type    int8      `gorm:"column:type" json:"type" form:"type" default:"-1"`
	UID     int64     `gorm:"column:uid" json:"uid"`
	Ctime   time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime   time.Time `gorm:"column:mtime" json:"mtime"`
}

// TokenBind .
type TokenBind struct {
	ID          int64     `gorm:"primary_key" json:"id"`
	Type        int8      `gorm:"column:type" json:"type"`
	ElementID   int64     `gorm:"column:element_id" json:"element_id"`
	TokenID     string    `gorm:"column:token_id" json:"token_id"`
	ChName      string    `gorm:"column:ch_name" json:"ch_name" `
	UID         int64     `gorm:"column:uid" json:"uid"`
	DisableTime time.Time `gorm:"column:disable_time" json:"disable_time"`
	Ctime       time.Time `gorm:"column:ctime" json:"ctime"`
	Mtime       time.Time `gorm:"column:mtime" json:"mtime"`
}

// TokenBindDetail .
type TokenBindDetail struct {
	ID          int64     `json:"id"`
	Type        int8      `json:"type"`
	ElementID   int64     `json:"element_id"`
	TokenID     string    `json:"token_id"`
	ChName      string    `json:"ch_name"`
	DisableTime time.Time `json:"disable_time"`
	Tokens      []*Token  `json:"tokens"`
}

// TableName .
func (tk *Token) TableName() string {
	return TableToken
}

// IsAssign .
func (tk *Token) IsAssign() bool {
	return tk.Compare == TokenCompareAssign
}

// FormatValue .
func (tk *Token) FormatValue() (vt interface{}, err error) {
	var (
		it int64
		ib bool
	)

	switch tk.Type {
	case TokenTypeInt8:
		if it, err = strconv.ParseInt(tk.Value, 10, 8); err != nil {
			return
		}
		vt = int8(it)
	case TokenTypeInt16:
		if it, err = strconv.ParseInt(tk.Value, 10, 16); err != nil {
			return
		}
		vt = int16(it)
	case TokenTypeInt32:
		if it, err = strconv.ParseInt(tk.Value, 10, 32); err != nil {
			return
		}
		vt = int32(it)
	case TokenTypeInt64:
		if it, err = strconv.ParseInt(tk.Value, 10, 8); err != nil {
			return
		}
		vt = it
	case TokenTypeString:
		vt = tk.Value
	case TokenTypeBool:
		if ib, err = strconv.ParseBool(tk.Value); err != nil {
			return
		}
		vt = ib
	default:
		err = fmt.Errorf("token not support tp(%d)!", tk.Type)
	}

	return
}

//FormatLog 日志形式
func (tk *Token) FormatLog() string {
	var (
		cm, tp string
		exist  bool
	)
	if cm, exist = TokenCompareDesc[tk.Compare]; !exist {
		cm = strconv.Itoa(int(tk.Compare))
	}
	if tp, exist = TokenValueTypeDesc[tk.Type]; !exist {
		tp = strconv.Itoa(int(tk.Type))
	}
	return fmt.Sprintf("%s%s%s(%s)", tk.Name, cm, tk.Value, tp)
}

// TableName .
func (tb *TokenBind) TableName() string {
	return TableTokenBind
}

// IsAvailable .
func (tb *TokenBind) IsAvailable() bool {
	return tb.DisableTime.IsZero()
}

// IsBatch .
func (tb *TokenBind) IsBatch() bool {
	return tb.Type == BindTypeTranBatch || tb.Type == BindTypeTranHelpBatch
}

//FormatLog 日志形式
func (tb *TokenBindDetail) FormatLog() string {
	var (
		logs = []string{}
	)
	for _, tk := range tb.Tokens {
		logs = append(logs, tk.FormatLog())
	}
	return fmt.Sprintf("{%s: %s}", tb.ChName, strings.Join(logs, ","))
}

// IsAvailable .
func (tb *TokenBindDetail) IsAvailable() bool {
	return tb.DisableTime.IsZero()
}

// GetTokenCompare .
func GetTokenCompare(compare int8) string {
	v, exist := TokenCompareDesc[compare]
	if !exist {
		return ""
	}

	return v
}

// GetTokenValueType .
func GetTokenValueType(tp int8) string {
	v, exist := TokenValueTypeDesc[tp]
	if !exist {
		return ""
	}

	return v
}

// Int64Slice data.Less and data.Swap
type Int64Slice []int64

func (sl Int64Slice) Len() int {
	return len(sl)
}

func (sl Int64Slice) Less(i, j int) bool {
	return sl[i] < sl[j]
}

func (sl Int64Slice) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

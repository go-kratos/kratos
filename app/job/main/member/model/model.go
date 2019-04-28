package model

import (
	"encoding/json"
	"fmt"
	"time"

	xtime "go-common/library/time"

	"github.com/satori/go.uuid"
)

const (
	// FlagNo 否
	FlagNo = int32(0)
	// FlagYes 是
	FlagYes = int32(1)
	// ExpFlagOnLogin 每日登录经验
	ExpFlagOnLogin = int32(1)
	// ExpFlagOnShare 每日分享经验
	ExpFlagOnShare = int32(2)
	// ExpFlagOnView 每日播放经验
	ExpFlagOnView = int32(4)
	// ExpFlagOnEmail 一次性绑定邮箱
	ExpFlagOnEmail = int32(8)
	// ExpFlagOnPhone 一次性绑定手机
	ExpFlagOnPhone = int32(16)
	// ExpFlagOnSafe 一次性绑定密保
	ExpFlagOnSafe = int32(32)
	// ExpFlagOnIdentify 一次性实名认证
	ExpFlagOnIdentify = int32(64)

	// ExpActOnCoin 投币奖励动作
	ExpActOnCoin = int64(1)
	// ExpActOnLogin 登录奖励动作
	ExpActOnLogin = int64(2)
	// ExpActOnView 播放奖励动作
	ExpActOnView = int64(3)
	// ExpActOnShare 分享奖励动作
	ExpActOnShare = int64(4)
	// ExpActOnEmail 绑定邮箱动作
	ExpActOnEmail = int64(5)
	// ExpActOnPhone 绑定手机动作
	ExpActOnPhone = int64(6)
	// ExpActOnSafe 绑定密保动作
	ExpActOnSafe = int64(7)
	// ExpActOnIdentify 实名认证动作
	ExpActOnIdentify = int64(8)
)

var (
	login    = &ExpOper{ExpFlagOnLogin, 5, "login", "登录奖励"}
	share    = &ExpOper{ExpFlagOnShare, 5, "shareClick", "分享视频奖励"}
	view     = &ExpOper{ExpFlagOnView, 5, "watch", "观看视频奖励"}
	email    = &ExpOper{ExpFlagOnEmail, 20, "bindEmail", "绑定邮箱奖励"}
	safe     = &ExpOper{ExpFlagOnSafe, 30, "pwdPro", "设置密保奖励"}
	phone    = &ExpOper{ExpFlagOnPhone, 100, "bindPhone", "绑定手机奖励"}
	identify = &ExpOper{ExpFlagOnIdentify, 50, "realIdentity", "实名认证奖励"}
	// ExpFlagOper exp flag map for oper
	ExpFlagOper = map[string]*ExpOper{"login": login, "share": share, "view": view, "email": email, "phone": phone, "safe": safe, "identify": identify}
)

// Exp userexp for mysql scan.
type Exp struct {
	Mid   int64      `json:"mid"`
	Exp   float32    `json:"exp"`
	Mtime xtime.Time `json:"modify_time"`
}

// ExpLog user exp log for mysql
type ExpLog struct {
	Mid      int64      `json:"mid"`
	FromExp  float32    `json:"from_exp"`
	ToExp    float32    `json:"to_exp"`
	Operater string     `json:"operater"`
	Reason   string     `json:"reason"`
	Action   int64      `json:"actin_id"`
	Mtime    xtime.Time `json:"modify_time"`
}

// NewExp userexp for mysql scan.
type NewExp struct {
	Mid     int64      `json:"mid"`
	Exp     int64      `json:"exp"`
	Flag    int32      `json:"flag"`
	Addtime xtime.Time `json:"addtime"`
	Mtime   xtime.Time `json:"mtime"`
}

// ExpOper exp operation
type ExpOper struct {
	Flag   int32
	Count  int64
	Oper   string
	Reason string
}

// Message binlog databus msg.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

//Flag is.
type Flag struct {
	Mid  int64 `json:"mid,omitempty"`
	Flag int64 `json:"finish_action"`
	Exp  int64 `json:"modify_exp"`
}

// UserLog user log.
type UserLog struct {
	Mid     int64             `json:"mid"`
	IP      string            `json:"ip"`
	TS      int64             `json:"ts"`
	LogID   string            `json:"log_id"`
	Content map[string]string `json:"content"`
}

// AddExp databus add exp arg.
type AddExp struct {
	Event string `json:"event,omitempty"`
	Mid   int64  `json:"mid,omitempty"`
	IP    string `json:"ip,omitempty"`
	Ts    int64  `json:"ts,omitempty"`
}

// LoginLog user login log.
type LoginLog struct {
	Mid       int64 `json:"mid,omitempty"`
	Loginip   int64 `json:"loginip,omitempty"`
	Timestamp int64 `json:"timestamp,omitempty"`
}

// LoginLogIPString user login log message with string ip.
type LoginLogIPString struct {
	Mid       int64  `json:"mid,omitempty"`
	Loginip   string `json:"loginip,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// FaceCheckRes face check result.
type FaceCheckRes struct {
	FileName string  `json:"file_name,omitempty"`
	Bucket   string  `json:"bucket,omitempty"`
	Sex      float64 `json:"sex,omitempty"`
	Politics float64 `json:"politics,omitempty"`
}

// FacePath is
func (fcr *FaceCheckRes) FacePath() string {
	return fmt.Sprintf("/bfs/%s/%s", fcr.Bucket, fcr.FileName)
}

//String is.
func (fcr *FaceCheckRes) String() string {
	return fmt.Sprintf("Sex: %.4f, Politics: %.4f", fcr.Sex, fcr.Politics)
}

// MemberMid member mid
type MemberMid struct {
	Mid int64 `json:"mid"`
}

// MemberAso member aso
type MemberAso struct {
	Email    string `json:"email"`
	Telphone string `json:"telphone"`
	SafeQs   int8   `json:"safe_question"`
	Spacesta int8   `json:"spacesta"`
}

// ExpMessage exp msg
type ExpMessage struct {
	Mid int64 `json:"mid"`
	Exp int64 `json:"exp"`
}

// FlagDailyReset reset daily flag with ts.
func (e *NewExp) FlagDailyReset(now time.Time) {
	e.Flag = e.Flag & ^0x7
	e.Addtime = RestrictDate(xtime.Time(now.Unix()))
}

// RestrictDate restric user brithday
func RestrictDate(xt xtime.Time) xtime.Time {
	t := xt.Time()
	year := t.Year()
	if year < 0 {
		year = 0
	}
	return xtime.Time(time.Date(year, t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix())
}

// UUID4 is generate uuid
func UUID4() string {
	return uuid.NewV4().String()
}

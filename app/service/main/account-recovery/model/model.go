package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	gotime "time"

	"go-common/app/service/main/account-recovery/conf"
	"go-common/library/log"
	"go-common/library/time"
)

// QueryInfoReq query mid info req
type QueryInfoReq struct {
	QType  string `form:"q_type" validate:"required"`
	QValue string `form:"q_value" validate:"required"`
	CToken string `form:"token" validate:"required"`
	Code   string `form:"code" validate:"required"`
}

// QueryInfoResp query info response
type QueryInfoResp struct {
	Status int64 `json:"status"`
	UID    int64 `json:"uid"`
}

// MIDInfo mid info data
type MIDInfo struct {
	Mids  string `json:"mids"`
	Count int64  `json:"count"`
}

// User user info data
type User struct {
	UserID string `json:"userid"`
	Pwd    string `json:"pwd"`
}

// Check check safe: Flag , reg: CheckInfo
type Check struct {
	CheckInfo string `json:"checkInfo"`
}

// UserInfoReq 用户申诉提交信息
type UserInfoReq struct {
	LoginAddrs   string    `json:"login_addrs" form:"login_addrs" validate:"required,min=1"`
	RegTime      time.Time `json:"reg_time" form:"reg_time" validate:"required"`
	RegType      int8      `json:"reg_type" form:"reg_type" validate:"required"`
	RegAddr      string    `json:"reg_addr" form:"reg_addr" validate:"required"`
	Unames       string    `json:"unames" form:"unames"`
	Pwds         string    `json:"pwds" form:"pwds" validate:"required"`
	Phones       string    `json:"phones" form:"phones"`
	Emails       string    `json:"emails" form:"emails"`
	SafeQuestion int8      `json:"safe_question" form:"safe_question" default:"99"`
	SafeAnswer   string    `json:"safe_answer" form:"safe_answer"`
	CardID       string    `json:"card_id" form:"card_id"`
	CardType     int8      `json:"card_type" form:"card_type" default:"99"`
	Captcha      string    `json:"captcha" form:"captcha" validate:"required"`
	LinkMail     string    `json:"link_mail" form:"link_mail" validate:"required,email"`
	Mid          int64     `json:"mid" form:"mid" validate:"required,min=1"`
	Business     string    `json:"business" form:"business" default:"account"`
	BusinessMap  map[string]string
	Files        []string  `json:"files" form:"files,split"`
	LastSucCount int64     //成功找回的总次数
	LastSucCTime time.Time //上次成功找回的提交时间
}

// DBAccountRecoveryAddit is
type DBAccountRecoveryAddit struct {
	Rid   int64     `json:"rid"`
	Files string    `json:"files"`
	Extra string    `json:"extra"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

// AsRecoveryAddit parse DBAccountRecoveryAddit to RecoveryAddit
func (dbAddit *DBAccountRecoveryAddit) AsRecoveryAddit() *RecoveryAddit {
	addit := &RecoveryAddit{
		Files: []string{},
		Extra: map[string]interface{}{},
	}
	if err := json.Unmarshal([]byte(dbAddit.Extra), &addit.Extra); err != nil {
		log.Error("QueryRecoveryAddit: json.Unmarshal(%s) error(%v)", addit.Extra, err)
	}
	addit.Files = strings.Split(dbAddit.Files, ",")
	for i, v := range addit.Files {
		addit.Files[i] = BuildFileURL(v)
	}

	return addit
}

// BuildFileURL  build bfs download url
func BuildFileURL(raw string) string {
	if raw == "" {
		return ""
	}
	ori, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if strings.HasPrefix(ori.Path, "/bfs/account") {
		//println("filename=====",bfsFilename(ori.Path, "account"))
		token := authorize(conf.Conf.Bfs.Key, conf.Conf.Bfs.Secret, "GET", "account", bfsFilename(ori.Path, "account"), gotime.Now().Unix())
		p := url.Values{}
		p.Set("token", token)
		ori.RawQuery = p.Encode()
	}
	if ori.Hostname() == "" {
		ori.Host = fmt.Sprintf("i%d.hdslb.com", rand.Int63n(3))
		ori.Scheme = "http"
	}
	return ori.String()
}

// authorize returns authorization for upload file to bfs
func authorize(key, secret, method, bucket, filename string, expire int64) string {
	content := fmt.Sprintf("%s\n%s\n%s\n%d\n", method, bucket, filename, expire)
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s:%s:%d", key, signature, expire)
}

func bfsFilename(path string, bucket string) string {
	return strings.TrimLeft(path, fmt.Sprintf("/bfs/%s", bucket))
}

// QueryRecoveryInfoReq 信息查询请求
type QueryRecoveryInfoReq struct {
	// FirstRid int64 `json:"first_rid" form:"first_rid"`
	// LastRid  int64 `json:"last_rid" form:"last_rid"`
	// RID account recovery info id
	RID int64 `json:"rid" form:"rid"`
	// UID user id
	UID int64 `json:"uid" form:"uid"`
	// Status account recovery status
	// status value, example: default "0", thought "1", reject "2"
	Status *int64 `json:"status" form:"-"`
	// Game is game user
	// game value, example: "1"/"0"
	Game *int64 `json:"game" form:"-"`
	// Size default size 10
	Size       int64     `json:"size" form:"size"`
	StartTime  time.Time `json:"start_time" form:"start_time"`
	EndTime    time.Time `json:"end_time" form:"end_time"`
	IsAdvanced bool      `json:"-"`
	// CurrPage   int64     `form:"curr_page"`
	Page      int64  `form:"page"`
	Bussiness string `json:"-"`
}

// DBRecoveryInfoParams DBRecoveryInfoParams
type DBRecoveryInfoParams struct {
	ExistGame   bool
	ExistStatus bool
	ExistMid    bool
	Mid         int64
	Game        int64
	Status      int64
	FirstRid    int64
	LastRid     int64
	Size        int64
	StartTime   time.Time
	EndTime     time.Time
	SubNum      int64
	CurrPage    int64
}

// AccountRecoveryInfo account recovery db info
type AccountRecoveryInfo struct {
	Rid             int64
	Mid             int64
	UserType        int64
	Status          int64
	LoginAddr       string
	UNames          string
	RegTime         time.Time `json:"-"`
	RegTimeStr      string    `json:"RegTime"`
	RegType         int64     `json:"-"`
	RegTypeStr      string    `json:"RegType"`
	RegAddr         string
	Pwd             string
	Phones          string
	Emails          string
	SafeQuestion    int64  `json:"-"`
	SafeQuestionStr string `json:"SafeQuestion"`
	SafeAnswer      string
	CardType        int64  `json:"-"`
	CardTypeStr     string `json:"CardType"`
	CardID          string
	SysLoginAddr    string
	SysReg          string
	SysUNames       string
	SysPwds         string
	SysPhones       string
	SysEmails       string
	SysSafe         string
	SysCard         string
	LinkEmail       string
	Operator        string
	OptTime         time.Time
	Remark          string
	CTime           time.Time
	MTime           time.Time
	Bussiness       string
	LastSucCount    int64
	LastSucCTime    time.Time
}

// RecoveryResInfo RecoveryResInfo
type RecoveryResInfo struct {
	AccountRecoveryInfo
	RecoverySuccess
	RecoveryAddit
	LastSuccessData
}

// RecoverySuccess recovery success info
type RecoverySuccess struct {
	SuccessMID       int64 `json:"-"`
	SuccessCount     int64
	FirstSuccessTime time.Time
	LastSuccessTime  time.Time
}

// LastSuccessData last recovery success info
type LastSuccessData struct {
	LastApplyMID  int64 `json:"-"`
	LastApplyTime time.Time
}

// RecoveryAddit common business field
type RecoveryAddit struct {
	Files []string
	Extra map[string]interface{}
}

// MultiQueryRes MultiQueryRes
type MultiQueryRes struct {
	Info []*RecoveryResInfo
	Page *Page
}

// UserInfoRes response of userInfo
type UserInfoRes struct {
	LoginAddrs string `json:"login_addrs"`
	Unames     string `json:"unames"`
	Pwds       string `json:"pwds"`
	Phones     string `json:"phones"`
	Emails     string `json:"emails"`

	RegInfo  *RegInfo  `json:"reg_info"`
	SafeInfo *SafeInfo `json:"safe_info"`
	CardInfo *CardInfo `json:"card_info"`
}

// RegInfo reg info
type RegInfo struct {
	RegTime time.Time `json:"reg_time"`
	RegType int64     `json:"reg_type"`
	RegAddr string    `json:"reg_addr"`
}

// SafeInfo safe info
type SafeInfo struct {
	SafeQuestion int8   `json:"safe_question"`
	SafeAnswer   string `json:"safe_answer"`
}

// CardInfo card info
type CardInfo struct {
	CardID   string `json:"card_id"`
	CardType int8   `json:"card_type"`
}

// SysInfo sys check info
type SysInfo struct {
	SysLoginAddrs string `json:"sys_login_addrs"`
	SysReg        string `json:"sys_reg"`
	SysUNames     string `json:"sys_nick"`
	SysPwds       string `json:"sys_pwds"`
	SysPhones     string `json:"sys_phones"`
	SysEmails     string `json:"sys_emails"`
	SysSafe       string `json:"sys_safe"`
	SysCard       string `json:"sys_card"`
}

// OperInfo operate info
type OperInfo struct {
	Operator string    `json:"operator"`
	OperTime time.Time `json:"oper_time"`
}

// AppealRes 信息查询返回结果
type AppealRes struct {
	Rid   int64     `json:"rid"`
	Mid   int64     `json:"uid"`
	Ctime time.Time `json:"ctime"`
	Count int64     `json:"count"`

	LinkEmail string `json:"link_email"`
	Remark    string `json:"remark"`
	Status    int64  `json:"status"`

	UserInfoRes *UserInfoRes `json:"user_info"`
	SysInfo     *SysInfo     `json:"sys_info"`
	OperInfo    *OperInfo    `json:"oper_info"`
}

// Page page
type Page struct {
	//Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

//JudgeReq appeal judge
type JudgeReq struct {
	Status   int64     `form:"status" validate:"required"`
	Rid      int64     `form:"rid" validate:"required"`
	Operator string    `form:"operator" validate:"required"`
	OptTime  time.Time `form:"opt_time" validate:"required"`
	Remark   string    `form:"remark"`
}

//BatchJudgeReq appeal judge
type BatchJudgeReq struct {
	Status   int64     `form:"status" validate:"required"`
	Rids     string    `form:"rids" validate:"required"`
	Operator string    `form:"operator" validate:"required"`
	OptTime  time.Time `form:"opt_time" validate:"required"`
	Remark   string    `form:"remark"`
	RidsAry  []int64
}

// CommonResq common response
type CommonResq struct {
	Code    int64  `json:"code"`
	TS      int64  `json:"ts"`
	Message string `json:"message"`
}

// Token captcha token
type Token struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

// TokenResq response of capthca
type TokenResq struct {
	CommonResq
	Data *Token `json:"data"`
}

// Game data of game
type Game struct {
	GameBaseID int    `json:"id"`
	GameName   string `json:"name"`
	LastLogin  string `json:"lastLogin"`
}

// GameListRes game list res
type GameListRes struct {
	Mid   int64   `json:"mid"`
	Items []*Game `json:"items"`
}

// UserBindLogRes UserBindLogRes
type UserBindLogRes struct {
	Page   Page           `json:"page"`
	Result []*UserBindLog `json:"result"`
}

// UserBindLog UserBindLog
type UserBindLog struct {
	Mid   int64  `json:"mid"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Time  string `json:"time"`
}

// BatchAppeal some of appeal info
type BatchAppeal struct {
	Rid      string
	Mid      string
	LinkMail string
	Ctime    time.Time
}

// NickNameReq  request es params
type NickNameReq struct {
	Mid  int64 `form:"mid"`
	Page int   `form:"page"`
	Size int   `form:"size"`
	From int64 `form:"from"`
	To   int64 `form:"to"`
}

// NickNameLogEs  query field form es
type NickNameLogEs struct {
	OldName string `json:"str_0"`
	NewName string `json:"str_1"`
}

// NickESRes the result of es query
type NickESRes struct {
	Page   Page             `json:"page"`
	Result []*NickNameLogEs `json:"result"`
}

// NickNameLogRes make NickESRes result become we need
type NickNameLogRes struct {
	Page   Page            `json:"page"`
	Result []*NickNameInfo `json:"result"`
}

// NickNameInfo nickName
type NickNameInfo struct {
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

// CheckEmailPhone check email and phone
type CheckEmailPhone struct {
	PhonesCheck string
	EmailCheck  string
}

//CaptchaMailReq get mail captcha
type CaptchaMailReq struct {
	Mid      int64  `form:"mid" validate:"required"`
	LinkMail string `form:"link_mail" validate:"required"`
}

// AddrInfo AddrInfo
type AddrInfo struct {
	OftenAddrs string `json:"oftenAddrs"`
}

// SendMailReq send mail according to status and rid
type SendMailReq struct {
	RID    int64 `form:"rid" validate:"required"`
	Status int64 `form:"status" validate:"required"`
}

// LoginIPInfo login ip info
type LoginIPInfo struct {
	LoginIP string `json:"loginip"`
}

// UserInfo user detail info
type UserInfo struct {
	Mid      int64     `json:"mid"`
	Phone    string    `json:"telphone"`
	Email    string    `json:"email"`
	JoinTime time.Time `json:"join_time"`
}

// UserBindLogReq UserBindLogReq
type UserBindLogReq struct {
	// Action value : telBindLog or emailBindLog
	Action string `form:"action"`
	Mid    int64  `form:"mid"`
	//Query search tel or email
	Query string `form:"query"`
	Page  int    `form:"page"`
	Size  int    `form:"size"`
	From  int64  `form:"from"`
	To    int64  `form:"to"`
}

// EsRes EsRes
type EsRes struct {
	Page   Page            `json:"page"`
	Result []*UserActLogEs `json:"result"`
}

// UserActLogEs UserActLogEs
type UserActLogEs struct {
	Mid       int64  `json:"mid"`
	Str0      string `json:"str_0"`
	ExtraData string `json:"extra_data"`
	CTime     string `json:"ctime"`
}

const (
	// VerifyMail verify mail
	VerifyMail = 1 //验证码邮件
	// CommitMail commit mail
	CommitMail = 2 //申诉信息提交邮件
	// RejectMail reject mail
	RejectMail = 3 //账号申诉驳回邮件
	// AgreeMail agree mail
	AgreeMail = 4 //账号审核通过邮件

	// DOAgree agree this appeal
	DOAgree = 1 //申诉通过
	// DOReject reject this appeal
	DOReject = 2 //申诉驳回

	// HIDEALL hide all
	HIDEALL = "*********"

	// BizGame game
	BizGame = "game"
)

// BusinessExtraArgs business args
func BusinessExtraArgs(business string) []string {
	switch business {
	case BizGame:
		return []string{"GamePlay", "GameArea", "GameRoleCTime", "GameRoleCreatePhoneType", "GameUsedPhoneType", "GameNames"}

	}
	return []string{}

}

package model

import (
	"regexp"
	"strings"
)

const (
	//UserTypeUp up主
	UserTypeUp = int8(1)
	//UserTypeBili b站用户
	UserTypeBili = int8(2)
	//UserTypeNew 新注册用户
	UserTypeNew = int8(3)
	//DegreeUncomp 未完成状态
	DegreeUncomp = int8(0)
	//DegreeComp 完成状态
	DegreeComp = int8(1)
	//SexMan 男
	SexMan = int8(1)
	//SexWoman 女
	SexWoman = int8(2)
	//SexAnimal 不明生物
	SexAnimal = int8(0)
)

// UserListType 用于指定列表类型
type UserListType int8

// UserListType的列表类型
const (
	FollowListType UserListType = 1
	FanListType    UserListType = 2
	BlackListType  UserListType = 4

	//ForbiddenStatus .
	ForbiddenStatus = 1
	//NormalStatus .
	NormalStatus = 0
)

const (
	// SpaceListLen 空间长度
	SpaceListLen = 20
	// BatchUserLen 批量请求用户信息时最大数量
	BatchUserLen = 50
	// MaxBlacklistLen 黑名单最大长度
	MaxBlacklistLen = 200
	// MaxFollowListLen 关注最大数
	MaxFollowListLen = 1000
)

// UserCard 主站返回的用户信息
type UserCard struct {
	MID     int64   `json:"mid"`
	Name    string  `json:"name"`
	Uname   string  `json:"uname"` // TODO: to delete
	Sex     string  `json:"sex"`
	Rank    int32   `json:"rank"`
	Face    string  `json:"face"`
	Sign    string  `json:"sign"`
	Level   int32   `json:"level"`
	VIPInfo VIPInfo `json:"vip_info"`
}

// UserInfoConfig 用于请求UserInfo的时候携带的参数
type UserInfoConfig struct {
	//needBase        bool // 必须基于UserBase信息
	NeedDesc        bool // 注意：desc和region_name一起，可能被降级，因为用户统计信息被认为是不重要信息
	NeedStatistic   bool // 注意：可能被降级，因为用户统计信息被认为是不重要信息
	NeedFollowState bool // 注意：可能被降级，因为关注关系信息被认为是不重要信息
}

//UpUserInfoRes account服务返回信息
type UpUserInfoRes struct {
	MID  int64  `json:"mid"`
	Name string `json:"name"`
	Sex  string `json:"sex"`
	Face string `json:"face"`
	Sign string `json:"sign"`
	Rank int64  `json:"rank"`
}

//VIPInfo .
type VIPInfo struct {
	Type    int32 `json:"type"`
	Status  int32 `json:"status"`
	DueDate int64 `json:"due_date"`
}

// CheckUnameSpecial 验证是否含有特殊字符
func CheckUnameSpecial(uname string) (matched bool) {
	matched, _ = regexp.MatchString("^[A-Za-z0-9\uAC00-\uD788\u3041-\u309E\u30A1-\u30FE\u3131-\u3163\u4E00-\u9FA5\uF92C-\uFA29_-]{1,}$", uname)
	return
}

//CheckUnameLength 验证长度
func CheckUnameLength(uname string) (matched bool) {
	lu := strings.Count(uname, "") - 1
	if lu < 3 || lu > 16 {
		return false
	}
	bt := []byte(uname)
	if len(bt) < 3 || len(bt) > 30 {
		return false
	}
	return true
}

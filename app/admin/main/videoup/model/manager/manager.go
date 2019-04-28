package manager

const (
	// UpperTypeWhite 白名单UP主
	UpperTypeWhite int8 = 1
	// UpperTypeBlack 黑名单UP主
	UpperTypeBlack int8 = 2
	// UpperTypePGC PGC UP主
	UpperTypePGC int8 = 3
	// UpperTypeUGCX UGC UP主
	UpperTypeUGCX int8 = 4
	// UpperTypePolity 时政 Up主
	UpperTypePolity int8 = 5
	// UpperTypeDanger 高危UP主
	UpperTypeDanger int8 = 6
	// UpperTypeTwoForbid 二禁UP主
	UpperTypeTwoForbid int8 = 10
	// UpperTypePGCWhite 视频自动锁定 PGC白名单
	UpperTypePGCWhite int8 = 11
	// ReasonLogTypeArc 审核理由类型：稿件
	ReasonLogTypeArc int8 = 1
	// ReasonLogTypeVideo 审核理由类型：视频
	ReasonLogTypeVideo int8 = 2
	// TaskLeader 组长
	TaskLeader int8 = 1
	// TaskMember 组员
	TaskMember int8 = 2
)

// User manager user struct
type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	State    int8   `json:"state"`
}

// UpGroupData uper group api return data
type UpGroupData struct {
	Code int        `json:"code"`
	Data []*UpGroup `json:"data"`
}

// UpGroup UP user group struct
type UpGroup struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	ShortTag string `json:"short_tag"`
	//Remark   string `json:"remark"`
	//State    int8   `json:"state"`
}

// UpGroup2 为了统一列表返回的up_group，
type UpGroup2 struct {
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
	GroupTag  string `json:"group_tag"`
}

package model

import (
	"strconv"

	v1 "go-common/app/service/main/account/api"
	mmodel "go-common/app/service/main/member/model"
)

// AccJavaInfo thin infomartion
type AccJavaInfo struct {
	Mid            int64 `json:"mid"`
	Scores         int32 `json:"scores"`
	JoinTime       int32 `json:"jointime"`
	Silence        int32 `json:"silence"`
	EmailStatus    int32 `json:"email_status"`
	TelStatus      int32 `json:"tel_status"`
	Identification int32 `json:"identification"`
	Moral          int32 `json:"moral"`
	Nameplate      struct {
		Nid        int    `json:"nid"`
		Name       string `json:"name"`
		Image      string `json:"image"`
		ImageSmall string `json:"image_small"`
		Level      string `json:"level"`
		Condition  string `json:"condition"`
	} `json:"nameplate"`
}

// OldInfo old info.
type OldInfo struct {
	Mid         string           `json:"mid"`
	Name        string           `json:"uname"`
	Sex         string           `json:"sex"`
	Sign        string           `json:"sign"`
	Avatar      string           `json:"avatar"`
	Rank        string           `json:"rank"`
	DisplayRank string           `json:"DisplayRank"`
	LevelInfo   mmodel.LevelInfo `json:"level_info"`
	Official    OldOfficial      `json:"official_verify"`
	Vip         v1.VipInfo       `json:"vip"`
}

// OldOfficial old official.
type OldOfficial struct {
	Type int8   `json:"type"`
	Desc string `json:"desc"`
}

// CvtOfficial is used to convert to old official.
func CvtOfficial(o v1.OfficialInfo) OldOfficial {
	old := OldOfficial{}
	if o.Role == 0 {
		old.Type = -1
	} else {
		if o.Role <= 2 {
			old.Type = 0
		} else {
			old.Type = 1
		}
		old.Desc = o.Title
	}
	return old
}

// Info old info -> info.
func (oi *OldInfo) Info() *v1.Info {
	mid, _ := strconv.ParseInt(oi.Mid, 10, 64)
	rank, _ := strconv.ParseInt(oi.Rank, 10, 64)
	i := &v1.Info{
		Mid:  mid,
		Name: oi.Name,
		Sex:  oi.Sex,
		Face: oi.Avatar,
		Sign: oi.Sign,
		Rank: int32(rank),
	}
	return i
}

// Relation relation.
type Relation struct {
	Following bool `json:"following"`
}

// ProfileStat profile with stat.
type ProfileStat struct {
	*v1.Profile
	LevelExp  mmodel.LevelInfo `json:"level_exp"`
	Coins     float64          `json:"coins"`
	Following int64            `json:"following"`
	Follower  int64            `json:"follower"`
}

// SearchMemberResult is.
type SearchMemberResult struct {
	Order  string `json:"order"`
	Sort   string `json:"sort"`
	Result []struct {
		Mid int64 `json:"mid"`
	} `json:"result"`
	Page Page `json:"page"`
}

// Privacy .
type Privacy struct {
	Realname     string `json:"realname"`
	IdentityCard string `json:"identity_card"`
	IdentitySex  string `json:"identity_sex"`
	Tel          string `json:"tel"`
	RegIP        string `json:"reg_ip"`
	RegTS        int64  `json:"reg_ts"`
	HandIMG      string `json:"hand_img"`
}

// Page page.
type Page struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// Mids is.
func (r *SearchMemberResult) Mids() []int64 {
	mids := make([]int64, 0, len(r.Result))
	for _, r := range r.Result {
		mids = append(mids, r.Mid)
	}
	return mids
}

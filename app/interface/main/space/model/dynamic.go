package model

import (
	"encoding/json"
	"strconv"

	v1 "go-common/app/service/main/archive/api"
)

// DyCard dynamic data.
type DyCard struct {
	Card json.RawMessage `json:"card"`
	Desc struct {
		UID         int64        `json:"uid"`
		Type        int          `json:"type"`
		ACL         int          `json:"acl"`
		Rid         int64        `json:"rid"`
		View        int32        `json:"view"`
		Repost      int32        `json:"repost"`
		Comment     int32        `json:"comment"`
		Like        int32        `json:"like"`
		IsLiked     int32        `json:"is_liked"`
		DynamicID   int64        `json:"dynamic_id"`
		CommentID   int64        `json:"comment_id"`
		Timestamp   int64        `json:"timestamp"`
		PreDyID     int64        `json:"pre_dy_id"`
		OrigDyID    int64        `json:"orig_dy_id"`
		OrigType    int          `json:"orig_type"`
		RType       int          `json:"r_type"`
		InnerID     int64        `json:"inner_id"`
		UserProfile *UserProfile `json:"user_profile,omitempty"`
	} `json:"desc"`
	Extension *DyExtension `json:"extension,omitempty"`
}

// DyResult dynamic result
type DyResult struct {
	Card json.RawMessage `json:"card"`
	Desc struct {
		UID         int64        `json:"uid"`
		Type        int          `json:"type"`
		ACL         int          `json:"acl"`
		Rid         int64        `json:"rid"`
		View        int32        `json:"view"`
		Repost      int32        `json:"repost"`
		Comment     int32        `json:"comment"`
		Like        int32        `json:"like"`
		IsLiked     int32        `json:"is_liked"`
		DynamicID   string       `json:"dynamic_id"`
		CommentID   int64        `json:"comment_id"`
		Timestamp   int64        `json:"timestamp"`
		PreDyID     string       `json:"pre_dy_id"`
		OrigDyID    string       `json:"orig_dy_id"`
		OrigType    int          `json:"orig_type"`
		RType       int          `json:"r_type"`
		InnerID     int64        `json:"inner_id"`
		UserProfile *UserProfile `json:"user_profile,omitempty"`
	} `json:"desc"`
	Extension *DyExtension `json:"extension,omitempty"`
}

// DyExtension .
type DyExtension struct {
	VoteCfg *struct {
		VoteID  int64  `json:"vote_id"`
		Desc    string `json:"desc"`
		JoinNum int64  `json:"join_num"`
	} `json:"vote_cfg,omitempty"`
}

// DyTotal out dynamic total.
type DyTotal struct {
	HasMore bool      `json:"has_more"`
	List    []*DyItem `json:"list"`
}

// DyItem out dynamic item.
type DyItem struct {
	Type    int        `json:"type"`
	Top     bool       `json:"top"`
	Card    *DyResult  `json:"card,omitempty"`
	Archive *VideoItem `json:"video,omitempty"`
	Fold    []*DyItem  `json:"fold,omitempty"`
	Ctime   int64      `json:"ctime"`
	Privacy bool       `json:"privacy"`
}

// UserProfile dynamic item user profile.
type UserProfile struct {
	Pendant struct {
		Pid    int64  `json:"pid"`
		Name   string `json:"name"`
		Image  string `json:"image"`
		Expire int64  `json:"expire"`
	} `json:"pendant,omitempty"`
	DecorateCard struct {
		Mid          int64  `json:"mid"`
		ID           int64  `json:"id"`
		CardURL      string `json:"card_url"`
		CardType     int    `json:"card_type"`
		Name         string `json:"name"`
		ExpireTime   int64  `json:"expire_time"`
		CardTypeName string `json:"card_type_name"`
		UID          int64  `json:"uid"`
	} `json:"decorate_card,omitempty"`
}

// FromCard format dynamic card.
func (d *DyResult) FromCard(c *DyCard) {
	d.Card = c.Card
	d.Desc.UID = c.Desc.UID
	d.Desc.Type = c.Desc.Type
	d.Desc.ACL = c.Desc.ACL
	d.Desc.Rid = c.Desc.Rid
	d.Desc.View = c.Desc.View
	d.Desc.Repost = c.Desc.Repost
	d.Desc.Comment = c.Desc.Comment
	d.Desc.Like = c.Desc.Like
	d.Desc.IsLiked = c.Desc.IsLiked
	d.Desc.DynamicID = strconv.FormatInt(c.Desc.DynamicID, 10)
	d.Desc.CommentID = c.Desc.CommentID
	d.Desc.Timestamp = c.Desc.Timestamp
	d.Desc.PreDyID = strconv.FormatInt(c.Desc.PreDyID, 10)
	d.Desc.OrigDyID = strconv.FormatInt(c.Desc.OrigDyID, 10)
	d.Desc.OrigType = c.Desc.OrigType
	d.Desc.RType = c.Desc.RType
	d.Desc.InnerID = c.Desc.InnerID
	if c.Extension != nil && c.Extension.VoteCfg != nil {
		d.Extension = c.Extension
	}
	if c.Desc.UserProfile != nil {
		d.Desc.UserProfile = c.Desc.UserProfile
	}
}

// DyList dynamic list.
type DyList struct {
	HasMore int       `json:"has_more"`
	Cards   []*DyCard `json:"cards"`
}

// DyActItem dynamic other action items.
type DyActItem struct {
	Aid        int64 `json:"aid"`
	Type       int   `json:"type"`
	ActionTime int64 `json:"action_time"`
	Privacy    bool  `json:"privacy"`
}

// VideoItem user action video item.
type VideoItem struct {
	Aid      int64  `json:"aid"`
	Pic      string `json:"pic"`
	Title    string `json:"title"`
	Duration int64  `json:"duration"`
	Author   struct {
		Mid  int64  `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	}
	Stat struct {
		View    int32 `json:"view"`
		Danmaku int32 `json:"danmaku"`
		Reply   int32 `json:"reply"`
		Fav     int32 `json:"favorite"`
		Coin    int32 `json:"coin"`
		Share   int32 `json:"share"`
		Like    int32 `json:"like"`
	}
	Rights     v1.Rights `json:"rights"`
	ActionTime int64     `json:"action_time"`
}

// DyListArg .
type DyListArg struct {
	Mid      int64 `form:"-"`
	Vmid     int64 `form:"mid" validate:"min=1"`
	DyID     int64 `form:"dy_id"`
	Qn       int   `form:"qn" default:"16" validate:"min=1"`
	Pn       int   `form:"pn" default:"1" validate:"min=1"`
	LastTime int64 `form:"last_time"`
}

// FromArchive from archive to video item.
func (v *VideoItem) FromArchive(arc *v1.Arc) {
	v.Aid = arc.Aid
	v.Pic = arc.Pic
	v.Title = arc.Title
	v.Duration = arc.Duration
	v.Author.Mid = arc.Author.Mid
	v.Author.Name = arc.Author.Name
	v.Author.Face = arc.Author.Face
	v.Stat.View = arc.Stat.View
	v.Stat.Danmaku = arc.Stat.Danmaku
	v.Stat.Reply = arc.Stat.Reply
	v.Stat.Fav = arc.Stat.Fav
	v.Stat.Share = arc.Stat.Share
	v.Stat.Like = arc.Stat.Like
	v.Rights = arc.Rights
}

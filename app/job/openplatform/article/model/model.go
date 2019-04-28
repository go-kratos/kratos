package model

import (
	"encoding/json"
	"strconv"

	artmdl "go-common/app/interface/openplatform/article/model"
)

const (
	// ActUpdate ...
	ActUpdate = "update"
	// ActInsert ...
	ActInsert = "insert"
	// ActDelete ...
	ActDelete = "delete"
)

// Message canal binlog message.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Article db struction.
type Article struct {
	ID              int64  `json:"article_id"`
	CTime           string `json:"ctime"`
	CategoryID      int64  `json:"category_id"`
	Title           string `json:"title"`
	Summary         string `json:"summary"`
	BannerURL       string `json:"banner_url"`
	TemplateID      int    `json:"template_id"`
	State           int    `json:"state"`
	Mid             int64  `json:"mid"`
	Reprint         int    `json:"reprint"`
	ImageURLs       string `json:"image_urls"`
	OriginImageURLs string `json:"origin_image_urls"`
	PublishTime     int    `json:"publish_time"`
	DeletedTime     int    `json:"deleted_time"`
	Attributes      int32  `json:"attributes,omitempty"`
	Reason          string `json:"reject_reason,omitempty"`
	Words           int64  `json:"words"`
	DynamicIntro    string `json:"dynamic_intro"`
	MediaID         int64  `json:"media_id"`
}

// SearchArticle .
type SearchArticle struct {
	Article
	Tags          string `json:"tags"`
	Content       string `json:"content"`
	StatsView     int64  `json:"stats_view"`
	StatsFavorite int64  `json:"stats_favorite"`
	StatsLikes    int64  `json:"stats_likes"`
	StatsDisLike  int64  `json:"stats_dislike"`
	StatsReply    int64  `json:"stats_reply"`
	StatsShare    int64  `json:"stats_share"`
	StatsCoin     int64  `json:"stats_coin"`
	Keywords      string `json:"keywords"`
}

// Author db struction.
type Author struct {
	ID         int64 `json:"id"`
	State      int   `json:"state"`
	Mid        int64 `json:"mid"`
	DailyLimit int   `json:"daily_limit"`
}

// Merge merges stat.
func Merge(last, m *artmdl.StatMsg) (changed [][2]int64) {
	if m.View != nil && *m.View >= 0 {
		*last.View += *m.View
		changed = append(changed, [2]int64{int64(artmdl.FieldView), *last.View})
	}
	if m.Like != nil {
		*last.Like = *m.Like
		changed = append(changed, [2]int64{int64(artmdl.FieldLike), *last.Like})
	}
	if m.Dislike != nil {
		*last.Dislike = *m.Dislike
	}
	if m.Share != nil && *m.Share >= 0 {
		*last.Share += *m.Share
	}
	if m.Favorite != nil && *m.Favorite >= 0 {
		*last.Favorite = *m.Favorite
		changed = append(changed, [2]int64{int64(artmdl.FieldFav), *last.Favorite})
	}
	if m.Reply != nil && *m.Reply >= 0 {
		*last.Reply = *m.Reply
		changed = append(changed, [2]int64{int64(artmdl.FieldReply), *last.Reply})
	}
	if m.Coin != nil && *m.Coin >= 0 {
		*last.Coin = *m.Coin
	}
	return
}

// ReadURLs returns article's read urls.
func ReadURLs(aid int64) []string {
	aidStr := strconv.FormatInt(aid, 10)
	return []string{
		"http://www.bilibili.com/read/cv/" + aidStr,
		"https://www.bilibili.com/read/cv/" + aidStr,
		"http://www.bilibili.com/read/app/" + aidStr,
		"https://www.bilibili.com/read/app/" + aidStr,
	}
}

// GameCacheRetry .
type GameCacheRetry struct {
	Action string `json:"action"`
	Aid    int64  `json:"aid"`
}

// FlowCacheRetry .
type FlowCacheRetry struct {
	Aid int64 `json:"aid"`
	Mid int64 `json:"mid"`
}

// DynamicCacheRetry .
type DynamicCacheRetry struct {
	Aid          int64
	Mid          int64
	Show         bool
	Comment      string
	Ts           int64
	DynamicIntro string
}

// LikeMsg msg
type LikeMsg struct {
	BusinessID    int64 `json:"business_id"`
	MessageID     int64 `json:"message_id"`
	LikesCount    int64 `json:"likes_count"`
	DislikesCount int64 `json:"dislikes_count"`
}

// DynamicMsg msg
type DynamicMsg struct {
	Card struct {
		Comment string `json:"comment"`
		Dynamic string `json:"dynamic"`
		OwnerID int64  `json:"owner_id"`
		Rid     int64  `json:"rid"`
		Show    int64  `json:"show"`
		Stype   int64  `json:"stype"`
		Ts      int64  `json:"ts"`
		Type    int64  `json:"type"`
	} `json:"card"`
}

// Setting the setting struct
type Setting struct {
	Recheck *Recheck
}

// Recheck setting struct
type Recheck struct {
	Day  int64 `json:"day"`
	View int64 `json:"view"`
}

// Read presents user reading duration struct
type Read struct {
	Buvid     string
	Aid       int64
	Mid       int64
	IP        string
	From      string
	StartTime int64
	EndTime   int64
}

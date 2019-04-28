package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	account "go-common/app/service/main/account/model"
	xtime "go-common/library/time"
)

// Const .
const (
	// State
	StateAutoLock    = -11
	StateLock        = -10
	StateReject      = -3
	StatePending     = -2
	StateOpen        = 0
	StateOpenPending = 2
	StateOpenReject  = 3
	StateAutoPass    = 4
	StateRePending   = 5 // 重复编辑待审
	StateReReject    = 6 // 重复编辑未通过
	StateRePass      = 7 // 重复编辑通过

	// groups for creation center.
	GroupAll      = 0 // except draft and deleted. 0 2 -2 3 -3 -10
	GroupPending  = 1 // -2 2
	GroupPassed   = 2 // 0
	GroupUnpassed = 3 // -10 -3 3

	NoLikeState  = 0
	LikeState    = 1
	DislikeState = 2

	// Templates
	TemplateText         = 1
	TemplateSingleImg    = 2
	TemplateMultiImg     = 3
	TemplateSingleBigImg = 4

	// Attributes
	//AttrBitNoDistribute 禁止分发(空间/分区/动态)
	AttrBitNoDistribute = uint(1)
	//AttrBitNoRegion 禁止在分区页显示
	AttrBitNoRegion = uint(2)
	//AttrBitNoRank 禁止排行
	AttrBitNoRank = uint(3)

	// Author

	// AuthorStatePass 过审
	AuthorStateReject  = -1
	AuthorStatePending = 0
	AuthorStatePass    = 1
	AuthorStateClose   = 2
	// 	AuthorStateIgnore = 3
)

var cleanURLRegexp = regexp.MustCompile(`^.+hdslb.com`)
var bfsRegexp = regexp.MustCompile(`^https?://.{1,6}\.hdslb+\.com/.+(?:jpg|gif|png|webp|jpeg)$`)

// Categories for sorting category.
type Categories []*Category

func (as Categories) Len() int { return len(as) }
func (as Categories) Less(i, j int) bool {
	return as[i].Position < as[j].Position
}
func (as Categories) Swap(i, j int) { as[i], as[j] = as[j], as[i] }

// StatMsg means article's stat message in databus.
type StatMsg struct {
	View      *int64     `json:"view"`
	Like      *int64     `json:"like"`
	Dislike   *int64     `json:"dislike"`
	Favorite  *int64     `json:"fav"`
	Reply     *int64     `json:"reply"`
	Share     *int64     `json:"share"`
	Coin      *int64     `json:"coin"`
	Aid       int64      `json:"aid"`
	Mid       int64      `json:"mid"`
	IP        string     `json:"ip"`
	CheatInfo *CheatInfo `json:"cheat_info"`
}

func (sm *StatMsg) String() (res string) {
	if sm == nil {
		res = "<nil>"
		return
	}
	res = fmt.Sprintf("aid: %v, mid: %v, ip: %v, view(%s) likes(%s) dislike(%s) favorite(%s) reply(%s) share(%s) coin(%s)", sm.Aid, sm.Mid, sm.IP, formatPInt(sm.View), formatPInt(sm.Like), formatPInt(sm.Dislike), formatPInt(sm.Favorite), formatPInt(sm.Reply), formatPInt(sm.Share), formatPInt(sm.Coin))
	return
}

// CheatInfo .
type CheatInfo struct {
	Valid    string `json:"valid"`
	Client   string `json:"client"`
	Cvid     string `json:"cvid"`
	Mid      string `json:"mid"`
	Lv       string `json:"lv"`
	Ts       string `json:"ts"`
	IP       string `json:"ip"`
	UA       string `json:"ua"`
	Refer    string `json:"refer"`
	Sid      string `json:"sid"`
	Buvid    string `json:"buvid"`
	DeviceID string `json:"device_id"`
	Build    string `json:"build"`
	Reason   string `json:"reason"`
}

func formatPInt(s *int64) (res string) {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%d", *s)
}

// DraftMsg means article's draft message in databus.
type DraftMsg struct {
	Aid int64 `json:"aid"`
	Mid int64 `json:"mid"`
}

// Draft draft struct.
type Draft struct {
	*Article
	Tags   []string `json:"tags"`
	ListID int64    `json:"list_id"`
	List   *List    `json:"list"`
}

// Metas Metas
type Metas []*Meta

func (as Metas) Len() int { return len(as) }
func (as Metas) Less(i, j int) bool {
	var it, jt xtime.Time
	if as[i] != nil {
		it = as[i].PublishTime
	}
	if as[j] != nil {
		jt = as[j].PublishTime
	}
	return it > jt
}
func (as Metas) Swap(i, j int) { as[i], as[j] = as[j], as[i] }

// CreationArtsType creation article-list type's count.
type CreationArtsType struct {
	All       int `json:"all"`
	Audit     int `json:"audit"`
	Passed    int `json:"passed"`
	NotPassed int `json:"not_passed"`
}

// ArtPage article page.
type ArtPage struct {
	Pn    int `json:"pn"`
	Ps    int `json:"ps"`
	Total int `json:"total"`
}

// CreationArts creation article list.
type CreationArts struct {
	Articles []*Meta           `json:"articles"`
	Type     *CreationArtsType `json:"type"`
	Page     *ArtPage          `json:"page"`
}

// Drafts draft list.
type Drafts struct {
	Drafts []*Draft `json:"drafts"`
	Page   *ArtPage `json:"page"`
}

// UpArtMetas article list.
type UpArtMetas struct {
	Articles []*Meta `json:"articles"`
	Pn       int     `json:"pn"`
	Ps       int     `json:"ps"`
	Count    int     `json:"count"`
}

// UpArtMetasLists .
type UpArtMetasLists struct {
	*UpArtMetas
	UpLists UpLists `json:"up_lists"`
}

// IsNormal judge whether article's state is normal.
func (a *Meta) IsNormal() bool {
	return (a != nil) && (a.State >= StateOpen)
}

// IsNormal judge article state.
func (a *Article) IsNormal() bool {
	if (a == nil) || (a.Meta == nil) {
		return false
	}
	return a.Meta.IsNormal()
}

// AttrVal gets attr val by bit.
func (a *Meta) AttrVal(bit uint) bool {
	return ((a.Attributes>>bit)&int32(1) == 1)
}

// AttrSet sets attr value by bit.
func (a *Meta) AttrSet(v int32, bit uint) {
	a.Attributes = a.Attributes&(^(1 << bit)) | (v << bit)
}

// Strong fill blank images and tags
func (a *Meta) Strong() *Meta {
	if a.ImageURLs == nil {
		a.ImageURLs = []string{}
	}
	if a.OriginImageURLs == nil {
		a.OriginImageURLs = []string{}
	}
	if a.Tags == nil {
		a.Tags = []*Tag{}
	}
	return a
}

// AuthorPermission recode of article_authors table.
type AuthorPermission struct {
	State int        `json:"state"`
	Rtime xtime.Time `json:"rtime"`
}

// Favorite user favorite list.
type Favorite struct {
	*Meta
	FavoriteTime int64 `json:"favorite_time"`
	Valid        bool  `json:"valid"`
}

// Page model
type Page struct {
	Pn    int `json:"pn"`
	Ps    int `json:"ps"`
	Total int `json:"total"`
}

// RecommendArt model
type RecommendArt struct {
	Meta
	Recommend
}

// RecommendArtWithLike model
type RecommendArtWithLike struct {
	RecommendArt
	LikeState int `json:"like_state"`
}

// MetaWithLike meta with like
type MetaWithLike struct {
	Meta
	LikeState int `json:"like_state"`
}

// Recommend model
type Recommend struct {
	ArticleID         int64  `json:"article_id,omitempty"`
	Position          int    `json:"-"`
	EndTime           int64  `json:"-"`
	Rec               bool   `json:"rec"`
	RecFlag           bool   `json:"rec_flag"`
	RecText           string `json:"rec_text"`
	RecImageURL       string `json:"rec_image_url"`
	RecImageStartTime int64  `json:"-"`
	RecImageEndTime   int64  `json:"-"`
}

// ViewInfo model
type ViewInfo struct {
	Like            int8     `json:"like"`
	Attention       bool     `json:"attention"`
	Favorite        bool     `json:"favorite"`
	Coin            int64    `json:"coin"`
	Stats           Stats    `json:"stats"`
	Title           string   `json:"title"`
	BannerURL       string   `json:"banner_url"`
	Mid             int64    `json:"mid"`
	AuthorName      string   `json:"author_name"`
	IsAuthor        bool     `json:"is_author"`
	ImageURLs       []string `json:"image_urls"`
	OriginImageURLs []string `json:"origin_image_urls"`
	Shareable       bool     `json:"shareable"`
	ShowLaterWatch  bool     `json:"show_later_watch"`
	ShowSmallWindow bool     `json:"show_small_window"`
	InList          bool     `json:"in_list"`
	Pre             int64    `json:"pre"`
	Next            int64    `json:"next"`
}

// Group2State mapping creation group to
func Group2State(group int) (states []int64) {
	switch group {
	case GroupPassed:
		states = []int64{StateOpen, StateAutoPass, StateRePass, StateReReject}
	case GroupPending:
		states = []int64{StatePending, StateOpenPending, StateRePending}
	case GroupUnpassed:
		states = []int64{StateReject, StateOpenReject, StateLock, StateAutoLock}
	case GroupAll:
		fallthrough
	default:
		states = []int64{StateOpen, StatePending, StateOpenPending, StateReject, StateOpenReject, StateLock, StateAutoPass, StateAutoLock, StateRePending, StateRePass, StateReReject}
	}
	return
}

// CompleteURL adds host on path.
func CompleteURL(path string) (url string) {
	if path == "" {
		// url = "http://static.hdslb.com/images/transparent.gif"
		return
	}
	url = path
	if strings.Index(path, "//") == 0 || strings.Index(path, "http://") == 0 || strings.Index(path, "https://") == 0 {
		return
	}
	url = "https://i0.hdslb.com" + url
	return
}

// CleanURL cuts host.
func CleanURL(url string) (path string) {
	path = string(cleanURLRegexp.ReplaceAll([]byte(url), nil))
	return
}

// CompleteURLs .
func CompleteURLs(paths []string) (urls []string) {
	for _, v := range paths {
		urls = append(urls, CompleteURL(v))
	}
	return
}

// CleanURLs .
func CleanURLs(urls []string) (paths []string) {
	for _, v := range urls {
		paths = append(paths, CleanURL(v))
	}
	return
}

// Recommends model
type Recommends [][]*Recommend

func (as Recommends) Len() int { return len(as) }
func (as Recommends) Less(i, j int) bool {
	return as[i][0].Position > as[j][0].Position
}
func (as Recommends) Swap(i, j int) { as[i], as[j] = as[j], as[i] }

// RecommendHome .
type RecommendHome struct {
	RecommendPlus
	Categories []*Category `json:"categories"`
	IP         string      `json:"ip"`
}

// RecommendPlus .
type RecommendPlus struct {
	Banners  []*Banner               `json:"banners"`
	Articles []*RecommendArtWithLike `json:"articles"`
	Ranks    []*RankMeta             `json:"ranks"`
	Hotspots []*Hotspot              `json:"hotspots"`
}

// Banner struct
type Banner struct {
	ID         int    `json:"id"`
	Plat       int8   `json:"-"`
	Position   int    `json:"index"`
	Title      string `json:"title"`
	Image      string `json:"image"`
	URL        string `json:"url"`
	Build      int    `json:"-"`
	Condition  string `json:"-"`
	Rule       string `json:"-"`
	ResID      int    `json:"resource_id"`
	ServerType int    `json:"server_type"`
	CmMark     int    `json:"cm_mark"`
	IsAd       bool   `json:"is_ad"`
	RequestID  string `json:"request_id"`
}

// ConvertPlat convert plat from resource
func ConvertPlat(p int8) (plat int8) {
	switch p {
	case 0: // resource iphone
		plat = PlatPC
	case 1: // resource iphone
		plat = PlatIPhone
	case 2: // resource android
		plat = PlatAndroid
	case 3: // resource pad
		plat = PlatIPad
	case 4: // resource iphoneg
		plat = PlatIPhoneI
	case 5: // resource androidg
		plat = PlatAndroidG
	case 6: // resource padg
		plat = PlatIPadI
	case 7: // resource h5
		plat = PlatH5
	case 8: // resource androidi
		plat = PlatAndroidI
	}
	return
}

// BannerRule .
type BannerRule struct {
	Area      string `json:"area"`
	Hash      string `json:"hash"`
	Build     int    `json:"build"`
	Condition string `json:"conditions"`
	Channel   string `json:"channel"`
}

// NoDistributeAttr check if no distribute
func NoDistributeAttr(attr int32) bool {
	meta := Meta{Attributes: attr}
	return meta.AttrVal(AttrBitNoDistribute)
}

// NoRegionAttr check if no region
func NoRegionAttr(attr int32) bool {
	meta := Meta{Attributes: attr}
	return meta.AttrVal(AttrBitNoRegion)
}

// AuthorLimit .
type AuthorLimit struct {
	Limit int        `json:"limit"`
	State int        `json:"state"`
	Rtime xtime.Time `json:"rtime"`
}

// Forbid .
func (a *AuthorLimit) Forbid() bool {
	// state -1 未通过 0待审 1通过 2关闭 3忽略
	if a == nil {
		return false
	}
	if a.State == AuthorStatePass {
		return false
	}
	return true
}

// Pass .
func (a *AuthorLimit) Pass() bool {
	return (a != nil) && (a.State == AuthorStatePass)
}

// Notice notice
type Notice struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	URL       string `json:"url"`
	Plat      int    `json:"-"`
	Condition int    `json:"-"`
	Build     int    `json:"-"`
}

// MoreArts .
type MoreArts struct {
	Articles  []*Meta      `json:"articles"`
	Total     int          `json:"total"`
	ReadCount int64        `json:"read_count"`
	Author    *AccountCard `json:"author"`
	Attention bool         `json:"attention"`
}

// AccountCard .
type AccountCard struct {
	Mid         string  `json:"mid"`
	Name        string  `json:"name"`
	Approve     bool    `json:"approve"`
	Sex         string  `json:"sex"`
	Rank        string  `json:"rank"`
	Face        string  `json:"face"`
	DisplayRank string  `json:"DisplayRank"`
	Regtime     int64   `json:"regtime"`
	Spacesta    int     `json:"spacesta"`
	Birthday    string  `json:"birthday"`
	Place       string  `json:"place"`
	Description string  `json:"description"`
	Article     int     `json:"article"`
	Attentions  []int64 `json:"attentions"`
	Fans        int     `json:"fans"`
	Friend      int     `json:"friend"`
	Attention   int     `json:"attention"`
	Sign        string  `json:"sign"`
	LevelInfo   struct {
		Cur     int         `json:"current_level"`
		Min     int         `json:"current_min"`
		NowExp  int         `json:"current_exp"`
		NextExp interface{} `json:"next_exp"`
	} `json:"level_info"`
	Pendant struct {
		Pid    int    `json:"pid"`
		Name   string `json:"name"`
		Image  string `json:"image"`
		Expire int    `json:"expire"`
	} `json:"pendant"`
	Nameplate struct {
		Nid        int    `json:"nid"`
		Name       string `json:"name"`
		Image      string `json:"image"`
		ImageSmall string `json:"image_small"`
		Level      string `json:"level"`
		Condition  string `json:"condition"`
	} `json:"nameplate"`
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Vip struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
}

// FromCard from account card.
func (ac *AccountCard) FromCard(card *account.Card) {
	ac.Mid = strconv.FormatInt(card.Mid, 10)
	ac.Name = card.Name
	// ac.Approve =
	ac.Sex = card.Sex
	ac.Rank = strconv.FormatInt(int64(card.Rank), 10)
	ac.Face = card.Face
	ac.DisplayRank = strconv.FormatInt(int64(card.Rank), 10)
	// ac.Regtime =
	// ac.Spacesta =
	// ac.Birthday =
	// ac.Place =
	// ac.Description =
	// ac.Article =
	ac.Attentions = []int64{}
	// ac.Fans =
	// ac.Friend =
	// ac.Attention =
	ac.Sign = card.Sign
	ac.LevelInfo.Cur = int(card.Level)
	ac.Pendant.Pid = card.Pendant.Pid
	ac.Pendant.Name = card.Pendant.Name
	ac.Pendant.Image = card.Pendant.Image
	ac.Pendant.Expire = card.Pendant.Expire
	ac.Nameplate.Nid = card.Nameplate.Nid
	ac.Nameplate.Name = card.Nameplate.Name
	ac.Nameplate.Image = card.Nameplate.Image
	ac.Nameplate.ImageSmall = card.Nameplate.ImageSmall
	ac.Nameplate.Level = card.Nameplate.Level
	ac.Nameplate.Condition = card.Nameplate.Condition
	if card.Official.Role == 0 {
		ac.OfficialVerify.Type = -1
	} else {
		if card.Official.Role <= 2 {
			ac.OfficialVerify.Type = 0
		} else {
			ac.OfficialVerify.Type = 1
		}
		ac.OfficialVerify.Desc = card.Official.Title
	}
	ac.Vip.Type = int(card.Vip.Type)
	ac.Vip.VipStatus = int(card.Vip.Status)
	ac.Vip.DueDate = card.Vip.DueDate
}

// FromProfileStat .
func (ac *AccountCard) FromProfileStat(card *account.ProfileStat) {
	ac.Mid = strconv.FormatInt(card.Mid, 10)
	ac.Name = card.Name
	// ac.Approve =
	ac.Sex = card.Sex
	ac.Rank = strconv.FormatInt(int64(card.Rank), 10)
	ac.Face = card.Face
	ac.DisplayRank = strconv.FormatInt(int64(card.Rank), 10)
	// ac.Regtime =
	// ac.Spacesta =
	// ac.Birthday =
	// ac.Place =
	// ac.Description =
	// ac.Article =
	ac.Attentions = []int64{}
	ac.Fans = int(card.Follower)
	// ac.Friend =
	// ac.Attention =
	ac.Sign = card.Sign
	ac.LevelInfo.Cur = int(card.Level)
	ac.Pendant.Pid = card.Pendant.Pid
	ac.Pendant.Name = card.Pendant.Name
	ac.Pendant.Image = card.Pendant.Image
	ac.Pendant.Expire = card.Pendant.Expire
	ac.Nameplate.Nid = card.Nameplate.Nid
	ac.Nameplate.Name = card.Nameplate.Name
	ac.Nameplate.Image = card.Nameplate.Image
	ac.Nameplate.ImageSmall = card.Nameplate.ImageSmall
	ac.Nameplate.Level = card.Nameplate.Level
	ac.Nameplate.Condition = card.Nameplate.Condition
	if card.Official.Role == 0 {
		ac.OfficialVerify.Type = -1
	} else {
		if card.Official.Role <= 2 {
			ac.OfficialVerify.Type = 0
		} else {
			ac.OfficialVerify.Type = 1
		}
		ac.OfficialVerify.Desc = card.Official.Title
	}
	ac.Vip.Type = int(card.Vip.Type)
	ac.Vip.VipStatus = int(card.Vip.Status)
	ac.Vip.DueDate = card.Vip.DueDate
}

// NoticeState .
type NoticeState map[string]bool

// 数据库为tinyint 长度必须小于7 字段只能追加
var _noticeStates = []string{"lead", "new"}

// NewNoticeState .
func NewNoticeState(value int64) (res NoticeState) {
	res = make(map[string]bool)
	for i, name := range _noticeStates {
		res[name] = ((value>>uint(i))&int64(1) == 1)
	}
	return
}

// ToInt64 .
func (n NoticeState) ToInt64() (res int64) {
	for i, name := range _noticeStates {
		if n[name] {
			res = res | (1 << uint(i))
		}
	}
	return
}

// Activity .
type Activity struct {
	ActURL  string `json:"act_url"`
	Author  string `json:"author"`
	Cover   string `json:"cover"`
	Ctime   string `json:"ctime"`
	Dic     string `json:"dic"`
	Etime   string `json:"etime"`
	Flag    string `json:"flag"`
	H5Cover string `json:"h5_cover"`
	ID      int64  `json:"id"`
	Letime  string `json:"letime"`
	Level   string `json:"level"`
	Lstime  string `json:"lstime"`
	Mtime   string `json:"mtime"`
	Name    string `json:"name"`
	Oid     int64  `json:"oid"`
	State   int64  `json:"state"`
	Stime   string `json:"stime"`
	Tags    string `json:"tags"`
	Type    int64  `json:"type"`
	Uetime  string `json:"uetime"`
	Ustime  string `json:"ustime"`
}

// SkyHorseResp response
type SkyHorseResp struct {
	Code int `json:"code"`
	Data []struct {
		ID        int64  `json:"id"`
		AvFeature string `json:"av_feature"`
	} `json:"data"`
	UserFeature string `json:"user_feature"`
}

// CheckBFSImage check bfs file
func CheckBFSImage(src string) bool {
	return bfsRegexp.MatchString(src)
}

// FillDefaultImage .
func (l *List) FillDefaultImage(image string) {
	if l != nil && l.ImageURL == "" {
		l.ImageURL = image
	}
}

// Articles .
type Articles struct {
	*Article
	Pre  int64 `json:"pre"`
	Next int64 `json:"next"`
}

// ArticleViewList .
type ArticleViewList struct {
	Position int     `json:"position"`
	Aids     []int64 `json:"articles_id"`
	From     string  `json:"from"`
	Mid      int64   `json:"mid"`
	Build    int     `json:"build"`
	Buvid    string  `json:"buvid"`
	Plat     int8    `json:"plat"`
}

// TagArts .
type TagArts struct {
	Tid  int64   `json:"tid"`
	Aids []int64 `json:"aids"`
}

// MediaResp .
type MediaResp struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Result  *MediaResult `json:"result"`
}

// MediaResult .
type MediaResult struct {
	Score int32 `json:"score"`
	Media struct {
		MediaID  int64  `json:"media_id"`
		Title    string `json:"title"`
		Cover    string `json:"cover"`
		Area     string `json:"area"`
		TypeID   int32  `json:"type_id"`
		TypeName string `json:"type_name"`
	} `json:"media"`
}

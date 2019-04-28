package web

import (
	"fmt"
	"math/rand"
	"time"

	"go-common/app/interface/main/web-goblin/conf"
	"go-common/app/service/main/archive/api"
)

const (
	_deal = 100
)

// Mi mi common .
type Mi struct {
	ID               int64       `json:"id"`
	Name             string      `json:"name"`
	Op               string      `json:"op"`
	AlternativeNames string      `json:"alternative_names"`
	Cover            Img         `json:"cover"`
	Thumbnail        Img         `json:"thumbnail"`
	Description      string      `json:"description"`
	Tags             string      `json:"tags"`
	CreateTime       string      `json:"create_time"`
	ModifyTime       string      `json:"modify_time"`
	PublishTime      string      `json:"publish_time"`
	Author           string      `json:"author"`
	Category         string      `json:"category"`
	Rating           float32     `json:"rating"`
	PlayCount        int32       `json:"play_count"`
	PlayCountMonth   int         `json:"play_count_month"`
	PlayCountWeek    int         `json:"play_count_week"`
	PlayLength       int64       `json:"play_length"`
	Language         string      `json:"language"`
	Images           Img         `json:"images"`
	Weburl           string      `json:"weburl"`
	Appurl           string      `json:"appurl"`
	MinVersion       int         `json:"min_version"`
	Pages            []*PageInfo `json:"pages"`
	CommentCount     int32       `json:"comment_count"`
	LikeCount        int32       `json:"like_count"`
}

// PageInfo page3 .
type PageInfo struct {
	Cid  int64 `json:"cid"`
	Page int32 `json:"page"`
}

// Img img .
type Img struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// SearchAids return aids .
type SearchAids struct {
	Aid    int64  `json:"aid"`
	Action string `json:"action"`
}

// FromArchive .
func (f *Mi) FromArchive(a *api.Arc, p []*api.Page, op, source string) {
	f.ID = a.Aid
	f.Name = a.Title
	f.Author = a.Author.Name
	f.Category = a.TypeName
	f.Weburl = fmt.Sprintf("https://www.bilibili.com/video/av%d%s", a.Aid, source)
	f.Appurl = fmt.Sprintf("bilibili://video/%d", a.Aid)
	f.ModifyTime = a.PubDate.Time().Format("2006-01-02 15:04:05")
	f.Description = a.Desc
	f.CreateTime = a.Ctime.Time().Format("2006-01-02 15:04:05")
	f.Images.URL = a.Pic
	f.PlayLength = a.Duration
	f.Cover.URL = a.Pic
	f.Thumbnail.URL = a.Pic
	f.PlayCount = a.Stat.View
	f.PublishTime = a.PubDate.Time().Format("2006-01-02 15:04:05")
	f.MinVersion = 1
	f.Op = op
	f.CommentCount = a.Stat.Reply
	f.LikeCount = a.Stat.Like
	pLen := len(p)
	if pLen > 0 {
		f.Pages = make([]*PageInfo, pLen)
		for i := 0; i < pLen; i++ {
			f.Pages[i] = &PageInfo{}
			f.Pages[i].Page = p[i].Page
			f.Pages[i].Cid = p[i].Cid
		}
	} else {
		f.Pages = []*PageInfo{}
	}
}

// UgcFullDeal .
func (f *Mi) UgcFullDeal() {
	var (
		lCount    = conf.Conf.OutSearch.DealLikeFull
		commCount = conf.Conf.OutSearch.DealCommFull
		commRes   = f.CommentCount + commCount
		likeRes   = f.LikeCount + lCount
	)
	if f.PlayCount+f.CommentCount > 0 {
		commRes = commRes + (f.PlayCount/(f.PlayCount+f.CommentCount))*f.CommentCount
	}
	if f.PlayCount+f.LikeCount > 0 {
		likeRes = likeRes + (f.PlayCount/(f.PlayCount+f.LikeCount))*f.LikeCount
	}
	f.CommentCount = commRes
	f.LikeCount = likeRes
}

// UgcIncreDeal .
func (f *Mi) UgcIncreDeal() {
	rand.Seed(time.Now().UnixNano())
	f.CommentCount = int32(rand.Intn(_deal))
	f.LikeCount = int32(rand.Intn(_deal)) + _deal
}

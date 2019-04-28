package favorite

import (
	"strconv"

	"go-common/app/interface/main/app-interface/model"
	"go-common/app/interface/main/app-interface/model/audio"
	"go-common/app/interface/main/app-interface/model/bplus"
	"go-common/app/interface/main/app-interface/model/sp"
	"go-common/app/interface/main/app-interface/model/topic"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/model/archive"
	"time"
)

type MyFavorite struct {
	Tab      *Tab     `json:"tab,omitempty"`
	Favorite *FavList `json:"favorite,omitempty"`
}

type Tab struct {
	Fav     bool `json:"favorite"`
	Topic   bool `json:"topic"`
	Article bool `json:"article"`
	Clips   bool `json:"clips"`
	Albums  bool `json:"albums"`
	Specil  bool `json:"specil"`
	Cinema  bool `json:"cinema"`
	Audios  bool `json:"audios"`
	Menu    bool `json:"menu"`
	PGCMenu bool `json:"pgc_menu"`
	Ticket  bool `json:"ticket"`
	Product bool `json:"product"`
}

type FavList struct {
	Count int        `json:"count"`
	Items []*FavItem `json:"items"`
}

type FavideoList struct {
	Count int            `json:"count"`
	Items []*FavideoItem `json:"items"`
}

type TopicList struct {
	Count int          `json:"count"`
	Items []*TopicItem `json:"items"`
}

type ArticleList struct {
	Count int            `json:"count"`
	Items []*ArticleItem `json:"items"`
}

type ClipsList struct {
	*bplus.PageInfo
	Items []*ClipsItem `json:"items"`
}

type AlbumsList struct {
	*bplus.PageInfo
	Items []*AlbumItem `json:"items"`
}

type SpList struct {
	Count int       `json:"count"`
	Items []*SpItem `json:"items"`
}

type AudioList struct {
	Count int          `json:"count"`
	Items []*AudioItem `json:"items"`
}

func (i *FavItem) FromFav(f *Folder) {
	i.MediaID = f.MediaID
	i.Fid = f.Fid
	i.Mid = f.Mid
	i.Name = f.Name
	if f.Cover != nil {
		i.Cover = f.Cover
	}
	i.CurCount = f.CurCount
	i.State = f.State
}

type FavItem struct {
	MediaID  int64   `json:"media_id"`
	Fid      int     `json:"fid"`
	Mid      int     `json:"mid"`
	Name     string  `json:"name"`
	CurCount int     `json:"cur_count"`
	State    int     `json:"state"`
	Cover    []Cover `json:"cover"`
}

func (i *FavideoItem) FromFavideo(fv *Archive) {
	i.Aid = fv.Aid
	i.Title = fv.Title
	i.Pic = fv.Pic
	i.Name = fv.Author.Name
	i.PlayNum = int(fv.Stat.View)
	i.Danmaku = int(fv.Stat.Danmaku)
	i.Param = strconv.FormatInt(int64(fv.Aid), 10)
	i.Goto = model.GotoAv
	i.URI = model.FillURI(i.Goto, i.Param, model.AvHandler(archive.BuildArchive3(fv.Arc)))
	i.UGCPay = fv.Rights.UGCPay

}

type FavideoItem struct {
	Aid     int64  `json:"aid"`
	Title   string `json:"title"`
	Pic     string `json:"pic"`
	Name    string `json:"name"`
	PlayNum int    `json:"play_num"`
	Danmaku int    `json:"danmaku"`
	Goto    string `json:"goto"`
	Param   string `json:"param"`
	URI     string `json:"uri"`
	UGCPay  int32  `json:"ugc_pay"`
}

func (i *TopicItem) FromTopic(tp *topic.List) {
	i.ID = tp.ID
	i.MID = tp.MID
	i.Name = tp.Name
	i.PCCover = tp.PCCover
	i.H5Cover = tp.H5Cover
	i.FavAt = tp.FavAt
	i.PCUrl = tp.PCUrl
	i.H5Url = tp.H5Url
	i.Desc = tp.Desc
	i.Param = strconv.FormatInt(int64(tp.ID), 10)
	i.Goto = model.GotoWeb
	i.URI = model.FillURI(i.Goto, i.Param, nil)
}

type TopicItem struct {
	ID      int64  `json:"id"`
	MID     int64  `json:"mid"`
	Name    string `json:"name"`
	PCCover string `json:"pc_cover"`
	H5Cover string `json:"h5_cover"`
	FavAt   int64  `json:"fav_at"`
	PCUrl   string `json:"pc_url"`
	H5Url   string `json:"h5_url"`
	Desc    string `json:"desc"`
	Goto    string `json:"goto"`
	Param   string `json:"param"`
	URI     string `json:"uri"`
}

func (i *ArticleItem) FromArticle(af *article.Favorite) {
	i.ID = af.ID
	i.Title = af.Title
	i.BannerURL = af.BannerURL
	i.TemplateID = int(af.TemplateID)
	i.Name = af.Author.Name
	i.ImageURLs = af.ImageURLs
	i.Summary = af.Summary
	i.FTime = af.FavoriteTime
	i.Param = strconv.FormatInt(int64(af.ID), 10)
	i.Goto = model.GotoArticle
	i.URI = model.FillURI(i.Goto, i.Param, nil)
}

type ArticleItem struct {
	ID         int64    `json:"id"`
	Title      string   `json:"title"`
	TemplateID int      `json:"template_id"`
	BannerURL  string   `json:"banner_url"`
	Name       string   `json:"name"`
	ImageURLs  []string `json:"image_urls"`
	Summary    string   `json:"summary"`
	FTime      int64    `json:"favorite_time"`
	Goto       string   `json:"goto"`
	Param      string   `json:"param"`
	URI        string   `json:"uri"`
}

func (i *ClipsItem) FromClips(c *bplus.ClipList) {
	i.ID = c.Content.Item.ID
	i.Name = c.Content.User.Name
	i.UID = c.Content.User.UID
	i.HeadURL = c.Content.User.HeadURL
	i.IsVIP = c.Content.User.IsVIP
	i.IsFollowed = c.Content.User.IsFollowed
	i.UploadTimeText = c.Content.Item.UploadTimeText
	i.Tags = c.Content.Item.Tags
	i.Cover = c.Content.Item.Cover
	i.VideoTime = c.Content.Item.VideoTime
	i.Desc = c.Content.Item.Desc
	i.DanakuNum = c.Content.Item.DanakuNum
	i.WatchedNum = c.Content.Item.WatchedNum
	i.Param = strconv.FormatInt(int64(c.Content.Item.ID), 10)
	i.Goto = model.GotoClip
	i.URI = model.FillURI(i.Goto, i.Param, nil)
	i.Status = c.Content.Item.ShowStatus
	i.Reply = c.Content.Item.Reply
	i.UploadTime = c.Content.Item.UploadTime
	i.Width = c.Content.Item.Width
	i.Height = c.Content.Item.Height
	i.FirstPic = c.Content.Item.FirstPic
	i.VideoPlayURL = c.Content.Item.VideoPlayURL
	i.BackupPlayURL = c.Content.Item.BackupPlayURL
	i.LikeNum = c.Content.Item.LikeNum
}

type ClipsItem struct {
	ID             int64    `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	UID            int64    `json:"uid,omitempty"`
	HeadURL        string   `json:"head_url,omitempty"`
	IsVIP          int      `json:"is_vip,omitempty"`
	UploadTimeText string   `json:"upload_time_text,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Cover          struct {
		Def string `json:"default,omitempty"`
	} `json:"cover,omitempty"`
	VideoTime     int      `json:"video_time,omitempty"`
	Desc          string   `json:"description,omitempty"`
	DanakuNum     int      `json:"damaku_num,omitempty"`
	WatchedNum    int      `json:"watched_num,omitempty"`
	Goto          string   `json:"goto,omitempty"`
	Param         string   `json:"param,omitempty"`
	URI           string   `json:"uri,omitempty"`
	Status        int      `json:"status,omitempty"`
	Reply         int      `json:"reply,omitempty"`
	FirstPic      string   `json:"first_pic,omitempty"`
	BackupPlayURL []string `json:"backup_playurl,omitempty"`
	IsFollowed    bool     `json:"is_followed,omitempty"`
	UploadTime    string   `json:"upload_time,omitempty"`
	Width         int      `json:"width,omitempty"`
	Height        int      `json:"height,omitempty"`
	VideoPlayURL  string   `json:"video_playurl,omitempty"`
	LikeNum       int      `json:"like_num,omitempty"`
}

func (i *AlbumItem) FromAlbum(bp *bplus.AlbumList) {
	i.ID = bp.Content.ID
	i.Pic = bp.Content.Pic
	i.PicCount = bp.Content.PicCount
	i.ShowStatus = bp.Content.ShowStatus
	i.Param = strconv.FormatInt(int64(bp.Content.ID), 10)
	i.Goto = model.GotoAlbum
	i.URI = model.FillURI(i.Goto, i.Param, nil)
}

type AlbumItem struct {
	ID         int64             `json:"id"`
	Pic        []*bplus.Pictures `json:"pictures"`
	ShowStatus int               `json:"show_status"`
	PicCount   int               `json:"pictures_count"`
	Goto       string            `json:"goto"`
	Param      string            `json:"param"`
	URI        string            `json:"uri"`
}

func (i *SpItem) FromSp(s *sp.Item) {
	i.SpID = s.SpID
	i.Title = s.Title
	i.Cover = s.Cover
	i.MCover = s.MCover
	i.SCover = s.SCover
	timeTmp, _ := time.Parse("2006-01-02 15:04", s.CTime)
	i.CTime = timeTmp.Unix()
	i.Param = strconv.FormatInt(int64(s.SpID), 10)
	i.Goto = model.GotoSp
	i.URI = model.FillURI(i.Goto, i.Param, nil)
}

type SpItem struct {
	SpID   int64  `json:"spid"`
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	MCover string `json:"m_cover"`
	SCover string `json:"s_cover"`
	CTime  int64  `json:"create_at"`
	Goto   string `json:"goto"`
	Param  string `json:"param"`
	URI    string `json:"uri"`
}

func (i *AudioItem) FromAudio(a *audio.FavAudio) {
	i.ID = a.ID
	i.Title = a.Title
	i.IsOpen = a.IsOpen
	i.Cover = a.ImgURL
	i.Count = a.RecordsNum
}

type AudioItem struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	IsOpen int    `json:"is_open"`
	Count  int    `json:"count"`
}

type TabItem struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
	Tab  string `json:"tab"`
}

type TabParam struct {
	MobiApp   string `form:"mobi_app"`
	Device    string `form:"device"`
	Build     int    `form:"build"`
	Platform  string `form:"platform"`
	Mid       int64  `form:"mid"`
	Business  string `form:"business"`
	AccessKey string `form:"access_key"`
	ActionKey string `form:"actionKey"`
	Filtered  string `form:"filtered"`
}

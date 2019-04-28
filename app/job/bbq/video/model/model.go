package model

import (
	"encoding/json"
	"go-common/library/time"
)

//静态变量
const (
	IsDeletedFalse = 0     //未删除标识
	DefaultVer     = "1.0" // DefaultVer 默认初始化版本
	FromMain       = 1     // FromMain 渠道来自主站

	//版本状态（1-草稿 2-待审核 3-待上架 4-已上架 -1-已下架 -2-强制下架）
	OnShelf = 4

	//Tag类型
	OPTag                 = 0 //运营标签
	TIDTag                = 1 //一级分区标签
	SubTIDTag             = 2 //二级分区标签
	NormalTag             = 3 //普通标签
	AllowSyncOperVideoTag = int64(1)
	DenySyncOperVideoTag  = int64(2)
	JobFinishNotice       = 1 //运营导入脚本完成邮件推送类型

	VideoStCheckBack         = 2  //视频状态回查
	VideoStPassReview        = 1  //审核通过
	VideoStPendingPassReview = 0  //原始稿件状态，等待安全审核
	VideoStPassReviewReject  = -1 //回查不通过
	VideoStCanPlay           = 3  //可放出
	VideoStHighGrade         = 4  //优质
	VideoStRecommend         = 5  //推荐
	VideoStInactive          = -3 //视频下架
	VideoStDeleted           = -4 //视频硬删除

	//origin sync st abandon
	VideoRepSyncStOrigin = 0
	//sub bvc commit
	VideoRepSyncStBvcCommit = 10
	//receive bvc resource
	VideoRepSyncStInsertBvcInfo = 20
	//video onshelf
	VideoRepSyncStOnshelf = 30

	UVStOpAdd = 1  //add
	UVStOpDel = -1 //delete
	//StateActive 评论状态
	StateActive = int16(0)
	//DefaultType ..
	DefaultType = int16(23)
	UserTypeUp  = 1
	//VideoFromBILI ..
	VideoFromBILI = 0
	//VideoFromBBQ ..
	VideoFromBBQ = 1
	//VideoFromCMS ..
	VideoFromCMS = 2

	//SourceRequest video_repository.sync_status source request
	SourceRequest = 1
	//SourceXcodeCover video_repository.sync_status xcode/cover
	SourceXcodeCover = 2
	//SourceAI video_repository.sync_status ai source
	SourceAI = 4
	//SourceOnshelf video_repository.sync_status video on shelf
	SourceOnshelf = 8

	VideoUploadProcessStatusFailed    = -1
	VideoUploadProcessStatusPending   = 0
	VideoUploadProcessStatusSuccessed = 1
)

//Tag .
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type int64  `json:"type"`
}

//VideoHiveInfo struct
type VideoHiveInfo struct {
	AVID           int64  `json:"avid"`
	CID            int64  `json:"cid"`
	MID            int64  `json:"mid"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Original       int16  `json:"original"`
	Report         int64  `json:"report"`
	DurationAll    int64  `json:"duration_all"`
	Play           int64  `json:"play"`
	PlayGuest      int64  `json:"play_guest"`
	PlayFans       int64  `json:"play_fans"`
	Access         int64  `json:"access"`
	Reply          int64  `json:"reply"`
	Fav            int64  `json:"fav"`
	Likes          int64  `json:"likes"`
	Coin           int64  `json:"coin"`
	Share          int64  `json:"share"`
	Danmu          int64  `json:"danmu"`
	ElecPay        int64  `json:"elec_pay"`
	ElecNum        int64  `json:"elec_num"`
	ElecUser       int64  `json:"elec_user"`
	Duration       int64  `json:"duration"`
	State          int64  `json:"state"`
	Tag            string `json:"tag"`
	ShareDaily     int64  `json:"share_daily"`
	PlayDaily      int64  `json:"play_daily"`
	FavDaily       int64  `json:"fav_daily"`
	ReplyDaily     int64  `json:"reply_daily"`
	DanmuDaily     int64  `json:"danmu_daily"`
	LikesDaily     int64  `json:"likes_daily"`
	DurationDaily  int64  `json:"duration_daily"`
	Pubtime        string `json:"pubtime"`
	LogDate        string `json:"log_date"`
	TID            int64  `json:"tid"`
	SubTID         int64  `json:"sub_tid"`
	Ctime          string `json:"ctime"`
	DispatchStatus int64  `json:"dispatch_status"`
	IsFullScreen   int16  `json:"is_full_screen"`
}

// VideoInfo 一般视频信息
type VideoInfo struct {
	SVID     int64     `json:"svid"`
	TID      int64     `json:"tid"`
	SubTID   int64     `json:"sub_tid"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	MID      int64     `json:"mid"`
	Report   int64     `json:"report"`
	Duration int64     `json:"duration"`
	Pubtime  string    `json:"pubtime"`
	Ctime    time.Time `json:"ctime"`
	AVID     int64     `json:"avid"`
	CID      int64     `json:"cid"`
	State    int16     `json:"state"`
	Original int16     `json:"original"`
	From     int16     `json:"from"`
	VerID    int64     `json:"ver_id"`
	Ver      int64     `json:"ver"`
}

// VideoStHive 视频hive统计数据
type VideoStHive struct {
	SVID           int64 `json:"svid"`
	Play           int64 `json:"play"`
	Report         int64 `json:"report"`
	DurationAll    int64 `json:"duration_all"`
	Access         int64 `json:"access"`
	Reply          int64 `json:"reply"`
	Fav            int64 `json:"fav"`
	Likes          int64 `json:"likes"`
	Coin           int64 `json:"coin"`
	Share          int64 `json:"share"`
	Subtitles      int64 `json:"subtitles"`
	ElecPay        int64 `json:"elec_pay"`
	ElecNum        int64 `json:"elec_num"`
	ElecUser       int64 `json:"elec_user"`
	DurationDaily  int64 `json:"duration_daily"`
	ShareDaily     int64 `json:"share_daily"`
	PlayDaily      int64 `json:"play_daily"`
	FavDaily       int64 `json:"fav_daily"`
	ReplyDaily     int64 `json:"reply_daily"`
	SubtitlesDaily int64 `json:"subtitles_daily"`
	LikesDaily     int64 `json:"likes_daily"`
}

// UserBaseDB 用户基础表字段
type UserBaseDB struct {
	ID        int64     `json:"id"`
	MID       int64     `json:"mid"`
	Uname     string    `json:"uname"`
	Face      string    `json:"face"`
	Birthday  string    `json:"birthday"`
	Exp       int64     `json:"exp"`
	Level     int64     `json:"level"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
	Signature string    `json:"signature"`
	Region    int64     `json:"region"`
	Sex       int16     `json:"sex"`
}

//UserDmg 用户画像
type UserDmg struct {
	MID          string           `json:"mid"`
	Gender       string           `json:"gender"`
	Age          string           `json:"age"`
	Geo          string           `json:"geo"`
	ContentTag   string           `json:"content_tag"`
	ViewedVideo  map[int64]string `json:"viewed_video"`
	ContentZone  string           `json:"content_zone"`
	ContentCount string           `json:"content_count"`
	FollowUps    string           `json:"follow_ups"`
}

//UserBbqDmg 用户画像
type UserBbqDmg struct {
	MID  string   `json:"mid"`
	Tag2 []string `json:"tag2"`
	Tag3 []string `json:"tag3"`
	Up   []string `json:"up"`
}

//UserBbqBuvidDmg 用户画像buvid
type UserBbqBuvidDmg struct {
	Buvid string   `json:"mid"`
	Tag2  []string `json:"tag2"`
	Tag3  []string `json:"tag3"`
	Up    []string `json:"up"`
}

//UpUserDmg 主站up主用户画像
type UpUserDmg struct {
	MID   int64  `json:"mid"`
	Uname string `json:"uname"`
	Play  int64  `json:"play"`
	Fans  int64  `json:"fans"`
	AVs   int64  `json:"avs"`
	Likes int64  `json:"likes"`
}

// CheckTask .
type CheckTask struct {
	TaskID    int64  `json:"task_id"`
	TaskName  string `json:"task_name"`
	LastCheck int64  `json:"last_check"`
}

// DatabusRes canal standary message
type DatabusRes struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

//DatabusBVCTransSub ...
type DatabusBVCTransSub struct {
	SVID int64 `json:"svid"`
}

// VideoDB 视频表数据库字段
type VideoDB struct {
	AutoID      int64     `json:"auto_id"`
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	MID         int64     `json:"mid"`
	CID         int64     `json:"cid"`
	Pubtime     time.Time `json:"pubtime"`
	Ctime       string    `json:"ctime"`
	Duration    int64     `json:"duration"`
	Original    int16     `json:"original"`
	State       int16     `json:"state"`
	IsFull      int16     `json:"is_full_screen"`
	VerID       int64     `json:"ver_id"`
	Ver         string    `json:"ver"`
	From        int16     `json:"from"`
	AVID        int64     `json:"avid"`
	TID         int64     `json:"tid"`
	SubTID      int64     `json:"sub_tid"`
	Score       int64     `json:"score"`
	CoverURL    string    `json:"cover_url"`
	CoverWidth  int64     `json:"cover_width"`
	CoverHeight int64     `json:"cover_height"`
}

// VideoRaw 视频原生表数据库字段
type VideoRaw struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MID         int64  `json:"mid"`
	CID         int64  `json:"cid"`
	Pubtime     string `json:"pubtime"`
	Ctime       string `json:"ctime"`
	Duration    int64  `json:"duration"`
	Original    int16  `json:"original"`
	State       int16  `json:"state"`
	IsFull      int16  `json:"is_full_screen"`
	VerID       int64  `json:"ver_id"`
	Ver         string `json:"ver"`
	From        int16  `json:"from"`
	AVID        int64  `json:"avid"`
	TID         int64  `json:"tid"`
	SubTID      int64  `json:"sub_tid"`
	Score       int64  `json:"score"`
	CoverURL    string `json:"cover_url"`
	CoverWidth  int64  `json:"cover_width"`
	CoverHeight int64  `json:"cover_height"`
	SVID        int64  `json:"svid"`
}

// VideoRepRaw 视频原生表数据库字段
type VideoRepRaw struct {
	ID            int64  `json:"id"`
	SVID          int64  `json:"svid"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	MID           int64  `json:"mid"`
	CID           int64  `json:"cid"`
	Pubtime       string `json:"pubtime"`
	Duration      int64  `json:"duration"`
	Original      int16  `json:"original"`
	IsFull        int16  `json:"is_full_screen"`
	From          int16  `json:"from"`
	AVID          int64  `json:"avid"`
	TID           int64  `json:"tid"`
	SubTID        int64  `json:"sub_tid"`
	Score         int64  `json:"score"`
	CoverURL      string `json:"cover_url"`
	CoverWidth    int64  `json:"cover_width"`
	CoverHeight   int64  `json:"cover_height"`
	Tag           string `json:"tag"`
	SyncStatus    int64  `json:"sync_status"`
	HomeImgURL    string `json:"home_img_url" form:"home_img_url"`
	HomeImgWidth  int64  `json:"home_img_width" form:"home_img_width"`
	HomeImgHeight int64  `json:"home_img_height" form:"home_img_height"`
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

// UserBase .
type UserBase struct {
	Mid   int64  `json:"mid"`
	Uname string `json:"uname"`
	Face  string `json:"face"`
}

//CmsVideo ..
type CmsVideo struct {
	ID       int64  `json:"id"`
	SVStatus int64  `json:"sv_status"`
	Pubtime  string `json:"pubtime"`
	Mid      int64  `json:"mid"`
	Title    string `json:"title"`
	From     int64  `json:"from"`
}

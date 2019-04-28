package model

import "go-common/library/time"

const (
	//FromBILI from bilibili
	FromBILI = 0
	//FromBBQ from bbq
	FromBBQ = 1
	//FromCMS from cms
	FromCMS = 2
	//SourceRequest video_repository.sync_status source request
	SourceRequest = 1
	//SourceXcodeCover video_repository.sync_status xcode/cover
	SourceXcodeCover = 2
	//SourceAI video_repository.sync_status ai source
	SourceAI = 4
	//SourceOnshelf video_repository.sync_status video on shelf
	SourceOnshelf = 8

	//UploadStatusFailed video_upload_process.upload_status
	UploadStatusFailed = -1
	//UploadStatusSuccessed video_upload_process.upload_status
	UploadStatusSuccessed = 1
	//UploadStatusWaiting video_upload_process.upload_status
	UploadStatusWaiting = 0
	//VideoUploadProcessStatusFailed .
	VideoUploadProcessStatusFailed = -1
	//VideoUploadProcessStatusPending .
	VideoUploadProcessStatusPending = 0
	//VideoUploadProcessStatusSuccessed .
	VideoUploadProcessStatusSuccessed = 1
)

//视频状态集合
const (
	//VideoStRecommend 推荐
	VideoStRecommend = 5
	//VideoStHighGrade 优质
	VideoStHighGrade = 4
	//VideoStCanPlay 可放出
	VideoStCanPlay = 3
	//VideoStCheckBack 视频状态回查
	VideoStCheckBack = 2
	//VideoStPassReview 审核通过
	VideoStPassReview = 1
	//VideoStPendingPassReview 原始稿件状态，等待安全审核
	VideoStPendingPassReview = 0
	//VideoStPassReviewReject 回查不通过，仅自见
	VideoStPassReviewReject = -1
	//VideoStCheckBackPatialPlay 回查不放出，在APP部分放出
	VideoStCheckBackPatialPlay = -2
	//VideoUnshelf 下架
	VideoUnshelf = -3
	//VideoDelete 删除
	VideoDelete = -4
)

//Tag .
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type int32  `json:"type"`
}

// VideoInfo 一般视频信息
type VideoInfo struct {
	SVID          int64     `json:"svid"`
	TID           int64     `json:"tid"`
	SubTID        int64     `json:"sub_tid"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	MID           int64     `json:"mid"`
	Report        int64     `json:"report"`
	Duration      int64     `json:"duration"`
	Pubtime       string    `json:"pubtime"`
	Ctime         time.Time `json:"ctime"`
	AVID          int64     `json:"avid"`
	CID           int64     `json:"cid"`
	State         int16     `json:"state"`
	Original      int64     `json:"original"`
	From          int16     `json:"from"`
	IsFullScreen  int16     `json:"is_full_screen"`
	CoverURL      string    `json:"cover_url"`
	CoverWidth    int64     `json:"cover_width"`
	CoverHeight   int64     `json:"cover_height"`
	HomeImgURL    string    `json:"home_img_url" form:"home_img_url"`
	HomeImgWidth  int64     `json:"home_img_width" form:"home_img_width"`
	HomeImgHeight int64     `json:"home_img_height" form:"home_img_height"`
}

//VideoUploadProcess .
type VideoUploadProcess struct {
	SVID          int64  `json:"svid"`
	Title         string `json:"Title"`
	Mid           int64  `json:"mid"`
	UploadStatus  int64  `json:"upload_status"`
	RetryTimes    int64  `json:"retry_times"`
	HomeImgURL    string `json:"home_img_url"`
	HomeImgWidth  int64  `json:"home_img_width"`
	HomeImgHeight int64  `json:"home_img_height"`
}

//VideoRepository ...
type VideoRepository struct {
	AVID          int64  `json:"avid"`
	CID           int64  `json:"cid"`
	MID           int64  `json:"mid"`
	SVID          int64  `json:"svid"`
	From          int64  `json:"from"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Original      int64  `json:"original"`
	Duration      int64  `json:"duration"`
	Pubtime       string `json:"pubtime"`
	TID           int64  `json:"tid"`
	SubTID        int64  `json:"sub_tid"`
	IsFullScreen  int64  `json:"is_full_screen"`
	CoverURL      string `json:"cover_url"`
	CoverWidth    string `json:"cover_width"`
	CoverHeight   string `json:"cover_height"`
	HomeImgURL    string `json:"home_img_url"`
	HomeImgWidth  int64  `json:"home_img_width"`
	HomeImgHeight int64  `json:"home_img_height"`
	SyncStatus    int64  `json:"sync_status"`
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

//VideoHiveInfo struct
type VideoHiveInfo struct {
	AVID          int64  `json:"avid"`
	CID           int64  `json:"cid"`
	MID           int64  `json:"mid"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Original      int16  `json:"original"`
	Report        int64  `json:"report"`
	DurationAll   int64  `json:"duration_all"`
	Play          int64  `json:"play"`
	PlayGuest     int64  `json:"play_guest"`
	PlayFans      int64  `json:"play_fans"`
	Access        int64  `json:"access"`
	Reply         int64  `json:"reply"`
	Fav           int64  `json:"fav"`
	Likes         int64  `json:"likes"`
	Coin          int64  `json:"coin"`
	Share         int64  `json:"share"`
	Danmu         int64  `json:"danmu"`
	ElecPay       int64  `json:"elec_pay"`
	ElecNum       int64  `json:"elec_num"`
	ElecUser      int64  `json:"elec_user"`
	Duration      int64  `json:"duration"`
	State         int64  `json:"state"`
	Tag           string `json:"tag"`
	ShareDaily    int64  `json:"share_daily"`
	PlayDaily     int64  `json:"play_daily"`
	FavDaily      int64  `json:"fav_daily"`
	ReplyDaily    int64  `json:"reply_daily"`
	DanmuDaily    int64  `json:"danmu_daily"`
	LikesDaily    int64  `json:"likes_daily"`
	DurationDaily int64  `json:"duration_daily"`
	Pubtime       string `json:"pubtime"`
	LogDate       string `json:"log_date"`
	TID           int64  `json:"tid"`
	SubTID        int64  `json:"sub_tid"`
	Ctime         string `json:"ctime"`
}

//UserBase .
type UserBase struct {
	Mid  int64  `json:"mid"`
	Name string `json:"uname"`
	Sex  string `json:"sex"`
	Face string `json:"face"`
	Sign string `json:"sign"`
	Rank int32  `json:"rank"`
}

// VideoBVC 视频转码信息
type VideoBVC struct {
	SVID            int64  `json:"svid"`
	Path            string `json:"path"`
	ResolutionRetio string `json:"resolution_retio"`
	CodeRate        int64  `json:"code_rate"`
	VideoCode       string `json:"video_code"`
	Duration        int64  `json:"duration"`
	FileSize        int64  `json:"file_size"`
}

// SvStInfo 视频统计
type SvStInfo struct {
	SVID      int64 `json:"svid"`
	Play      int64 `json:"view"` //和上层的play重复，因此改成view
	Subtitles int64 `json:"subtitles"`
	Like      int64 `json:"like"`
	Share     int64 `json:"share"`
	Reply     int64 `json:"reply"`
	Report    int64 `json:"report"`
}

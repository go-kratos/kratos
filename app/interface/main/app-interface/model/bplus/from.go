package bplus

import "go-common/library/time"

type Clips struct {
	List     []*ClipList `json:"list,omitempty"`
	PageInfo *PageInfo   `json:"pageinfo,omitempty"`
}

type ClipList struct {
	BizType    int       `json:"biz_type,omitempty"`
	Favid      int       `json:"fav_id,omitempty"`
	UTime      time.Time `json:"utime,omitempty"`
	HasMore    int       `json:"has_more,omitempty"`
	NextFffset int       `json:"next_offse,omitempty"`
	Content    struct {
		User struct {
			UID        int64  `json:"uid,omitempty"`
			HeadURL    string `json:"head_url,omitempty"`
			IsVIP      int    `json:"is_vip,omitempty"`
			Name       string `json:"name,omitempty"`
			IsFollowed bool   `json:"is_followed,omitempty"`
		} `json:"user,omitempty"`
		Item struct {
			Type  int   `json:"type,omitempty"`
			ID    int64 `json:"id,omitempty"`
			Cover struct {
				Def string `json:"default,omitempty"`
			} `json:"cover,omitempty"`
			Desc             string   `json:"description,omitempty"`
			Tags             []string `json:"tags,omitempty"`
			VideoTime        int      `json:"video_time,omitempty"`
			UploadTime       string   `json:"upload_time,omitempty"`
			Width            int      `json:"width,omitempty"`
			Height           int      `json:"height,omitempty"`
			UploadTimeText   string   `json:"upload_time_text,omitempty"`
			VerifyStatusText string   `json:"verify_status_text,omitempty"`
			ShareURL         string   `json:"share_url,omitempty"`
			JumpURL          string   `json:"jump_url,omitempty"`
			DanakuNum        int      `json:"damaku_num,omitempty"`
			WatchedNum       int      `json:"watched_num,omitempty"`
			VideoPlayURL     string   `json:"video_playurl,omitempty"`
			ShowStatus       int      `json:"show_status,omitempty"`
			ShareNum         int      `json:"share_num,omitempty"`
			EnshrineNum      int      `json:"enshrine_num,omitempty"`
			Reply            int      `json:"reply,omitempty"`
			FirstPic         string   `json:"first_pic,omitempty"`
			BackupPlayURL    []string `json:"backup_playurl,omitempty"`
			LikeNum          int      `json:"like_num,omitempty"`
		} `json:"item,omitempty"`
	} `json:"content,omitempty"`
}

type Albums struct {
	List     []*AlbumList `json:"list,omitempty"`
	PageInfo *PageInfo    `json:"pageinfo,omitempty"`
}

type AlbumList struct {
	BizType int       `json:"biz_type,omitempty"`
	Favid   int       `json:"fav_id,omitempty"`
	UTime   time.Time `json:"utime,omitempty"`
	Content struct {
		ID         int64       `json:"id,omitempty"`
		Pic        []*Pictures `json:"pictures,omitempty"`
		ShowStatus int         `json:"show_status,omitempty"`
		PicCount   int         `json:"pictures_count,omitempty"`
	} `json:"content,omitempty"`
}

type PageInfo struct {
	Page      string `json:"page,omitempty"`
	PageSize  string `json:"pagesize,omitempty"`
	TotalPage int    `json:"totalpage,omitempty"`
	Count     int    `json:"count,omitempty"`
}

// Detail struct
type Detail struct {
	ID              int64  `json:"dynamic_id,omitempty"`
	Mid             int64  `json:"mid,omitempty"`
	FaceImg         string `json:"face_img,omitempty"`
	NickName        string `json:"nick_name,omitempty"`
	PublishTimeText string `json:"publish_time_text	,omitempty"`
	ImgCount        int    `json:"img_count,omitempty"`
	ViewCount       int    `json:"view_count,omitempty"`
	CommentCount    int    `json:"comment_count,omitempty"`
	LikeCount       int    `json:"like_count,omitempty"`
}

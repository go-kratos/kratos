package archive

import (
	"go-common/app/admin/main/videoup/model/utils"
)

const (
	//QATypeVideo 视频质检任务
	QATypeVideo = int8(1)
)

//QAVideo 质检视频详情
type QAVideo struct {
	UID          int64   `json:"uid"`
	Oname        string  `json:"username"`
	AID          int64   `json:"aid"`
	CID          int64   `json:"cid"`
	TaskID       int64   `json:"task_id"`
	TaskUTime    int64   `json:"task_utime"`
	Attribute    int32   `json:"attribute"`
	TagID        int64   `json:"tag_id"`
	ArcTitle     string  `json:"arc_title"`
	ArcTypeid    int64   `json:"arc_typeid"`
	AuditStatus  int16   `json:"audit_status"`
	AuditSubmit  string  `json:"audit_submit"`
	AuditDetails string  `json:"audit_details"`
	MID          int64   `json:"mid"`
	UPGroups     []int64 `json:"up_groups"`
	Fans         int64   `json:"fans"`
}

//AuditSubmit 提交的审核内容
type AuditSubmit struct {
	Encoding  string `json:"encoding"`
	Attribute string `json:"attribute"`
	ReasonID  string `json:"reason_id"`
	Reason    string `json:"reason"`
	Note      string `json:"note"`
}

//AuditDetails 提交详情
type AuditDetails struct {
	UserInfo       map[string]interface{} `json:"user_info"`
	RelationVideos []*RelationVideo       `json:"relation_videos"`
	Task           []*Task                `json:"task"`
	Video          *VideoInfo             `json:"video"`
	Watermark      []*Watermark           `json:"watermark"`
	Mosaic         []*Mosaic              `json:"mosaic"`
}

//RelationVideo related video
type RelationVideo struct {
	Filename   string           `json:"filename"`
	Status     int16            `json:"status"`
	AID        int64            `json:"aid"`
	IndexOrder int              `json:"index_order"`
	Title      string           `json:"title"`
	Ctime      utils.FormatTime `json:"ctime"`
}

//VideoInfo video info
type VideoInfo struct {
	ID             int64            `json:"id"`
	MID            int64            `json:"mid"`
	CID            int64            `json:"cid"`
	Eptitle        string           `json:"eptitle"`
	Filename       string           `json:"filename"`
	Epctime        utils.FormatTime `json:"epctime"`
	AID            int64            `json:"aid"`
	Ctime          utils.FormatTime `json:"ctime"`
	Description    string           `json:"description"`
	Title          string           `json:"-"`
	Tag            string           `json:"tag"`
	Content        string           `json:"content"`
	Dynamic        string           `json:"dynamic"`
	Author         string           `json:"author"`
	Copyright      string           `json:"copyright"`
	Source         string           `json:"source"`
	Typename       string           `json:"typename"`
	Cover          string           `json:"cover"`
	XcodeState     int8             `json:"xcode_state"`
	XcodeStateName string           `json:"xcode_state_name"`
	Playurl        string           `json:"playurl"`
	Typeid         int64            `json:"-"`
}

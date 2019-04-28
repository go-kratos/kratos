package archive

import (
	"encoding/json"

	"go-common/library/time"
)

// ArcVideos 稿件及其所属视频
type ArcVideos struct {
	Archive *Archive `json:"archive"`
	Videos  []*Video `json:"videos"`
}

//UpArchives upper主的稿件ID和个数
type UpArchives struct {
	Count int64   `json:"count"`
	Aids  []int64 `json:"aids"`
}

// ArcMissionParam str
type ArcMissionParam struct {
	AID       int64  `form:"aid" validate:"required"`
	MID       int64  `form:"mid" validate:"required"`
	MissionID int64  `form:"mission_id" validate:"required"`
	Tag       string `form:"tag"`
}

// ArcDynamicParam str
type ArcDynamicParam struct {
	AID int64 `form:"aid" validate:"required"`
	MID int64 `form:"mid" validate:"required"`
}

//ArcParam 提交的稿件参数
type ArcParam struct {
	Aid          int64         `json:"aid"`
	Mid          int64         `json:"mid"`
	Author       string        `json:"author"`
	TypeID       int16         `json:"tid"`
	Title        string        `json:"title"`
	Cover        string        `json:"cover"`
	Tag          string        `json:"tag"`
	Copyright    int8          `json:"copyright"`
	Desc         string        `json:"desc"`
	AllowTag     int32         `json:"allow_tag"`
	NoReprint    int32         `json:"no_reprint"`
	UGCPay       int32         `json:"ugcpay"`
	MissionID    int64         `json:"mission_id"`
	FromIP       int64         `json:"from_ip"`
	IPv6         []byte        `json:"ipv6"`
	UpFrom       int8          `json:"up_from"`
	Source       string        `json:"source"`
	DTime        time.Time     `json:"dtime"`
	Videos       []*VideoParam `json:"videos"`
	Staffs       []*StaffParam `json:"staffs"`
	HandleStaff  bool          `json:"handle_staff"`
	CodeMode     bool          `json:"code_mode"`
	OrderID      int64         `json:"order_id"`
	FlowRemark   string        `json:"flow_remark"`
	Dynamic      string        `json:"dynamic"`
	IsDRM        int8          `json:"is_drm"`
	DescFormatID int64         `json:"desc_format_id"`
	Porder       *Porder       `json:"porder"`
	POI          *PoiObj       `json:"poi_object"`
	Vote         *Vote         `json:"vote"`
	Lang         string        `json:"lang"`
}

//Porder str
type Porder struct {
	// for user operation
	FlowID     int64  `json:"flow_id"`
	IndustryID int64  `json:"industry_id"`
	BrandID    int64  `json:"brand_id"`
	BrandName  string `json:"brand_name"`
	Official   int8   `json:"official"`
	ShowType   string `json:"show_type"`
	// for admin operation
	Advertiser string `json:"advertiser"`
	Agent      string `json:"agent"`
	//state 0 自首  1  审核添加
	State int8 `json:"state"`
}

//VideoParam 提交的视频参数
type VideoParam struct {
	Title    string  `json:"title"`
	Desc     string  `json:"desc"`
	Filename string  `json:"filename"`
	Cid      int64   `json:"cid"`
	Sid      int64   `json:"sid"`
	SrcType  string  `json:"src_type"`
	IsDRM    int8    `json:"is_drm"`
	Editor   *Editor `json:"editor"`
}

// Editor str
type Editor struct {
	CID    int64 `json:"cid"`
	UpFrom int8  `json:"upfrom"` // filled by backend
	// ids set
	Filters         interface{} `json:"filters"`          // 滤镜
	Fonts           interface{} `json:"fonts"`            //字体
	Subtitles       interface{} `json:"subtitles"`        //字幕
	Bgms            interface{} `json:"bgms"`             //bgm
	Stickers        interface{} `json:"stickers"`         //3d拍摄贴纸
	VideoupStickers interface{} `json:"videoup_stickers"` //2d投稿贴纸
	Transitions     interface{} `json:"trans"`            //视频转场特效
	// switch env 0/1
	Split        int8 `json:"split"`         //视频切片
	Cut          int8 `json:"cut"`           //拿时间窗口切子集
	VideoRotate  int8 `json:"rotate"`        //画面坐标轴变换
	AudioRecord  int8 `json:"audio_record"`  //录音
	Camera       int8 `json:"camera"`        //拍摄
	Speed        int8 `json:"speed"`         //变速
	Beauty       int8 `json:"beauty"`        //美颜特效
	Flashlight   int8 `json:"flashlight"`    //闪光灯
	CameraRotate int8 `json:"camera_rotate"` //摄像头翻转
	CountDown    int8 `json:"countdown"`     //拍摄倒计时
}

// UnmarshalJSON fn
func (vp *VideoParam) UnmarshalJSON(data []byte) (err error) {
	type VpAlias VideoParam
	tmp := &VpAlias{SrcType: "vupload"}
	if err = json.Unmarshal(data, tmp); err != nil {
		return err
	}
	*vp = VideoParam(*tmp)
	return
}

//PubAgentParam 提交的视频参数
type PubAgentParam struct {
	Route       string `json:"route"`
	Timestamp   string `json:"timestamp"`
	Filename    string `json:"filename"`
	Xcode       int8   `json:"xcode"`
	VideoDesign string `json:"video_design"`
	Submit      int8   `json:"submit"`
}

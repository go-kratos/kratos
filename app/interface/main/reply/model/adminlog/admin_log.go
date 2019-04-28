package adminlog

// AdminLog AdminLog
type AdminLog struct {
	ReplyID      int64  `json:"reply_id"`       //操作人
	AdminID      int64  `json:"adminid"`        //操作人
	Operator     string `json:"operator"`       //操作人昵称
	State        int32  `json:"operator_type"`  //操作人身份
	ReplyMid     int64  `json:"replymid"`       //评论人
	ReplyUser    string `json:"reply_user"`     //评论人昵称
	ReplyFacePic string `json:"reply_face_pic"` //评论人头像
	CTime        string `json:"operation_time"` //删除时间
}

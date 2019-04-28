package model

const (
	// MaxUploadSize max upload file size
	MaxUploadSize = 20 * 1024 * 1024
)

// UploadActionType report action type
type UploadActionType int

// report action type
const (
	UploadInternal      UploadActionType = iota + 1 // 内网用户
	UploadInternalAdmin                             // 内网管理员
	UploadPublic                                    // 外网公用
	UploadApp                                       // 外网 app
	UploadWeb                                       // 外网web
)

func (a UploadActionType) String() (s string) {
	switch a {
	case UploadInternal:
		s = "internal_upload"
	case UploadInternalAdmin:
		s = "internal_admin_upload"
	case UploadPublic:
		s = "outer_public_upload"
	case UploadApp:
		s = "outer_app_upload"
	case UploadWeb:
		s = "outer_web_upload"
	default:
		s = "undefined_upload"
	}
	return
}

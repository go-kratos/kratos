package model

const (
	// MobileUserAgentFlag is
	MobileUserAgentFlag = "Mobile"
	// PlatH5 is
	PlatH5 = "h5"
	// PlatPC is
	PlatPC = "pc"
)

// ProcessRes 请求获取极验验证响应参数
type ProcessRes struct {
	Success    int64  `json:"success"`
	CaptchaID  string `json:"gt"`
	Challenge  string `json:"challenge"`
	NewCaptcha int    `json:"new_captcha"`
}

//ValidateRes 验证返回值
type ValidateRes struct {
	Seccode string `json:"seccode"`
}

// GeeCaptchaRequest 获取极验验证请求参数
type GeeCaptchaRequest struct {
	MID int64  `json:"mid" form:"mid"`
	IP  string `json:"ip" form:"ip"`
	// h5  web  native
	ClientType string `json:"client_type" form:"client_type" default:"web"`
}

// GeeCheckRequest 校验极验验证码请求参数
type GeeCheckRequest struct {
	MID        int64
	Challenge  string `json:"challenge" form:"challenge" validate:"required"`
	Validate   string `json:"validate" form:"validate" validate:"required"`
	Seccode    string `json:"seccode" form:"seccode" validate:"required"`
	ClientType string `json:"client_type" form:"client_type" default:"web"`
}

// GeeCheckResponse 校验极验验证码响应参数
type GeeCheckResponse struct {
	NewVoucher string `json:"new_voucher"`
}

package model

// var
var (
	// Geetest captcha_type .
	Geetest     = "gt"
	BiliCaptcha = "bili"
)

// ProcessRes geetest Captcha resp info.
type ProcessRes struct {
	Success     int8   `json:"success"`
	CaptchaID   string `json:"gt"`
	Challenge   string `json:"challenge"`
	NewCaptcha  int    `json:"new_captcha"`
	CaptchaType string `json:"type,omitempty"`
	Token       string `json:"token,omitempty"`
	URL         string `json:"url,omitempty"`
}

// ValidateRes info.
type ValidateRes struct {
	Seccode string `json:"seccode"`
}

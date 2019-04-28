package geetest

// ProcessRes str
type ProcessRes struct {
	Success    int8            `json:"success"`
	CaptchaID  string          `json:"gt"`
	Challenge  string          `json:"challenge"`
	NewCaptcha int             `json:"new_captcha"`
	Limit      map[string]bool `json:"limit"`
}

// ValidateRes str
type ValidateRes struct {
	Seccode string `json:"seccode"`
}

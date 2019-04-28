package geetest

type ProcessRes struct {
	Success    int8   `json:"success"`
	CaptchaID  string `json:"gt"`
	Challenge  string `json:"challenge"`
	NewCaptcha int    `json:"new_captcha"`
}

type ValidateRes struct {
	Seccode string `json:"seccode"`
}

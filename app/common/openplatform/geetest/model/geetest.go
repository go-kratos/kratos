package model

//ProcessRes 拉起返回值
type ProcessRes struct {
	Success    int    `json:"success"`
	CaptchaID  string `json:"gt"`
	Challenge  string `json:"challenge"`
	NewCaptcha int    `json:"new_captcha"`
}

//ValidateRes 验证返回值
type ValidateRes struct {
	Seccode string `json:"seccode"`
}

package model

// RealnameBool is.
type RealnameBool uint8

// RealnameBool enum
const (
	RealnameFalse RealnameBool = iota
	RealnameTrue
)

// RealnameCountry .
type RealnameCountry struct {
	ID    int    `json:"id"`
	CName string `json:"cname"`
}

// RealnameCardType .
type RealnameCardType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RealnameChannel .
type RealnameChannel struct {
	Name string
	Flag RealnameBool
}

// ParamRealnameAlipayApply .
type ParamRealnameAlipayApply struct {
	Realname string `form:"real_name"`
	CardNum  string `form:"card_num"`
	Capture  int    `form:"capture"`
	ImgToken string `form:"img_token"`
}

// RealnameAlipayApply .
type RealnameAlipayApply struct {
	URL string `json:"url"`
	// Bizno string `json:"biz_no"`
}

// RealnameAlipayConfirm .
type RealnameAlipayConfirm struct {
	Passed RealnameBool `json:"passed"`
	Reason string       `json:"reason"`
}

// ParamRealnameCaptchaGTCheck .
type ParamRealnameCaptchaGTCheck struct {
	Remote    RealnameBool `form:"remote"`
	Challenge string       `form:"challenge"`
	Validate  string       `form:"validate"`
	Seccode   string       `form:"seccode"`
}

// ParamRealnameCaptchaGTRefresh .
type ParamRealnameCaptchaGTRefresh struct {
	Hash string `form:"hash"`
}

// RealnameCaptchaGTRegister .
type RealnameCaptchaGTRegister struct {
	Remote    RealnameBool `json:"remote"`
	ID        string       `json:"gt_id"`
	Challenge string       `json:"gt_challenge"`
}

// RealnameCaptchaGTValidate .
type RealnameCaptchaGTValidate struct {
	State RealnameBool `json:"state"`
}

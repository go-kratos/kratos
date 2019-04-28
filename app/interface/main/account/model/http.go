package model

// TokenResq .
type TokenResq struct {
	CommonResq
	Data *Token `json:"data"`
}

// Token .
type Token struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

// CommonResq .
type CommonResq struct {
	Code    int64  `json:"code"`
	TS      int64  `json:"ts"`
	Message string `json:"message"`
}

// ResourceCodeResq .
type ResourceCodeResq struct {
	CommonResq
	Data *ResourceCode `json:"data"`
}

// ResourceCode .
type ResourceCode struct {
	Code        string `json:"code"`
	Status      int8   `json:"status"`
	Days        int32  `json:"days"`
	BatchCodeID int64  `json:"batch_code_id"`
}

// CMAccountInfo is
type CMAccountInfo struct {
	Nickname           string `json:"nickname"`
	CertificationTitle string `json:"certification_title"`
	CreditCode         string `json:"credit_code"`
	CompanyName        string `json:"company_name"`
	Organization       string `json:"organization"`
	OrganizationType   string `json:"organization_type"`
}

// OfficialPermissionResponse is
type OfficialPermissionResponse struct {
	DeniedRoles []int8                 `json:"denied_roles"`
	Metadata    map[string]interface{} `json:"metadata"`
}

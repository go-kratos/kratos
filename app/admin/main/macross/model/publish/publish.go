package publish

// Dashboard for dashboard.
type Dashboard struct {
	Name          string `json:"name"`
	Label         string `json:"label"`
	Commit        string `json:"commit"`
	OutURL        string `json:"out_url"`
	CoverageURL   string `json:"coverage_url"`
	TextSizeArm64 int64  `json:"text_size_arm64"`
	ResSize       int64  `json:"res_size"`
	Logs          []*Log `json:"logs"`
	Extra         string `json:"extra"`
}

// Log from Dashboard.
type Log struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

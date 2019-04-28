package model

// UserConf user configruation
type UserConf struct {
	ModuleKey int   `json:"module_key,omitempty"`
	Mid       int64 `json:"mid,omitempty"`
	CheckSum  int64 `json:"check_sum"`
	Timestamp int64 `json:"timestamp"`
}

// Document data store
type Document struct {
	CheckSum int64  `json:"check_sum"`
	Doc      string `json:"doc"`
}

package model

// Secret secret
type Secret struct {
	KeyType int8   `json:"key_type"`
	Key     string `json:"key"`
}

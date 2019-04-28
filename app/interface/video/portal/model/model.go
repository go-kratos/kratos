package model

// StreamForwardingConfResponse 转推白名单返回结构
type StreamForwardingConfResponse struct {
	Global  []string            `json:"global,omitempty"`
	Address map[string]string   `json:"address,omitempty"`
	List    map[string][]string `json:"list,omitempty"`
}

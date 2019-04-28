package tree

import "time"

// Token token
type Token struct {
	Token    string `json:"token"`
	UserName string `json:"user_name"`
	Secret   string `json:"secret"`
	Expired  int64  `json:"expired"`
}

// TokenResult token result
type TokenResult struct {
	Code    int    `json:"code"`
	Data    *Token `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Resp tree resp
type Resp struct {
	Data []*Node `json:"data"`
}

// Node node
type Node struct {
	TreeID      int    `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        int    `json:"type"`
	Role        int    `json:"role"`
	DiscoveryID string `json:"discovery_id"`
}

// Tree tree model
type Tree struct {
	Project string  `json:"project"`
	Subs    []*Tree `json:"subs"`
}

// Rest tree rest
type Rest struct {
	Data []*Info `json:"data"`
}

// Info tree info
type Info struct {
	AppTreeID int    `json:"app_tree_id"`
	AppID     string `json:"app_id"`
}

// Resd tree resd
type Resd struct {
	Data  []*DiscoveryID `json:"data"`
	CTime time.Time      `json:"ctime"`
}

// DiscoveryID node
type DiscoveryID struct {
	TreeID      int    `json:"app_tree_id"`
	AppID       string `json:"app_id"`
	AppAuth     string `json:"app_auth"`
	DiscoveryID string `json:"discovery_id"`
}

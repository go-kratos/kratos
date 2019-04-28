package model

import "time"

// Res res.
type Res struct {
	Count   int         `json:"count"`
	Data    []*TreeNode `json:"data"`
	Page    int         `json:"page"`
	Results int         `json:"results"`
}

// TreeNode TreeNode.
type TreeNode struct {
	Alias     string      `json:"alias"`
	CreatedAt string      `json:"created_at"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Tags      interface{} `json:"tags"`
	Type      int         `json:"type"`
}

// Node node.
type Node struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	TreeID int64  `json:"tree_id"`
}

//CacheData ...
type CacheData struct {
	Data  map[int64]*RoleNode `json:"data"`
	CTime time.Time           `json:"ctime"`
}

//RoleNode roleNode .
type RoleNode struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Type int8   `json:"type"`
	Role int8   `json:"role"`
}

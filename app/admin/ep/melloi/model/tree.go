package model

// TreeResponse service tree response model
type TreeResponse struct {
	Code int      `json:"code"`
	Data UserTree `json:"data"`
}

// UserTree user tree model
type UserTree struct {
	Bilibili map[string]interface{} `json:"bilibili"`
}

// TreeSonResponse  tree son response model
type TreeSonResponse struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
}

// TreeConf service tree conf model
type TreeConf struct {
	Host string
}

// TokenResponse service tree token response model
type TokenResponse struct {
	Token    string `json:"token"`
	UserName string `json:"user_name"`
	Secret   string `json:"secret"`
	Expired  int64  `json:"expired"`
}

// TreeAdminResponse service tree admin response model
type TreeAdminResponse struct {
	Code int         `json:"code"`
	Data []*TreeRole `json:"data"`
}

// TreeRole service tree role
type TreeRole struct {
	UserName string `json:"user_name"`
	Role     int    `json:"role"`
	OldRole  int    `json:"old_role"`
}

// TreeRoleApp tree role app
type TreeRoleApp struct {
	Code int        `json:"code"`
	Data []*RoleApp `json:"data"`
}

// RoleApp role app
type RoleApp struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Type int    `json:"type"` // Type 1.公司 2.部门 3.项目 4. 应用 5.环境 6.挂载点
	Role int    `json:"role"` // Role 1:管理员 2:研发 3:测试 4:运维 5:访客
}

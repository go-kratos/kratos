package model

//UserRole role
type UserRole struct {
	UID  int64  `json:"uid"`
	Name string `json:"username"`
	Role int8   `json:"role"`
}

//UserDepart department
type UserDepart struct {
	UID        int64  `json:"uid"`
	Name       string `json:"username"`
	Department string `json:"department"`
}

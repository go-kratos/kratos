package model

// official role const.
const (
	OfficialRoleUnauth = iota
	OfficialRoleUp
	OfficialRoleIdentify
	OfficialRoleBusiness
	OfficialRoleGov
	OfficialRoleMedia
	OfficialRoleOther
)

// Official is.
type Official struct {
	Role  int8   `json:"role"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

// FromCert is.
func FromCert(v *MemberVerify) Official {
	of := Official{}
	switch v.Type {
	case -1:
		of.Role = 0
		of.Title = ""
	case 0:
		of.Role = 2
		of.Title = v.Desc
	case 1:
		of.Role = 3
		of.Title = v.Desc
	}
	of.Desc = ""
	return of
}

package model

import "go-common/library/ecode"

// ProjectInfo def
type ProjectInfo struct {
	ProjectID     int    `json:"project_id" gorm:"column:project_id"`
	Name          string `json:"name" gorm:"column:name"`
	Description   string `json:"description" gorm:"column:description"`
	WebURL        string `json:"web_url" gorm:"column:web_url"`
	Repo          string `json:"repo" gorm:"column:repo"`
	DefaultBranch string `json:"default_branch" gorm:"column:default_branch"`
	Owner         string `json:"owner" gorm:"column:owner"`
	SpaceName     string `json:"namespace_name" gorm:"column:namespace_name"`
	SpaceKind     string `json:"namespace_kind" gorm:"column:namespace_kind"`
	Saga          bool   `json:"saga" gorm:"column:saga"`
	Runner        bool   `json:"runner" gorm:"column:runner"`
	Department    string `json:"department" gorm:"column:department"`
	Business      string `json:"business" gorm:"column:business"`
	Language      string `json:"language" gorm:"column:language"`
}

// ProjectInfoRequest Project Info Request.
type ProjectInfoRequest struct {
	Pagination
	TeamParam
	Username string `form:"username"`
	Name     string `form:"name"`
}

// ProjectInfoResp ...
type ProjectInfoResp struct {
	Total       int              `json:"total"`
	Saga        int              `json:"saga"`
	Runner      int              `json:"runner"`
	SagaScale   int              `json:"saga_scale"`
	RunnerScale int              `json:"runner_scale"`
	PageNum     int              `json:"page_num"`
	PageSize    int              `json:"page_size"`
	ProjectInfo []*MyProjectInfo `json:"ProjectInfo"`
}

// FavoriteProjectsResp resp for favorite projects
type FavoriteProjectsResp struct {
	Total int `json:"total"`
	Pagination
	Projects []*MyProjectInfo `json:"projects,omitempty"`
}

// ProjectFavorite def
type ProjectFavorite struct {
	ID       int64  `json:"id,omitempty" gorm:"column:id"`
	UserName string `json:"user_name,omitempty" gorm:"column:user_name"`
	ProjID   int    `json:"proj_id,omitempty" gorm:"column:proj_id"`
}

// EditFavoriteReq params for edit favorite
type EditFavoriteReq struct {
	ProjID int  `json:"proj_id" validate:"required"`
	Star   bool `json:"star"`
}

// MyProjectInfo mask as star
type MyProjectInfo struct {
	*ProjectInfo
	Star bool `json:"is_star"`
}

// ProjectsInfoResp resp for query projects with start mark
type ProjectsInfoResp struct {
	Projects []*MyProjectInfo `json:"projects"`
}

// Verify verify the value of pageNum and pageSize.
func (p *Pagination) Verify() error {
	if p.PageNum < 0 {
		return ecode.MerlinIllegalPageNumErr
	} else if p.PageNum == 0 {
		p.PageNum = DefaultPageNum
	}
	if p.PageSize < 0 {
		return ecode.MerlinIllegalPageSizeErr
	} else if p.PageSize == 0 {
		p.PageSize = DefaultPageSize
	}
	return nil
}

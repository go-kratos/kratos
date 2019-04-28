package model

// SearchHubResponse Search Hub Response.
type SearchHubResponse struct {
	Repository []*HubRepo `json:"repository"`
}

// HubRepo HubRepo.
type HubRepo struct {
	ProjectID      int    `json:"project_id"`
	ProjectName    string `json:"project_name"`
	ProjectPublic  bool   `json:"project_public"`
	RepositoryName string `json:"repository_name"`
	TagsCount      int    `json:"tags_count"`
}

// HubProject HubProject.
type HubProject struct {
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
	RepoCount int    `json:"repo_count"`
}

// GetHubProjectDetailResponse GetHubProjectDetailResponse.
type GetHubProjectDetailResponse struct {
	ProjectID   int    `json:"project_id"`
	ProjectName string `json:"name"`
	RepoCount   int    `json:"repo_count"`
}

// PaginateProjectRepoRecord PaginateProjectRepoRecord.
type PaginateProjectRepoRecord struct {
	Total             int                  `json:"total"`
	PageNum           int                  `json:"page_num"`
	PageSize          int                  `json:"page_size"`
	ProjectRepository []*ProjectRepository `json:"project_repositories"`
}

// ProjectRepositoryRequest ProjectRepositoryRequest.
type ProjectRepositoryRequest struct {
	ProjectRepository []*ProjectRepository
}

// ProjectRepository ProjectRepository.
type ProjectRepository struct {
	RepositoryID   int    `json:"id"`
	RepositoryName string `json:"name"`
	TagCount       int    `json:"tags_count"`
	CreateTime     string `json:"creation_time"`
	UpdateTime     string `json:"update_time"`
}

// RepositoryTagResponse Repository Tag Response.
type RepositoryTagResponse struct {
	Digest  string `json:"digest"`
	Name    string `json:"name"`
	OS      string `json:"os"`
	Size    int64  `json:"size"`
	Created string `json:"created"`
}

// RepositoryTag Repository Tag.
type RepositoryTag struct {
	RepositoryTagResponse
	ImageFullName string `json:"image_full_name"`
}

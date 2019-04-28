package model

// FollowersListPage is the model for followers list result
type FollowersListPage struct {
	TotalCount int           `json:"total_count"`
	PN         int           `json:"pn"`
	PS         int           `json:"ps"`
	List       FollowersList `json:"list"`

	Order string `json:"order"`
	Sort  string `json:"sort"`
}

// FollowingsListPage is the model for followings list result
type FollowingsListPage struct {
	TotalCount int            `json:"total_count"`
	PN         int            `json:"pn"`
	PS         int            `json:"ps"`
	List       FollowingsList `json:"list"`

	Order string `json:"order"`
	Sort  string `json:"sort"`
}

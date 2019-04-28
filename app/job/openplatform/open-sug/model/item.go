package model

// Item ...
type Item struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	IPRightID    int    `json:"ip_right_id"`
	Brief        string `json:"brief"`
	Keywords     string `json:"keywords"`
	SalesCount   int    `json:"sales_count"`
	WishCount    int    `json:"wish_count"`
	CommentCount int    `json:"comment_count"`
	HeadImg      string `json:"head_img"`
	SugImg       string `json:"sug_img"`
}

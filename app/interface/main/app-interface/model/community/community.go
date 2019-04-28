package community

type Community struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Desc           string `json:"desc"`
	Thumb          string `json:"thumb"`
	PostCount      int    `json:"post_count"`
	MemberCount    int    `json:"member_count"`
	PostNickname   string `json:"post_nickname"`
	MemberNickname string `json:"member_nickname"`
}

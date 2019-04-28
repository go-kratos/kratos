package common

//CardMap .
type CardMap struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

//WebSearch .
type WebSearch struct {
	CardList []*CardMap `json:"web_search"`
}

//CardType .
func (s *Service) CardType() (m *WebSearch) {
	c := &CardMap{
		Name: "特殊小卡",
		ID:   1,
	}
	cards := make([]*CardMap, 0)
	cards = append(cards, c)
	return &WebSearch{
		CardList: cards,
	}
}

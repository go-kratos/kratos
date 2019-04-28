package model

// Score ...
type Score struct {
	SeasonID   string
	SeasonName string
	Score      float64
}

// Sug ...
type Sug struct {
	Item     *Item
	SeasonID int64
	Score    int64
}

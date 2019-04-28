package v1

// ScoreRequest .
type ScoreRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

// ScoreResponse .
type ScoreResponse struct {
	Name  string  `json:"name"`
	Score float32 `json:"score"`
}

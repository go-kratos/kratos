package model

// HTTPFilterRes struct .
type HTTPFilterRes struct {
	MSG    string   `json:"msg"`
	Level  int8     `json:"level"`
	TypeID []int64  `json:"typeid"`
	Hit    []string `json:"hit"`
	Limit  int      `json:"limit"`
	AI     *AiScore `json:"ai"`
}

// HTTPAreaFilterRes struct .
type HTTPAreaFilterRes struct {
	MSG    string  `json:"msg"`
	Level  int8    `json:"level"`
	TypeID []int64 `json:"typeid"`
}

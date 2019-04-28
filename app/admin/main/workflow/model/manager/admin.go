package manager

// UNameSearchResult .
type UNameSearchResult struct {
	Code    int              `json:"code"`
	Data    map[int64]string `json:"data"`
	Message string           `json:"message"`
	TTL     int32            `json:"ttl"`
}

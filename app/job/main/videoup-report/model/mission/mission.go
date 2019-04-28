package mission

import "time"

type Mission struct {
	ID    int       `json:"id"`
	Name  string    `json:"name"`
	Tags  string    `json:"tags"`
	ETime time.Time `json:"etime"`
}

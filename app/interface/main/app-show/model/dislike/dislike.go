package dislike

type DisklikeReason struct {
	ReasonID   int    `json:"reason_id,omitempty"`
	ReasonName string `json:"reason_name,omitempty"`
}

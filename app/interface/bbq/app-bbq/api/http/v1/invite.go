package v1

// InvitationRequest .
type InvitationRequest struct {
	MID int64 `json:"mid" form:"mid" validate:"required"`
}

// Download .
type Download struct {
	URL  string `json:"url"`
	PWD  string `json:"pwd"`
	Type int8   `json:"type"`
}

// InvitationResponse .
type InvitationResponse struct {
	Data       []*Download `json:"download,omitempty"`
	InviteCode int32       `json:"invite_code"`
}

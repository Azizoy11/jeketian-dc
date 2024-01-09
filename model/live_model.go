package model

type Live struct {
	MemberUsername    string `json:"member_username,omitempty"`
	MemberDisplayName string `json:"member_display_name,omitempty"`
	Platform          string `json:"platform,omitempty"`
	Title             string `json:"title,omitempty"`
	StreamUrl         string `json:"stream_url,omitempty"`
	Views             uint   `json:"views,omitempty"`
	ImageUrl          string `json:"image_url"`
	StartedAt         uint   `json:"started_at"`
}

type OnLives struct {
	MemberOnLives []string `json:"member_onlives"`
}

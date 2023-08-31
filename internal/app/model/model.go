package model

type Segment struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserExperiment struct {
	ID      int64   `json:"id"`
	UserID  int64   `json:"user_id"`
	Segment Segment `json:"segment"`
}

type UserExperimentList struct {
	UserID   int64     `json:"user_id"`
	Segments []Segment `json:"segments"`
}

type LogInfo struct {
	UserID int64  `json:"user_id"`
	From   string `json:"from"`
	Path   string `json:"url"`
}

type UserExperimentItem struct {
	Name      string `json:"name" validate:"required"`
	ExpiresAt string `json:"expired_at,omitempty"`
}

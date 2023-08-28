package model

// type User struct {
// 	ID int64
// }

type Segment struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserExperiment struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	SegmentID int64 `json:"segment_id"`
}

type UserExperimentList struct {
	UserID   int64     `json:"user_id"`
	Segments []Segment `json:"segments"`
}

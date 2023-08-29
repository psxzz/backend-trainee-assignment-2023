package model

// type User struct {
// 	ID int64
// }

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

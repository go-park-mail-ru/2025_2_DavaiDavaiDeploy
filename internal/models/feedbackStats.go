package models

type FeedbackStats struct {
	Total       int64 `json:"total" db:"total"`
	Open        int64 `json:"open" db:"open"`
	InProgress  int64 `json:"in_progress" db:"in_progress"`
	Closed      int64 `json:"closed" db:"closed"`
	Bugs        int64 `json:"bugs" db:"bugs"`
	FeatureReqs int64 `json:"feature_requests" db:"feature_requests"`
	Complaints  int64 `json:"complaints" db:"complaints"`
	Questions   int64 `json:"questions" db:"questions"`
}

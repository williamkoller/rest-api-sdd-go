package entity

import "time"

type FeedPost struct {
	ID          string
	SchoolID    string
	UnitID      string
	AuthorID    string
	Body        string
	ImageURL    string
	PublishedAt time.Time
	UpdatedAt   time.Time
}

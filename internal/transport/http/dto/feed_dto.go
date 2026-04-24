package dto

import (
	"time"

	"github.com/williamkoller/rest-api-sdd-go/internal/domain/entity"
)

type FeedPostResponse struct {
	ID          string    `json:"id"`
	SchoolID    string    `json:"schoolId"`
	UnitID      string    `json:"unitId"`
	AuthorID    string    `json:"authorId"`
	Body        string    `json:"body"`
	ImageURL    string    `json:"imageUrl"`
	PublishedAt time.Time `json:"publishedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func MapFeedPost(e *entity.FeedPost) *FeedPostResponse {
	return &FeedPostResponse{
		ID:          e.ID,
		SchoolID:    e.SchoolID,
		UnitID:      e.UnitID,
		AuthorID:    e.AuthorID,
		Body:        e.Body,
		ImageURL:    e.ImageURL,
		PublishedAt: e.PublishedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func MapFeedPosts(list []*entity.FeedPost) []*FeedPostResponse {
	res := make([]*FeedPostResponse, len(list))
	for i, e := range list {
		res[i] = MapFeedPost(e)
	}
	return res
}

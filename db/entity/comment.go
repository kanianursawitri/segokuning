package entity

import "time"

type Comment struct {
	Id        int       `json:"id"`
	Comment   string    `json:"comment"`
	PostID    int       `json:"postId"`
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

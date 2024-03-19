package entity

import "time"

type Post struct {
	Id         int       `json:"id"`
	PostInHtml string    `json:"postInHtml"`
	Tags       []string  `json:"tags"`
	UserID     int       `json:"userId"`
	CreatedAt  time.Time `json:"createdAt"`
}

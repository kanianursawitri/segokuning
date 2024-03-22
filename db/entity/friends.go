package entity

import "time"

type Friend struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	FriendID  int       `json:"friendId"`
	CreatedAt time.Time `json:"createdAt"`
}

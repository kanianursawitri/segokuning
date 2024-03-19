package entity

import "time"

type Friend struct {
	FirstUserID  int       `json:"firstUserId"`
	SecondUserID int       `json:"secondUserId"`
	CreatedAt    time.Time `json:"createdAt"`
}

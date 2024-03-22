package entity

import "time"

type (
	Friend struct {
		ID        int       `json:"id"`
		UserID    int       `json:"userId"`
		FriendID  int       `json:"friendId"`
		Name      string    `json:"name"`
		ImageUrl  *string   `json:"imageUrl"`
		CreatedAt time.Time `json:"createdAt"`
	}

	QueryGetFriends struct {
		UserID      int    `query:"userId"`
		Limit       int    `query:"limit"`
		Offset      int    `query:"offset"`
		SortBy      string `query:"sortBy"`
		OrderBy     string `query:"orderBy"`
		OnlyFriends bool   `query:"onlyFriends"`
		Search      string `query:"search"`
	}
)

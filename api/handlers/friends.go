package handlers

import (
	"errors"
	"strconv"

	"segokuning/api/responses"
	"segokuning/db/entity"
	"segokuning/db/functions"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
)

type (
	Friend struct {
		Database *functions.Friend
	}

	FriendRequest struct {
		FriendID string `json:"userId"`
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

func (fr FriendRequest) Validate() error {
	return validation.ValidateStruct(&fr,
		// PostInHtml is not null, minLength 2, maxLength 500, no need to validate this for HTML
		validation.Field(&fr.FriendID, validation.Required, validation.Length(2, 500)),
	)
}

func (qgf QueryGetFriends) Validate() error {
	return validation.ValidateStruct(&qgf,
		validation.Field(&qgf.Limit, validation.Min(1)),
		validation.Field(&qgf.Offset, validation.Min(0)),
		validation.Field(&qgf.OnlyFriends, validation.In(true, false)),
		validation.Field(&qgf.SortBy, validation.In("friendCount", "createdAt")),
		validation.Field(&qgf.OrderBy, validation.In("asc", "desc")),
	)
}

func (f *Friend) GetFriends(ctx *fiber.Ctx) error {
	var (
		err         error
		friendsData entity.FriendData
		userID      int
	)

	// Get user ID from context
	userIDClaim := ctx.Locals("user_id").(string)
	userID, err = strconv.Atoi(userIDClaim)
	if err != nil {
		return err
	}

	// Parse query parameters
	queryParams := new(QueryGetFriends)
	if err := ctx.QueryParser(queryParams); err != nil {
		return err
	}

	err = queryParams.Validate()
	if err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	// Set default values if not provided
	if queryParams.Limit == 0 {
		queryParams.Limit = 5
	}
	if queryParams.Offset == 0 {
		queryParams.Offset = 0
	}
	if queryParams.OrderBy == "" {
		queryParams.OrderBy = "desc"
	}
	if queryParams.SortBy == "" {
		queryParams.SortBy = "createdAt"
	}

	// Call the Get method to fetch friends data
	friendsData, err = f.Database.Get(ctx.Context(), entity.QueryGetFriends{
		UserID:      userID,
		Limit:       queryParams.Limit,
		Offset:      queryParams.Offset,
		SortBy:      queryParams.SortBy,
		OrderBy:     queryParams.OrderBy,
		OnlyFriends: queryParams.OnlyFriends,
		Search:      queryParams.Search,
	})
	if err != nil {
		return err
	}

	// Send the response
	return responses.SuccessMeta(ctx, friendsData.Data, friendsData.Meta)
}

// AddFriend is a handler to add a friend
func (f *Friend) AddFriend(ctx *fiber.Ctx) error {
	var (
		req      FriendRequest
		err      error
		userID   int
		friendID int
	)

	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	// Validate request
	if req.FriendID == "" {
		return responses.ErrorBadRequest(ctx, "userId is required")
	}

	// Get user ID from context
	userIDClaim := ctx.Locals("user_id").(string)
	userID, err = strconv.Atoi(userIDClaim)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	// Convert friend ID to integer
	friendID, err = strconv.Atoi(req.FriendID)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	// Check if user is trying to add self as friend
	if userID == friendID {
		return responses.ErrorBadRequest(ctx, "Cannot add self as friend")
	}

	// Add friend
	err = f.Database.AddFriend(ctx.Context(), userID, friendID)
	if err != nil {
		if errors.Is(err, errors.New("NO_ADD_SELF")) {
			return responses.ErrorBadRequest(ctx, err.Error())
		}
		if errors.Is(err, errors.New("FRIENDSHIP_EXISTS")) {
			return responses.ErrorBadRequest(ctx, err.Error())
		}
		if errors.Is(err, errors.New("FRIEND_NOT_FOUND")) {
			return responses.ErrorNotFound(ctx, err.Error())
		}
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return responses.Success(ctx, map[string]interface{}{
		"message": "Successfully added friend",
	})
}

// DeleteFriend is a handler to delete a friend
func (f *Friend) DeleteFriend(ctx *fiber.Ctx) error {
	var (
		req      FriendRequest
		err      error
		userID   int
		friendID int
	)

	// Parse request body
	if err := ctx.BodyParser(&req); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	// Validate request
	if req.FriendID == "" {
		return responses.ErrorBadRequest(ctx, "userId is required")
	}

	// Get user ID from context
	userIDClaim := ctx.Locals("user_id").(string)
	userID, err = strconv.Atoi(userIDClaim)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	// Convert friend ID to integer
	friendID, err = strconv.Atoi(req.FriendID)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	// Delete friend
	err = f.Database.DeleteFriend(ctx.Context(), userID, friendID)
	if err != nil {
		if errors.Is(err, errors.New("NO_ADD_SELF")) {
			return responses.ErrorBadRequest(ctx, err.Error())
		}
		if errors.Is(err, errors.New("FRIENDSHIP_EXISTS")) {
			return responses.ErrorBadRequest(ctx, err.Error())
		}
		if errors.Is(err, errors.New("FRIEND_NOT_FOUND")) {
			return responses.ErrorNotFound(ctx, err.Error())
		}
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return responses.Success(ctx, map[string]interface{}{
		"message": "Successfully deleted friend",
	})
}

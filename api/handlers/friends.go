package handlers

import (
	"errors"
	"strconv"

	"segokuning/api/responses"
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
)

func (fr FriendRequest) Validate() error {
	return validation.ValidateStruct(&fr,
		// PostInHtml is not null, minLength 2, maxLength 500, no need to validate this for HTML
		validation.Field(&fr.FriendID, validation.Required, validation.Length(2, 500)),
	)
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

package handlers

import (
	"shopifyx/api/responses"
	"shopifyx/db/entity"
	"shopifyx/db/functions"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
)

type (
	Comment struct {
		PostDatabase   *functions.Post
		FriendDatabase *functions.Friend
	}

	AddCommentRequest struct {
		Comment string `json:"comment"`
		PostID  string `json:"postId"`
	}
)

func (acr AddCommentRequest) Validate() error {
	//Comment cannot be empty, minimum length is 2 maximum length is 500
	return validation.ValidateStruct(&acr,
		validation.Field(&acr.Comment, validation.Required, validation.Length(2, 500)),
		validation.Field(&acr.PostID, validation.Required),
	)
}

func (c *Comment) AddComment(ctx *fiber.Ctx) error {
	var acr AddCommentRequest
	if err := ctx.BodyParser(&acr); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	if err := acr.Validate(); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	userIDClaim := ctx.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	postID, err := strconv.Atoi(acr.PostID)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	// If post is not found return 404
	post, err := c.PostDatabase.GetByID(ctx.Context(), postID)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}
	if post.Id == 0 {
		return responses.ErrorNotFound(ctx, "Post not found")
	}

	// if post is found but not comes from the user's friend return 400
	isFriend, err := c.FriendDatabase.IsFriend(ctx.Context(), post.UserID, userID)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}
	if !isFriend {
		return responses.ErrorBadRequest(ctx, "You can only comment on your friend's post")
	}

	comment := entity.CommentPerPost{
		Comment:   acr.Comment,
		CreatedAt: time.Now(),
		Creator: entity.Creator{
			UserId: userID,
		},
	}

	comment, err = c.PostDatabase.AddComment(ctx.Context(), postID, comment)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return responses.Success(ctx, comment)
}

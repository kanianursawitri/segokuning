package handlers

import (
	"shopifyx/api/responses"
	"shopifyx/db/entity"
	"shopifyx/db/functions"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
)

type (
	Post struct {
		Database *functions.Post
	}

	AddPostRequest struct {
		PostInHtml string   `json:"postInHtml"`
		Tags       []string `json:"tags"`
	}
)

func (ap AddPostRequest) Validate() error {
	return validation.ValidateStruct(&ap,
		// PostInHtml is not null, minLength 2, maxLength 500, no need to validate this for HTML
		validation.Field(&ap.PostInHtml, validation.Required, validation.Length(2, 500)),
		// Tags is not null
		validation.Field(&ap.Tags, validation.Required),
	)
}

// AddPost is a handler to add a post
func (p *Post) AddPost(ctx *fiber.Ctx) error {
	var (
		req    AddPostRequest
		err    error
		userID int
	)

	if err := ctx.BodyParser(&req); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	userIDClaim := ctx.Locals("user_id").(string)
	userID, err = strconv.Atoi(userIDClaim)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	post := entity.Post{
		PostInHtml: req.PostInHtml,
		Tags:       req.Tags,
		UserID:     userID,
	}

	post, err = p.Database.Add(ctx.Context(), post)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return responses.Success(ctx, post)
}

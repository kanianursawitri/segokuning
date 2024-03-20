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

	QueryGetPosts struct {
		Limit      int      `query:"limit"`
		Offset     int      `query:"offset"`
		Search     string   `query:"search"`
		SearchTags []string `query:"searchTags"`
	}

	GetPostsResponse struct {
		Message string `json:"message"`
		Data    []struct {
			Post struct {
				PostInHtml string   `json:"postInHtml"`
				Tags       []string `json:"tags"`
				CreatedAt  string   `json:"createdAt"`
			} `json:"post"`
			Comments []struct {
				Comment string `json:"comment"`
				Creator struct {
					UserId      string `json:"userId"`
					Name        string `json:"name"`
					ImageUrl    string `json:"imageUrl"`
					FriendCount int    `json:"friendCount"`
					CreatedAt   string `json:"createdAt"`
				} `json:"creator"`
			} `json:"comments"`
			Creator struct {
				UserId      string `json:"userId"`
				Name        string `json:"name"`
				ImageUrl    string `json:"imageUrl"`
				FriendCount int    `json:"friendCount"`
				CreatedAt   string `json:"createdAt"`
			} `json:"creator"`
		} `json:"data"`
		Meta struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
			Total  int `json:"total"`
		} `json:"meta"`
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

func (qgp QueryGetPosts) Validate() error {
	return validation.ValidateStruct(&qgp,
		// Limit is optional, default 5
		validation.Field(&qgp.Limit, validation.Min(1)),
		// Offset is optional, default 0
		validation.Field(&qgp.Offset, validation.Min(0)),
	)
}

func (qgp QueryGetPosts) ToEntity() entity.QueryGetPosts {
	return entity.QueryGetPosts{
		Limit:      qgp.Limit,
		Offset:     qgp.Offset,
		Search:     qgp.Search,
		SearchTags: qgp.SearchTags,
	}
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

// GetPosts is a handler to get posts
func (p *Post) GetPosts(ctx *fiber.Ctx) error {
	var (
		req QueryGetPosts
		err error
	)

	if err := ctx.QueryParser(&req); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	if req.Limit == 0 {
		req.Limit = 5
	}

	posts, err := p.Database.Get(ctx.Context(), req.ToEntity())
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return responses.Success(ctx, posts)
}

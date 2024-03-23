package handlers

import (
	"segokuning/api/responses"
	"segokuning/db/entity"
	"segokuning/db/functions"
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

	PostData struct {
		PostInHtml string   `json:"postInHtml"`
		Tags       []string `json:"tags"`
		CreatedAt  string   `json:"createdAt"`
	}

	Creator struct {
		UserId      string  `json:"userId"`
		Name        string  `json:"name"`
		ImageUrl    *string `json:"imageUrl"`
		FriendCount int     `json:"friendCount"`
	}

	CreatorPost struct {
		Creator
		CreatedAt string `json:"createdAt"`
	}
	CommentPerPost struct {
		Comment   string  `json:"comment"`
		Creator   Creator `json:"creator"`
		CreatedAt string  `json:"createdAt"`
	}

	ElemData struct {
		PostId   int              `json:"postId"`
		Post     PostData         `json:"post"`
		Comments []CommentPerPost `json:"comments"`
		Creator  CreatorPost      `json:"creator"`
	}

	Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	GetPostsResponse struct {
		Message string     `json:"message"`
		Data    []ElemData `json:"data"`
		Meta    Meta       `json:"meta"`
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

func (qgp QueryGetPosts) ToEntity(userID int) entity.QueryGetPosts {
	return entity.QueryGetPosts{
		UserId:     userID,
		Limit:      qgp.Limit,
		Offset:     qgp.Offset,
		Search:     qgp.Search,
		SearchTags: qgp.SearchTags,
	}
}

func (p *Post) convertEntityPostsToResponse(posts []entity.Post) []ElemData {
	var elemData []ElemData
	for _, post := range posts {
		var comments []CommentPerPost
		for _, comment := range post.Comments {
			comments = append(comments, CommentPerPost{
				Comment:   comment.Comment,
				Creator:   Creator{UserId: strconv.Itoa(comment.Creator.UserId), Name: comment.Creator.Name, ImageUrl: comment.Creator.ImageUrl, FriendCount: comment.Creator.FriendCount},
				CreatedAt: comment.CreatedAt.String(),
			})
		}

		elemData = append(elemData, ElemData{
			PostId: post.Id,
			Post: PostData{
				PostInHtml: post.PostInHtml,
				Tags:       post.Tags,
				CreatedAt:  post.CreatedAt.String(),
			},
			Comments: comments,
			Creator: CreatorPost{
				Creator: Creator{
					UserId:      strconv.Itoa(post.Creator.UserId),
					Name:        post.Creator.Name,
					ImageUrl:    post.Creator.ImageUrl,
					FriendCount: post.Creator.FriendCount,
				},
				CreatedAt: post.CreatedAt.String(),
			},
		})
	}

	return elemData
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

	userIDClaim := ctx.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	filter := req.ToEntity(userID)
	posts, err := p.Database.Get(ctx.Context(), filter)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	count, err := p.Database.Count(ctx.Context(), filter)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	response := GetPostsResponse{
		Message: "Success",
		Data:    p.convertEntityPostsToResponse(posts),
		Meta: Meta{
			Limit:  req.Limit,
			Offset: req.Offset,
			Total:  count,
		},
	}

	return responses.Success(ctx, response)
}

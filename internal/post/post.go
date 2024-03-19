package post

import (
	"context"

	"github.com/rrab-0/its-gram/internal"
)

type PostIdUriRequest struct {
	PostId string `uri:"id" binding:"required"`
}

type CreatePostRequest struct {
	PictureLink string `json:"picture_link" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type Repository interface {
	CreatePost(ctx context.Context, userId string, post internal.Post) (internal.Post, error)
	GetUserPosts(ctx context.Context, userId string) ([]internal.Post, error)
	GetPostById(ctx context.Context, id string) (internal.Post, error)
}

type Service interface {
	CreatePost(ctx context.Context, reqUri internal.UserIdUriRequest, reqBody CreatePostRequest) (internal.Post, error)
	GetUserPosts(ctx context.Context, reqUri internal.UserIdUriRequest) ([]internal.Post, error)
	GetPostById(ctx context.Context, reqUri PostIdUriRequest) (internal.Post, error)
}

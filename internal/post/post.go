package post

import (
	"context"

	"github.com/google/uuid"
	"github.com/rrab-0/its-gram/internal"
)

type PostIdUriRequest struct {
	PostId string `uri:"id" binding:"required,uuid"`
}

type CreatePostRequest struct {
	PictureLink string `json:"picture_link" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type PostAndUserUriRequest struct {
	UserId string `uri:"id" binding:"required"`
	PostId string `uri:"postId" binding:"required,uuid"`
}

type Repository interface {
	GetUserPosts(ctx context.Context, userId string) ([]internal.Post, error)
	GetPostById(ctx context.Context, id string) (internal.Post, error)
	CreatePost(ctx context.Context, userId string, post internal.Post) (internal.Post, error)
	DeletePost(ctx context.Context, userId string, postId uuid.UUID) error
	LikePost(ctx context.Context, userId string, postId uuid.UUID) error
	UnlikePost(ctx context.Context, userId string, postId uuid.UUID) error
	CommentPost(ctx context.Context, userId string, postId uuid.UUID) error
	UncommentPost(ctx context.Context, userId string, postId uuid.UUID) error
}

type Service interface {
	CreatePost(ctx context.Context, reqUri internal.UserIdUriRequest, reqBody CreatePostRequest) (internal.Post, error)
	GetUserPosts(ctx context.Context, reqUri internal.UserIdUriRequest) ([]internal.Post, error)
	GetPostById(ctx context.Context, reqUri PostIdUriRequest) (internal.Post, error)
	DeletePost(ctx context.Context, reqUri PostAndUserUriRequest) error
	LikePost(ctx context.Context, reqUri PostAndUserUriRequest) error
	UnlikePost(ctx context.Context, reqUri PostAndUserUriRequest) error
	CommentPost(ctx context.Context, reqUri PostAndUserUriRequest) error
	UncommentPost(ctx context.Context, reqUri PostAndUserUriRequest) error
}

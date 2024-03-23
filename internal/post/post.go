package post

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrab-0/its-gram/internal"
	"gorm.io/gorm"
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

type CreateCommentRequest struct {
	Description string `json:"description" binding:"required"`
}

type CommentAndUserUriRequest struct {
	UserId    string `uri:"id" binding:"required"`
	CommentId string `uri:"commentId" binding:"required,uuid"`
}

type ReplyCommentRequest struct {
	UserId    string `uri:"id" binding:"required"`
	PostId    string `uri:"postId" binding:"required,uuid"`
	CommentId string `uri:"commentId" binding:"required,uuid"`
}

type GetCommentRequest struct {
	CommentId string `uri:"commentId" binding:"required,uuid"`
}

type GetCommentResponse struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	CreatedBy   internal.User   `json:"created_by"`
	Description string          `json:"description"`
	Likes       []internal.User `json:"likes"`
	Replies     []interface{}   `json:"replies"`
}

type Repository interface {
	GetPostById(ctx context.Context, id string) (post internal.Post, totalComments int, err error)
	CreatePost(ctx context.Context, userId string, post internal.Post) (internal.Post, error)
	DeletePost(ctx context.Context, userId string, postId uuid.UUID) error
	LikePost(ctx context.Context, userId string, postId uuid.UUID) error
	UnlikePost(ctx context.Context, userId string, postId uuid.UUID) error

	GetComment(ctx context.Context, commentId uuid.UUID) (internal.Comment, error)
	CommentPost(ctx context.Context, userId, description string, postId uuid.UUID) error
	UncommentPost(ctx context.Context, userId string, commentId uuid.UUID) error
	ReplyComment(ctx context.Context, userId, description string, postId, commentId uuid.UUID) error
	RemoveReplyFromComment(ctx context.Context, userId string, commentId uuid.UUID) error
	LikeComment(ctx context.Context, userId string, commentId uuid.UUID) error
	UnlikeComment(ctx context.Context, userId string, commentId uuid.UUID) error
}

type Service interface {
	GetPostById(ctx context.Context, reqUri PostIdUriRequest) (post internal.Post, totalComments int, err error)
	CreatePost(ctx context.Context, reqUri internal.UserIdUriRequest, reqBody CreatePostRequest) (internal.Post, error)
	DeletePost(ctx context.Context, reqUri PostAndUserUriRequest) error
	LikePost(ctx context.Context, reqUri PostAndUserUriRequest) error
	UnlikePost(ctx context.Context, reqUri PostAndUserUriRequest) error

	GetComment(ctx context.Context, reqUri GetCommentRequest) (internal.Comment, error)
	CommentPost(ctx context.Context, reqUri PostAndUserUriRequest, reqBody CreateCommentRequest) error
	UncommentPost(ctx context.Context, reqUri CommentAndUserUriRequest) error
	ReplyComment(ctx context.Context, reqUri ReplyCommentRequest, reqBody CreateCommentRequest) error
	RemoveReplyFromComment(ctx context.Context, reqUri CommentAndUserUriRequest) error
	LikeComment(ctx context.Context, reqUri CommentAndUserUriRequest) error
	UnlikeComment(ctx context.Context, reqUri CommentAndUserUriRequest) error
}

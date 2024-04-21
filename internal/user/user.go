package user

import (
	"context"

	"github.com/rrab-0/its-gram/internal"
)

type UpdateUserProfileRequest struct {
	Username    string `json:"username" binding:"required"`
	PictureLink string `json:"picture_link" binding:"required"`
	Description string `json:"description"`
}

type UserSearchRequest struct {
	Username string `form:"username" binding:"required"`
}

type FollowOtherUserRequest struct {
	UserId      string `uri:"id" binding:"required"`
	OtherUserId string `uri:"otherUserId" binding:"required"`
}

type GetLikesResponse struct {
	Likes []any `json:"likes"`
}

type GetLikesCommentQueryRes struct {
	Type    string           `json:"type"`
	Comment internal.Comment `json:"comment"`
}

type GetLikesPostQueryRes struct {
	Type string        `json:"type"`
	Post internal.Post `json:"post"`
}

type GetHomepageQueryRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type GetHomepageQueryRes struct {
	TotalPage int             `json:"total_page"`
	Posts     []internal.Post `json:"posts"`
}

type GetUserHomepageCursorQueryRes struct {
	NextCursor string          `json:"next_cursor"`
	Posts      []internal.Post `json:"posts"`
}

type GetUserHomepageCursorQueryRequest struct {
	Cursor string `form:"cursor" binding:"required"`
	Limit  int    `form:"limit" binding:"required"`
}

type GetUserHomepageInitialCursorQueryRequest struct {
	Limit int `form:"limit" binding:"required"`
}

type Repository interface {
	GetUser(ctx context.Context, id string) (internal.User, error)
	SearchUser(ctx context.Context, username string) ([]internal.User, error)
	GetUserHomepage(ctx context.Context, page, limit int, id string) (GetHomepageQueryRes, error)
	GetUserHomepageInitialCursor(ctx context.Context, limit int, id string) (*GetUserHomepageCursorQueryRes, error)
	GetUserHomepageCursor(ctx context.Context, cursor string, limit int, id string) (*GetUserHomepageCursorQueryRes, error)

	CreateUser(ctx context.Context, user internal.User) (internal.User, error)
	UpdateUserProfile(ctx context.Context, id, username, picture, description string) (internal.User, error)
	DeleteUser(ctx context.Context, id string) (internal.User, error)

	FollowOtherUser(ctx context.Context, userId, otherUserId string) error
	UnfollowOtherUser(ctx context.Context, userId, otherUserId string) error

	GetLikes(ctx context.Context, userId string) ([]any, error)
	GetPosts(ctx context.Context, userId string) ([]internal.Post, []int, error)
	GetComments(ctx context.Context, userId string) ([]internal.Comment, error)
}

type Service interface {
	GetUser(ctx context.Context, reqUri internal.UserIdUriRequest) (internal.User, error)
	SearchUser(ctx context.Context, reqQuery UserSearchRequest) ([]internal.User, error)
	GetUserHomepage(ctx context.Context, reqUri internal.UserIdUriRequest, reqQuery GetHomepageQueryRequest) (GetHomepageQueryRes, error)
	GetUserHomepageInitialCursor(ctx context.Context, reqUri internal.UserIdUriRequest, reqQuery GetUserHomepageInitialCursorQueryRequest) (*GetUserHomepageCursorQueryRes, error)
	GetUserHomepageCursor(ctx context.Context, reqUri internal.UserIdUriRequest, reqQuery GetUserHomepageCursorQueryRequest) (*GetUserHomepageCursorQueryRes, error)

	CreateUser(ctx context.Context, firebaseId, username, email, picture string) (internal.User, error)
	UpdateUserProfile(ctx context.Context, reqUri internal.UserIdUriRequest, reqBody UpdateUserProfileRequest) (internal.User, error)
	DeleteUser(ctx context.Context, reqUri internal.UserIdUriRequest) (internal.User, error)

	FollowOtherUser(ctx context.Context, reqUri FollowOtherUserRequest) error
	UnfollowOtherUser(ctx context.Context, reqUri FollowOtherUserRequest) error

	GetLikes(ctx context.Context, reqUri internal.UserIdUriRequest) ([]any, error)
	GetPosts(ctx context.Context, reqUri internal.UserIdUriRequest) ([]internal.Post, []int, error)
	GetComments(ctx context.Context, reqUri internal.UserIdUriRequest) ([]internal.Comment, error)
}

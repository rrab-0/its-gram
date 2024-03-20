package user

import (
	"context"

	"github.com/rrab-0/its-gram/internal"
)

type UpdateUserProfileRequest struct {
	Username    string `json:"username" binding:"required"`
	PictureLink string `json:"picture_link" binding:"required"`
}

type UserSearchRequest struct {
	Username string `form:"username" binding:"required"`
}

type FollowOtherUserRequest struct {
	UserId      string `uri:"id" binding:"required"`
	OtherUserId string `uri:"otherUserId" binding:"required"`
}

type Repository interface {
	GetUser(ctx context.Context, id string) (internal.User, error)
	SearchUser(ctx context.Context, username string) ([]internal.User, error)
	GetUserHomepage(ctx context.Context, id string) ([]internal.Post, error)

	CreateUser(ctx context.Context, user internal.User) (internal.User, error)
	UpdateUserProfile(ctx context.Context, id, username, picture string) (internal.User, error)
	DeleteUser(ctx context.Context, id string) (internal.User, error)

	FollowOtherUser(ctx context.Context, userId, otherUserId string) error
	UnfollowOtherUser(ctx context.Context, userId, otherUserId string) error

	GetLikes(ctx context.Context, userId string) ([]internal.Post, error)
}

type Service interface {
	GetUser(ctx context.Context, reqUri internal.UserIdUriRequest) (internal.User, error)
	SearchUser(ctx context.Context, reqQuery UserSearchRequest) ([]internal.User, error)
	GetUserHomepage(ctx context.Context, reqUri internal.UserIdUriRequest) ([]internal.Post, error)

	CreateUser(ctx context.Context, firebaseId, username, email, picture string) (internal.User, error)
	UpdateUserProfile(ctx context.Context, reqUri internal.UserIdUriRequest, reqBody UpdateUserProfileRequest) (internal.User, error)
	DeleteUser(ctx context.Context, reqUri internal.UserIdUriRequest) (internal.User, error)

	FollowOtherUser(ctx context.Context, reqUri FollowOtherUserRequest) error
	UnfollowOtherUser(ctx context.Context, reqUri FollowOtherUserRequest) error

	GetLikes(ctx context.Context, reqUri internal.UserIdUriRequest) ([]internal.Post, error)
}

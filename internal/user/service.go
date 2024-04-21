package user

import (
	"context"

	"github.com/rrab-0/its-gram/internal"
)

type userService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return userService{
		repo: repo,
	}
}

func (s userService) CreateUser(ctx context.Context, firebaseId, username, email, picture string) (internal.User, error) {
	user := internal.User{
		ID:          firebaseId,
		Username:    username,
		Email:       email,
		PictureLink: picture,
	}

	user, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return internal.User{}, err
	}

	return user, nil
}

func (s userService) GetUser(ctx context.Context, reqUri internal.UserIdUriRequest) (internal.User, error) {
	user, err := s.repo.GetUser(ctx, reqUri.UserId)
	if err != nil {
		return internal.User{}, err
	}

	return user, nil
}

func (s userService) SearchUser(ctx context.Context, reqQuery UserSearchRequest) ([]internal.User, error) {
	users, err := s.repo.SearchUser(ctx, reqQuery.Username)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s userService) GetUserHomepage(ctx context.Context, reqUri internal.UserIdUriRequest, reqQuery GetHomepageQueryRequest) (GetHomepageQueryRes, error) {
	posts, err := s.repo.GetUserHomepage(ctx, reqQuery.Page, reqQuery.Limit, reqUri.UserId)
	if err != nil {
		return GetHomepageQueryRes{}, err
	}

	return posts, nil
}

func (s userService) GetUserHomepageInitialCursor(ctx context.Context, reqUri internal.UserIdUriRequest, reqQuery GetUserHomepageInitialCursorQueryRequest) (*GetUserHomepageCursorQueryRes, error) {
	posts, err := s.repo.GetUserHomepageInitialCursor(ctx, reqQuery.Limit, reqUri.UserId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s userService) GetUserHomepageCursor(ctx context.Context, reqUri internal.UserIdUriRequest, reqQuery GetUserHomepageCursorQueryRequest) (*GetUserHomepageCursorQueryRes, error) {
	posts, err := s.repo.GetUserHomepageCursor(ctx, reqQuery.Cursor, reqQuery.Limit, reqUri.UserId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s userService) UpdateUserProfile(ctx context.Context, reqUri internal.UserIdUriRequest, reqBody UpdateUserProfileRequest) (internal.User, error) {
	user, err := s.repo.UpdateUserProfile(ctx, reqUri.UserId, reqBody.Username, reqBody.PictureLink, reqBody.Description)
	if err != nil {
		return internal.User{}, err
	}

	return user, nil
}

func (s userService) DeleteUser(ctx context.Context, reqUri internal.UserIdUriRequest) (internal.User, error) {
	user, err := s.repo.DeleteUser(ctx, reqUri.UserId)
	if err != nil {
		return internal.User{}, err
	}

	return user, nil
}

func (s userService) FollowOtherUser(ctx context.Context, reqUri FollowOtherUserRequest) error {
	return s.repo.FollowOtherUser(ctx, reqUri.UserId, reqUri.OtherUserId)
}

func (s userService) UnfollowOtherUser(ctx context.Context, reqUri FollowOtherUserRequest) error {
	return s.repo.UnfollowOtherUser(ctx, reqUri.UserId, reqUri.OtherUserId)
}

func (s userService) GetPosts(ctx context.Context, reqUri internal.UserIdUriRequest) ([]internal.Post, []int, error) {
	posts, totalComments, err := s.repo.GetPosts(ctx, reqUri.UserId)
	if err != nil {
		return nil, nil, err
	}

	return posts, totalComments, nil
}

func (s userService) GetLikes(ctx context.Context, reqUri internal.UserIdUriRequest) ([]any, error) {
	likes, err := s.repo.GetLikes(ctx, reqUri.UserId)
	if err != nil {
		return nil, err
	}

	return likes, nil
}

func (s userService) GetComments(ctx context.Context, reqUri internal.UserIdUriRequest) ([]internal.Comment, error) {
	comments, err := s.repo.GetComments(ctx, reqUri.UserId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

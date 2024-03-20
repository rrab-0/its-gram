package post

import (
	"context"

	"github.com/google/uuid"
	"github.com/rrab-0/its-gram/internal"
)

type postService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return postService{
		repo: repo,
	}
}

func (s postService) CreatePost(ctx context.Context, reqUri internal.UserIdUriRequest, reqBody CreatePostRequest) (internal.Post, error) {
	post := internal.Post{
		PictureLink: reqBody.PictureLink,
		Title:       reqBody.Title,
		Description: reqBody.Description,
	}

	post, err := s.repo.CreatePost(ctx, reqUri.UserId, post)
	if err != nil {
		return internal.Post{}, err
	}

	return post, nil
}

func (s postService) GetUserPosts(ctx context.Context, reqUri internal.UserIdUriRequest) ([]internal.Post, error) {
	posts, err := s.repo.GetUserPosts(ctx, reqUri.UserId)
	if err != nil {
		return []internal.Post{}, err
	}

	return posts, nil
}

func (s postService) GetPostById(ctx context.Context, reqUri PostIdUriRequest) (internal.Post, error) {
	post, err := s.repo.GetPostById(ctx, reqUri.PostId)
	if err != nil {
		return internal.Post{}, err
	}

	return post, nil
}

func (s postService) DeletePost(ctx context.Context, reqUri PostAndUserUriRequest) error {
	postId, _ := uuid.Parse(reqUri.PostId)
	return s.repo.DeletePost(ctx, reqUri.UserId, postId)
}

func (s postService) LikePost(ctx context.Context, reqUri PostAndUserUriRequest) error {
	postId, _ := uuid.Parse(reqUri.PostId)
	return s.repo.LikePost(ctx, reqUri.UserId, postId)
}

func (s postService) UnlikePost(ctx context.Context, reqUri PostAndUserUriRequest) error {
	postId, _ := uuid.Parse(reqUri.PostId)
	return s.repo.LikePost(ctx, reqUri.UserId, postId)
}

func (s postService) CommentPost(ctx context.Context, reqUri PostAndUserUriRequest) error {
	postId, _ := uuid.Parse(reqUri.PostId)
	return s.repo.LikePost(ctx, reqUri.UserId, postId)
}

func (s postService) UncommentPost(ctx context.Context, reqUri PostAndUserUriRequest) error {
	postId, _ := uuid.Parse(reqUri.PostId)
	return s.repo.LikePost(ctx, reqUri.UserId, postId)
}

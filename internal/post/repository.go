package post

import (
	"context"

	"github.com/google/uuid"
	"github.com/rrab-0/its-gram/internal"
	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return gormRepository{
		db: db,
	}
}

func (r gormRepository) CreatePost(ctx context.Context, userId string, post internal.Post) (internal.Post, error) {
	var (
		user internal.User
		tx   = r.db.WithContext(ctx).Begin()
	)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Where("id = ?", userId).First(&user).Error
	if err != nil {
		tx.Rollback()
		return internal.Post{}, err
	}

	post.CreatedBy = user
	post.UserID = user.ID
	err = tx.Create(&post).Error
	if err != nil {
		tx.Rollback()
		return internal.Post{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return internal.Post{}, err
	}

	return post, nil
}

func (r gormRepository) GetUserPosts(ctx context.Context, userId string) ([]internal.Post, error) {
	var (
		post []internal.Post
		user internal.User
	)

	err := r.db.WithContext(ctx).Preload("Posts").Where("id = ?", userId).First(&user).Error
	if err != nil {
		return []internal.Post{}, err
	}

	post = user.Posts
	return post, nil
}

func (r gormRepository) GetPostById(ctx context.Context, id string) (internal.Post, error) {
	var post internal.Post

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return internal.Post{}, err
	}

	return post, nil
}

func (r gormRepository) DeletePost(ctx context.Context, userId string, postId uuid.UUID) error {
	var (
		user internal.User
		post internal.Post
	)

	user.ID = userId
	post.ID = postId
	return r.db.Model(&user).Association("Posts").Unscoped().Delete(&post)
}

func (r gormRepository) LikePost(ctx context.Context, userId string, postId uuid.UUID) error {
	var (
		user internal.User
		post internal.Post
	)

	user.ID = userId
	post.ID = postId
	return r.db.Model(&user).Association("LikedPosts").Append(&post)
}

func (r gormRepository) UnlikePost(ctx context.Context, userId string, postId uuid.UUID) error {
	var (
		user internal.User
		post internal.Post
	)

	user.ID = userId
	post.ID = postId
	return r.db.Model(&user).Association("LikedPosts").Delete(&post)
}

func (r gormRepository) CommentPost(ctx context.Context, userId string, postId uuid.UUID) error {
	var (
		user    internal.User
		post    internal.Post
		comment internal.Comment
		tx      = r.db.WithContext(ctx).Begin()
	)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Where("id = ?", userId).First(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Where("id = ?", postId).First(&post).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	comment.CreatedBy = user
	comment.UserID = user.ID

	comment.CreatedIn = post
	comment.PostCreatedInID = post.ID

	err = tx.Create(&comment).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

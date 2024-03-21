package post

import (
	"context"

	"github.com/google/uuid"
	"github.com/rrab-0/its-gram/internal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r gormRepository) GetPostById(ctx context.Context, id string) (internal.Post, error) {
	var post internal.Post

	err := r.db.
		WithContext(ctx).
		Unscoped().
		Preload(clause.Associations).
		Preload("Comments", "parent_id IS NULL").
		Preload("Comments.CreatedBy").
		Where("id = ?", id).
		First(&post).
		Error
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
	return r.db.WithContext(ctx).Model(&user).Association("Posts").Unscoped().Delete(&post)
}

func (r gormRepository) LikePost(ctx context.Context, userId string, postId uuid.UUID) error {
	var (
		user internal.User
		post internal.Post
	)

	user.ID = userId
	post.ID = postId
	return r.db.WithContext(ctx).Model(&user).Association("LikedPosts").Append(&post)
}

func (r gormRepository) UnlikePost(ctx context.Context, userId string, postId uuid.UUID) error {
	var (
		user internal.User
		post internal.Post
	)

	user.ID = userId
	post.ID = postId
	return r.db.WithContext(ctx).Model(&user).Association("LikedPosts").Delete(&post)
}

func (r gormRepository) GetComment(ctx context.Context, commentId uuid.UUID) (internal.Comment, error) {
	var (
		comment internal.Comment
	)

	err := r.db.WithContext(ctx).Unscoped().Preload("Replies").Where("id = ?", commentId).First(&comment).Error
	if err != nil {
		return internal.Comment{}, err
	}

	return comment, nil
}

func (r gormRepository) CommentPost(ctx context.Context, userId, description string, postId uuid.UUID) error {
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

	comment.PostID = post.ID

	comment.Description = description

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

func (r gormRepository) UncommentPost(ctx context.Context, userId string, commentId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", commentId).Delete(&internal.Comment{}).Error
}

func (r gormRepository) ReplyComment(ctx context.Context, userId, description string, commentId uuid.UUID) error {
	var (
		user       internal.User
		comment    internal.Comment
		newComment internal.Comment
		tx         = r.db.WithContext(ctx).Begin()
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

	err = tx.Where("id = ?", commentId).First(&comment).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	newComment.CreatedBy = user
	newComment.UserID = user.ID

	newComment.CreatedIn = comment.CreatedIn
	newComment.PostCreatedInID = comment.PostCreatedInID

	newComment.PostID = comment.PostID

	newComment.Description = description
	newComment.ParentID = &comment.ID

	err = tx.Create(&newComment).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&comment).Association("Replies").Append(&newComment)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r gormRepository) RemoveReplyFromComment(ctx context.Context, userId string, commentId uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", commentId).Delete(&internal.Comment{}).Error
}

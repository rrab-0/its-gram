package user

import (
	"context"

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

func (r gormRepository) CreateUser(ctx context.Context, user internal.User) (internal.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return internal.User{}, err
	}

	return user, nil
}

func (r gormRepository) GetUser(ctx context.Context, id string) (internal.User, error) {
	var user internal.User
	err := r.db.WithContext(ctx).Preload(clause.Associations).Preload("Posts.CreatedBy").Where("id = ?", id).First(&user).Error
	if err != nil {
		return internal.User{}, err
	}

	return user, nil
}

func (r gormRepository) SearchUser(ctx context.Context, username string) ([]internal.User, error) {
	var users []internal.User
	err := r.db.WithContext(ctx).Where("username LIKE ?", "%"+username+"%").Find(&users).Error
	if err != nil {
		return []internal.User{}, err
	}

	return users, nil
}

func (r gormRepository) GetUserHomepage(ctx context.Context, id string) ([]internal.Post, error) {
	var posts []internal.Post
	var user internal.User

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Model(&user).
		Select("Followers.Posts").
		Find(&posts).
		Error
	if err != nil {
		return []internal.Post{}, err
	}

	return posts, nil
}

func (r gormRepository) UpdateUserProfile(ctx context.Context, id, username, picture string) (internal.User, error) {
	var user internal.User
	user.Username = username
	user.PictureLink = picture

	err := r.db.
		WithContext(ctx).
		Clauses(clause.Returning{}).
		Model(&user).
		Where("id = ?", id).
		Updates(user).
		Error
	if err != nil {
		return internal.User{}, err
	}

	return user, nil
}

// Remove user's references of:
// - LikedPosts
// - Followers
// - Followings
// Remove user's actual:
// - Posts
func (r gormRepository) DeleteUser(ctx context.Context, id string) (internal.User, error) {
	var (
		user internal.User
		tx   = r.db.WithContext(ctx).Begin()
	)

	err := tx.Model(&user).Where("id = ?", id).Association("Followers").Clear()
	if err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	err = tx.Model(&user).Where("id = ?", id).Association("Followings").Clear()
	if err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	err = tx.Model(&user).Where("id = ?", id).Association("LikedPosts").Clear()
	if err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	err = tx.Unscoped().Model(&user).Where("id = ?", id).Association("Posts").Unscoped().Clear()
	if err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	if err := tx.Unscoped().Where("id = ?", id).Delete(&user).Error; err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return internal.User{}, err
	}

	return user, nil
}

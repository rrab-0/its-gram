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
	var (
		followingsPosts []internal.Post
		user            internal.User
	)

	user.ID = id

	res := r.db.
		WithContext(ctx).
		Preload("Followings.Posts").
		First(&user)
	if res.Error != nil {
		return []internal.Post{}, res.Error
	}

	if res.RowsAffected == 0 || len(user.Followings) == 0 {
		return []internal.Post{}, gorm.ErrRecordNotFound
	}

	for _, following := range user.Followings {
		followingsPosts = append(followingsPosts, following.Posts...)
	}
	return followingsPosts, nil
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
// Then hard delete the user
func (r gormRepository) DeleteUser(ctx context.Context, id string) (internal.User, error) {
	var (
		user internal.User
		tx   = r.db.WithContext(ctx).Begin()
	)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user.ID = id

	err := tx.Model(&user).Association("Followers").Clear()
	if err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	err = tx.Model(&user).Association("Followings").Clear()
	if err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	err = tx.Model(&user).Association("LikedPosts").Clear()
	if err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	err = tx.Unscoped().Model(&user).Association("Posts").Unscoped().Clear()
	if err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	if err := tx.Unscoped().Delete(&user).Error; err != nil {
		tx.Rollback()
		return internal.User{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return internal.User{}, err
	}

	return user, nil
}

func (r gormRepository) FollowOtherUser(ctx context.Context, userId, otherUserId string) error {
	var (
		user      internal.User
		otherUser internal.User
		tx        = r.db.WithContext(ctx).Begin()
	)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user.ID = userId
	otherUser.ID = otherUserId

	if err := tx.Model(&user).Association("Followings").Append(&otherUser); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&otherUser).Association("Followers").Append(&user); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r gormRepository) UnfollowOtherUser(ctx context.Context, userId, otherUserId string) error {
	var (
		user      internal.User
		otherUser internal.User
		tx        = r.db.WithContext(ctx).Begin()
	)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user.ID = userId
	otherUser.ID = otherUserId

	if err := tx.Model(&user).Association("Followings").Delete(&otherUser); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&otherUser).Association("Followers").Delete(&user); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r gormRepository) GetLikes(ctx context.Context, userId string) ([]internal.Post, error) {
	var (
		user       internal.User
		likedPosts []internal.Post
	)

	user.ID = userId

	if err := r.db.Model(&user).Association("LikedPosts").Find(&likedPosts); err != nil {
		return []internal.Post{}, err
	}

	return likedPosts, nil
}

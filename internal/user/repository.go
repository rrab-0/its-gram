package user

import (
	"context"
	"math"
	"net/url"
	"sort"
	"time"

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
	err := r.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "email"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"id":           user.ID,
				"created_at":   time.Now(),
				"updated_at":   time.Now(),
				"deleted_at":   nil,
				"username":     user.Username,
				"email":        user.Email,
				"picture_link": user.PictureLink,
			}),
		}).
		Create(&user).
		Error
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

func (r gormRepository) GetUserHomepage(ctx context.Context, page, limit int, id string) (GetHomepageQueryRes, error) {
	var (
		followingsPosts GetHomepageQueryRes
		user            internal.User
		tx              = r.db.WithContext(ctx).Begin()
		totalPosts      int
	)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get total posts to validate page request
	user.ID = id
	err := tx.Preload("Followings.Posts").First(&user).Error
	if err != nil {
		tx.Rollback()
		return GetHomepageQueryRes{}, err
	}

	for _, following := range user.Followings {
		totalPosts += len(following.Posts)
	}

	totalPage := math.Ceil(float64(totalPosts) / float64(limit))
	if page > int(totalPage) {
		page = int(totalPage)
	}

	// Get offset from page and limit then query posts
	offset := (page - 1) * limit
	res := tx.
		Preload("Followings.Posts", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("created_at DESC").
				Offset(offset).
				Limit(limit)
		}).
		Preload("Followings.Posts.CreatedBy").
		Preload("Followings.Posts.Likes").
		Preload("Followings.Posts.Comments").
		First(&user)
	if res.Error != nil {
		tx.Rollback()
		return GetHomepageQueryRes{}, res.Error
	}

	if err := tx.Commit().Error; err != nil {
		return GetHomepageQueryRes{}, err
	}

	if res.RowsAffected == 0 || len(user.Followings) == 0 {
		return GetHomepageQueryRes{}, gorm.ErrRecordNotFound
	}

	for _, following := range user.Followings {
		followingsPosts.Posts = append(followingsPosts.Posts, following.Posts...)
	}

	followingsPosts.TotalPage = int(totalPage)
	return followingsPosts, nil
}

func (r gormRepository) GetUserHomepageInitialCursor(ctx context.Context, limit int, id string) (*GetUserHomepageCursorQueryRes, error) {
	var (
		user  internal.User
		posts []internal.Post
	)

	user.ID = id
	res := r.db.
		WithContext(ctx).
		Preload("Followings.Posts", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("created_at DESC").
				Limit(limit)
		}).
		Preload("Followings.Posts.CreatedBy").
		Preload("Followings.Posts.Likes").
		Preload("Followings.Posts.Comments").
		First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 || len(user.Followings) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	for _, following := range user.Followings {
		posts = append(posts, following.Posts...)
	}

	return &GetUserHomepageCursorQueryRes{
		NextCursor: url.QueryEscape(posts[len(posts)-1].CreatedAt.Format(time.RFC3339Nano)),
		Posts:      posts,
	}, nil
}

func (r gormRepository) GetUserHomepageCursor(ctx context.Context, cursor string, limit int, id string) (*GetUserHomepageCursorQueryRes, error) {
	var (
		user  internal.User
		posts []internal.Post
	)

	cursorTime, err := time.Parse(time.RFC3339Nano, cursor)
	if err != nil {
		return nil, err
	}

	user.ID = id
	res := r.db.WithContext(ctx).
		Preload("Followings.Posts", func(db *gorm.DB) *gorm.DB {
			return db.
				Where("created_at < ?", cursorTime).
				Order("created_at DESC").
				Limit(limit)
		}).
		Preload("Followings.Posts.CreatedBy").
		Preload("Followings.Posts.Likes").
		Preload("Followings.Posts.Comments").
		First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 || len(user.Followings) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	for _, following := range user.Followings {
		posts = append(posts, following.Posts...)
	}

	return &GetUserHomepageCursorQueryRes{
		NextCursor: url.QueryEscape(posts[len(posts)-1].CreatedAt.Format(time.RFC3339Nano)),
		Posts:      posts,
	}, nil
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
// Then soft delete the user
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

	if err := tx.Delete(&user).Error; err != nil {
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

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
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

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r gormRepository) GetPosts(ctx context.Context, userId string) ([]internal.Post, []int, error) {
	var (
		totalComments []int
		posts         []internal.Post
		tx            = r.db.WithContext(ctx).Begin()
	)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Unscoped().Preload("Comments").Where("user_id = ?", userId).Find(&posts).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	for _, post := range posts {
		totalComments = append(totalComments, len(post.Comments))
	}

	posts = nil

	err = tx.
		Unscoped().
		Preload(clause.Associations).
		Preload("Comments", "parent_id IS NULL").
		Preload("Comments.CreatedBy").
		Preload("Comments.Likes").
		Preload("Comments.Replies").
		Where("user_id = ?", userId).
		Find(&posts).
		Error

	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	return posts, totalComments, nil
}

func getCreatedAt(like interface{}) time.Time {
	switch v := like.(type) {
	case internal.Post:
		return v.CreatedAt
	case internal.Comment:
		return v.CreatedAt
	default:
		return time.Time{}
	}
}

func (r gormRepository) GetLikes(ctx context.Context, userId string) ([]any, error) {
	var (
		user  internal.User
		likes []any
	)

	user.ID = userId

	err := r.db.
		WithContext(ctx).
		Model(&user).
		Preload("LikedPosts", func(db *gorm.DB) *gorm.DB {
			return db.Preload(clause.Associations)
		}).
		Preload("LikedComments", func(db *gorm.DB) *gorm.DB {
			return db.Preload(clause.Associations)
		}).
		First(&user).
		Error
	if err != nil {
		return []any{}, err
	}

	for _, comment := range user.LikedComments {
		likes = append(likes, GetLikesCommentQueryRes{
			Type:    "comment",
			Comment: comment,
		})
	}

	for _, post := range user.LikedPosts {
		likes = append(likes, GetLikesPostQueryRes{
			Type: "post",
			Post: post,
		})
	}

	sort.Slice(likes, func(i, j int) bool {
		iCreatedAt := getCreatedAt(likes[i])
		jCreatedAt := getCreatedAt(likes[j])
		return iCreatedAt.Before(jCreatedAt)
	})

	return likes, nil
}

func (r gormRepository) GetComments(ctx context.Context, userId string) ([]internal.Comment, error) {
	var comments []internal.Comment

	if err := r.db.WithContext(ctx).Preload(clause.Associations).Where("user_id = ?", userId).Find(&comments).Error; err != nil {
		return []internal.Comment{}, err
	}

	return comments, nil
}

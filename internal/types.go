package internal

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uuid.UUID      `json:"id" uri:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// "comment" belongs to "user"
	CreatedBy User      `json:"created_by" gorm:"foreignKey:UserID;not null"`
	UserID    uuid.UUID `json:"-" gorm:"type:uuid;not null"`

	// "comment" belongs to "post" (comment "created in" post)
	CreatedIn       Post      `json:"created_in" gorm:"foreignKey:PostCreatedInID;not null"`
	PostCreatedInID uuid.UUID `json:"-" gorm:"type:uuid;not null"`

	// "post" has many "comments"
	PostID uuid.UUID `json:"-" gorm:"type:uuid;not null"`

	Description string `json:"description"`
	Likes       int64  `json:"likes"`

	// "comments" has many "comments"
	Replies  []Comment  `json:"replies" gorm:"foreignKey:ParentID;"`
	ParentID *uuid.UUID `json:"-" gorm:"type:uuid"`
}

type Post struct {
	ID        uuid.UUID      `json:"id" uri:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// // "user" has many "posts"
	// UserPostsRefer uuid.UUID `json:"-" gorm:"type:uuid;not null"`

	// // "user" has many "(liked) posts"
	// UserLikedPostsRefer uuid.UUID `json:"-" gorm:"type:uuid;not null"`

	// "post" belongs to "user"
	CreatedBy User      `json:"created_by" gorm:"foreignKey:UserID;references:ID;not null"`
	UserID    uuid.UUID `json:"-" gorm:"type:uuid;not null"`

	PictureLink string `json:"picture_link" gorm:"not null"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Likes       int64  `json:"likes"`

	// "post" has many "comments"
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID;"`
}

type User struct {
	ID        uuid.UUID      `json:"id" uri:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Username    string `json:"username" gorm:"not null"`
	Email       string `json:"-" gorm:"unique; not null"`
	PictureLink string `json:"picture_link"`

	Posts      []Post  `json:"posts" gorm:"foreignKey:UserID;references:ID"`   // "user" has many "posts"
	LikedPosts []Post  `json:"liked_posts" gorm:"many2many:user_liked_posts;"` // "user" has many "(liked) posts"
	Followers  []*User `json:"followers" gorm:"many2many:user_followers;"`     // "user" many to many "user"
	Followings []*User `json:"followings" gorm:"many2many:user_followings;"`   // "user" many to many "user"

	// Posts      []Post  `json:"posts" gorm:"foreignKey:UserPostsRefer"`            // "user" has many "posts"
	// LikedPosts []Post  `json:"liked_posts" gorm:"foreignKey:UserLikedPostsRefer"` // "user" has many "(liked) posts"
}

type UserIdUriRequest struct {
	UserId string `uri:"id" binding:"required"`
}

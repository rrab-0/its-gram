package internal

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// "comment" belongs to "user"
	CreatedBy User   `json:"created_by" gorm:"foreignKey:UserID;references:ID;not null"`
	UserID    string `json:"-" gorm:"type:unique;not null"`

	// "comment" belongs to "post" (comment "created in" post)
	CreatedIn       Post      `json:"created_in" gorm:"foreignKey:PostCreatedInID;not null"`
	PostCreatedInID uuid.UUID `json:"-" gorm:"type:uuid;not null"`

	// "post" has many "comments"
	PostID uuid.UUID `json:"-" gorm:"type:uuid;not null"`

	Description string `json:"description"`
	Likes       []User `json:"likes" gorm:"many2many:user_liked_comments;"`

	// "comments" has many "comments"
	Replies  []Comment  `json:"replies" gorm:"foreignKey:ParentID;"`
	ParentID *uuid.UUID `json:"-" gorm:"type:uuid"`
}

type Post struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// "post" belongs to "user"
	CreatedBy User   `json:"created_by" gorm:"foreignKey:UserID;references:ID;not null"`
	UserID    string `json:"-" gorm:"type:unique;not null"`

	PictureLink string `json:"picture_link" gorm:"not null"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`

	// "user" many to many "(liked) posts"
	Likes []User `json:"likes" gorm:"many2many:user_liked_posts;"`

	// "post" has many "comments"
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID;"`
}

type User struct {
	ID        string         `json:"id" gorm:"primaryKey;"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Username    string `json:"username" gorm:"not null"`
	Email       string `json:"-" gorm:"unique; not null"`
	PictureLink string `json:"picture_link"`

	Posts      []Post `json:"posts" gorm:"foreignKey:UserID;references:ID"`   // "user" has many "posts"
	LikedPosts []Post `json:"liked_posts" gorm:"many2many:user_liked_posts;"` // "user" many to many "(liked) posts"

	Comments      []Comment `json:"comments" gorm:"foreignKey:UserID;references:ID"`      // "user" has many "comments"
	LikedComments []Comment `json:"liked_comments" gorm:"many2many:user_liked_comments;"` // "user" many to many "(liked) comments"

	Followers  []*User `json:"followers" gorm:"many2many:user_followers;"`   // "user" many to many "user"
	Followings []*User `json:"followings" gorm:"many2many:user_followings;"` // "user" many to many "user"
}

type UserIdUriRequest struct {
	UserId string `uri:"id" binding:"required"`
}

type DeletedCommentOrPostResponse struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	IsDeleted bool           `json:"is_deleted"`
	CreatedBy User           `json:"created_by"`
}

type GetPostResponse struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`

	CreatedBy     User          `json:"created_by"`
	PictureLink   string        `json:"picture_link"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	Likes         []User        `json:"likes"`
	Comments      []interface{} `json:"comments"`
	TotalComments int           `json:"total_comments"`
}

package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrab-0/its-gram/internal"
	"gorm.io/gorm"
)

type Handler struct {
	Service
}

func NewHandler(db *gorm.DB) Handler {
	return Handler{
		Service: NewService(NewRepository(db)),
	}
}

// TODO: check if after deleting, the same user (email) can register again or not
func (h Handler) CreateUser(ctx *gin.Context) {
	firebaseId, fExists := ctx.Get("user_id")
	username, uExists := ctx.Get("username")
	email, eExists := ctx.Get("email")
	picture, pExists := ctx.Get("picture")

	if !fExists || !uExists || !eExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &internal.ErrorResponse{
			Message: "Failed to register.",
			Error:   "invalid token",
		})
		return
	}

	var profilePicture string
	if !pExists {
		profilePicture = ""
	} else {
		profilePicture = picture.(string)
	}

	user, err := h.Service.CreateUser(ctx.Request.Context(), firebaseId.(string), username.(string), email.(string), profilePicture)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			ctx.AbortWithStatusJSON(http.StatusConflict, internal.ErrorResponse{
				Message: "User already exists.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to create User.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, internal.SuccessResponse{
		Message: "User created successfully",
		Data:    user,
	})
}

// GetUser is a method to get a user
// @Summary Get a user
// @Description Returns a user and their followers, followings, posts, comments, liked posts, and liked comments.
// @Tags user
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} internal.SuccessResponse
// @Failure 400 {object} internal.ErrorResponse
// @Failure 404 {object} internal.ErrorResponse
// @Failure 500 {object} internal.ErrorResponse
// @Router /user/{id} [get]
func (h Handler) GetUser(ctx *gin.Context) {
	var reqUri internal.UserIdUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	user, err := h.Service.GetUser(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch user, user not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch User.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "User fetched successfully.",
		Data:    user,
	})
}

func (h Handler) SearchUser(ctx *gin.Context) {
	var reqQuery UserSearchRequest
	if err := ctx.Bind(&reqQuery); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	users, err := h.Service.SearchUser(ctx.Request.Context(), reqQuery)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch users, users not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch Users.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "Users fetched successfully.",
		Data:    users,
	})
}

const (
	MAXIMUM_LIMIT = 50
	MINIMUM_LIMIT = 10
	MINIMUM_PAGE  = 1
)

func (h Handler) GetUserHomepage(ctx *gin.Context) {
	var (
		reqUri   internal.UserIdUriRequest
		reqQuery GetHomepageQueryRequest
	)

	if err := ctx.Bind(&reqQuery); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	if reqQuery.Limit < MINIMUM_LIMIT {
		reqQuery.Limit = MINIMUM_LIMIT
	} else if reqQuery.Limit > MAXIMUM_LIMIT {
		reqQuery.Limit = MAXIMUM_LIMIT
	}

	if reqQuery.Page < MINIMUM_PAGE {
		reqQuery.Page = MINIMUM_PAGE
	}

	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	userId, idExists := ctx.Get("user_id")
	if !idExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &internal.ErrorResponse{
			Message: "Failed to fetch user's homepage.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to fetch user's homepage.",
			Error:   "invalid token",
		})
		return
	}

	posts, err := h.Service.GetUserHomepage(ctx.Request.Context(), reqUri, reqQuery)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch user's homempage.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch user's homempage.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "User's homepage fetched successfully.",
		Data:    posts,
	})
}

func (h Handler) GetUserHomepageInitialCursor(ctx *gin.Context) {
	var (
		reqUri   internal.UserIdUriRequest
		reqQuery GetUserHomepageInitialCursorQueryRequest
	)

	if err := ctx.Bind(&reqQuery); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	if reqQuery.Limit < MINIMUM_LIMIT {
		reqQuery.Limit = MINIMUM_LIMIT
	} else if reqQuery.Limit > MAXIMUM_LIMIT {
		reqQuery.Limit = MAXIMUM_LIMIT
	}

	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	userId, idExists := ctx.Get("user_id")
	if !idExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &internal.ErrorResponse{
			Message: "Failed to fetch user's initial cursor homepage.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to fetch user's initial cursor homepage.",
			Error:   "invalid token",
		})
		return
	}

	posts, err := h.Service.GetUserHomepageInitialCursor(ctx.Request.Context(), reqUri, reqQuery)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch user's initial cursor homempage.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch user's initial cursor homempage.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "User's initial cursor homepage fetched successfully.",
		Data:    posts,
	})
}

func (h Handler) GetUserHomepageCursor(ctx *gin.Context) {
	var (
		reqUri   internal.UserIdUriRequest
		reqQuery GetUserHomepageCursorQueryRequest
	)

	if err := ctx.Bind(&reqQuery); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	if reqQuery.Limit < MINIMUM_LIMIT {
		reqQuery.Limit = MINIMUM_LIMIT
	} else if reqQuery.Limit > MAXIMUM_LIMIT {
		reqQuery.Limit = MAXIMUM_LIMIT
	}

	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	userId, idExists := ctx.Get("user_id")
	if !idExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &internal.ErrorResponse{
			Message: "Failed to fetch user's cursor homepage.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to fetch user's cursor homepage.",
			Error:   "invalid token",
		})
		return
	}

	posts, err := h.Service.GetUserHomepageCursor(ctx.Request.Context(), reqUri, reqQuery)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch user's cursor homempage.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch user's cursor homempage.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "User's cursor homepage fetched successfully.",
		Data:    posts,
	})
}

func (h Handler) UpdateUserProfile(ctx *gin.Context) {
	var reqUri internal.UserIdUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	var reqBody UpdateUserProfileRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	userId, idExists := ctx.Get("user_id")
	if !idExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &internal.ErrorResponse{
			Message: "Failed to update user.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to update user.",
			Error:   "invalid token",
		})
		return
	}

	user, err := h.Service.UpdateUserProfile(ctx.Request.Context(), reqUri, reqBody)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to update user, user not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to update User.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "User updated successfully",
		Data:    user,
	})
}

func (h Handler) DeleteUser(ctx *gin.Context) {
	var reqUri internal.UserIdUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	userId, idExists := ctx.Get("user_id")
	if !idExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &internal.ErrorResponse{
			Message: "Failed to delete user.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to delete user.",
			Error:   "invalid token",
		})
		return
	}

	_, err := h.Service.DeleteUser(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to delete user, user not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to delete user.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: fmt.Sprintf("User with id %v deleted successfully.", reqUri.UserId),
	})
}

func (h Handler) FollowOtherUser(ctx *gin.Context) {
	var reqUri FollowOtherUserRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	userId, idExists := ctx.Get("user_id")
	if !idExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &internal.ErrorResponse{
			Message: "Failed to follow user.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to follow user.",
			Error:   "invalid token",
		})
		return
	}

	err := h.Service.FollowOtherUser(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to follow user, user not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to follow user.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: fmt.Sprintf("User with id %v successfully followed user with id %v.", reqUri.UserId, reqUri.OtherUserId),
	})
}

func (h Handler) UnfollowOtherUser(ctx *gin.Context) {
	var reqUri FollowOtherUserRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	userId, idExists := ctx.Get("user_id")
	if !idExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &internal.ErrorResponse{
			Message: "Failed to unfollow user.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to unfollow user.",
			Error:   "invalid token",
		})
		return
	}

	err := h.Service.UnfollowOtherUser(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to unfollow user, user not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to unfollow user.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: fmt.Sprintf("User with id %v successfully unfollowed user with id %v.", reqUri.UserId, reqUri.OtherUserId),
	})
}

func (h Handler) GetPosts(ctx *gin.Context) {
	var (
		reqUri  internal.UserIdUriRequest
		postRes []any
	)
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	posts, totalComments, err := h.Service.GetPosts(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch user's posts, posts not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch user's posts.",
			Error:   err.Error(),
		})
		return
	}

	for i, post := range posts {
		if post.DeletedAt.Valid {
			var deletedPost internal.DeletedCommentOrPostResponse
			deletedPost.ID = post.ID
			deletedPost.CreatedAt = post.CreatedAt
			deletedPost.UpdatedAt = post.UpdatedAt
			deletedPost.DeletedAt = post.DeletedAt
			deletedPost.IsDeleted = true
			deletedPost.CreatedBy = post.CreatedBy

			postRes = append(postRes, deletedPost)
			continue
		}

		var newPost internal.GetPostResponse
		newPost.ID = post.ID
		newPost.CreatedAt = post.CreatedAt
		newPost.UpdatedAt = post.UpdatedAt
		newPost.DeletedAt = post.DeletedAt

		newPost.CreatedBy = post.CreatedBy
		newPost.PictureLink = post.PictureLink
		newPost.Title = post.Title
		newPost.Description = post.Description
		newPost.Likes = post.Likes
		newPost.TotalComments = totalComments[i]

		for _, comment := range post.Comments {
			if !comment.DeletedAt.Valid {
				newPost.Comments = append(newPost.Comments, comment)
				continue
			}

			var deletedCommentRes internal.DeletedCommentOrPostResponse
			deletedCommentRes.ID = comment.ID
			deletedCommentRes.CreatedAt = comment.CreatedAt
			deletedCommentRes.UpdatedAt = comment.UpdatedAt
			deletedCommentRes.DeletedAt = comment.DeletedAt
			deletedCommentRes.IsDeleted = true
			deletedCommentRes.CreatedBy = comment.CreatedBy
			newPost.Comments = append(newPost.Comments, deletedCommentRes)
		}

		postRes = append(postRes, newPost)
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "User's posts fetched successfully.",
		Data:    postRes,
	})
}

func (h Handler) GetLikes(ctx *gin.Context) {
	var (
		reqUri internal.UserIdUriRequest
	)
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	likes, err := h.Service.GetLikes(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch user's likes, likes not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch user's likes.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "User's likes fetched successfully.",
		Data:    likes,
	})
}

func (h Handler) GetComments(ctx *gin.Context) {
	var reqUri internal.UserIdUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	comments, err := h.Service.GetComments(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch user's comments, comments not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch user's comments.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "User's comments fetched successfully.",
		Data:    comments,
	})
}

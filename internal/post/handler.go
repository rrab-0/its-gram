package post

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

func (h Handler) CreatePost(ctx *gin.Context) {
	var reqUri internal.UserIdUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	var reqBody CreatePostRequest
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
			Message: "Failed to create post.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to create post.",
			Error:   "invalid token",
		})
		return
	}

	post, err := h.Service.CreatePost(ctx.Request.Context(), reqUri, reqBody)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to create post, user not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to create post.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, internal.SuccessResponse{
		Message: "Post created successfully",
		Data:    post,
	})
}

func (h Handler) GetUserPosts(ctx *gin.Context) {
	var reqUri internal.UserIdUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	post, err := h.Service.GetUserPosts(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch posts, user or posts not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch posts.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "Posts fetched successfully",
		Data:    post,
	})
}

func (h Handler) GetPostById(ctx *gin.Context) {
	var reqUri PostIdUriRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	post, err := h.Service.GetPostById(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch post, post not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch post.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "Post fetched successfully",
		Data:    post,
	})
}

func (h Handler) DeletePost(ctx *gin.Context) {
	var reqUri PostAndUserUriRequest
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
			Message: "Failed to delete post.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to delete post.",
			Error:   "invalid token",
		})
		return
	}

	err := h.Service.DeletePost(ctx.Request.Context(), reqUri)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to delete post.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: fmt.Sprintf("Post with id %v deleted successfully", reqUri.PostId),
	})
}

func (h Handler) LikePost(ctx *gin.Context) {
	var reqUri PostAndUserUriRequest
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
			Message: "Failed to like post.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to like post.",
			Error:   "invalid token",
		})
		return
	}

	err := h.Service.LikePost(ctx.Request.Context(), reqUri)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to like post.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: fmt.Sprintf("Post with id %v liked successfully", reqUri.PostId),
	})
}

func (h Handler) UnlikePost(ctx *gin.Context) {
	var reqUri PostAndUserUriRequest
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
			Message: "Failed to unlike post.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to unlike post.",
			Error:   "invalid token",
		})
		return
	}

	err := h.Service.UnlikePost(ctx.Request.Context(), reqUri)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to unlike post.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: fmt.Sprintf("Post with id %v unliked successfully", reqUri.PostId),
	})
}

func (h Handler) CommentPost(ctx *gin.Context) {
	var reqUri PostAndUserUriRequest
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
			Message: "Failed to comment on post.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to comment on post.",
			Error:   "invalid token",
		})
		return
	}

	err := h.Service.CommentPost(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to comment on post, post not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to comment on post.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "Comment posted on post successfully.",
	})
}

func (h Handler) UncommentPost(ctx *gin.Context) {
	var reqUri PostAndUserUriRequest
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
			Message: "Failed to uncomment on post.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to uncomment on post.",
			Error:   "invalid token",
		})
		return
	}

	err := h.Service.UncommentPost(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to uncomment on post, post not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to uncomment on post.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "Comment removed from post successfully.",
	})
}

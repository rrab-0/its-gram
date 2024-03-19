package post

import (
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

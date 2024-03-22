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
	var (
		reqUri  internal.UserIdUriRequest
		reqBody CreatePostRequest
	)

	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

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

func (h Handler) GetPostById(ctx *gin.Context) {
	var (
		reqUri  PostIdUriRequest
		postRes GetPostByIdResponse
	)
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	post, totalComments, err := h.Service.GetPostById(ctx.Request.Context(), reqUri)
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

	postRes.ID = post.ID
	postRes.CreatedAt = post.CreatedAt
	postRes.CreatedBy = post.CreatedBy
	postRes.PictureLink = post.PictureLink
	postRes.Title = post.Title
	postRes.Description = post.Description
	postRes.Likes = post.Likes
	postRes.TotalComments = totalComments

	for _, comment := range post.Comments {
		if !comment.DeletedAt.Valid {
			postRes.Comments = append(postRes.Comments, comment)
			continue
		}

		var deletedCommentRes DeletedCommentResponse
		deletedCommentRes.ID = comment.ID
		deletedCommentRes.CreatedAt = comment.CreatedAt
		deletedCommentRes.UpdatedAt = comment.UpdatedAt
		deletedCommentRes.DeletedAt = comment.DeletedAt
		deletedCommentRes.IsDeleted = true
		deletedCommentRes.CreatedBy = comment.CreatedBy
		postRes.Comments = append(postRes.Comments, deletedCommentRes)
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "Post fetched successfully.",
		Data:    postRes,
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

func (h Handler) GetComment(ctx *gin.Context) {
	var (
		reqUri     GetCommentRequest
		commentRes GetCommentResponse
	)
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	comment, err := h.Service.GetComment(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to fetch comment, comment not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to fetch comment.",
			Error:   err.Error(),
		})
		return
	}

	commentRes.ID = comment.ID
	commentRes.CreatedAt = comment.CreatedAt
	commentRes.UpdatedAt = comment.UpdatedAt
	commentRes.DeletedAt = comment.DeletedAt
	commentRes.CreatedBy = comment.CreatedBy
	commentRes.Description = comment.Description
	commentRes.Likes = comment.Likes

	for _, reply := range comment.Replies {
		if !reply.DeletedAt.Valid {
			commentRes.Replies = append(commentRes.Replies, reply)
			continue
		}

		var deletedCommentRes DeletedCommentResponse
		deletedCommentRes.ID = reply.ID
		deletedCommentRes.CreatedAt = reply.CreatedAt
		deletedCommentRes.UpdatedAt = reply.UpdatedAt
		deletedCommentRes.DeletedAt = reply.DeletedAt
		deletedCommentRes.IsDeleted = true
		deletedCommentRes.CreatedBy = reply.CreatedBy
		commentRes.Replies = append(commentRes.Replies, deletedCommentRes)
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "Fetched comment successfully.",
		Data:    commentRes,
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

	var reqBody CreateCommentRequest
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

	err := h.Service.CommentPost(ctx.Request.Context(), reqUri, reqBody)
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

	ctx.JSON(http.StatusCreated, internal.SuccessResponse{
		Message: "Comment posted on post successfully.",
	})
}

func (h Handler) UncommentPost(ctx *gin.Context) {
	var reqUri CommentAndUserUriRequest
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

func (h Handler) ReplyComment(ctx *gin.Context) {
	var reqUri ReplyCommentRequest
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &internal.ErrorResponse{
			Message: "Invalid request.",
			Error:   internal.GenerateRequestValidatorError(err).Error(),
		})
		return
	}

	var reqBody CreateCommentRequest
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
			Message: "Failed to reply on comment.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to reply on comment.",
			Error:   "invalid token",
		})
		return
	}

	err := h.Service.ReplyComment(ctx.Request.Context(), reqUri, reqBody)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to reply on comment, post not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to reply on comment.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, internal.SuccessResponse{
		Message: "Reply posted on comment successfully.",
	})
}

func (h Handler) RemoveReplyFromComment(ctx *gin.Context) {
	var reqUri CommentAndUserUriRequest
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
			Message: "Failed to remove reply from comment.",
			Error:   "invalid token",
		})
		return
	}

	if reqUri.UserId != userId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, &internal.ErrorResponse{
			Message: "Failed to remove reply from comment.",
			Error:   "invalid token",
		})
		return
	}

	err := h.Service.RemoveReplyFromComment(ctx.Request.Context(), reqUri)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, internal.ErrorResponse{
				Message: "Failed to remove reply from comment, comment not found.",
				Error:   err.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, internal.ErrorResponse{
			Message: "Failed to remove reply from comment.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: "Reply removed from comment successfully.",
	})
}

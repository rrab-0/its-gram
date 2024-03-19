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

func (h Handler) CreateUser(ctx *gin.Context) {
	username, uExists := ctx.Get("username")
	email, eExists := ctx.Get("email")
	picture, pExists := ctx.Get("picture")

	if !uExists || !eExists {
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

	user, err := h.Service.CreateUser(ctx.Request.Context(), username.(string), email.(string), profilePicture)
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
	if err := ctx.ShouldBindUri(&reqQuery); err != nil {
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

func (h Handler) GetUserHomepage(ctx *gin.Context) {
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

	posts, err := h.Service.GetUserHomepage(ctx.Request.Context(), reqUri)
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
		Message: "User fetched successfully.",
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
			Message: "Failed to delete User.",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, internal.SuccessResponse{
		Message: fmt.Sprintf("User with id %v deleted successfully", reqUri.UserId),
	})
}

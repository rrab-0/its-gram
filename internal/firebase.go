package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type FirebaseApp struct {
	App *firebase.App
}

func NewFirebaseApp(serviceAccPath string) (*FirebaseApp, error) {
	firebaseOpt := option.WithCredentialsFile(serviceAccPath)
	app, err := firebase.NewApp(context.Background(), nil, firebaseOpt)
	if err != nil {
		return nil, err
	}

	return &FirebaseApp{
		App: app,
	}, nil
}

type FirebaseAuth struct {
	auth *auth.Client
}

func NewFirebaseAuth(app *firebase.App) (*FirebaseAuth, error) {
	auth, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &FirebaseAuth{
		auth: auth,
	}, nil
}

func (f *FirebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := f.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id token: %v", err)
	}

	return token, nil
}

func (f *FirebaseAuth) ValidateToken(funcType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: "Failed to authenticate.",
				Error:   "token is empty",
			})
			return
		}

		fields := strings.Fields(tokenStr)
		token := fields[1]

		idToken, err := f.auth.VerifyIDToken(ctx, token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: "Failed to authenticate.",
				Error:   err.Error(),
			})
			return
		}

		// // 2024/03/18 21:09:12 map[
		// 	auth_time:1.71077095e+09
		// 	email:jasapedia2024@gmail.com
		// 	email_verified:true
		// 	firebase:map[
		// 		identities:map[
		// 			email:[jasapedia2024@gmail.com]
		// 			google.com:[102893898261984639861]]
		// 			sign_in_provider:google.com]
		// 			name:jasa pedia
		// 			picture:https://lh3.googleusercontent.com/a/ACg8ocL3ls0l2jiZ6rRuYV8NDhtGF2O2QOHHtbFTjwu63u9X=s96-c
		// 			user_id:shxCXE3jYfYuQAUCqHCZS4g22xz2]
		if funcType == "" {
			ctx.Set("user_id", idToken.Claims["user_id"])
		}

		if funcType == "REGISTER" {
			ctx.Set("user_id", idToken.Claims["user_id"])
			ctx.Set("username", idToken.Claims["name"])
			ctx.Set("email", idToken.Claims["email"])
			ctx.Set("picture", idToken.Claims["picture"])
		}
		ctx.Next()
	}
}

func (f *FirebaseAuth) ValidateNgrokDevToken(funcType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization-Bearer")
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: "Failed to authenticate.",
				Error:   "token is empty",
			})
			return
		}

		fields := strings.Fields(tokenStr)
		token := fields[1]

		idToken, err := f.auth.VerifyIDToken(ctx, token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: "Failed to authenticate.",
				Error:   err.Error(),
			})
			return
		}

		if funcType == "" {
			ctx.Set("user_id", idToken.Claims["user_id"])
		}

		if funcType == "REGISTER" {
			ctx.Set("user_id", idToken.Claims["user_id"])
			ctx.Set("username", idToken.Claims["name"])
			ctx.Set("email", idToken.Claims["email"])
			ctx.Set("picture", idToken.Claims["picture"])
		}
		ctx.Next()
	}
}

var dummyUserCount = 3

// Checks if userId (doesn't have to be valid) is present in URI request or not,
// if present sets random values in contexts needed.
func (f *FirebaseAuth) ValidateDevToken(funcType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			firebaseId = fmt.Sprintf("%d", dummyUserCount)
			username   = fmt.Sprintf("dummy%v", dummyUserCount)
			email      = fmt.Sprintf("dummy%v@gmail.com", dummyUserCount)
			picture    = ""
		)

		dummyUserCount++

		id := ctx.Param("id")
		if id == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: "Failed to authenticate.",
				Error:   "id uri is empty",
			})
			return
		}

		if funcType == "" {
			ctx.Set("user_id", id)
		}

		if funcType == "REGISTER" {
			ctx.Set("user_id", firebaseId)
			ctx.Set("username", username)
			ctx.Set("email", email)
			ctx.Set("picture", picture)
		}
		ctx.Next()
	}
}

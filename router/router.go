package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rrab-0/its-gram/internal"
	"github.com/rrab-0/its-gram/internal/post"
	"github.com/rrab-0/its-gram/internal/user"
	"github.com/spf13/viper"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine, firebaseAuth *internal.FirebaseAuth, userHandler user.Handler, postHandler post.Handler) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{ /* "http://localhost:5173" */ "*"},
		AllowMethods:     []string{"OPTIONS", "POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	var (
		validateRegisterToken gin.HandlerFunc
		validateToken         gin.HandlerFunc
	)

	if viper.GetString("ENV") == "LOCAL_DEV" {
		validateRegisterToken = firebaseAuth.ValidateDevToken("REGISTER")
		validateToken = firebaseAuth.ValidateDevToken("")
	} else if viper.GetString("ENV") == "NGROK_DEV" {
		validateRegisterToken = firebaseAuth.ValidateNgrokDevToken("REGISTER")
		validateToken = firebaseAuth.ValidateNgrokDevToken("")
	} else {
		validateRegisterToken = firebaseAuth.ValidateToken("REGISTER")
		validateToken = firebaseAuth.ValidateToken("")
	}

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, "world")
	})

	// NOTE: ":id" here is
	// user_id from context which is,
	// idToken.Claims["user_id"] from firebase auth idToken.
	v1 := r.Group("/api/v1")
	user := v1.Group("/user")
	{
		user.POST("/register/:id", validateRegisterToken, userHandler.CreateUser)

		user.GET("/:id", userHandler.GetUser)
		user.GET("/:id/posts", userHandler.GetPosts)
		user.GET("/:id/comments", userHandler.GetComments)
		user.GET("/:id/likes", userHandler.GetLikes)
		user.GET("/search", userHandler.SearchUser)

		user.Use(validateToken)
		user.GET("/:id/homepage", userHandler.GetUserHomepage)
		user.GET("/:id/homepage/cursor/initial", userHandler.GetUserHomepageInitialCursor)
		user.GET("/:id/homepage/cursor", userHandler.GetUserHomepageCursor)
		user.PATCH("/profile/update/:id", userHandler.UpdateUserProfile)
		user.DELETE("/delete/:id", userHandler.DeleteUser)

		user.POST("/:id/follow/:otherUserId", userHandler.FollowOtherUser)
		user.DELETE("/:id/unfollow/:otherUserId", userHandler.UnfollowOtherUser)
	}

	post := v1.Group("/post")
	{
		post.GET("/:id", postHandler.GetPostById)
		post.GET("/comment/:commentId", postHandler.GetComment)

		post.Use(validateToken)
		post.POST("/create/:id", postHandler.CreatePost)
		post.DELETE("/:postId/user/:id/delete", postHandler.DeletePost)

		post.POST("/:postId/user/:id/like", postHandler.LikePost)
		post.DELETE("/:postId/user/:id/unlike", postHandler.UnlikePost)

		post.POST("/:postId/user/:id/comment", postHandler.CommentPost)
		post.DELETE("/user/:id/comment/remove/:commentId", postHandler.UncommentPost)

		post.POST("/:postId/user/:id/comment/reply/:commentId", postHandler.ReplyComment)
		post.DELETE("/user/:id/comment/remove/reply/:commentId", postHandler.RemoveReplyFromComment)

		post.POST("/user/:id/comment/like/:commentId", postHandler.LikeComment)
		post.DELETE("/user/:id/comment/unlike/:commentId", postHandler.UnlikeComment)
	}
}

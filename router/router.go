package router

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rrab-0/its-gram/internal"
	"github.com/rrab-0/its-gram/internal/post"
	"github.com/rrab-0/its-gram/internal/user"
)

func Setup(r *gin.Engine, firebaseAuth *internal.FirebaseAuth, userHandler user.Handler, postHandler post.Handler) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{ /* "http://localhost:5173" */ "*"},
		AllowMethods:     []string{"OPTIONS", "POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	var (
		validateRegisterToken gin.HandlerFunc
		validateToken         gin.HandlerFunc
	)

	if os.Getenv("ENV") == "LOCAL_DEV" {
		validateRegisterToken = firebaseAuth.ValidateDevToken("REGISTER")
		validateToken = firebaseAuth.ValidateDevToken("")
	} else if os.Getenv("ENV") == "NGROK_DEV" {
		validateRegisterToken = firebaseAuth.ValidateNgrokDevToken("REGISTER")
		validateToken = firebaseAuth.ValidateNgrokDevToken("")
	} else {
		validateRegisterToken = firebaseAuth.ValidateToken("REGISTER")
		validateToken = firebaseAuth.ValidateToken("")
	}

	// TODO:
	// - need to clean all response
	// - need to handle query responses (not found, already liked, etc)

	// NOTE: ":id" here is user_id from context which is from firebase auth idToken
	v1 := r.Group("/api/v1")
	user := v1.Group("/user")
	{
		user.POST("/register/:id", validateRegisterToken, userHandler.CreateUser)

		user.GET("/:id", userHandler.GetUser)
		user.GET("/:id/likes", userHandler.GetLikes)       // TODO: check if working
		user.GET("/:id/comments", userHandler.GetComments) // TODO: check if working
		user.GET("/search", userHandler.SearchUser)

		user.Use(validateToken)
		user.GET("/:id/homepage", userHandler.GetUserHomepage)
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
		post.DELETE("/user/:id/delete/:postId", postHandler.DeletePost)

		post.POST("/user/:id/like/:postId", postHandler.LikePost)
		post.DELETE("/user/:id/unlike/:postId", postHandler.UnlikePost)

		post.POST("/user/:id/comment/:postId", postHandler.CommentPost)
		post.DELETE("/user/:id/comment/remove/:commentId", postHandler.UncommentPost)

		post.POST("/:postId/user/:id/comment/reply/:commentId", postHandler.ReplyComment)
		post.DELETE("/user/:id/comment/remove/reply/:commentId", postHandler.RemoveReplyFromComment)

		// TODO: implement this also
		// post.POST("/user/:id/comment/like/:commentId", postHandler.LikeComment) // TODO: check if working
		// post.DELETE("/user/:id/comment/unlike/:commentId", postHandler.UnlikeComment) // TODO: check if working
	}
}

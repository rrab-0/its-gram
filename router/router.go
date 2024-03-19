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

	v1 := r.Group("/api/v1")
	user := v1.Group("/user")
	{
		user.GET("/:id", userHandler.GetUser)
		user.GET("/search", userHandler.SearchUser) // TODO: check if working
		user.POST("/register/:id", validateRegisterToken, userHandler.CreateUser)

		user.Use(validateToken)
		user.GET("/:id/homepage", userHandler.GetUserHomepage) // TODO: check if working
		user.PATCH("/profile/update/:id", userHandler.UpdateUserProfile)
		user.DELETE("/delete/:id", userHandler.DeleteUser)
	}

	post := v1.Group("/post")
	{
		post.GET("/:id", postHandler.GetPostById)
		post.GET("/user/:id", postHandler.GetUserPosts)

		post.Use(validateToken)
		post.POST("/create/:id", postHandler.CreatePost)
	}
}

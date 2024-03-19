package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rrab-0/its-gram/db"
	"github.com/rrab-0/its-gram/internal"
	"github.com/rrab-0/its-gram/internal/post"
	"github.com/rrab-0/its-gram/internal/user"
	"github.com/rrab-0/its-gram/router"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("ERROR: Failed to load .env")
	}

	pgsql, err := db.NewPostgreSQL()
	if err != nil {
		log.Fatalf("ERROR: Failed to connect to PostgreSQL: %v", err.Error())
	}

	if err := pgsql.Migrate(); err != nil {
		log.Fatalf("ERROR: Failed to migrate PostgreSQL: %v", err.Error())
	}

	firebaseAuth, err := internal.NewFirebaseAuth(os.Getenv("SERVICE_ACCOUNT_KEY_PATH"))
	if err != nil {
		log.Fatalf("ERROR: Failed to initialize firebase auth: %v", err.Error())
	}

	userHandler := user.NewHandler(pgsql.DB)
	postHandler := post.NewHandler(pgsql.DB)

	gin.ForceConsoleColor()
	r := gin.Default()
	router.Setup(
		r,
		firebaseAuth,
		userHandler,
		postHandler,
	)

	if err := r.Run(os.Getenv("DEV_HOST") + ":" + os.Getenv("DEV_PORT")); err != nil {
		panic("ERROR: Failed to start server: " + err.Error())
	}
}

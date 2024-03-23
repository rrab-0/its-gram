package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/rrab-0/its-gram/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rrab-0/its-gram/db"
	"github.com/rrab-0/its-gram/internal"
	"github.com/rrab-0/its-gram/internal/post"
	"github.com/rrab-0/its-gram/internal/user"
	"github.com/rrab-0/its-gram/router"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

// @title           its-gram api docs
// @version         1.0
// @description     © Layanan Aplikasi its-gram
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /api/v1
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

	firebase, err := internal.NewFirebaseApp(os.Getenv("SERVICE_ACCOUNT_KEY_PATH"))
	if err != nil {
		log.Fatalf("ERROR: Failed to initialize firebase app: %v", err.Error())
	}

	firebaseAuth, err := internal.NewFirebaseAuth(firebase.App)
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

func runServer(ctx context.Context, r *gin.Engine) error {
	env := os.Getenv("ENV")
	if env == "" {
		return fmt.Errorf("ENV in .env is not set")
	}

	if os.Getenv("ENV") == "NGROK_DEV" {
		listener, err := ngrokListener(ctx)
		if err != nil {
			return err
		}

		fmt.Printf("\n")
		log.Printf("NGROK: Ingress established with %v at: https://%v\n\n", listener.Addr().Network(), listener.Addr())

		return r.RunListener(listener)
	}

	if os.Getenv("ENV") == "LOCAL_DEV" {
		return r.Run(os.Getenv("DEV_HOST") + ":" + os.Getenv("DEV_PORT"))
	}

	return nil
}

func ngrokListener(ctx context.Context) (net.Listener, error) {
	return ngrok.Listen(ctx,
		config.HTTPEndpoint(
			config.WithBasicAuth(os.Getenv("NGROK_BASIC_AUTH_USERNAME"), os.Getenv("NGROK_BASIC_AUTH_PASSWORD")),
		),
		ngrok.WithAuthtokenFromEnv(),
	)
}

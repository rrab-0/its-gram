package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/rrab-0/its-gram/docs"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
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
	if len(os.Args) == 1 || (os.Args[1] != "DEV" && os.Args[1] != "AWS") {
		log.Fatalf("ERROR: Need to specify command-line argument, only 'DEV' and 'AWS' are allowed")
	}

	if err := internal.LoadConfig(os.Args[1]); err != nil {
		log.Fatalf("ERROR: Failed to load configs: %v", err.Error())
	}

	pgsql, err := db.NewPostgreSQL()
	if err != nil {
		log.Fatalf("ERROR: Failed to connect to PostgreSQL: %v", err.Error())
	}

	if err := pgsql.Migrate(); err != nil {
		log.Fatalf("ERROR: Failed to migrate PostgreSQL: %v", err.Error())
	}

	firebase, err := internal.NewFirebaseApp(viper.GetString("SERVICE_ACCOUNT_KEY"))
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

	if err := runServer(context.Background(), r); err != nil {
		log.Fatalf("ERROR: Failed to start server: %v", err.Error())
	}
}

func runServer(ctx context.Context, r *gin.Engine) error {
	env := viper.GetString("ENV")

	if env == "NGROK_DEV" {
		listener, err := ngrokListener(ctx)
		if err != nil {
			return err
		}

		fmt.Printf("\n")
		log.Printf("NGROK: Ingress established with %v at: https://%v\n\n", listener.Addr().Network(), listener.Addr())

		return r.RunListener(listener)
	}

	if env == "LOCAL_DEV" {
		return r.Run(viper.GetString("DEV_HOST") + ":" + viper.GetString("DEV_PORT"))
	}

	if env == "AWS" {
		return r.Run(":" + viper.GetString("DEV_PORT"))
	}

	return nil
}

func ngrokListener(ctx context.Context) (net.Listener, error) {
	return ngrok.Listen(ctx,
		config.HTTPEndpoint(
			config.WithBasicAuth(viper.GetString("NGROK_BASIC_AUTH_USERNAME"), viper.GetString("NGROK_BASIC_AUTH_PASSWORD")),
		),
		ngrok.WithAuthtokenFromEnv(),
	)
}

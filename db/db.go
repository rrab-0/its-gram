package db

import (
	"fmt"
	"log"
	"os"

	"github.com/rrab-0/its-gram/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	Migrate() error
}

type postgreSQL struct {
	DB *gorm.DB
}

func NewPostgreSQL() (postgreSQL, error) {
	var (
		DB_HOST     = os.Getenv("DB_HOST")
		DB_USER     = os.Getenv("DB_USER")
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
		DB_NAME     = os.Getenv("DB_NAME")
		DB_PORT     = os.Getenv("DB_PORT")
	)

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		DB_HOST,
		DB_USER,
		DB_PASSWORD,
		DB_NAME,
		DB_PORT,
	)

	var gormConfig gorm.Config
	if os.Getenv("ENV") == "LOCAL_DEV" {
		gormConfig = gorm.Config{}
	} else {
		gormConfig = gorm.Config{
			TranslateError: true,
		}
	}

	db, err := gorm.Open(postgres.Open(connStr), &gormConfig)
	if err != nil {
		return postgreSQL{}, err
	}

	log.Println("SUCCESS: Connected to PostgreSQL database.")
	return postgreSQL{DB: db}, nil
}

func (p postgreSQL) Migrate() error {
	err := p.DB.AutoMigrate(
		internal.Comment{},
		internal.Post{},
		internal.User{},
	)
	if err != nil {
		return err
	}

	log.Println("SUCCESS: PostgreSQL migration completed (Some tables won't be created if they already exist but new fields will be appended).")
	return nil
}

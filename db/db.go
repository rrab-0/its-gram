package db

import (
	"fmt"
	"log"

	"github.com/rrab-0/its-gram/internal"
	"github.com/spf13/viper"
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
		DB_HOST     = viper.GetString("DB_HOST")
		DB_USER     = viper.GetString("DB_USER")
		DB_PASSWORD = viper.GetString("DB_PASSWORD")
		DB_NAME     = viper.GetString("DB_NAME")
		DB_PORT     = viper.GetString("DB_PORT")
	)

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		DB_HOST,
		DB_USER,
		DB_PASSWORD,
		DB_NAME,
		DB_PORT,
	)

	var gormConfig gorm.Config
	if viper.GetString("ENV") == "LOCAL_DEV" {
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

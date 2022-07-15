package db

import (
	"context"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func Init() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		host,
		user,
		password,
		name,
		port,
	)
	var err error
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false, // disables implicit prepared statement usage
	}), &gorm.Config{})
	db.Use(otelgorm.NewPlugin())
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func GetDB(c context.Context) *gorm.DB {
	return db.WithContext(c)
}

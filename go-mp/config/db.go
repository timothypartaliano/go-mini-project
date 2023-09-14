package config

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "os"
)

var db *gorm.DB

func InitDatabase() *gorm.DB {
    err := godotenv.Load()
    if err != nil {
        logrus.Fatalf("Error loading .env file: %v", err)
    }

    dsn := os.Getenv("DB_CONNECTION_STRING")
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        logrus.Fatalf("Error connecting to the database: %v", err)
    }

    return db
}
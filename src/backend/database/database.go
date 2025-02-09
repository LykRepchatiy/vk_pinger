package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	logger = log.New(os.Stdout, "docker-pinger: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
)

type DBContainer struct {
	ID          uint      `gorm:"primaryKey"`
	ContainerID string    `gorm:"uniqueIndex;not null"`
	IP          string    `gorm:"type:varchar(255);not null"`
	Status      string    `gorm:"type:varchar(255);not null"`
	Timestamp   time.Time `gorm:"not null"`
	Datestamp   time.Time `gorm:"not null"`
}

type Env struct {
	Port       string
	DBHost     string
	DBUser     string
	DBPort     string
	DBName     string
	DBPassword string
}

func ParseEnv() Env {
	var env Env
	env.Port = os.Getenv("PORT")
	if env.Port == "" {
		logger.Fatal("PORT environment variable is required")
		os.Exit(1)
	}
	env.DBHost = os.Getenv("DATABASE_HOST")
	if env.DBHost == "" {
		logger.Fatal("DATABASE_HOST environment variable is required")
		os.Exit(1)
	}
	env.DBUser = os.Getenv("DATABASE_USER")
	if env.DBUser == "" {
		logger.Fatal("DATABASE_USER environment variable is required")
		os.Exit(1)
	}
	env.DBPort = os.Getenv("DATABASE_PORT")
	if env.DBPort == "" {
		logger.Fatal("DATABASE_PORT environment variable is required")
		os.Exit(1)
	}
	env.DBName = os.Getenv("DATABASE_NAME")
	if env.DBName == "" {
		logger.Fatal("DATABASE_NAME environment variable is required")
		os.Exit(1)
	}
	env.DBPassword = os.Getenv("DATABASE_PASSWORD")
	if env.DBPassword == "" {
		logger.Fatal("DATABASE_PASSWORD environment variable is required")
		os.Exit(1)
	}
	return env
}

func DBConnect() (*gorm.DB, error) {
	env := ParseEnv()
	dsn := fmt.Sprintf("host=%s user=%s port=%s dbname=%s password=%s", env.DBHost, env.DBUser, env.DBPort, env.DBName, env.DBPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&DBContainer{})
	return db, nil
}
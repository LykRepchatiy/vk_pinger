package main

import (
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBContainer struct {
	ID        uint      `gorm:"primaryKey"`
	IP        string    `gorm:"uniqueIndex;not null"`
	Status    string    `gorm:"type:varchar(255);not null"`
	TimeStamp time.Time `gorm:"not null"`
	DateStamp string    `gorm:"type:varchar(255)"`
}

func dbConnect() {
	container := DBContainer{
		IP:        "0.0.0.0",
		Status:    "down",
		TimeStamp: time.Now(),
		DateStamp: time.Now().String(),
	}
	log.Println(container.TimeStamp)
	time.Sleep(2 * time.Second)
	newContainer := DBContainer{
		IP:        "0.0.0.0",
		Status:    "ok",
		TimeStamp: time.Now(),
		DateStamp: time.Now().String(),
	}
	dsn := "host=localhost user=myuser port=5433 dbname=mydatabase password='mypassword'"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&DBContainer{})
	db.Create(&container)
	db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "ip"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":     newContainer.Status,
			"time_stamp": newContainer.TimeStamp,
			"date_stamp": gorm.Expr("CASE WHEN ? = 'ok' THEN ? END", newContainer.Status, newContainer.DateStamp),
		}),
	}).Create(&newContainer)
}

func main() {
	// http.HandleFunc("/ping", pingFunc)
	dbConnect()
	http.ListenAndServe(":8080", nil)
}

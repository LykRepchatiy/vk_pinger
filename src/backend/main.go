package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBContainer struct {
	ID        uint      `gorm:"primaryKey"`
	IP        string    `json:"ip" gorm:"uniqueIndex;not null"`
	Status    string    `json:"status" gorm:"type:varchar(255);not null"`
	Timestamp time.Time `json:"timestamp" gorm:"not null"`
	Datestamp string    `json:"datestamp" gorm:"type:varchar(255)"`
}

func PutStatus(w http.ResponseWriter, r *http.Request) {
	var dbContainer DBContainer
	if r.Method != http.MethodPost {
		log.Println("wrong method")
		return
	}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(byteReq, &dbContainer)
	if err != nil {
		log.Println(err)
		return
	}
	db, err := dbConnect()
	if err != nil {
		log.Println(err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
		return
	}
	defer sqlDB.Close()
	err = db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "ip"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":    dbContainer.Status,
			"timestamp": dbContainer.Timestamp,
			"datestamp": gorm.Expr("CASE WHEN ? = 'ok' THEN ? ELSE db_containers.datestamp END", dbContainer.Status, dbContainer.Datestamp),
		}),
	}).Create(&dbContainer).Error
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func dbConnect() (*gorm.DB, error) {
	dsn := "host=postgres user=myuser port=5432 dbname=mydatabase password='mypassword'"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&DBContainer{})
	return db, nil
}

func main() {
	http.HandleFunc("/putStatus", PutStatus)
	dbConnect()
	http.ListenAndServe(":8080", nil)
}

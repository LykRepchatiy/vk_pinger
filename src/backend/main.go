package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Request struct {
	ContainerID string            `json:"containerID"`
	IP          map[string]string `json:"ip"`
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Datestamp   time.Time         `json:"datestamp"`
}

type ToFront struct {
	ContainerID string    `json:"containerID"`
	IP          string    `json:"ip"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Datestamp   time.Time `json:"datestamp"`
}

type DBContainer struct {
	ID          uint      `gorm:"primaryKey"`
	ContainerID string    `gorm:"uniqueIndex;not null"`
	IP          string    `gorm:"type:varchar(255);not null"`
	Status      string    `gorm:"type:varchar(255);not null"`
	Timestamp   time.Time `gorm:"not null"`
	Datestamp   time.Time `gorm:"not null"`
}

var (
	logger = log.New(os.Stdout, "backend: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
)

func ContainerList(w http.ResponseWriter, r *http.Request) {
	reqs := []ToFront{}
	conts := []DBContainer{}
	db, err := dbConnect()
	if err != nil {
		logger.Println(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Println(err)
	}
	defer sqlDB.Close()
	query := db.Order("timestamp ASC").Find(&conts)
	if query.Error != nil {
		logger.Println(query.Error)
	}

	for i := range conts {
		reqs = append(reqs, ToFront{
			ContainerID: conts[i].ContainerID,
			IP:          conts[i].IP,
			Status:      conts[i].Status,
			Timestamp:   conts[i].Timestamp,
			Datestamp:   conts[i].Datestamp,
		})
	}

	json, err := json.Marshal(reqs)
	if err != nil {
		logger.Println(err)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Write(json)
	w.WriteHeader(http.StatusOK)
}

func PutStatus(w http.ResponseWriter, r *http.Request) {
	reqs := make([]Request, 1)
	if r.Method != http.MethodPost {
		logger.Println("wrong method")
		return
	}
	byteReq, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Println(err)
		return
	}
	err = json.Unmarshal(byteReq, &reqs)
	if err != nil {
		logger.Println(err)
		return
	}
	DBconts := make([]DBContainer, len(reqs))
	for i, req := range reqs {
		for net, ip := range req.IP {
			DBconts[i].IP = net + ", " + ip + "\n"
		}
		DBconts[i].ContainerID = req.ContainerID
		DBconts[i].Status = req.Status
		DBconts[i].Timestamp = req.Timestamp
		DBconts[i].Datestamp = req.Datestamp
	}
	db, err := dbConnect()
	if err != nil {
		logger.Println(err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Println(err)
		return
	}
	defer sqlDB.Close()

	for _, dbContainer := range DBconts {
		err = db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "container_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"ip":        gorm.Expr("CASE WHEN ? = 'running' THEN ? ELSE db_containers.ip END", dbContainer.Status, dbContainer.IP),
				"status":    dbContainer.Status,
				"timestamp": dbContainer.Timestamp,
				"datestamp": gorm.Expr("CASE WHEN ? = 'running' THEN ? ELSE db_containers.datestamp END", dbContainer.Status, dbContainer.Datestamp),
			}),
		}).Create(&dbContainer).Error
		if err != nil {
			logger.Println(err)
			return
		}
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
	mux := http.NewServeMux()
	mux.HandleFunc("/putStatus", PutStatus)
	mux.HandleFunc("/containerList", ContainerList)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":8080", handler)
}

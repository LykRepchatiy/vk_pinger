package handlers

import (
	database "backend/database"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	logger = log.New(os.Stdout, "backend: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
)

type Request struct {
	ContainerID string            `json:"containerID"`
	IP          map[string]string `json:"ip"`
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Datestamp   time.Time         `json:"datestamp"`
}

type ToFront struct {
	IP        string    `json:"ip"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Datestamp time.Time `json:"datestamp"`
}

func ContainerList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Println("Error: wrong http method")
		return
	}
	reqs := []ToFront{}
	conts := []database.DBContainer{}
	db, err := database.DBConnect()
	if err != nil {
		logger.Println(err)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Println(err)
		return
	}
	defer sqlDB.Close()
	query := db.Order("timestamp ASC").Find(&conts)
	if query.Error != nil {
		logger.Println(query.Error)
		return
	}

	for i := range conts {
		reqs = append(reqs, ToFront{
			IP:        conts[i].IP,
			Status:    conts[i].Status,
			Timestamp: conts[i].Timestamp,
			Datestamp: conts[i].Datestamp,
		})
	}

	json, err := json.Marshal(reqs)
	if err != nil {
		logger.Println(err)
		return
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Write(json)
	w.WriteHeader(http.StatusOK)
}

func PutStatus(w http.ResponseWriter, r *http.Request) {
	reqs := []Request{}
	if r.Method != http.MethodPost {
		logger.Println("Error: wrong http method")
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
	DBconts := make([]database.DBContainer, len(reqs))
	for i, req := range reqs {
		for net, ip := range req.IP {
			DBconts[i].IP = net + ", " + ip + "\n"
		}
		DBconts[i].ContainerID = req.ContainerID
		DBconts[i].Status = req.Status
		DBconts[i].Timestamp = req.Timestamp
		DBconts[i].Datestamp = req.Datestamp
	}
	db, err := database.DBConnect()
	if err != nil {
		logger.Println(err)
		http.Error(w, `{"Error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Println(err)
		return
	}
	defer sqlDB.Close()

	for _, dbCont := range DBconts {
		err = db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "container_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"ip":        gorm.Expr("CASE WHEN ? = 'running' THEN ? ELSE db_containers.ip END", dbCont.Status, dbCont.IP),
				"status":    dbCont.Status,
				"timestamp": dbCont.Timestamp,
				"datestamp": gorm.Expr("CASE WHEN ? = 'running' THEN ? ELSE db_containers.datestamp END", dbCont.Status, dbCont.Datestamp),
			}),
		}).Create(&dbCont).Error
		if err != nil {
			logger.Println(err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

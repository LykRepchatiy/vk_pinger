package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type DBContainer struct {
	IP        string    `json:"ip"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Datestamp string    `json:"datestamp"`
}

func createRequest(container DBContainer) (*http.Request, error){
	byteSl, err := json.Marshal(container)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/putStatus", bytes.NewBuffer(byteSl))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-type", "application/json")
	return req, nil
}

func main() {
	client := http.Client{}
	container := DBContainer{
			IP: "0.0.0.0",
			Status: "down",
			Timestamp: time.Now(),
			Datestamp: "2025-02-07T01:06:25.601011292+03:00",
	}
	req, err := createRequest(container)
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode == http.StatusOK {
		log.Println("OK")
	}
}

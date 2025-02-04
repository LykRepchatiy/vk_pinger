package main

import (
	"log"
	"net/http"
)

func pingFunc(w http.ResponseWriter, r *http.Request) {

	log.Println("Hello world!")
}

func main() {
	http.HandleFunc("/ping", pingFunc)
	http.ListenAndServe(":8080", nil)
}
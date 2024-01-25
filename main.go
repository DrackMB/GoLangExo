package main

import (
	"estiam/dictionary"
	"estiam/middleware"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var d *dictionary.Dictionary

func main() {
	d = dictionary.New()
	r := mux.NewRouter()
	r.Use(middleware.LogMiddleware)
	r.Use(middleware.AuthenticationMiddleware)
	fmt.Println("Starting dictionary server on port 8083")

	r.HandleFunc("/", AddWord).Methods("POST")

	r.HandleFunc("/define/{word}", GetDefin).Methods("GET")

	r.HandleFunc("/remove/{word}", DeleteWord).Methods("DELETE")

	r.HandleFunc("/list", GetALL).Methods("GET")

	http.ListenAndServe(":8083", r)
}

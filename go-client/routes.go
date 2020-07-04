package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var router *mux.Router = mux.NewRouter()

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success")
}

func DefaultHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "default")
}

func init() {
	dir := "./static"
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer((http.Dir(dir)))))
	web := router.PathPrefix("/web").Subrouter()
	api := router.PathPrefix("/api").Subrouter()

	web.HandleFunc("/", DefaultHandle).Methods("GET")
	api.HandleFunc("/", DefaultHandle).Methods("GET")
}

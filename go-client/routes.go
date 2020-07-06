package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-client/api"
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
	rWeb := router.PathPrefix("/web").Subrouter()
	rApi := router.PathPrefix("/api").Subrouter()

	rWeb.HandleFunc("/", DefaultHandle).Methods("GET")
	rApi.HandleFunc("/", DefaultHandle).Methods("GET")

	rApi.HandleFunc("/wx_user_list", api.GetWxUserList).Methods("GET")
}

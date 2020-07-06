package main

import (
	"go-client/config"
	"go-client/ws"
	"log"
	"net/http"
	"time"
)

func init() {
	config.SetUp()
	ws.SetUp()
}
func main() {

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + config.Config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	ws.RecvLog()
	log.Fatal(srv.ListenAndServe())
}

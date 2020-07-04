package main

import (
	"go-client/config"
	"go-client/ws"
	"log"
	"net/http"
	"time"
)

func main() {

	cfg := config.GetConfig()
	log.Println((cfg))

	wsClient := ws.NewClient(cfg.DllServer)
	defer wsClient.Conn.Close()
	wsClient.Start()

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + cfg.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

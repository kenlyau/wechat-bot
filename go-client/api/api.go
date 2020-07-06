package api

import (
	"fmt"
	"go-client/ws"
	"net/http"
)

func GetWxUserList(w http.ResponseWriter, r *http.Request) {
	ws.GetWxUserList()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "get wx user list")
}

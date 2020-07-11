package api

import (
	"encoding/json"
	"net/http"
)

func jsonWrite(w http.ResponseWriter, d interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func jsonBind(r *http.Request, d interface{}) error {
	return json.NewDecoder(r.Body).Decode(d)
}
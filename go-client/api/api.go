package api

import (
	"go-client/ws"
	"log"
	"net/http"
	"time"
)

func GetWxUserList(w http.ResponseWriter, r *http.Request) {
	client := ws.GetWxClient()
	client.GetWxUserList()
	w.WriteHeader(http.StatusOK)
	time.Sleep(5 * time.Second)
	jsonWrite(w, client.Users)
}

func PostTxtMessage(w http.ResponseWriter, r *http.Request) {
	params := &ws.Params{}
	client := ws.GetWxClient()
	log.Println(r.Header.Get("Content-Type"))
	if r.Header.Get("Content-Type") == "application/json" {

		if err := jsonBind(r, params); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			jsonWrite(w, err.Error())
			return
		}
	} else {
		r.ParseForm()
		params.Content = r.Form.Get("content")
		params.Wxid = r.Form.Get("wxid")
	}
	client.PostTxtMessage(params.Content, params.Wxid)
	w.WriteHeader(http.StatusOK)
	jsonWrite(w, params)
}

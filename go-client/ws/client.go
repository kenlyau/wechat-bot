package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type Client struct {
	Url  url.URL
	Conn *websocket.Conn
}

func NewClient(addr string) *Client {
	log.Println("new ws client")
	u := url.URL{Scheme: "ws", Host: addr, Path: ""}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("websocket client error:", err)
	}
	return &Client{
		Url:  u,
		Conn: conn,
	}
}

func (c *Client) Start() {

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.Conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv %s", message)
		}
	}()
}

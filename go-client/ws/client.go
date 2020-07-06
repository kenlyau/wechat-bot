package ws

import (
	"encoding/json"
	"go-client/config"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var client *Client

type Params struct {
	Id      string `json:"id"`
	Type    int    `json:"type"`
	Content string `json:"content"`
	Wxid    string `json:"wxid"`
}

type User struct {
	Name string `json:"name"`
	Wxid string `json:"wxid"`
	Type string `json:"type"`
}

type Message struct {
	Id      string      `json:"id"`
	Content interface{} `json:"content"`
	Sender  string      `json:"sender"`
	Srvid   int         `json:"srvid"`
	Time    string      `json:"time"`
	Type    int         `json:"type"`
}

func (p *Params) Json() []byte {
	p.Id = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	j, _ := json.Marshal(p)
	return j
}

type Client struct {
	Url   url.URL
	Conn  *websocket.Conn
	Users []User
}

func SetUp() {
	log.Println("new ws client")
	u := url.URL{Scheme: "ws", Host: config.Config.DllServer, Path: ""}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("websocket client error:", err)
	}
	client = &Client{
		Url:  u,
		Conn: conn,
	}
	//初始化发起获取用户列表请求
	GetWxUserList()

}

func RecvLog() {
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, msg, err := client.Conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv %s", msg)
			message := &Message{}
			err = json.Unmarshal(msg, message)
			if err != nil {
				log.Println("json parse message error:", err)
				return
			}

			SwitchMessage(message)
		}
	}()
}

func SwitchMessage(message *Message) {
	log.Printf("%+v", message)
	switch message.Type {
	case GET_USER_LIST_SUCCESS:
		//强转slice
		sli := reflect.ValueOf(message.Content)
		users := make([]User, 0)
		for i := 0; i < sli.Len(); i++ {
			j, _ := json.Marshal(sli.Index(i).Interface())
			user := &User{}
			json.Unmarshal(j, user)
			users = append(users, *user)
		}
		client.Users = users
		log.Printf("%+v", users)
	case RECV_TEXT_MSG:
		if message.Sender == "newall" && message.Content == "Ping" {
			PostTxtMessage("success", "newall")
		}
	}
}

func Close() {
	client.Conn.Close()
}

func WriteMessage(p *Params) {
	err := client.Conn.WriteMessage(websocket.TextMessage, p.Json())
	log.Println(string(p.Json()))
	if err != nil {
		log.Println("websocket write message error:", err)
	} else {
		log.Printf("websocket write message success")
	}
}
func GetWxUserList() {
	params := &Params{
		Type:    USER_LIST,
		Content: "user list",
		Wxid:    "null",
	}
	WriteMessage(params)
}

func PostTxtMessage(content string, wxid string) {
	params := &Params{
		Type:    TXT_MSG,
		Content: content,
		Wxid:    wxid,
	}
	WriteMessage(params)
}

// func (c *Client) Start() {

// 	done := make(chan struct{})

// 	go func() {
// 		defer close(done)
// 		for {
// 			_, message, err := c.Conn.ReadMessage()
// 			if err != nil {
// 				log.Println("read:", err)
// 				return
// 			}
// 			log.Printf("recv %s", message)
// 		}
// 	}()
// }

// func (c *Client) GetWxUserList() {
// 	p := Params{Id: time.Now().UnixNano(), Type: USER_LIST, Content: "user list", Wxid: "null"}
// 	c.Conn.WriteJSON(p)
// }

package ws

import (
	"encoding/json"
	"fmt"
	"go-client/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var client *Client

func GetWxClient() *Client {
	return client
}

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

func (c *Client) UpdateUsers(message *Message) {
	sytemWxids := "@floatbottle@medianote@weixin@newsapp@fmessage"
	//强转slice
	sli := reflect.ValueOf(message.Content)
	users := make([]User, 0)
	for i := 0; i < sli.Len(); i++ {
		j, _ := json.Marshal(sli.Index(i).Interface())
		user := &User{}
		json.Unmarshal(j, user)
		if strings.Contains(user.Wxid, "@chatroom") {
			user.Type = "chatroom"
		} else if strings.Contains(sytemWxids, "@"+user.Wxid) {
			user.Type = "system"
		} else if strings.HasPrefix(user.Wxid, "gh_") {
			user.Type = "mp"
		} else {
			user.Type = "user"
		}
		users = append(users, *user)
	}
	c.Users = users
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
	client.GetWxUserList()

}

func RecvLog() {
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, msg, err := client.Conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				continue
			}
			message := &Message{}
			err = json.Unmarshal(msg, message)

			if err != nil {
				log.Println("json parse message error:", err)
				continue
			}
			if message.Type == HEART_BEAT {
				continue
			}
			SwitchMessage(message)
		}
	}()
}

func SwitchMessage(message *Message) {
	log.Printf("%+v", message)
	commands := config.GetCommands()
	switch message.Type {
	case GET_USER_LIST_SUCCESS:
		client.UpdateUsers(message)
	case RECV_TEXT_MSG:
		log.Println(commands)
		log.Println(commands[message.Sender])
		if commands[message.Sender] != nil {

			execUserCommands(message.Sender, message)
		}
		// if message.Sender == "newall" && message.Content == "Ping" {
		// 	client.PostTxtMessage("success", "newall")
		// }
	}
}

func execUserCommands(user string, message *Message) {
	commands := config.GetCommands()
	myCommands := commands[user]
	for k, v := range myCommands {
		if strings.HasPrefix(message.Content.(string), "#"+k) {
			switch v.Classify {
			case "template":
				execTemplateCommand(user, message, v)
				continue
			case "hook":
				execHookCommand(user, message, v)
				continue
			}
			break
		}
	}
}

func execTemplateCommand(user string, message *Message, command config.Command) {
	templates := config.GetTemplates()
	templateString := "%s"
	if templates[command.Variate] != "" {
		templateString = templates[command.Variate]

	}
	msg := fmt.Sprintf(templateString, "success")
	client.PostTxtMessage(msg, user)
}

func execHookCommand(user string, message *Message, command config.Command) {

	str := fmt.Sprintf(
		"user=%s&id=%s&content=%s&sender=%s&srvid=%s&time=%s",
		user, message.Id, message.Content, message.Sender, message.Srvid, message.Time,
	)
	res, err := http.Post(command.Variate, "application/x-www-form-urlencoded", strings.NewReader(str))
	if err != nil {
		log.Println("exec hook command error:", err)
		return
	}
	log.Println("res:", res)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("exec hook command response error:", err)
		return
	}
	client.PostTxtMessage(string(body), user)
}

func Close() {
	client.Conn.Close()
}

func (c *Client) WriteMessage(p *Params) {
	if len(p.Content) > 1024 {
		log.Println("消息文字最大1024", p.Id)
		return
	}
	err := c.Conn.WriteMessage(websocket.TextMessage, p.Json())
	log.Println(string(p.Json()))
	if err != nil {
		log.Println("websocket write message error:", err)
	} else {
		log.Printf("websocket write message success")
	}
}
func (c *Client) GetWxUserList() {
	params := &Params{
		Type:    USER_LIST,
		Content: "user list",
		Wxid:    "null",
	}
	c.WriteMessage(params)
}

func (c *Client) PostTxtMessage(content string, wxid string) {
	params := &Params{
		Type:    TXT_MSG,
		Content: content,
		Wxid:    wxid,
	}
	c.WriteMessage(params)
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

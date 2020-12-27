package api

import (
	"encoding/json"
	"fmt"
	user "go-redis-chatroom/user"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type commandMsg struct {
	Content string `json:"content,omitempty"`
	Channel string `json:"channel,omitempty"`
	Command string `json:"command,omitempty"`
	Err     string `json:"err,omitempty"`
}

// WsHandler used to handle the websocket request
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	u, err := Connect(r, conn)
	if err != nil {
		http.Error(w, "Connect to channel error", 400)
		return
	}
	closeCh := Disconnect(r, conn, u)

	for {
		select {
		case <-closeCh:
			return
		default:
			onUserMessage(u, conn, r)
		}
	}
}

// TestHandler is to test api
func TestHandler(w http.ResponseWriter, r *http.Request) {
	content := map[string]string{
		"key": "value",
	}
	json.NewEncoder(w).Encode(content)
}

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// Connect is to add user to channel
func Connect(r *http.Request, conn *websocket.Conn) (*user.User, error) {
	params := r.URL.Query()
	var username string
	if params["username"] != nil {
		username = params["username"][0]
	} else {
		username = "visitor_" + randomString(6)
	}
	fmt.Println("connected from:", conn.RemoteAddr(), "user:", username)
	u, err := user.Connect(username)
	if err != nil {
		return nil, err
	}

	onChannelMessage(conn, r, u)
	return u, nil
}

func onChannelMessage(conn *websocket.Conn, r *http.Request, u *user.User) {
	go func() {
		for m := range u.MessageChan {
			if err := conn.WriteJSON(m); err != nil {
				fmt.Println(err)
				break
			}
		}

	}()
}

func onUserMessage(u *user.User, conn *websocket.Conn, r *http.Request) {

	var commandMsg commandMsg
	if err := conn.ReadJSON(&commandMsg); err != nil {
		conn.WriteJSON(err.Error())
		return
	}

	// username := r.URL.Query()["username"][0]
	switch commandMsg.Command {
	// case "subscribe":
	// 	if err := u.Subscribe(rdb, msg.Channel); err != nil {
	// 		conn.WriteJSON(err.Error())
	// 	}
	// case "unsubscribe":
	// 	if err := u.Unsubscribe(rdb, msg.Channel); err != nil {
	// 		conn.WriteJSON(err.Error())
	// 	}
	case "chat":
		if err := user.Chat(commandMsg.Channel, commandMsg.Content, u); err != nil {
			conn.WriteJSON(err.Error())
		}
	}
}

// Disconnect close the channels
func Disconnect(r *http.Request, conn *websocket.Conn, u *user.User) chan struct{} {

	closeCh := make(chan struct{})
	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("connection closed for user", u)
		if err := u.Disconnect(); err != nil {
			return err
		}
		close(closeCh)
		return nil
	})
	return closeCh
}

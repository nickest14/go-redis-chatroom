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
var connectedUsers = make(map[string]*user.User)

// WsHandler used to handle the websocket request
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Println(conn)
	err = Connect(r, conn)
	if err != nil {
		http.Error(w, "Connect to channel error", 400)
		return
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
func Connect(r *http.Request, conn *websocket.Conn) error {
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
		return err
	}
	connectedUsers[username] = u
	return nil
}

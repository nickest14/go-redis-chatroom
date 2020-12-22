package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go-redis-chatroom/redis"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// WsHandler used to handle the websocket request
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Println(conn)
	ctx := context.Background()
	rdb := redis.Client
	fmt.Println(rdb.Get(ctx, "test"))
}

// TestHandler is to test api
func TestHandler(w http.ResponseWriter, r *http.Request) {
	content := map[string]string{
		"key": "value",
	}
	json.NewEncoder(w).Encode(content)
}

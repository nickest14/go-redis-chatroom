package main

import (
	api "go-redis-chatroom/api"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Path("/chat").Methods("GET").HandlerFunc(api.WsHandler)
	r.Path("/users").Methods("GET").HandlerFunc(api.UsersHandler)
	r.Path("/user/{user}/channels").Methods("GET").HandlerFunc(api.UserChannelsHandler)
	r.Path("/test").Methods("GET").HandlerFunc(api.TestHandler)
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":9000"
	}
	http.ListenAndServe(port, r)
}

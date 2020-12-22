package main

import (
	api "go-redis-chatroom/api"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Path("/chat").Methods("GET").HandlerFunc(api.WsHandler)
	r.Path("/test").Methods("GET").HandlerFunc(api.TestHandler)
	http.ListenAndServe(":9000", r)
}

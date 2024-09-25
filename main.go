package main

import (
	"fmt"
	"log"
	"net/http"

	"socket-v1/src/http/middlewares"
	"socket-v1/src/http/routes"
	"socket-v1/src/infra/database"
	messageQueue "socket-v1/src/infra/message-queue"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()

	conn, ch := messageQueue.ConnectMQ()
	defer conn.Close()
	defer ch.Close()

	r := mux.NewRouter()
	routes.RegisterWebsocketRoute(r)

	wrapMiddleware := middlewares.WrapRequest()
	wrapRoute := wrapMiddleware(r)
	handler := middlewares.Cors(wrapRoute)

	PORT := 3000

	log.Printf("Start server successfully on port %d\n", PORT)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), handler)
}
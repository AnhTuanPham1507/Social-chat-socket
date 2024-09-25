package routes

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"

	messageQueue "socket-v1/src/infra/message-queue"
	"socket-v1/src/services/websocket"
)

var RegisterWebsocketRoute = func(router *mux.Router) {
	pool := websocket.NewPool()
	go pool.Start()

	sb := router.PathPrefix("/v1").Subrouter()
	sb.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(pool, w, r)
	})
}

func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r)
	messageChannel := make(chan []byte)

	if(err != nil) {
		log.Println("Error: ", err)
	}
	
	randomID := rand.Intn(100)
	client := &websocket.Client{
		ID:     uint(randomID),
		Connection: conn,
		Pool:       pool,
		Name:      "Phạm Anh Tuấn",
	}
	go client.InitConnection(messageChannel)

	// lắng nghe và đẩy sự kiện từ queue
	br := messageQueue.GetRabbitMQBroker()
	go br.PublishMessage(messageChannel)
	go br.ReadMessages(pool)
}


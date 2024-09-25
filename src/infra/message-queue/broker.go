package messageQueue

import (
	"context"
	"log"
	"time"

	"socket-v1/src/utils"

	amqp "github.com/rabbitmq/amqp091-go"

	"socket-v1/src/services/websocket"
)

type Broker struct {
	Queue amqp.Queue
	Channel *amqp.Channel
}

type QueueResponse struct {
	RoomID uint `json:"RoomID"`
	Message string `json:"Message"`
	Owner string `json:"Owner"`
}

func (b *Broker) SetUp(ch *amqp.Channel) {
	QUEUE_NAME := "CHAT"
	
	q1, err := ch.QueueDeclare(
		QUEUE_NAME, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if(err != nil){
		log.Fatalf("Failed to declare queue: %s", err)
	}

	b.Queue = q1
	b.Channel = ch
}

func (b *Broker) PublishMessage(requestBody chan []byte) {
	FIVE_SECONDS := 5 * time.Second

	for body := range requestBody {	
		ctx, cancel := context.WithTimeout(context.Background(), FIVE_SECONDS)
		
		err := b.Channel.PublishWithContext(ctx,
			"",
			b.Queue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body: body,
			},
		)

		cancel()

		if err != nil {
			log.Printf("PublishMessage Error occured %s\n", err)
			continue
		}
		
		log.Printf(" [x] Sent %s\n", body)
	}
}

func (b *Broker) ReadMessages(pool *websocket.Pool) {
	msgs, err := b.Channel.Consume(
		b.Queue.Name, // queue name
        "",       // consumer tag
        false,     // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,      // args
	)
	if  err != nil {
		log.Printf("ReadMessages Error occured %s\n", err)
		return
	}

	rsvdMsgs := make(chan QueueResponse)
	go messageTransformer(msgs, rsvdMsgs)
	go processResponse(rsvdMsgs, pool)
	
	select {}
}

func messageTransformer(entries <-chan amqp.Delivery, receivedMessages chan QueueResponse ){
	var response QueueResponse
	for d := range entries {
		err := utils.ParseByteArray(d.Body, &response)

		if err != nil {
			// log.Printf("ParseByteArray Error occured %s\n", err)
            continue
		}
		
		receivedMessages <- response
	}
}

func processResponse(s <-chan QueueResponse, pool *websocket.Pool) {
	for r := range s {
		log.Println("processing stock response for ", r)
		
		message := websocket.Message{
			RoomID:  int32(r.RoomID),
			Message: r.Message,
			Owner: r.Owner,
		}
		
		pool.Broadcast <- message
	}
}
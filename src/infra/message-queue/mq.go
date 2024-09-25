package messageQueue

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var br Broker 

func ConnectMQ() (*amqp.Connection, *amqp.Channel) {
	rmqHost := os.Getenv("MQ_URI") 

	conn, _ := amqp.Dial(rmqHost)

	ch, _ := conn.Channel()

	br.SetUp(ch)

	log.Printf("Connected to rabbitmq server")
	
	return conn, ch
}

func GetRabbitMQBroker() *Broker {
	return &br
}
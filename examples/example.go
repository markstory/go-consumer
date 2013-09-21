package main

import (
	"consumer"
	"flag"
	"log"
	"github.com/streadway/amqp"
)

func main() {
	var ini = flag.String("config", "./test.ini", "The configuration file to use.")

	c, err := consumer.Create(*ini)
	if err != nil {
		log.Fatalf("Unable to create consumer: Error: %v", err)
	}
	c.Consume(worker)
}

// Example worker function
// Should eventually get type coming out of consumer
// instead of an amqp.Delivery value.
func worker (msg amqp.Delivery) {
	log.Printf("msg received")
	msg.Ack(true)
}
